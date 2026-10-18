//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

var fpClassLUT = [16]byte{0, 1, 2, 3, 3, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5}

func quantizeFPBlockSIMD(coeff, iscan, qcoeff, dqcoeff []int16, n int, r, q, d [2]int16, classes []byte) int {

	if n >= 64 && n <= 1024 && n&15 == 0 && Enabled && archsimd.X86.AVX2() && fp16WindowsDisjoint(coeff, iscan, qcoeff, dqcoeff, classes, n) {
		cp, ip, qp, dp := unsafe.Pointer(unsafe.SliceData(coeff)), unsafe.Pointer(unsafe.SliceData(iscan)), unsafe.Pointer(unsafe.SliceData(qcoeff)), unsafe.Pointer(unsafe.SliceData(dqcoeff))
		clp := unsafe.Pointer(unsafe.SliceData(classes))
		max16 := archsimd.BroadcastUint16x16(32767)
		zero := archsimd.Int16x16{}
		last := zero
		var lut archsimd.Uint8x16
		if classes != nil {
			lut = archsimd.LoadUint8x16Array(&fpClassLUT)
		}
		eleven := archsimd.BroadcastInt16x8(11)
		// First packet owns DC/AC vectors. They die before AC-only loop tables.
		{
			i := 0
			vr, vq, vd := archsimd.BroadcastInt16x16(r[1]).SetLo(archsimd.BroadcastInt16x8(r[1]).SetElem(0, r[0])), archsimd.BroadcastInt16x16(q[1]).SetLo(archsimd.BroadcastInt16x8(q[1]).SetElem(0, q[0])), archsimd.BroadcastInt16x16(d[1]).SetLo(archsimd.BroadcastInt16x8(d[1]).SetElem(0, d[0]))
			c := archsimd.LoadInt16x16Array((*[16]int16)(unsafe.Add(cp, uintptr(i)*2)))
			a := c.Abs().ToBits().Add(vr.ToBits()).Min(max16).BitsToInt16()
			magnitude := a.MulHigh(vq).Masked(a.GreaterEqual(vd))
			sign := c.ShiftAllRight(15)
			t := magnitude.Xor(sign).Sub(sign)
			if qcoeff != nil {
				t.StoreArray((*[16]int16)(unsafe.Add(qp, uintptr(i)*2)))
			}
			t.Mul(vd).StoreArray((*[16]int16)(unsafe.Add(dp, uintptr(i)*2)))
			last = last.Max(archsimd.LoadInt16x16Array((*[16]int16)(unsafe.Add(ip, uintptr(i)*2))).Masked(t.NotEqual(zero)))
			if classes != nil {
				lo := fpClassIndices(magnitude.GetLo().Min(eleven))
				hi := fpClassIndices(magnitude.GetHi().Min(eleven))
				*(*uint64)(unsafe.Add(clp, uintptr(i))) = permute16(lut, lo).ReshapeToUint64s().GetElem(0)
				*(*uint64)(unsafe.Add(clp, uintptr(i+8))) = permute16(lut, hi).ReshapeToUint64s().GetElem(0)
			}
		}
		acR, acQ, acD := archsimd.BroadcastInt16x16(r[1]), archsimd.BroadcastInt16x16(q[1]), archsimd.BroadcastInt16x16(d[1])
		for i := 16; i < n; i += 16 {
			{

				vr, vq, vd := acR, acQ, acD
				c := archsimd.LoadInt16x16Array((*[16]int16)(unsafe.Add(cp, uintptr(i)*2)))
				a := c.Abs().ToBits().Add(vr.ToBits()).Min(max16).BitsToInt16()
				magnitude := a.MulHigh(vq).Masked(a.GreaterEqual(vd))
				sign := c.ShiftAllRight(15)
				t := magnitude.Xor(sign).Sub(sign)
				if qcoeff != nil {
					t.StoreArray((*[16]int16)(unsafe.Add(qp, uintptr(i)*2)))
				}
				t.Mul(vd).StoreArray((*[16]int16)(unsafe.Add(dp, uintptr(i)*2)))
				last = last.Max(archsimd.LoadInt16x16Array((*[16]int16)(unsafe.Add(ip, uintptr(i)*2))).Masked(t.NotEqual(zero)))
				if classes != nil {
					lo := fpClassIndices(magnitude.GetLo().Min(eleven))
					hi := fpClassIndices(magnitude.GetHi().Min(eleven))
					*(*uint64)(unsafe.Add(clp, uintptr(i))) = permute16(lut, lo).ReshapeToUint64s().GetElem(0)
					*(*uint64)(unsafe.Add(clp, uintptr(i+8))) = permute16(lut, hi).ReshapeToUint64s().GetElem(0)
				}
			}
		}
		result := fpReduceMax16(last.GetLo().Max(last.GetHi()))
		archsimd.ClearAVXUpperBits()
		return result
	}
	acR, acQ, acD := archsimd.BroadcastInt16x8(r[1]), archsimd.BroadcastInt16x8(q[1]), archsimd.BroadcastInt16x8(d[1])
	dcR, dcQ, dcD := acR.SetElem(0, r[0]), acQ.SetElem(0, q[0]), acD.SetElem(0, d[0])
	max16 := archsimd.BroadcastUint16x8(32767)
	zero := archsimd.Int16x8{}
	last := zero
	lut := archsimd.LoadUint8x16Array(&fpClassLUT)
	i := 0
	for ; i+8 <= n; i += 8 {
		vr, vq, vd := acR, acQ, acD
		if i == 0 {
			vr, vq, vd = dcR, dcQ, dcD
		}
		c := archsimd.LoadInt16x8(coeff[i:])
		// Unsigned magnitude handles -32768; adding a nonnegative int16 round
		// fits uint16 exactly, then clips to the scalar signed-int16 bound.
		a := c.Abs().ToBits().Add(vr.ToBits()).Min(max16).BitsToInt16()
		t := fpMulHigh16(a, vq).Masked(a.GreaterEqual(vd))
		magnitude := t
		t = zero.Sub(t).IfElse(c.Less(zero), t)
		if qcoeff != nil {
			t.Store(qcoeff[i:])
		}
		t.Mul(vd).Store(dqcoeff[i:])
		last = last.Max(archsimd.LoadInt16x8(iscan[i:]).Masked(t.NotEqual(zero)))
		if classes != nil {
			indices := fpClassIndices(magnitude.Min(archsimd.BroadcastInt16x8(11)))
			storeBytes8(classes[i:], permute16(lut, indices))
		}
	}
	return max(fpReduceMax16(last), quantizeFPBlockTail(coeff, iscan, qcoeff, dqcoeff, i, n, r, q, d, classes))
}

