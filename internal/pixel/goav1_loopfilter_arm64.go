// SPDX-License-Identifier: BSD-2-Clause
// Byte-domain dataflow extracted from goav1 ff46346e73f063ff1ea3a8d7f06ee26cd9baaad4.
// See UPSTREAM.md and ../../third_party/goav1/LICENSE.
//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

var lfByteRight1 = archsimd.BroadcastInt8x16(-1)
var lfByteRight3 = archsimd.BroadcastInt8x16(-3)
var lfByteUnsignedRight1 = archsimd.BroadcastInt8x16(-1)

// Saturating the threshold sum is exact for block limits below 255. The
// complete int16 path retains the full public threshold domain.
func loopFilterWide(p []byte, q0, step, outer, n, blimit, limit, thresh, kind int) bool {
	if n <= 0 || step <= 0 || outer <= 0 || n%8 != 0 || (kind != 2 && kind != 4 && kind != 6) || blimit < 0 || blimit >= 255 || limit < 0 || limit > 255 || thresh < 0 || thresh > 255 || (step != 1 && outer != 1) {
		return false
	}
	// With positive strides these endpoints cover every exact 4/8/16-byte
	// row load and each strictly smaller modified store window below.
	left, right := 4, 3
	if kind == 2 {
		left, right = 2, 1
	}
	_ = p[q0-left*step]
	_ = p[q0+right*step+(n-1)*outer]
	chunk := 16
	// Preserve the original eight-position ordering for overlapping layouts.
	if (outer == 1 && step < 16) || (step == 1 && outer < left+right+1) {
		chunk = 8
	}
	for base := 0; base < n; base += chunk {
		lanes := min(chunk, n-base)
		if lanes != 8 && lanes != 16 {
			return false
		}
		loopFilterBytes(p, q0+base*outer, step, outer, lanes, blimit, limit, thresh, kind)
	}
	return true
}
func loadLFByte(p []byte, off, lanes int) archsimd.Uint8x16 {
	ptr := unsafe.Add(unsafe.Pointer(unsafe.SliceData(p)), off)
	if lanes == 8 {
		return loadNPtr(ptr, 8)
	}
	return archsimd.LoadUint8x16Array((*[16]byte)(ptr))
}
func loadLFColumnTile(p []byte, off, pitch, width int) (archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16) {
	ptr := unsafe.Add(unsafe.Pointer(unsafe.SliceData(p)), off)
	return transposeBytes8(loadNPtr(ptr, width), loadNPtr(unsafe.Add(ptr, pitch), width), loadNPtr(unsafe.Add(ptr, 2*pitch), width), loadNPtr(unsafe.Add(ptr, 3*pitch), width), loadNPtr(unsafe.Add(ptr, 4*pitch), width), loadNPtr(unsafe.Add(ptr, 5*pitch), width), loadNPtr(unsafe.Add(ptr, 6*pitch), width), loadNPtr(unsafe.Add(ptr, 7*pitch), width))
}
func joinLFBytes(a, b archsimd.Uint8x16) archsimd.Uint8x16 {
	return a.ReshapeToUint64s().InterleaveLo(b.ReshapeToUint64s()).ReshapeToUint8s()
}
func loadLFColumns(p []byte, off, pitch, width, lanes int) (archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16, archsimd.Uint8x16) {
	a0, a1, a2, a3, a4, a5, a6, a7 := loadLFColumnTile(p, off, pitch, width)
	if lanes == 8 {
		return a0, a1, a2, a3, a4, a5, a6, a7
	}
	b0, b1, b2, b3, b4, b5, b6, b7 := loadLFColumnTile(p, off+8*pitch, pitch, width)
	return joinLFBytes(a0, b0), joinLFBytes(a1, b1), joinLFBytes(a2, b2), joinLFBytes(a3, b3), joinLFBytes(a4, b4), joinLFBytes(a5, b5), joinLFBytes(a6, b6), joinLFBytes(a7, b7)
}
func lfAbsByte(a, b archsimd.Uint8x16) archsimd.Uint8x16 { return a.Max(b).Sub(a.Min(b)) }
func lfJoinSigned(a, b archsimd.Int8x16) archsimd.Int8x16 {
	return joinLFBytes(a.ToBits(), b.ToBits()).BitsToInt8()
}
func lfMacroCorrection(f archsimd.Int8x16, c int16) archsimd.Int8x16 {
	k, bias := archsimd.BroadcastInt16x8(c), archsimd.BroadcastInt16x8(63)
	lo := filterShift7(f.ExtendLo8ToInt16().Mul(k).Add(bias)).SaturateToInt8()
	hi := filterShift7(f.HiToLo().ExtendLo8ToInt16().Mul(k).Add(bias)).SaturateToInt8()
	return lfJoinSigned(lo, hi)
}
func loopFilterBytes(p []byte, off, step, outer, lanes, blimit, limit, thresh, kind int) {
	var p3, p2, p1, p0, q0, q1, q2, q3 archsimd.Uint8x16
	if step == 1 && outer > 1 {
		start, width := off-4, 8
		if kind == 2 {
			start, width = off-2, 4
		}
		a0, a1, a2, a3, a4, a5, a6, a7 := loadLFColumns(p, start, outer, width, lanes)
		if kind == 2 {
			p1, p0, q0, q1 = a0, a1, a2, a3
		} else {
			p3, p2, p1, p0, q0, q1, q2, q3 = a0, a1, a2, a3, a4, a5, a6, a7
		}
	} else {
		p1, p0, q0, q1 = loadLFByte(p, off-2*step, lanes), loadLFByte(p, off-step, lanes), loadLFByte(p, off, lanes), loadLFByte(p, off+step, lanes)
		if kind != 2 {
			p3, p2, q2, q3 = loadLFByte(p, off-4*step, lanes), loadLFByte(p, off-3*step, lanes), loadLFByte(p, off+2*step, lanes), loadLFByte(p, off+3*step, lanes)
		}
	}
	d1, d2, d0, d3 := lfAbsByte(p1, p0), lfAbsByte(q1, q0), lfAbsByte(p0, q0), lfAbsByte(p1, q1)
	threshold := d0.AddSaturated(d0).AddSaturated(d3.Shift(lfByteUnsignedRight1))
	need := threshold.LessEqual(archsimd.BroadcastUint8x16(byte(blimit)))
	if kind != 2 {
		delta := d1.Max(d2).Max(lfAbsByte(p3, p2)).Max(lfAbsByte(p2, p1)).Max(lfAbsByte(q2, q1)).Max(lfAbsByte(q3, q2))
		need = need.And(delta.LessEqual(archsimd.BroadcastUint8x16(byte(limit))))
	}
	bits := need.ToInt8x16().ToBits().ReshapeToUint64s()
	if bits.GetElem(0)|bits.GetElem(1) == 0 {
		return
	}
	hev := d1.Max(d2).Greater(archsimd.BroadcastUint8x16(byte(thresh)))
	bias := archsimd.BroadcastUint8x16(128)
	p1c, p0c, q0c, q1c := p1.Xor(bias).BitsToInt8(), p0.Xor(bias).BitsToInt8(), q0.Xor(bias).BitsToInt8(), q1.Xor(bias).BitsToInt8()
	outerF := p1c.SubSaturated(q1c)
	if kind == 4 {
		outerF = outerF.Masked(hev)
	}
	lo := q0.ExtendLo8ToUint16().BitsToInt16().Sub(p0.ExtendLo8ToUint16().BitsToInt16())
	hi := q0.HiToLo().ExtendLo8ToUint16().BitsToInt16().Sub(p0.HiToLo().ExtendLo8ToUint16().BitsToInt16())
	three16 := archsimd.BroadcastInt16x8(3)
	f := lfJoinSigned(lo.Mul(three16).Add(outerF.ExtendLo8ToInt16()).SaturateToInt8(), hi.Mul(three16).Add(outerF.HiToLo().ExtendLo8ToInt16()).SaturateToInt8()).Masked(need)
	three, four := archsimd.BroadcastInt8x16(3), archsimd.BroadcastInt8x16(4)
	z := archsimd.Uint8x16{}
	if kind == 6 {
		fh := f.Masked(hev)
		ps0 := p0c.AddSaturated(fh.AddSaturated(three).Shift(lfByteRight3))
		qs0 := q0c.SubSaturated(fh.AddSaturated(four).Shift(lfByteRight3))
		ff := f.Masked(hev.Not())
		u := lfMacroCorrection(ff, 27)
		np0, nq0 := ps0.AddSaturated(u).ToBits().Xor(bias), qs0.SubSaturated(u).ToBits().Xor(bias)
		u = lfMacroCorrection(ff, 18)
		np1, nq1 := p1c.AddSaturated(u).ToBits().Xor(bias), q1c.SubSaturated(u).ToBits().Xor(bias)
		u = lfMacroCorrection(ff, 9)
		np2, nq2 := p2.Xor(bias).BitsToInt8().AddSaturated(u).ToBits().Xor(bias), q2.Xor(bias).BitsToInt8().SubSaturated(u).ToBits().Xor(bias)
		storeLFBytes(p, off-3*step, step, outer, lanes, 6, np2, np1, np0, nq0, nq1, nq2)
		return
	}
	f1, f2 := f.AddSaturated(four).Shift(lfByteRight3), f.AddSaturated(three).Shift(lfByteRight3)
	np0, nq0 := p0c.AddSaturated(f2).ToBits().Xor(bias), q0c.SubSaturated(f1).ToBits().Xor(bias)
	if kind == 2 {
		storeLFBytes(p, off-step, step, outer, lanes, 2, np0, nq0, z, z, z, z)
		return
	}
	u := f1.Add(archsimd.BroadcastInt8x16(1)).Shift(lfByteRight1).Masked(hev.Not())
	np1, nq1 := p1c.AddSaturated(u).ToBits().Xor(bias), q1c.SubSaturated(u).ToBits().Xor(bias)
	storeLFBytes(p, off-2*step, step, outer, lanes, 4, np1, np0, nq0, nq1, z, z)
}
func storeLFBytes(p []byte, off, step, outer, lanes, count int, a, b, c, d, e, f archsimd.Uint8x16) {
	if outer == 1 {
		ptr := unsafe.Pointer(unsafe.SliceData(p))
		store := func(at int, v archsimd.Uint8x16) {
			dest := unsafe.Add(ptr, at)
			if lanes == 16 {
				v.StoreArray((*[16]byte)(dest))
			} else {
				*(*uint64)(dest) = v.ReshapeToUint64s().GetElem(0)
			}
		}
		store(off, a)
		store(off+step, b)
		if count >= 4 {
			store(off+2*step, c)
			store(off+3*step, d)
		}
		if count == 6 {
			store(off+4*step, e)
			store(off+5*step, f)
		}
		return
	}
	storeLFByteTile(p, off, outer, count, a, b, c, d, e, f)
	if lanes == 16 {
		storeLFByteTile(p, off+8*outer, outer, count, a.HiToLo(), b.HiToLo(), c.HiToLo(), d.HiToLo(), e.HiToLo(), f.HiToLo())
	}
}
func storeLFByteTile(p []byte, off, pitch, count int, a, b, c, d, e, f archsimd.Uint8x16) {
	z := archsimd.Uint8x16{}
	r0, r1, r2, r3, r4, r5, r6, r7 := transposeBytes8(a, b, c, d, e, f, z, z)
	rows := [8]uint64{r0.ReshapeToUint64s().GetElem(0), r1.ReshapeToUint64s().GetElem(0), r2.ReshapeToUint64s().GetElem(0), r3.ReshapeToUint64s().GetElem(0), r4.ReshapeToUint64s().GetElem(0), r5.ReshapeToUint64s().GetElem(0), r6.ReshapeToUint64s().GetElem(0), r7.ReshapeToUint64s().GetElem(0)}
	ptr := unsafe.Add(unsafe.Pointer(unsafe.SliceData(p)), off)
	// These unaligned stores touch exactly the original modified columns.
	// The wide entry point checks their entire window before any load/store.
	switch count {
	case 2:
		for y, v := range rows {
			*(*uint16)(unsafe.Add(ptr, y*pitch)) = uint16(v)
		}
	case 4:
		for y, v := range rows {
			*(*uint32)(unsafe.Add(ptr, y*pitch)) = uint32(v)
		}
	case 6:
		for y, v := range rows {
			dest := unsafe.Add(ptr, y*pitch)
			*(*uint32)(dest) = uint32(v)
			*(*uint16)(unsafe.Add(dest, 4)) = uint16(v >> 32)
		}
	}
}
