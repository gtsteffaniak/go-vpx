//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// sse8RowsWide requires the existing bounded SSE8 complete-window checks.
// Two exact eight-byte rows share one AVX2 int16 subtraction/dot product.
// The maximum 64-row sum is 64*8*255*255, which fits int32.
const sse8WideAvailable = true

func sse8RowsWide(sp unsafe.Pointer, ss int, rp unsafe.Pointer, rs, h int) uint32 {
	defer archsimd.ClearAVXUpperBits()
	var acc archsimd.Int32x8
	so, ro, y := 0, 0, 0
	for ; y+2 <= h; y += 2 {
		a := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so))).InterleaveLoGrouped(archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so+ss)))).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
		b := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro))).InterleaveLoGrouped(archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro+rs)))).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
		d := a.Sub(b)
		acc = acc.Add(d.DotProductPairs(d))
		so += 2 * ss
		ro += 2 * rs
	}
	if y < h {
		a := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so))).InterleaveLoGrouped(archsimd.Uint64x4{}).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
		b := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro))).InterleaveLoGrouped(archsimd.Uint64x4{}).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
		d := a.Sub(b)
		acc = acc.Add(d.DotProductPairs(d))
	}
	return uint32(sum32(acc.GetLo().Add(acc.GetHi())))
}
