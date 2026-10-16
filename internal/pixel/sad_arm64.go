//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// SAD keeps separate low/high uint16 accumulators until the final widening
// reduction. At most 64*64/16 byte vectors contribute to each lane, so the
// maximum is 256*255=65280 and cannot wrap. The endpoint checks prove every
// unaligned fixed-width load below lies inside its source slice, including
// the final row. Pointer arithmetic never leaves the originating allocation.
func SAD(src []byte, ss int, ref []byte, rs, w, h int) int {
	if !Enabled {
		return SADScalar(src, ss, ref, rs, w, h)
	}

	if w <= 0 || h <= 0 {
		return 0
	}
	if w > 64 || h > 64 || ss < 0 || rs < 0 {
		return sadGeneric(src, ss, ref, rs, w, h)
	}
	_ = src[(h-1)*ss+w-1]
	_ = ref[(h-1)*rs+w-1]
	var lo, hi archsimd.Uint16x8
	tail := 0
	sp, rp := unsafe.Pointer(unsafe.SliceData(src)), unsafe.Pointer(unsafe.SliceData(ref))
	if w == 8 && (h == 8 || h == 16) {
		value := sad8x8Ptr(sp, ss, rp, rs)
		if h == 16 {
			value += sad8x8Ptr(step(sp, 8*ss), ss, step(rp, 8*rs), rs)
		}
		return value
	}
	for y := range h {
		a, b := unsafe.Add(sp, y*ss), unsafe.Add(rp, y*rs)
		x := 0
		for ; x+16 <= w; x += 16 {
			va := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(a, x)))
			vb := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(b, x)))
			d := va.Max(vb).Sub(va.Min(vb))
			lo = lo.Add(d.ExtendLo8ToUint16())
			hi = hi.Add(d.HiToLo().ExtendLo8ToUint16())
		}
		if x+8 <= w {
			va := archsimd.Uint64x2{}.SetElem(0, *(*uint64)(unsafe.Add(a, x))).ReshapeToUint8s()
			vb := archsimd.Uint64x2{}.SetElem(0, *(*uint64)(unsafe.Add(b, x))).ReshapeToUint8s()
			d := va.Max(vb).Sub(va.Min(vb))
			lo = lo.Add(d.ExtendLo8ToUint16())
			x += 8
		}
		if x+4 <= w {
			va := archsimd.Uint32x4{}.SetElem(0, *(*uint32)(unsafe.Add(a, x))).ReshapeToUint8s()
			vb := archsimd.Uint32x4{}.SetElem(0, *(*uint32)(unsafe.Add(b, x))).ReshapeToUint8s()
			d := va.Max(vb).Sub(va.Min(vb))
			lo = lo.Add(d.ExtendLo8ToUint16())
			x += 4
		}
		for ; x < w; x++ {
			d := int(*(*byte)(unsafe.Add(a, x))) - int(*(*byte)(unsafe.Add(b, x)))
			tail += max(d, -d)
		}
	}
	sum := lo.ExtendLo4ToUint32().Add(lo.HiToLo().ExtendLo4ToUint32()).Add(hi.ExtendLo4ToUint32()).Add(hi.HiToLo().ExtendLo4ToUint32())
	return int(sum.ReduceSum()) + tail
}
