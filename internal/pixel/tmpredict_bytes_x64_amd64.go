//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// Saturated byte differences express clip(above+left-top) exactly: at
// most one of positive and negative is nonzero. Capture above only after
// proving it cannot be changed by a destination store; left stays per-row.
func tmpredictBytesX64(dst []byte, ds int, above, left []byte, top byte, w int) bool {
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
	var v archsimd.Uint8x16
	if w == 8 {
		v = archsimd.Uint64x2{}.SetElem(0, *(*uint64)(ap)).ReshapeToUint8s()
	} else {
		v = archsimd.LoadUint8x16Array((*[16]byte)(ap))
	}
	t := archsimd.BroadcastUint8x16(top)
	for y := 0; y < w; y++ {
		l := archsimd.BroadcastUint8x16(left[y])
		row := v.AddSaturated(l.SubSaturated(t)).SubSaturated(t.SubSaturated(l))
		if w == 8 {
			*(*uint64)(unsafe.Add(dp, y*ds)) = row.ReshapeToUint64s().GetElem(0)
		} else {
			row.StoreArray((*[16]byte)(unsafe.Add(dp, y*ds)))
		}
	}
	return true
}
