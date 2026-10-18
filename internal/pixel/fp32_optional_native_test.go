//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"bytes"
	"fmt"
	"os"
	"simd/archsimd"
	"testing"
	"unsafe"
)

func TestFP32PublicHost(t *testing.T) {
	t.Logf("FP32PublicHost AVX512=%v AVX2=%v Enabled=%v", archsimd.X86.AVX512(), archsimd.X86.AVX2(), Enabled)
	if os.Getenv("GOVPX_REQUIRE_FP32_NATIVE") == "1" && (!Enabled || !archsimd.X86.AVX2() || !archsimd.X86.AVX512()) {
		t.Fatal("actual native AVX512 SIMD backend required")
	}
}

func fp32OptionalWindow(size, offset int) ([]byte, []byte) {
	p := bytes.Repeat([]byte{0xa7}, size+192)
	base := uintptr(unsafe.Pointer(unsafe.SliceData(p)))
	start := 64 + int((-base)&63) + offset
	return p, p[start : start+size : start+size]
}

// Independent fixed-point equations and scalar energy thresholds; no production
// FP, token, clipping or vector helpers. Optional q is deliberately supported.
func fp32OptionalExpected(c, scan, q, dq []int16, cl []byte, r, qt, d [2]int16) int {
	eob := 0
	for i, v := range c {
		k := 1
		if i == 0 {
			k = 0
		}
		a := int(v)
		if a < 0 {
			a = -a
		}
		a += int(r[k])
		mag := 0
		if a >= int(d[k]) {
			if a > 32767 {
				a = 32767
			}
			mag = a * int(qt[k]) / 65536
		}
		signed := mag
		if v < 0 {
			signed = -signed
		}
		if q != nil {
			q[i] = int16(signed)
		}
		dq[i] = int16(signed * int(d[k]))
		if signed != 0 && int(scan[i]) > eob {
			eob = int(scan[i])
		}
		if cl != nil {
			value := byte(0)
			for _, bound := range []int{0, 1, 2, 4, 10} {
				if mag > bound {
					value++
				}
			}
			cl[i] = value
		}
	}
	return eob
}

var fp32OptionalSink uint64

// This measures the real public pixel API, not an encoder qnil caller. The
// encoder APIs require qcoeff and are separately covered by unchanged fixtures.
func BenchmarkFP32PublicOptionalNative(b *testing.B) {
	for _, n := range []int{64, 256} {
		for _, qnil := range []bool{false, true} {
			for _, classes := range []bool{false, true} {
				for _, pattern := range []string{"zero", "sparse", "dense"} {
					for _, unaligned := range []bool{false, true} {
						alignment := "aligned64"
						if unaligned {
							alignment = "offset2-classes1"
						}
						b.Run(fmt.Sprintf("n%d/qnil%v/classes%v/%s/%s", n, qnil, classes, pattern, alignment), func(b *testing.B) {
							offset, coff := 0, 0
							if unaligned {
								offset, coff = 2, 1
							}
							var backing, want [5][]byte
							var off [5]int
							var views [5][]byte
							for k := 0; k < 5; k++ {
								size := 2 * n
								align := offset
								if k == 4 {
									size = n
									align = coff
								}
								backing[k], views[k] = fp32OptionalWindow(size, align)
								off[k] = int(uintptr(unsafe.Pointer(unsafe.SliceData(views[k]))) - uintptr(unsafe.Pointer(unsafe.SliceData(backing[k]))))
								if uintptr(unsafe.Pointer(unsafe.SliceData(views[k])))&63 != uintptr(align) {
									b.Fatal("fixture alignment")
								}
							}
							word := func(v []byte) []int16 { return unsafe.Slice((*int16)(unsafe.Pointer(unsafe.SliceData(v))), n) }
							c, s, q, dq, cl := word(views[0]), word(views[1]), word(views[2]), word(views[3]), views[4]
							for i := range c {
								if pattern == "zero" {
									c[i] = 0
								} else if pattern == "dense" {
									c[i] = int16((i*347)%65536 - 32768)
								} else if i%11 == 0 {
									c[i] = int16((i*271)%1024 - 512)
								} else {
									c[i] = 0
								}
								s[i] = int16((i*37)%n + 1)
							}
							for k := range backing {
								want[k] = bytes.Clone(backing[k])
							}
							wc, ws, wq, wd, wcl := word(want[0][off[0]:off[0]+2*n]), word(want[1][off[1]:off[1]+2*n]), word(want[2][off[2]:off[2]+2*n]), word(want[3][off[3]:off[3]+2*n]), want[4][off[4]:off[4]+n]
							if qnil {
								q, wq = nil, nil
							}
							if !classes {
								cl, wcl = nil, nil
							}
							r, qt, dequant := [2]int16{8, 13}, [2]int16{2849, 1598}, [2]int16{23, 41}
							expected := fp32OptionalExpected(wc, ws, wq, wd, wcl, r, qt, dequant)
							check := func(e int) {
								b.Helper()
								if e != expected {
									b.Fatal("EOB", e, expected)
								}
								for k := range backing {
									if !bytes.Equal(backing[k], want[k]) {
										b.Fatal("full input/output/canary backing", k)
									}
								}
							}
							eob := QuantizeFPBlock(c, s, q, dq, n, r, qt, dequant, cl)
							check(eob)
							b.ReportAllocs()
							for b.Loop() {
								eob = QuantizeFPBlock(c, s, q, dq, n, r, qt, dequant, cl)
							}
							check(eob)
							sum := uint64(eob)
							for _, v := range dq {
								sum += uint64(uint16(v))
							}
							if q != nil {
								for _, v := range q {
									sum += uint64(uint16(v))
								}
							}
							if cl != nil {
								for _, v := range cl {
									sum += uint64(v)
								}
							}
							fp32OptionalSink = sum
						})
					}
				}
			}
		}
	}
}
