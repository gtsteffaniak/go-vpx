//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// TMPredictBytes retains libvpx v1.16.0 vpx_dsp/intrapred.c true-motion
// rounding/clipping. It captures left[y] per row, as the original SIMD path
// does. Above may be hoisted only when the entire destination span is disjoint.
// Declines do not read input or write output, preserving the original fallback.
func TMPredictBytes(dst []byte, ds int, above, left []byte, top byte, w int) bool {
	if !Enabled || (w != 8 && w != 16) || uint(ds)>>30 != 0 {
		return false
	}
	span := (w-1)*ds + w
	if span > len(dst) || len(above) < w || len(left) < w {
		return false
	}
	dp, ap := unsafe.Pointer(unsafe.SliceData(dst)), unsafe.Pointer(unsafe.SliceData(above))
	d, a := uintptr(dp), uintptr(ap)
	if (d <= a && a-d < uintptr(span)) || (a < d && d-a < uintptr(w)) {
		return false
	}
	lp := unsafe.Pointer(unsafe.SliceData(left))
	l := uintptr(lp)
	if !((d <= l && l-d < uintptr(span)) || (l < d && d-l < uintptr(w))) {
		// Only disjoint left/output windows may capture all left rows up front.
		// The original per-row left read remains below for every alias.
		t := archsimd.BroadcastUint8x16(top)
		if w == 8 {
			a8 := *(*uint64)(ap)
			v := archsimd.Uint64x2{}.SetElem(0, a8).SetElem(1, a8).ReshapeToUint8s()
			lv := archsimd.Uint64x2{}.SetElem(0, *(*uint64)(lp)).ReshapeToUint8s()
			pos, neg := lv.SubSaturated(t), t.SubSaturated(lv)
			idx := archsimd.Uint64x2{}.SetElem(1, 0x0101010101010101).ReshapeToUint8s()
			step := archsimd.BroadcastUint8x16(2)
			for remaining := 4; ; {
				rows := v.AddSaturated(pos.LookupOrZero(idx)).SubSaturated(neg.LookupOrZero(idx)).ReshapeToUint64s()
				*(*uint64)(dp) = rows.GetElem(0)
				*(*uint64)(unsafe.Add(dp, ds)) = rows.GetElem(1)
				remaining--
				if remaining == 0 {
					break
				}
				dp = unsafe.Add(dp, 2*ds)
				idx = idx.Add(step)
			}
		} else {
			v := archsimd.LoadUint8x16Array((*[16]byte)(ap))
			lv := archsimd.LoadUint8x16Array((*[16]byte)(lp))
			pos, neg := lv.SubSaturated(t), t.SubSaturated(lv)
			idx := archsimd.Uint8x16{}
			step := archsimd.BroadcastUint8x16(1)
			for remaining := 16; ; {
				row := v.AddSaturated(pos.LookupOrZero(idx)).SubSaturated(neg.LookupOrZero(idx))
				row.StoreArray((*[16]byte)(dp))
				remaining--
				if remaining == 0 {
					break
				}
				dp = unsafe.Add(dp, ds)
				idx = idx.Add(step)
			}
		}
		return true
	}
	// Read left per row after the previous row stores, including aliases.
	// delta is in [-255,255]. Only one of plus/minus is nonzero, so these
	// two unsigned saturations equal clamp(above+left-top,0,255) exactly.
	if w == 8 {
		v := archsimd.Uint64x2{}.SetElem(0, *(*uint64)(ap)).ReshapeToUint8s()
		for y := 0; y < 8; y++ {
			delta := int(*(*byte)(unsafe.Add(lp, y))) - int(top)
			row := v.AddSaturated(archsimd.BroadcastUint8x16(byte(max(delta, 0)))).SubSaturated(archsimd.BroadcastUint8x16(byte(max(-delta, 0))))
			*(*uint64)(unsafe.Add(dp, y*ds)) = row.ReshapeToUint64s().GetElem(0)
		}
	} else {
		v := archsimd.LoadUint8x16Array((*[16]byte)(ap))
		for y := 0; y < 16; y++ {
			delta := int(*(*byte)(unsafe.Add(lp, y))) - int(top)
			row := v.AddSaturated(archsimd.BroadcastUint8x16(byte(max(delta, 0)))).SubSaturated(archsimd.BroadcastUint8x16(byte(max(-delta, 0))))
			row.StoreArray((*[16]byte)(unsafe.Add(dp, y*ds)))
		}
	}
	return true
}
