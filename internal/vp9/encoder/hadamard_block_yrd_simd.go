//go:build goexperiment.simd && (arm64 || amd64) && !purego

package encoder

import (
	"github.com/thesyncim/govpx/internal/pixel"
	"simd/archsimd"
	"unsafe"
)

// hadamardBlockYrdGo mirrors the scoring-only layout of libvpx v1.16.0
// vpx_dsp/arm/hadamard_neon.c. BlockYrd supplies disjoint residual/output
// windows and residuals in [-255,255]. Each 8x8 coefficient is therefore in
// [-16320,16320], so the 16x16 combine's pre-shift sums fit signed 16 bits.
func hadamardBlockYrdGo(src []int16, stride int, coeff []int16, n int) bool {
	if !pixel.Enabled || (n != 8 && n != 16) || stride < n || len(src) < n ||
		stride > (len(src)-n)/(n-1) || len(coeff) < n*n {
		return false
	}
	sp, cp := unsafe.Pointer(unsafe.SliceData(src)), unsafe.Pointer(unsafe.SliceData(coeff))
	if n == 8 {
		hadamardBlockYrd8(sp, stride, cp)
		return true
	}
	hadamardBlockYrd8(sp, stride, cp)
	hadamardBlockYrd8(unsafe.Add(sp, 16), stride, unsafe.Add(cp, 128))
	hadamardBlockYrd8(unsafe.Add(sp, 16*stride), stride, unsafe.Add(cp, 256))
	hadamardBlockYrd8(unsafe.Add(sp, 16*stride+16), stride, unsafe.Add(cp, 384))
	for i := 0; i < 128; i += 16 {
		a := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(cp, i)))
		b := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(cp, i+128)))
		c := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(cp, i+256)))
		d := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(cp, i+384)))
		u, v, w, z := a.Add(b).ShiftAllRight(1), a.Sub(b).ShiftAllRight(1), c.Add(d).ShiftAllRight(1), c.Sub(d).ShiftAllRight(1)
		u.Add(w).StoreArray((*[8]int16)(unsafe.Add(cp, i)))
		v.Add(z).StoreArray((*[8]int16)(unsafe.Add(cp, i+128)))
		u.Sub(w).StoreArray((*[8]int16)(unsafe.Add(cp, i+256)))
		v.Sub(z).StoreArray((*[8]int16)(unsafe.Add(cp, i+384)))
	}
	return true
}

// The caller proves all eight complete input rows and the 64 output samples.
// Named vectors keep both butterflies and the single transpose in registers.
func hadamardBlockYrd8(src unsafe.Pointer, stride int, coeff unsafe.Pointer) {
	v0 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 0*stride)))
	v1 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 2*stride)))
	v2 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 4*stride)))
	v3 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 6*stride)))
	v4 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 8*stride)))
	v5 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 10*stride)))
	v6 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 12*stride)))
	v7 := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(src, 14*stride)))
	{
		b0, b1 := v0.Add(v1), v0.Sub(v1)
		b2, b3 := v2.Add(v3), v2.Sub(v3)
		b4, b5 := v4.Add(v5), v4.Sub(v5)
		b6, b7 := v6.Add(v7), v6.Sub(v7)
		c0, c1, c2, c3 := b0.Add(b2), b1.Add(b3), b0.Sub(b2), b1.Sub(b3)
		c4, c5, c6, c7 := b4.Add(b6), b5.Add(b7), b4.Sub(b6), b5.Sub(b7)
		v0, v7, v3, v4, v2, v6, v1, v5 = c0.Add(c4), c1.Add(c5), c2.Add(c6), c3.Add(c7), c0.Sub(c4), c1.Sub(c5), c2.Sub(c6), c3.Sub(c7)
	}
	{
		a0, a1 := v0.InterleaveLo(v1).ToBits().ReshapeToUint32s(), v0.InterleaveHi(v1).ToBits().ReshapeToUint32s()
		a2, a3 := v2.InterleaveLo(v3).ToBits().ReshapeToUint32s(), v2.InterleaveHi(v3).ToBits().ReshapeToUint32s()
		a4, a5 := v4.InterleaveLo(v5).ToBits().ReshapeToUint32s(), v4.InterleaveHi(v5).ToBits().ReshapeToUint32s()
		a6, a7 := v6.InterleaveLo(v7).ToBits().ReshapeToUint32s(), v6.InterleaveHi(v7).ToBits().ReshapeToUint32s()
		b0, b1, b2, b3 := a0.InterleaveLo(a2).ReshapeToUint64s(), a0.InterleaveHi(a2).ReshapeToUint64s(), a1.InterleaveLo(a3).ReshapeToUint64s(), a1.InterleaveHi(a3).ReshapeToUint64s()
		b4, b5, b6, b7 := a4.InterleaveLo(a6).ReshapeToUint64s(), a4.InterleaveHi(a6).ReshapeToUint64s(), a5.InterleaveLo(a7).ReshapeToUint64s(), a5.InterleaveHi(a7).ReshapeToUint64s()
		v0, v1, v2, v3 = b0.InterleaveLo(b4).ReshapeToUint16s().BitsToInt16(), b0.InterleaveHi(b4).ReshapeToUint16s().BitsToInt16(), b1.InterleaveLo(b5).ReshapeToUint16s().BitsToInt16(), b1.InterleaveHi(b5).ReshapeToUint16s().BitsToInt16()
		v4, v5, v6, v7 = b2.InterleaveLo(b6).ReshapeToUint16s().BitsToInt16(), b2.InterleaveHi(b6).ReshapeToUint16s().BitsToInt16(), b3.InterleaveLo(b7).ReshapeToUint16s().BitsToInt16(), b3.InterleaveHi(b7).ReshapeToUint16s().BitsToInt16()
	}
	{
		b0, b1 := v0.Add(v1), v0.Sub(v1)
		b2, b3 := v2.Add(v3), v2.Sub(v3)
		b4, b5 := v4.Add(v5), v4.Sub(v5)
		b6, b7 := v6.Add(v7), v6.Sub(v7)
		c0, c1, c2, c3 := b0.Add(b2), b1.Add(b3), b0.Sub(b2), b1.Sub(b3)
		c4, c5, c6, c7 := b4.Add(b6), b5.Add(b7), b4.Sub(b6), b5.Sub(b7)
		v0, v7, v3, v4, v2, v6, v1, v5 = c0.Add(c4), c1.Add(c5), c2.Add(c6), c3.Add(c7), c0.Sub(c4), c1.Sub(c5), c2.Sub(c6), c3.Sub(c7)
	}
	v0.StoreArray((*[8]int16)(unsafe.Add(coeff, 0)))
	v1.StoreArray((*[8]int16)(unsafe.Add(coeff, 16)))
	v2.StoreArray((*[8]int16)(unsafe.Add(coeff, 32)))
	v3.StoreArray((*[8]int16)(unsafe.Add(coeff, 48)))
	v4.StoreArray((*[8]int16)(unsafe.Add(coeff, 64)))
	v5.StoreArray((*[8]int16)(unsafe.Add(coeff, 80)))
	v6.StoreArray((*[8]int16)(unsafe.Add(coeff, 96)))
	v7.StoreArray((*[8]int16)(unsafe.Add(coeff, 112)))
}
