//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// DCPairNative selects a single arithmetic owner for paired DC reconstruction.
// Complete x64 windows use byte SIMD; fallback access order stays unchanged.
const DCPairNative = true

// DCPairInt32Native owns rounding, byte arithmetic and fallback, so a public
// caller need not keep ten arguments live across a boolean helper call.
func DCPairInt32Native(c0, c1 int32, pred []byte, ps int, dst []byte, ds int) {
	if !Enabled {
		dcPairNativeScalar(c0, c1, pred, ps, dst, ds)
		return
	}
fast:
	for {
		if ps < 8 || ds < 8 || len(pred) < 8 || len(dst) < 8 || ps > (len(pred)-8)/3 || ds > (len(dst)-8)/3 {
			break fast
		}
		pp, dp := uintptr(unsafe.Pointer(unsafe.SliceData(pred))), uintptr(unsafe.Pointer(unsafe.SliceData(dst)))
		exact := pp == dp && ps == ds
		disjoint := (pp <= dp && dp-pp >= uintptr(len(pred))) || (dp < pp && pp-dp >= uintptr(len(dst)))
		if !exact && !disjoint {
			break fast
		}
		plus, minus := dcPairCoeffBytes(c0, c1)
		sp, out := unsafe.Pointer(unsafe.SliceData(pred)), unsafe.Pointer(unsafe.SliceData(dst))
		dcPairByteRows2(sp, ps, out, ds, plus, minus)
		dcPairByteRows2(unsafe.Add(sp, 2*ps), ps, unsafe.Add(out, 2*ds), ds, plus, minus)
		return
	}
	d0, d1 := dcPairNativeDelta(c0), dcPairNativeDelta(c1)
	// Exact enabled fallback: block0 is complete before slicing block1.
	// Capture each four-byte prediction before its store as AddConstant does.
	for block := range 2 {
		p, d, delta := pred, dst, d0
		if block != 0 {
			p = pred[4:]
			d = dst[4:]
			delta = d1
		}
		v := archsimd.BroadcastInt32x4(int32(delta))
		for y := range 4 {
			StorePixels4(d[y*ds:], LoadPixels4(p[y*ps:]).Add(v))
		}
	}
}
