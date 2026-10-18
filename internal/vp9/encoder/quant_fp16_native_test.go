package encoder

import (
	"bytes"
	"fmt"
	"testing"
	"unsafe"

	"github.com/thesyncim/govpx/internal/vp9/common"
)

var fp16NativeSink uint64

// fp16NativeWindow aligns the actual data address, retaining both canaries.
// Int16 windows use offset 2; byte classes use offset 1 in the unaligned mode.
func fp16NativeWindow(n, offset int) ([]byte, []byte) {
	backing := bytes.Repeat([]byte{0xa7}, n+128)
	base := uintptr(unsafe.Pointer(unsafe.SliceData(backing)))
	start := 32 + int((-base)&31) + offset
	return backing, backing[start : start+n : start+n]
}

// fp16NativeReference is an independent integer implementation of the pinned
// libvpx v1.16.0 FP operation and energy-class thresholds. It calls no SIMD,
// production quantization, token mapping, or clipping helpers.
func fp16NativeReference(coeff, iscan, qo, dq []int16, cl []byte, round, quant, dequant [2]int16) int {
	eob := 0
	for i, c16 := range coeff {
		k := 1
		if i == 0 {
			k = 0
		}
		c := int(c16)
		magnitude := c
		if magnitude < 0 {
			magnitude = -magnitude
		}
		rounded := magnitude + int(round[k])
		value := 0
		if rounded >= int(dequant[k]) {
			if rounded > 32767 {
				rounded = 32767
			}
			value = rounded * int(quant[k]) / 65536
		}
		signed := value
		if c < 0 {
			signed = -signed
		}
		qo[i] = int16(signed)
		dq[i] = int16(signed * int(dequant[k]))
		if value != 0 && int(iscan[i]) > eob {
			eob = int(iscan[i])
		}
		if cl != nil {
			class := byte(0)
			for _, bound := range [...]int{0, 1, 2, 4, 10} {
				if value > bound {
					class++
				}
			}
			cl[i] = class
		}
	}
	return eob
}

