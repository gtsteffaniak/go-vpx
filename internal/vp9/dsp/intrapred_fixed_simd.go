//go:build goexperiment.simd && (arm64 || amd64) && !purego

package dsp

import (
	"github.com/gtsteffaniak/go-vpx/internal/pixel"
	"simd/archsimd"
	"unsafe"
)

// Fixed-row kernels implement libvpx v1.16.0 vpx_dsp/intrapred.c arithmetic.
// Entry guards prove every complete row and preserve the scalar alias/panic path.
func intraTM32Fixed(dst []byte, stride int, above, left []byte, top byte) bool {
	if !pixel.Enabled {
		return false
	}
	dp, ap, lp := unsafe.Pointer(unsafe.SliceData(dst)), unsafe.Pointer(unsafe.SliceData(above)), unsafe.Pointer(unsafe.SliceData(left))
	a := archsimd.LoadUint8x16Array((*[16]byte)(ap))
	b := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(ap, 16)))
	for y := 0; y < 32; y++ {
		delta := int(*(*byte)(unsafe.Add(lp, y))) - int(top)
		pos, neg := archsimd.BroadcastUint8x16(byte(max(delta, 0))), archsimd.BroadcastUint8x16(byte(max(-delta, 0)))
		row := unsafe.Add(dp, y*stride)
		a.AddSaturated(pos).SubSaturated(neg).StoreArray((*[16]byte)(row))
		b.AddSaturated(pos).SubSaturated(neg).StoreArray((*[16]byte)(unsafe.Add(row, 16)))
	}
	return true
}
func intraVFixed(dst []byte, stride, n int, above []byte) bool {
	if !intraBridgeDisjoint(dst, stride, n, above, n) {
		return false
	}
	dp, ap := unsafe.Pointer(unsafe.SliceData(dst)), unsafe.Pointer(unsafe.SliceData(above))
	switch n {
	case 4:
		v := *(*uint32)(ap)
		for y := 0; y < 4; y++ {
			*(*uint32)(unsafe.Add(dp, y*stride)) = v
		}
	case 8:
		v := *(*uint64)(ap)
		for y := 0; y < 8; y++ {
			*(*uint64)(unsafe.Add(dp, y*stride)) = v
		}
	case 16:
		v := archsimd.LoadUint8x16Array((*[16]byte)(ap))
		for y := 0; y < 16; y++ {
			v.StoreArray((*[16]byte)(unsafe.Add(dp, y*stride)))
		}
	case 32:
		a := archsimd.LoadUint8x16Array((*[16]byte)(ap))
		b := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(ap, 16)))
		for y := 0; y < 32; y++ {
			row := unsafe.Add(dp, y*stride)
			a.StoreArray((*[16]byte)(row))
			b.StoreArray((*[16]byte)(unsafe.Add(row, 16)))
		}
	default:
		return false
	}
	return true
}
func intraHFixed(dst []byte, stride, n int, left []byte) bool {
	if !intraBridgeWindow(dst, stride, n) || len(left) < n {
		return false
	}
	dp, lp := unsafe.Pointer(unsafe.SliceData(dst)), unsafe.Pointer(unsafe.SliceData(left))
	switch n {
	case 4:
		for y := 0; y < 4; y++ {
			v := uint32(*(*byte)(unsafe.Add(lp, y))) * 0x01010101
			*(*uint32)(unsafe.Add(dp, y*stride)) = v
		}
	case 8:
		for y := 0; y < 8; y++ {
			v := uint64(*(*byte)(unsafe.Add(lp, y))) * 0x0101010101010101
			*(*uint64)(unsafe.Add(dp, y*stride)) = v
		}
	case 16:
		for y := 0; y < 16; y++ {
			v := archsimd.BroadcastUint8x16(*(*byte)(unsafe.Add(lp, y)))
			v.StoreArray((*[16]byte)(unsafe.Add(dp, y*stride)))
		}
	case 32:
		for y := 0; y < 32; y++ {
			v := archsimd.BroadcastUint8x16(*(*byte)(unsafe.Add(lp, y)))
			row := unsafe.Add(dp, y*stride)
			v.StoreArray((*[16]byte)(row))
			v.StoreArray((*[16]byte)(unsafe.Add(row, 16)))
		}
	default:
		return false
	}
	return true
}
func intraEdgeSum(p unsafe.Pointer, n int) int {
	switch n {
	case 4:
		return intraSumBytes(archsimd.Uint32x4{}.SetElem(0, *(*uint32)(p)).ReshapeToUint8s())
	case 8:
		return intraSumBytes(archsimd.Uint64x2{}.SetElem(0, *(*uint64)(p)).ReshapeToUint8s())
	case 16:
		return intraSumBytes(archsimd.LoadUint8x16Array((*[16]byte)(p)))
	default:
		return intraSumBytes(archsimd.LoadUint8x16Array((*[16]byte)(p))) + intraSumBytes(archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(p, 16))))
	}
}
func intraDCFixed(dst []byte, stride, n int, above, left []byte, up, lf bool) bool {
	if !intraBridgeWindow(dst, stride, n) || (up && len(above) < n) || (lf && len(left) < n) {
		return false
	}
	sum, count := 0, 0
	if up {
		sum += intraEdgeSum(unsafe.Pointer(unsafe.SliceData(above)), n)
		count += n
	}
	if lf {
		sum += intraEdgeSum(unsafe.Pointer(unsafe.SliceData(left)), n)
		count += n
	}
	value := byte(128)
	if count > 0 {
		value = byte((sum + count/2) / count)
	}
	dp := unsafe.Pointer(unsafe.SliceData(dst))
	v := archsimd.BroadcastUint8x16(value)
	switch n {
	case 4:
		word := v.ReshapeToUint32s().GetElem(0)
		for y := 0; y < 4; y++ {
			*(*uint32)(unsafe.Add(dp, y*stride)) = word
		}
	case 8:
		word := v.ReshapeToUint64s().GetElem(0)
		for y := 0; y < 8; y++ {
			*(*uint64)(unsafe.Add(dp, y*stride)) = word
		}
	case 16:
		for y := 0; y < 16; y++ {
			v.StoreArray((*[16]byte)(unsafe.Add(dp, y*stride)))
		}
	case 32:
		for y := 0; y < 32; y++ {
			row := unsafe.Add(dp, y*stride)
			v.StoreArray((*[16]byte)(row))
			v.StoreArray((*[16]byte)(unsafe.Add(row, 16)))
		}
	default:
		return false
	}
	return true
}
