//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// QuantizeVP8BatchOnly omits dequantized coefficients when the caller only
// consumes quantized coefficients and EOBs. Store order and fallback guards
// match QuantizeVP8Batch for its remaining inputs and outputs.
func QuantizeVP8BatchOnly(coeff []int16, round, quant *[16]int16, qcoeff []int16, iscan *[16]int16, eobs []uint8, count int) bool {
	if !Enabled || count <= 0 || count > min(len(coeff), len(qcoeff))/16 || count > len(eobs) {
		return false
	}
	// Validate element counts before multiplying; uintptr byte offsets also
	// avoid signed overflow for large accepted counts.
	n := count * 16

	// A conservative span covers the shared encoder BlockQuant tables.
	// Separately allocated tables may decline, preserving the general fallback.
	overlap := func(a uintptr, an uintptr, b uintptr, bn uintptr) bool {
		return (a <= b && b-a < an) || (b < a && a-b < bn)
	}
	rptr, qptr := uintptr(unsafe.Pointer(round)), uintptr(unsafe.Pointer(quant))
	lo, hi := min(rptr, qptr), max(rptr, qptr)+32
	scan := uintptr(unsafe.Pointer(iscan))
	for k, p := range [2]uintptr{uintptr(unsafe.Pointer(unsafe.SliceData(qcoeff))), uintptr(unsafe.Pointer(unsafe.SliceData(eobs)))} {
		size := uintptr(n) * 2
		if k == 1 {
			size = uintptr(count)
		}
		if overlap(p, size, lo, hi-lo) || overlap(p, size, scan, 32) {
			return false
		}
	}

	r := archsimd.LoadInt16x16Array(round)
	q := archsimd.LoadInt16x16Array(quant)
	rq := r.Or(q)
	// Keep this load VEX encoded while r/q YMM values are live. A literal
	// XMM zero makes Go 1.27 emit legacy MOVUPS and pay an AVX/SSE transition.
	if any16(rq.GetLo().Or(rq.GetHi()).Less(archsimd.LoadInt16x8Array(&packZero16))) {
		archsimd.ClearAVXUpperBits()
		return false
	}
	positions := archsimd.LoadInt16x16Array(iscan)
	z := archsimd.Int16x16{}
	cp, qp := unsafe.Pointer(unsafe.SliceData(coeff)), unsafe.Pointer(unsafe.SliceData(qcoeff))
	for i := range count {
		off := uintptr(i) * 32
		c := archsimd.LoadInt16x16Array((*[16]int16)(unsafe.Add(cp, off)))
		magnitude := c.Abs().ToBits().Add(r.ToBits()).MulHigh(q.ToBits()).BitsToInt16()
		sign := c.ShiftAllRight(15)
		t := magnitude.Xor(sign).Sub(sign)
		t.StoreArray((*[16]int16)(unsafe.Add(qp, off)))
		last := positions.Masked(t.NotEqual(z))
		eobs[i] = uint8(fpReduceMax16(last.GetLo().Max(last.GetHi())))
	}
	// Return to ordinary Go/SSE with no live wide vector results.
	archsimd.ClearAVXUpperBits()
	return true
}
