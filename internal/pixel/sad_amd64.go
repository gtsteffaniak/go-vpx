//go:build goexperiment.simd && amd64 && !purego

package pixel

import "unsafe"

func SAD(src []byte, ss int, ref []byte, rs, w, h int) int {
	if !Enabled {
		return SADScalar(src, ss, ref, rs, w, h)
	}
	if w == 8 && (h == 8 || h == 16) && ss >= 0 && rs >= 0 {
		_ = src[(h-1)*ss+7]
		_ = ref[(h-1)*rs+7]
		sp, rp := unsafe.Pointer(unsafe.SliceData(src)), unsafe.Pointer(unsafe.SliceData(ref))
		value := sad8x8Ptr(sp, ss, rp, rs)
		if h == 16 {
			value += sad8x8Ptr(step(sp, 8*ss), ss, step(rp, 8*rs), rs)
		}
		return value
	}
	return sadGeneric(src, ss, ref, rs, w, h)
}
