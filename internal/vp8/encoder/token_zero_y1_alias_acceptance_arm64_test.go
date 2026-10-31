//go:build goexperiment.simd && arm64 && !purego && !govpx_oracle_trace

package encoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/gtsteffaniak/go-vpx/internal/vp8/tables"
)

// Arena views remain inside one owned allocation. The helper must decline before
// any store for overlap cases: the legacy fallback is deliberately not executed
// with a payload intersecting its own live slice header.
func TestTokenZeroY1DirectAliasWindows(t *testing.T) {
	const prefix = 3
	var sample coefficientTokenCounts
	counterOffset := int(uintptr(unsafe.Pointer(&sample[3][0][0][tables.DCTEOBToken])) - uintptr(unsafe.Pointer(&sample)))
	sizes := map[string]int{"payload": 64, "header": 24, "counter": int(unsafe.Sizeof(int(0))), "coeff": int(unsafe.Sizeof(MacroblockCoefficients{})), "above": int(unsafe.Sizeof(TokenContextPlanes{})), "left": int(unsafe.Sizeof(TokenContextPlanes{}))}
	type pair struct{ a, b string }
	pairs := []pair{{"payload", "header"}, {"payload", "counter"}, {"header", "counter"}}
	for _, w := range []string{"payload", "header", "counter"} {
		for _, r := range []string{"coeff", "above", "left"} {
			pairs = append(pairs, pair{w, r})
		}
	}
	checks, declines, handled := 0, 0, 0
	for _, p := range pairs {
		for direction := range 2 {
			for _, overlap := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s-%s/order%d/overlap%v", p.a, p.b, direction, overlap), func(t *testing.T) {
					arena := make([]uint64, 16384)
					base := unsafe.Pointer(&arena[0])
					raw := unsafe.Slice((*byte)(base), len(arena)*8)
					pos := map[string]int{"payload": 8192, "header": 16384, "counter": 32768, "coeff": 49152, "above": 65536, "left": 81920}
					// Keep the slice header naturally aligned. Other pointer-free views permit
					// exact byte-boundary probes on ARM; no pointer arithmetic leaves arena.
					gap := 0
					if overlap {
						gap = -1
					}
					moving, fixed := p.a, p.b
					dir := direction
					if moving == "header" {
						moving, fixed = fixed, moving
						dir = 1 - dir
					}
					if dir == 0 {
						pos[moving] = pos[fixed] + sizes[fixed] + gap
					} else {
						pos[moving] = pos[fixed] - sizes[moving] - gap
					}
					ptr := func(name string) unsafe.Pointer { return unsafe.Add(base, pos[name]) }
					coeff := (*MacroblockCoefficients)(ptr("coeff"))
					above := (*TokenContextPlanes)(ptr("above"))
					left := (*TokenContextPlanes)(ptr("left"))
					records := (*InterCoefficientTokenRecords)(ptr("header"))
					counts := (*coefficientTokenCounts)(unsafe.Add(base, pos["counter"]-counterOffset))
					storage := unsafe.Slice((*CoefficientTokenRecord)(unsafe.Add(ptr("payload"), -prefix*4)), prefix+16)
					records.Records = storage[: prefix : prefix+16]
					// Assert the exact earlier predicates, so an EOB/context/capacity failure
					// cannot masquerade as an alias-guard test.
					if binary.LittleEndian.Uint32(above.Y1[:])|binary.LittleEndian.Uint32(left.Y1[:]) != 0 {
						t.Fatal("setup failed zero-context eligibility")
					}
					if cap(records.Records)-len(records.Records) != 16 {
						t.Fatal("setup failed capacity eligibility")
					}
					for _, e := range coeff.EOB[:16] {
						if e != 0 {
							t.Fatal("setup failed EOB eligibility")
						}
					}
					targetOverlap := func(a int, an int, b int, bn int) bool { return a < b+bn && b < a+an }
					if targetOverlap(pos[p.a], sizes[p.a], pos[p.b], sizes[p.b]) != overlap {
						t.Fatal("target interval setup")
					}
					before := bytes.Clone(raw)
					want := bytes.Clone(before)
					if !overlap {
						var one coefficientTokenCounts
						reference := InterCoefficientTokenRecords{Records: make([]CoefficientTokenRecord, 0, 1)}
						var qc [16]int16
						if err := countBlockCoefficientTokensAndRecords(&one, &reference, 3, 0, 0, &qc, 0); err != nil {
							t.Fatal(err)
						}
						if len(reference.Records) != 1 {
							t.Fatal("frozen scalar EOB record count")
						}
						for j := range 16 {
							binary.LittleEndian.PutUint32(want[pos["payload"]+j*4:], uint32(reference.Records[0]))
						}
						binary.LittleEndian.PutUint64(want[pos["header"]+8:], uint64(prefix+16))
						initial := *(*int)(ptr("counter"))
						expected := initial + 16
						binary.LittleEndian.PutUint64(want[pos["counter"]:], uint64(expected))
					}
					got := zeroY1TokenRecords(above, left, coeff, counts, records, 3, 0)
					if got == overlap {
						t.Fatalf("handled=%v, target overlap=%v", got, overlap)
					}
					checks++
					if !got {
						declines++
						if !bytes.Equal(raw, before) {
							t.Fatal("decline modified complete raw backing")
						}
					} else {
						handled++
						if !bytes.Equal(raw, want) {
							t.Fatal("handled adjacent case differs from independent scalar record/full-backing expectation")
						}
					}
					runtime.KeepAlive(arena)
				})
			}
		}
	}
	if checks != 48 || declines != 24 || handled != 24 {
		t.Fatalf("checks=%d decline=%d handled=%d", checks, declines, handled)
	}
	t.Logf("direct alias boundary cases=%d declined=%d adjacent handled=%d", checks, declines, handled)
}

