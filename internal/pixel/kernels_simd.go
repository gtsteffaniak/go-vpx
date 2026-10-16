//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import "simd/archsimd"

func SAD(src []byte, ss int, ref []byte, rs, w, h int) int {
	var acc archsimd.Uint32x4
	tail := 0
	for y := range h {
		a, b := src[y*ss:], ref[y*rs:]
		x := 0
		for ; x+16 <= w; x += 16 {
			acc = sadAdd(acc, archsimd.LoadUint8x16(a[x:]), archsimd.LoadUint8x16(b[x:]))
		}
		if x+8 <= w {
			acc = sadAdd(acc, load8(a[x:]), load8(b[x:]))
			x += 8
		}
		for ; x < w; x++ {
			d := int(a[x]) - int(b[x])
			tail += max(d, -d)
		}
	}
	return int(sumU32(acc)) + tail
}
func Stats(src []byte, ss int, ref []byte, rs, w, h int) (int, uint32) {
	var sums, squares archsimd.Int32x4
	tailSum := 0
	var tailSSE uint32
	for y := range h {
		a, b := src[y*ss:], ref[y*rs:]
		x := 0
		for ; x+8 <= w; x += 8 {
			d := widen8(load8(a[x:])).Sub(widen8(load8(b[x:])))
			lo := d.ExtendLo4ToInt32()
			hi := high4(d)
			sums = sums.Add(lo).Add(hi)
			squares = squares.Add(lo.Mul(lo)).Add(hi.Mul(hi))
		}
		for ; x < w; x++ {
			d := int(a[x]) - int(b[x])
			tailSum += d
			tailSSE += uint32(d * d)
		}
	}
	return int(sum32(sums)) + tailSum, uint32(sum32(squares)) + tailSSE
}
func Subtract(src []byte, ss int, pred []byte, ps int, out []int16, os, w, h int) bool {
	var bits archsimd.Int16x8
	tail := int16(0)
	for y := range h {
		a, b, dst := src[y*ss:], pred[y*ps:], out[y*os:]
		x := 0
		for ; x+8 <= w; x += 8 {
			d := widen8(load8(a[x:])).Sub(widen8(load8(b[x:])))
			d.Store(dst[x:])
			bits = bits.Or(d)
		}
		for ; x < w; x++ {
			d := int16(a[x]) - int16(b[x])
			dst[x] = d
			tail |= d
		}
	}
	return bits.ToBits().ReshapeToUint64s().GetElem(0)|bits.ToBits().ReshapeToUint64s().GetElem(1) != 0 || tail != 0
}

// Filter uses widening int32 products: negative filter lobes cannot overflow
// or saturate before the codec's single rounding and clipping operation.
func Filter(src []byte, ss int, dst []byte, ds, w, h, step int, taps []int16) {
	var c [8]archsimd.Int16x8
	for i, v := range taps {
		c[i] = archsimd.BroadcastInt16x8(v)
	}
	bias := archsimd.BroadcastInt32x4(64)
	for y := range h {
		a, b := src[y*ss:], dst[y*ds:]
		x := 0
		for ; x+8 <= w; x += 8 {
			lo, hi := bias, bias
			for k := range taps {
				v := widen8(load8(a[x+k*step:]))
				lo, hi = mac(lo, hi, v, c[k])
			}
			store8(b[x:], lo.ShiftAllRight(7), hi.ShiftAllRight(7))
		}
		for ; x < w; x++ {
			v := 64
			for k, c := range taps {
				v += int(a[x+k*step]) * int(c)
			}
			b[x] = byte(min(max(v>>7, 0), 255))
		}
	}
}
func Average(src []byte, ss int, dst []byte, ds, w, h int) {
	one := archsimd.BroadcastInt16x8(1)
	for y := range h {
		a, b := src[y*ss:], dst[y*ds:]
		x := 0
		for ; x+8 <= w; x += 8 {
			v := widen8(load8(a[x:])).Add(widen8(load8(b[x:]))).Add(one).ShiftAllRight(1)
			storeBytes8(b[x:], packBytes(v))
		}
		for ; x < w; x++ {
			b[x] = byte((int(a[x]) + int(b[x]) + 1) >> 1)
		}
	}
}

// Error widens before subtraction, preserving the full signed int16 domain.
// Four squares fit in uint64; no int32 horizontal accumulation can overflow.
func Error(a, b []int16, n int) (sum uint64) {
	i := 0
	for ; i+8 <= n; i += 8 {
		va, vb := archsimd.LoadInt16x8(a[i:]), archsimd.LoadInt16x8(b[i:])
		lo := va.ExtendLo4ToInt32().Sub(vb.ExtendLo4ToInt32())
		hi := high4(va).Sub(high4(vb))
		sum += sumSquares(lo) + sumSquares(hi)
	}
	for ; i < n; i++ {
		d := int64(a[i]) - int64(b[i])
		sum += uint64(d * d)
	}
	return
}

