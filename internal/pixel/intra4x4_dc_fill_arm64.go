//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// Intra4x4DCFill follows libvpx v1.16.0 vpx_dsp/intrapred.c DC rounding.
// The exact four-byte references are completely consumed before output stores.
func Intra4x4DCFill(dst []byte, ds int, above, left []byte) bool {
	if !Enabled || uint(ds)>>30 != 0 || 3*ds+4 > len(dst) || len(above) < 4 || len(left) < 4 {
		return false
	}
	a := *(*uint32)(unsafe.Pointer(unsafe.SliceData(above)))
	l := *(*uint32)(unsafe.Pointer(unsafe.SliceData(left)))
	v := archsimd.Uint32x4{}.SetElem(0, a).SetElem(1, l).ReshapeToUint8s().ExtendLo8ToUint16()
	value := byte((int(v.ReduceSum()) + 4) >> 3)
	word := archsimd.BroadcastUint8x16(value).ReshapeToUint32s().GetElem(0)
	p := unsafe.Pointer(unsafe.SliceData(dst))
	*(*uint32)(p) = word
	*(*uint32)(unsafe.Add(p, ds)) = word
	*(*uint32)(unsafe.Add(p, 2*ds)) = word
	*(*uint32)(unsafe.Add(p, 3*ds)) = word
	return true
}
