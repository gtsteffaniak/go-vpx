//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// SSE8Ptr requires caller-validated complete read-only 8x8 windows.
// Exact eight-byte loads support unaligned, overlapping and zero-stride rows.
// Each signed difference is in [-255, 255]; the full sum fits int32.
func SSE8Ptr(src *byte, ss int, ref *byte, rs int) uint32 {
	if !Enabled || !archsimd.X86.AVX2() || src == nil || ref == nil || (uint(ss)|uint(rs))>>30 != 0 {
		return SSE8(unsafe.Slice(src, 7*ss+8), ss, unsafe.Slice(ref, 7*rs+8), rs, 8)
	}
	defer archsimd.ClearAVXUpperBits()
	sp, rp := unsafe.Pointer(src), unsafe.Pointer(ref)
	var acc archsimd.Int32x8
	so, ro := 0, 0
	for y := 0; y < 8; y += 4 {
		{
			a := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so+0*ss))).InterleaveLoGrouped(archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so+1*ss)))).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
			b := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro+0*rs))).InterleaveLoGrouped(archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro+1*rs)))).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
			d := a.Sub(b)
			acc = acc.Add(d.DotProductPairs(d))
		}
		{
			a := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so+2*ss))).InterleaveLoGrouped(archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(sp, so+3*ss)))).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
			b := archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro+2*rs))).InterleaveLoGrouped(archsimd.BroadcastUint64x4(*(*uint64)(unsafe.Add(rp, ro+3*rs)))).GetLo().ReshapeToUint8s().ExtendToUint16().BitsToInt16()
			d := a.Sub(b)
			acc = acc.Add(d.DotProductPairs(d))
		}
		so += 4 * ss
		ro += 4 * rs
	}
	return uint32(sum32(acc.GetLo().Add(acc.GetHi())))
}
