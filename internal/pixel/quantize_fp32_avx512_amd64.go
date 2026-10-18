//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// All used indices are clamped to0..11; upper entries are unreachable.
var fpClassWordLUT32 = [32]uint16{0, 1, 2, 3, 3, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5}

// quantizeFPBlock32AVX512 requires the existing public length/table checks,
// complete disjoint windows, n256..1024 divisible by32, and real AVX512().
func quantizeFPBlock32AVX512(coeff, iscan, qcoeff, dqcoeff []int16, n int, r, q, d [2]int16, classes []byte) int {
	cp, ip, qp, dp := unsafe.Pointer(unsafe.SliceData(coeff)), unsafe.Pointer(unsafe.SliceData(iscan)), unsafe.Pointer(unsafe.SliceData(qcoeff)), unsafe.Pointer(unsafe.SliceData(dqcoeff))
	clp := unsafe.Pointer(unsafe.SliceData(classes))
	max16 := archsimd.BroadcastUint16x32(32767)
	zero := archsimd.Int16x32{}
	last := zero
	var lut archsimd.Uint16x32
	if classes != nil {
		lut = archsimd.LoadUint16x32Array(&fpClassWordLUT32)
	}
	eleven := archsimd.BroadcastInt16x32(11)
	{
		i := 0
		vr := archsimd.BroadcastInt16x32(r[1]).SetLo(archsimd.BroadcastInt16x16(r[1]).SetLo(archsimd.BroadcastInt16x8(r[1]).SetElem(0, r[0])))
		vq := archsimd.BroadcastInt16x32(q[1]).SetLo(archsimd.BroadcastInt16x16(q[1]).SetLo(archsimd.BroadcastInt16x8(q[1]).SetElem(0, q[0])))
		vd := archsimd.BroadcastInt16x32(d[1]).SetLo(archsimd.BroadcastInt16x16(d[1]).SetLo(archsimd.BroadcastInt16x8(d[1]).SetElem(0, d[0])))
		c := archsimd.LoadInt16x32Array((*[32]int16)(unsafe.Add(cp, uintptr(i)*2)))
		a := c.Abs().ToBits().Add(vr.ToBits()).Min(max16).BitsToInt16()
		active := a.GreaterEqual(vd)
		if active.ToBits() == 0 {
			if qcoeff != nil {
				zero.StoreArray((*[32]int16)(unsafe.Add(qp, uintptr(i)*2)))
			}
			zero.StoreArray((*[32]int16)(unsafe.Add(dp, uintptr(i)*2)))
			if classes != nil {
				archsimd.Uint8x32{}.StoreArray((*[32]byte)(unsafe.Add(clp, uintptr(i))))
			}
		} else {
			magnitude := a.MulHigh(vq).Masked(active)
			t := zero.Sub(magnitude).IfElse(c.Less(zero), magnitude)
			if qcoeff != nil {
				t.StoreArray((*[32]int16)(unsafe.Add(qp, uintptr(i)*2)))
			}
			t.Mul(vd).StoreArray((*[32]int16)(unsafe.Add(dp, uintptr(i)*2)))
			last = last.Max(archsimd.LoadInt16x32Array((*[32]int16)(unsafe.Add(ip, uintptr(i)*2))).Masked(t.NotEqual(zero)))
			if classes != nil {
				lut.Permute(magnitude.Min(eleven).ToBits()).TruncToUint8().StoreArray((*[32]byte)(unsafe.Add(clp, uintptr(i))))
			}
		}
	}
	acR, acQ, acD := archsimd.BroadcastInt16x32(r[1]), archsimd.BroadcastInt16x32(q[1]), archsimd.BroadcastInt16x32(d[1])
	for i := 32; i < n; i += 32 {
		c := archsimd.LoadInt16x32Array((*[32]int16)(unsafe.Add(cp, uintptr(i)*2)))
		a := c.Abs().ToBits().Add(acR.ToBits()).Min(max16).BitsToInt16()
		active := a.GreaterEqual(acD)
		if active.ToBits() == 0 {
			if qcoeff != nil {
				zero.StoreArray((*[32]int16)(unsafe.Add(qp, uintptr(i)*2)))
			}
			zero.StoreArray((*[32]int16)(unsafe.Add(dp, uintptr(i)*2)))
			if classes != nil {
				archsimd.Uint8x32{}.StoreArray((*[32]byte)(unsafe.Add(clp, uintptr(i))))
			}
			continue
		}
		magnitude := a.MulHigh(acQ).Masked(active)
		t := zero.Sub(magnitude).IfElse(c.Less(zero), magnitude)
		if qcoeff != nil {
			t.StoreArray((*[32]int16)(unsafe.Add(qp, uintptr(i)*2)))
		}
		t.Mul(acD).StoreArray((*[32]int16)(unsafe.Add(dp, uintptr(i)*2)))
		last = last.Max(archsimd.LoadInt16x32Array((*[32]int16)(unsafe.Add(ip, uintptr(i)*2))).Masked(t.NotEqual(zero)))
		if classes != nil {
			lut.Permute(magnitude.Min(eleven).ToBits()).TruncToUint8().StoreArray((*[32]byte)(unsafe.Add(clp, uintptr(i))))
		}
	}
	v := last.GetLo().Max(last.GetHi())
	result := fpReduceMax16(v.GetLo().Max(v.GetHi()))
	archsimd.ClearAVXUpperBits()
	return result
}
