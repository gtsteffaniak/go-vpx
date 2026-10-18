package dsp

import (
	"github.com/thesyncim/govpx/internal/pixel"
	"unsafe"
)

// Ported from libvpx v1.16.0 vp8/common/dequantize.c.

func DequantizeBlock(qcoeff *[16]int16, dequant *[16]int16, dqcoeff *[16]int16) {
	// Exact in-place multiplication is safe; partial overlap retains scalar ordering.
	q, d, o := uintptr(unsafe.Pointer(qcoeff)), uintptr(unsafe.Pointer(dequant)), uintptr(unsafe.Pointer(dqcoeff))
	if pixel.Enabled && (q == o || q < o && o-q >= 32 || o < q && q-o >= 32) && (d == o || d < o && o-d >= 32 || o < d && d-o >= 32) {
		pixel.Mul16(qcoeff[:], dequant[:], dqcoeff[:], 16)
		return
	}

	for i := range 16 {
		dqcoeff[i] = qcoeff[i] * dequant[i]
	}
}

func DequantIDCT4x4Add(input *[16]int16, dequant *[16]int16, dst []byte, stride int) {
	for i := range 16 {
		input[i] = dequant[i] * input[i]
	}

	IDCT4x4Add(input, dst, stride, dst, stride)

	for i := range 16 {
		input[i] = 0
	}
}