func Quantize(coeff, iscan, qcoeff, dqcoeff []int16, n, zbin, round, quant, shift, dequant int, fp bool) (eob int) {
	z, r, q, s, d := archsimd.BroadcastInt32x4(int32(zbin)), archsimd.BroadcastInt32x4(int32(round)), archsimd.BroadcastInt32x4(int32(quant)), archsimd.BroadcastInt32x4(int32(shift)), archsimd.BroadcastInt32x4(int32(dequant))
	zero := archsimd.Int32x4{}
	max16 := archsimd.BroadcastInt32x4(32767)
	min16 := archsimd.BroadcastInt32x4(-32768)
	var last archsimd.Int32x4
	i := 0
	for ; i+4 <= n; i += 4 {
		c := Load4(coeff[i:])
		a := c.Abs()
		v := a.Add(r)
		keep := a.GreaterEqual(z)
		if fp {
			keep = v.GreaterEqual(d)
		}
		v = v.Min(max16).Max(min16)
		t := v.Mul(q).ShiftAllRight(16)
		if !fp {
			t = t.Add(v).Mul(s).ShiftAllRight(16)
		}
		t = t.Masked(keep)
		t = zero.Sub(t).IfElse(c.Less(zero), t)
		if qcoeff != nil {
			Store4(qcoeff[i:], t)
		}
		Store4(dqcoeff[i:], t.Mul(d))
		scan := Load4(iscan[i:])
		last = last.Max(scan.Masked(t.NotEqual(zero)))
	}
	eob = int(Max4(last))
	for ; i < n; i++ {
		c := int(coeff[i])
		a := max(c, -c)
		t := 0
		if (fp && a+round >= dequant) || (!fp && a >= zbin) {
			v := min(max(a+round, -32768), 32767)
			if fp {
				t = v * quant >> 16
			} else {
				t = (((v * quant >> 16) + v) * shift) >> 16
			}
			if c < 0 {
				t = -t
			}
		}
		if qcoeff != nil {
			qcoeff[i] = int16(t)
		}
		dqcoeff[i] = int16(t * dequant)
		if t != 0 {
			eob = max(eob, int(iscan[i]))
		}
	}
	return
}

// Gather writes residuals in VP8's 4x4 block scan order, without staging
// complete macroblocks or allocating packing buffers.
func Gather(src []byte, ss int, pred []byte, ps int, out []int16, w int) {
	for y := range w {
		a, b := src[y*ss:], pred[y*ps:]
		for x := 0; x < w; x += 8 {
			v := widen8(load8(a[x:])).Sub(widen8(load8(b[x:]))).ToBits().ReshapeToUint64s()
			off := ((y/4)*(w/4)+x/4)*16 + (y%4)*4
			put4(out[off:], v.GetElem(0))
			put4(out[off+16:], v.GetElem(1))
		}
	}
}
func Temporal(src []byte, ss int, pred []byte, ps, w, strength, weight int, acc, count []uint32) {
	rounding := int32(0)
	if strength > 0 {
		rounding = 1 << (strength - 1)
	}
	r := archsimd.BroadcastInt32x4(rounding)
	three := archsimd.BroadcastInt32x4(3)
	sixteen := archsimd.BroadcastInt32x4(16)
	wt := archsimd.BroadcastInt32x4(int32(weight))
	for y := range w {
		for x := 0; x < w; x += 4 {
			a, b := LoadPixels4(src[y*ss+x:]), LoadPixels4(pred[y*ps+x:])
			d := a.Sub(b)
			m := sixteen.Sub(d.Mul(d).Mul(three).Add(r).ShiftAllRight(uint64(strength)).Min(sixteen)).Mul(wt).ToBits()
			i := y*w + x
			m.Add(archsimd.LoadUint32x4(count[i:])).Store(count[i:])
			m.Mul(b.ToBits()).Add(archsimd.LoadUint32x4(acc[i:])).Store(acc[i:])
		}
	}
}
func VP8Quant(coeff, round, quant, dequant []int16, qcoeff, dqcoeff []int16) {
	for i := 0; i < 16; i += 4 {
		c := Load4(coeff[i:])
		sign := c.ShiftAllRight(31)
		q := c.Abs().Add(Load4(round[i:])).Mul(Load4(quant[i:])).ShiftAllRight(16)
		q = q.Xor(sign).Sub(sign)
		Store4(qcoeff[i:], q)
		Store4(dqcoeff[i:], q.Mul(Load4(dequant[i:])))
	}
}
func AddConstant(pred []byte, ps int, dst []byte, ds, w, h, delta int) {
	v := archsimd.BroadcastInt32x4(int32(delta))
	for y := range h {
		for x := 0; x < w; x += 4 {
			StorePixels4(dst[y*ds+x:], LoadPixels4(pred[y*ps+x:]).Add(v))
		}
	}
}
func TMPredict(dst []byte, ds int, above, left []byte, top byte, w int) {
	t := archsimd.BroadcastInt32x4(int32(top))
	for y := range w {
		v := archsimd.BroadcastInt32x4(int32(left[y])).Sub(t)
		for x := 0; x < w; x += 4 {
			StorePixels4(dst[y*ds+x:], LoadPixels4(above[x:]).Add(v))
		}
	}
}
