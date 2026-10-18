//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// DCPairInPlaceInt32Native specializes libvpx v1.16.0's idct_dequant_0_2x_neon
// for one prediction/destination window, retaining every fallback access order.
func DCPairInPlaceInt32Native(c0, c1 int32, b []byte, stride int) {
	if !Enabled {
		dcPairNativeScalar(c0, c1, b, stride, b, stride)
		return
	}
	d0, d1 := dcPairNativeDelta(c0), dcPairNativeDelta(c1)
	if stride >= 8 && len(b) >= 8 && stride <= (len(b)-8)/3 {
		delta := dcPairSignedDelta8(d0, d1)
		sp := unsafe.Pointer(unsafe.SliceData(b))
		dcPairSignedRows2(sp, stride, sp, stride, delta)
		sp2 := unsafe.Add(sp, 2*stride)
		dcPairSignedRows2(sp2, stride, sp2, stride, delta)
		return
	}
	// Exact original enabled fallback: complete block0 before slicing block1,
	// and evaluate destination before capturing each four-byte prediction.
	for block := range 2 {
		p, d, delta := b, b, d0
		if block != 0 {
			p = b[4:]
			d = b[4:]
			delta = d1
		}
		v := archsimd.BroadcastInt32x4(int32(delta))
		for y := range 4 {
			StorePixels4(d[y*stride:], LoadPixels4(p[y*stride:]).Add(v))
		}
	}
}

// The low four lanes use d0, the high four d1. Callers clamp both first.
func dcPairSignedDelta8(d0, d1 int16) archsimd.Int16x8 {
	lo := archsimd.BroadcastInt16x8(d0).ToBits().ReshapeToUint64s()
	hi := archsimd.BroadcastInt16x8(d1).ToBits().ReshapeToUint64s()
	return lo.InterleaveLo(hi).ReshapeToUint16s().BitsToInt16()
}

// Callers prove both exact eight-byte windows and preserve the original
// two-row read packet before its two ordered stores. Pixel+delta fits int16.
func dcPairSignedRows2(sp unsafe.Pointer, ps int, out unsafe.Pointer, ds int, delta archsimd.Int16x8) {
	lo := *(*uint64)(sp)
	hi := *(*uint64)(unsafe.Add(sp, ps))
	l := archsimd.BroadcastUint64x2(lo).ReshapeToUint8s().ExtendLo8ToUint16().BitsToInt16().Add(delta).SaturateToUint8().ReshapeToUint64s()
	h := archsimd.BroadcastUint64x2(hi).ReshapeToUint8s().ExtendLo8ToUint16().BitsToInt16().Add(delta).SaturateToUint8().ReshapeToUint64s()
	v := l.InterleaveLo(h)
	*(*uint64)(out) = v.GetElem(0)
	*(*uint64)(unsafe.Add(out, ds)) = v.GetElem(1)
}
