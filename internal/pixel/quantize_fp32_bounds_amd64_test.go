//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"bytes"
	"testing"
	"unsafe"
)

func fp32WordView(p []byte, off, n int) []int16 {
	return unsafe.Slice((*int16)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(p)), off)), n)
}

func TestFP32EligibleExactWindows(t *testing.T) {
	const n = 256
	for align := 0; align < 64; align++ {
		for _, trail := range []int{0, 19} {
			for options := 0; options < 4; options++ {
				var got, want [5][]byte
				var offsets [5]int
				for k := 0; k < 5; k++ {
					size := 2 * n
					if k == 4 {
						size = n
					}
					got[k] = bytes.Repeat([]byte{0xa7}, 128+size+trail)
					base := uintptr(unsafe.Pointer(unsafe.SliceData(got[k]))) & 63
					desired := (align*(2*k+1) + k) % 64
					offsets[k] = (desired - int(base) + 64) % 64
					if uintptr(unsafe.Add(unsafe.Pointer(unsafe.SliceData(got[k])), offsets[k]))&63 != uintptr(desired) {
						t.Fatal("alignment fixture")
					}
				}
				c, s := fp32WordView(got[0], offsets[0], n), fp32WordView(got[1], offsets[1], n)
				for i := 0; i < n; i++ {
					c[i] = int16(i*771 + align*347 - 32768)
					s[i] = int16(i*127 - 16384)
				}
				for k := range got {
					want[k] = bytes.Clone(got[k])
				}
				qg, dg := fp32WordView(got[2], offsets[2], n), fp32WordView(got[3], offsets[3], n)
				qw, dw := fp32WordView(want[2], offsets[2], n), fp32WordView(want[3], offsets[3], n)
				cg, cw := got[4][offsets[4]:offsets[4]+n:offsets[4]+n], want[4][offsets[4]:offsets[4]+n:offsets[4]+n]
				if options&1 != 0 {
					qg, qw = nil, nil
				}
				if options&2 != 0 {
					cg, cw = nil, nil
				}
				if !fp16WindowsDisjoint(c, s, qg, dg, cg, n) {
					t.Fatal("eligible exact windows unexpectedly overlap")
				}
				r, q, d := [2]int16{32767, 17}, [2]int16{32767, 12345}, [2]int16{1, 23}
				we := fp16ScalarIndependent(fp32WordView(want[0], offsets[0], n), fp32WordView(want[1], offsets[1], n), qw, dw, cw, n, r, q, d)
				ge := QuantizeFPBlock(c, s, qg, dg, n, r, q, d, cg)
				if ge != we {
					t.Fatalf("align%d options%d EOB", align, options)
				}
				for k := range got {
					if !bytes.Equal(got[k], want[k]) {
						t.Fatalf("align%d trail%d options%d backing%d", align, trail, options, k)
					}
				}
			}
		}
	}
}

func TestFP32EligibleAliasedDeclines(t *testing.T) {
	const n = 256
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 5; j++ {
			for _, shift := range []int{0, 1, n - 1, 2*n - 1} {
				offsets := [5]int{64, 1088, 2112, 3136, 4160}
				offsets[j] = offsets[i] + shift
				got := make([]byte, 6144)
				for k := range got {
					got[k] = byte(k*37 + k/3)
				}
				want := bytes.Clone(got)
				c, s, q, d, cl := fp16Views(got, offsets, n)
				cl = cl[:n:n]
				if fp16WindowsDisjoint(c, s, q, d, cl, n) {
					t.Fatalf("pair%d/%d shift%d accepted overlap", i, j, shift)
				}
				wc, ws, wq, wd, wcl := fp16Views(want, offsets, n)
				wcl = wcl[:n:n]
				r, qt, dt := [2]int16{17, 32767}, [2]int16{12345, 32767}, [2]int16{23, 1}
				ge, gp := fp16Capture(func() int { return QuantizeFPBlock(c, s, q, d, n, r, qt, dt, cl) })
				we, wp := fp16Capture(func() int { return fp16OriginalPublic(wc, ws, wq, wd, n, r, qt, dt, wcl) })
				if ge != we || gp != wp || !bytes.Equal(got, want) {
					t.Fatalf("pair%d/%d shift%d changed frozen state/panic", i, j, shift)
				}
			}
		}
	}
}

func TestFP32EligiblePanicWindows(t *testing.T) {
	const n = 256
	for short := 0; short < 5; short++ {
		for _, length := range []int{0, 1, n - 1} {
			for _, clipped := range []bool{false, true} {
				initial := make([]byte, 6144)
				for k := range initial {
					initial[k] = byte(k*37 + 11)
				}
				got, want := bytes.Clone(initial), bytes.Clone(initial)
				offs := [5]int{64, 1088, 2112, 3136, 4160}
				run := func(storage []byte, old bool) (int, string) {
					c, s, q, d, cl := fp16Views(storage, offs, n)
					cl = cl[:n:n]
					slices := []*[]int16{&c, &s, &q, &d}
					if short < 4 {
						p := slices[short]
						if clipped {
							*p = (*p)[:length:length]
						} else {
							*p = (*p)[:length]
						}
					} else if clipped {
						cl = cl[:length:length]
					} else {
						cl = cl[:length]
					}
					return fp16Capture(func() int {
						if old {
							return fp16OriginalPublic(c, s, q, d, n, [2]int16{0, 17}, [2]int16{32767, 12345}, [2]int16{23, 41}, cl)
						}
						return QuantizeFPBlock(c, s, q, d, n, [2]int16{0, 17}, [2]int16{32767, 12345}, [2]int16{23, 41}, cl)
					})
				}
				ge, gp := run(got, false)
				we, wp := run(want, true)
				if gp == "" || ge != we || gp != wp || !bytes.Equal(got, want) {
					t.Fatalf("short%d length%d clipped%v panic%q/%q", short, length, clipped, gp, wp)
				}
			}
		}
	}
}