// fp16WindowsDisjoint compares touched byte ranges without adding end pointers.
// The public entry has already checked every required length; n is bounded.
// Even read-only input aliases decline, keeping this specialization conservative.
func fp16WindowsDisjoint(coeff, iscan, qcoeff, dqcoeff []int16, classes []byte, n int) bool {
	type span struct{ p, bytes uintptr }
	spans := [5]span{
		{uintptr(unsafe.Pointer(unsafe.SliceData(coeff))), uintptr(n) * 2},
		{uintptr(unsafe.Pointer(unsafe.SliceData(iscan))), uintptr(n) * 2},
		{uintptr(unsafe.Pointer(unsafe.SliceData(dqcoeff))), uintptr(n) * 2},
	}
	count := 3
	if qcoeff != nil {
		spans[count] = span{uintptr(unsafe.Pointer(unsafe.SliceData(qcoeff))), uintptr(n) * 2}
		count++
	}
	if classes != nil {
		spans[count] = span{uintptr(unsafe.Pointer(unsafe.SliceData(classes))), uintptr(n)}
		count++
	}
	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			a, b := spans[i], spans[j]
			if a.p <= b.p && b.p-a.p < a.bytes || b.p < a.p && a.p-b.p < b.bytes {
				return false
			}
		}
	}
	return true
}