// BenchmarkVP9FP16Native overlays unchanged on assembly, parent and current.
// Only coeff/q/dq/classes alignment is controlled; the real immutable
// ScanOrder's IScan is used by the actual public encoder operation.
func BenchmarkVP9FP16Native(b *testing.B) {
	for _, size := range []struct {
		n  int
		tx common.TxSize
	}{{16, common.Tx4x4}, {64, common.Tx8x8}, {256, common.Tx16x16}} {
		for _, classesOn := range []bool{false, true} {
			for _, pattern := range []string{"zero", "sparse", "dense"} {
				for _, unaligned := range []bool{false, true} {
					mode, alignment := "no-classes", "aligned32"
					if classesOn {
						mode = "classes"
					}
					if unaligned {
						alignment = "offset2-classes1"
					}
					b.Run(fmt.Sprintf("n%d/%s/%s/%s", size.n, mode, pattern, alignment), func(b *testing.B) {
						n := size.n
						offset := 0
						classOffset := 0
						if unaligned {
							offset = 2
							classOffset = 1
						}
						cb, cbytes := fp16NativeWindow(2*n, offset)
						qb, qbytes := fp16NativeWindow(2*n, offset)
						db, dbytes := fp16NativeWindow(2*n, offset)
						lb, classes := fp16NativeWindow(n, classOffset)
						coeff := unsafe.Slice((*int16)(unsafe.Pointer(unsafe.SliceData(cbytes))), n)
						qo := unsafe.Slice((*int16)(unsafe.Pointer(unsafe.SliceData(qbytes))), n)
						dq := unsafe.Slice((*int16)(unsafe.Pointer(unsafe.SliceData(dbytes))), n)
						for i := range coeff {
							switch pattern {
							case "zero":
								coeff[i] = 0
							case "sparse":
								coeff[i] = 0
								if i%11 == 0 {
									coeff[i] = int16((i*271)%1024 - 512)
								}
							case "dense":
								coeff[i] = int16((i*347)%65536 - 32768)
							}
						}
						for label, p := range map[string]unsafe.Pointer{"coeff": unsafe.Pointer(unsafe.SliceData(coeff)), "q": unsafe.Pointer(unsafe.SliceData(qo)), "dq": unsafe.Pointer(unsafe.SliceData(dq)), "classes": unsafe.Pointer(unsafe.SliceData(classes))} {
							want := uintptr(offset)
							if label == "classes" {
								want = uintptr(classOffset)
							}
							if uintptr(p)&31 != want {
								b.Fatalf("%s alignment=%d want%d", label, uintptr(p)&31, want)
							}
						}
						scan := &common.ScanOrders[size.tx][common.DctDct]
						scanBefore, iscanBefore := append([]int16(nil), scan.Scan[:n]...), append([]int16(nil), scan.IScan[:n]...)
						b.Logf("alignment coeff/q/dq=%d classes=%d actual-iscan-mod32=%d", offset, classOffset, uintptr(unsafe.Pointer(&scan.IScan[0]))&31)
						dequant := [2]int16{23, 41}
						tables := QuantizeFPTablesForDequant(dequant)
						cbefore := bytes.Clone(cb)
						qw, dw, lw := bytes.Clone(qb), bytes.Clone(db), bytes.Clone(lb)
						qOffset := int(uintptr(unsafe.Pointer(unsafe.SliceData(qbytes))) - uintptr(unsafe.Pointer(unsafe.SliceData(qb))))
						dOffset := int(uintptr(unsafe.Pointer(unsafe.SliceData(dbytes))) - uintptr(unsafe.Pointer(unsafe.SliceData(db))))
						lOffset := int(uintptr(unsafe.Pointer(unsafe.SliceData(classes))) - uintptr(unsafe.Pointer(unsafe.SliceData(lb))))
						qwant := unsafe.Slice((*int16)(unsafe.Pointer(&qw[qOffset])), n)
						dwant := unsafe.Slice((*int16)(unsafe.Pointer(&dw[dOffset])), n)
						var lwant []byte
						if classesOn {
							lwant = lw[lOffset : lOffset+n : lOffset+n]
						}
						wantEOB := fp16NativeReference(coeff, scan.IScan[:n], qwant, dwant, lwant, tables.RoundFP, tables.QuantFP, dequant)
						var eob int
						check := func() {
							b.Helper()
							if eob != wantEOB || !bytes.Equal(cb, cbefore) || !bytes.Equal(qb, qw) || !bytes.Equal(db, dw) || !bytes.Equal(lb, lw) {
								b.Fatalf("full backing/EOB mismatch got%d want%d", eob, wantEOB)
							}
							for i := 0; i < n; i++ {
								if scan.Scan[i] != scanBefore[i] || scan.IScan[i] != iscanBefore[i] {
									b.Fatal("ScanOrder mutated")
								}
							}
						}
						// Same actual production caller fallback as the existing composite fixture.
						fallbackClasses := func() {
							for i, c := range qo {
								m := int(c)
								if m < 0 {
									m = -m
								}
								token, _ := TokenForAbsCoeff(m)
								classes[i] = PtEnergyClass[token]
							}
						}
						if classesOn {
							var produced bool
							eob, produced = QuantizeFPWithQTablesScanOrderPtrTokenClasses(coeff, n, dequant, tables, scan, qo, dq, classes)
							if !produced {
								fallbackClasses()
							}
						} else {
							eob = QuantizeFPWithQTablesScanOrderPtrValidated(coeff, n, dequant, tables, scan, qo, dq)
						}
						check()
						b.ReportAllocs()
						for b.Loop() {
							if classesOn {
								var produced bool
								eob, produced = QuantizeFPWithQTablesScanOrderPtrTokenClasses(coeff, n, dequant, tables, scan, qo, dq, classes)
								if !produced {
									for i, c := range qo {
										m := int(c)
										if m < 0 {
											m = -m
										}
										token, _ := TokenForAbsCoeff(m)
										classes[i] = PtEnergyClass[token]
									}
								}
							} else {
								eob = QuantizeFPWithQTablesScanOrderPtrValidated(coeff, n, dequant, tables, scan, qo, dq)
							}
						}
						check()
						sum := uint64(eob)
						for i := range qo {
							sum += uint64(uint16(qo[i])) + uint64(uint16(dq[i]))
							if classesOn {
								sum += uint64(classes[i])
							}
						}
						fp16NativeSink = sum
					})
				}
			}
		}
	}
}