// Frozen original Y1 loop, copied from tokenCountFrozenMacro. It calls the
// unchanged scalar block/token API, never the new macro helper.
func tokenZeroY1FrozenWalk(above, left *TokenContextPlanes, coeffs *MacroblockCoefficients, counts *coefficientTokenCounts, records *InterCoefficientTokenRecords, bt, skip int) error {
	for block := range 16 {
		eob := coeffs.BlockEOB(block, skip)
		a, l := block&3, (block&0x0c)>>2
		ctx := int(above.Y1[a] + left.Y1[l])
		if ctx >= tables.PrevCoefContexts {
			return ErrInvalidPacketConfig
		}
		if err := countBlockCoefficientTokensAndRecords(counts, records, bt, ctx, skip, &coeffs.QCoeff[block], eob); err != nil {
			return err
		}
		has := uint8(0)
		if eob > skip {
			has = 1
		}
		above.Y1[a], left.Y1[l] = has, has
	}
	return nil
}

func TestTokenZeroY1CapacityCounterFrozen(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	checks := 0
	for _, remaining := range []int{15, 16, 17} {
		for _, initial := range []int{maxInt, maxInt - 15, maxInt - 16, minInt, minInt + 15, -1, 0, 1} {
			for _, is4 := range []bool{false, true} {
				for _, sameContext := range []bool{false, true} {
					bt, skip := 0, 1
					if is4 {
						bt, skip = 3, 0
					}
					var ca, cb MacroblockCoefficients
					for i := range 16 {
						ca.EOB[i] = uint8(skip)
						for k := range 16 {
							ca.QCoeff[i][k] = int16(37*i - 19*k)
						}
					}
					cb = ca
					aa := TokenContextPlanes{U: [2]byte{1, 2}, V: [2]byte{3, 4}, Y2: 1}
					la := TokenContextPlanes{U: [2]byte{5, 6}, V: [2]byte{7, 8}, Y2: 2}
					ab, lb := aa, la
					ap, lp, bp, mp := &aa, &la, &ab, &lb
					if sameContext {
						lp, mp = ap, bp
					}
					var ac, bc coefficientTokenCounts
					ac[bt][skip][0][tables.DCTEOBToken] = initial
					bc = ac
					ba, bb := make([]CoefficientTokenRecord, 40), make([]CoefficientTokenRecord, 40)
					for i := range ba {
						ba[i] = CoefficientTokenRecord(0x8bad0000 + uint32(i)*71)
					}
					copy(bb, ba)
					ra := InterCoefficientTokenRecords{Records: ba[2 : 5 : 5+remaining], Rows: 7}
					rb := InterCoefficientTokenRecords{Records: bb[2 : 5 : 5+remaining], Rows: 7}
					ra.rowStartsInline[17], rb.rowStartsInline[17] = 123, 123
					wantErr, wantPanic := tokenCountMacroCapture(func() error { return tokenZeroY1FrozenWalk(bp, mp, &cb, &bc, &rb, bt, skip) })
					used := false
					gotErr, gotPanic := tokenCountMacroCapture(func() error {
						used = zeroY1TokenRecords(ap, lp, &ca, &ac, &ra, bt, skip)
						if used {
							return nil
						}
						return tokenZeroY1FrozenWalk(ap, lp, &ca, &ac, &ra, bt, skip)
					})
					if used != (remaining >= 16) {
						t.Fatalf("remaining%d helper handled%v", remaining, used)
					}
					checks++
					if gotErr != wantErr || gotPanic != wantPanic || aa != ab || la != lb || ca != cb || ac != bc || !reflect.DeepEqual(ra, rb) || !reflect.DeepEqual(ba, bb) {
						t.Fatalf("remaining%d count%d is4%v contextsAlias%v err%q/%q panic%q/%q", remaining, initial, is4, sameContext, gotErr, wantErr, gotPanic, wantPanic)
					}
				}
			}
		}
	}
	if checks != 96 {
		t.Fatalf("checks=%d", checks)
	}
	t.Logf("frozen nonempty-prefix/capacity/counter-wrap cases=%d", checks)
}
