//go:build !goexperiment.simd || (!arm64 && !amd64) || purego

package pixel

func SAD(src []byte, ss int, ref []byte, rs, w, h int) (sum int) {
	for y := range h {
		for x := range w {
			d := int(src[y*ss+x]) - int(ref[y*rs+x])
			sum += max(d, -d)
		}
	}
	return
}
func Stats(src []byte, ss int, ref []byte, rs, w, h int) (sum int, sse uint32) {
	for y := range h {
		for x := range w {
			d := int(src[y*ss+x]) - int(ref[y*rs+x])
			sum += d
			sse += uint32(d * d)
		}
	}
	return
}
func Subtract(src []byte, ss int, pred []byte, ps int, out []int16, os, w, h int) (nonzero bool) {
	for y := range h {
		for x := range w {
			d := int16(src[y*ss+x]) - int16(pred[y*ps+x])
			out[y*os+x] = d
			nonzero = nonzero || d != 0
		}
	}
	return
}
func Filter(src []byte, ss int, dst []byte, ds, w, h, step int, taps []int16) {
	for y := range h {
		for x := range w {
			v := 64
			for k, c := range taps {
				v += int(src[y*ss+x+k*step]) * int(c)
			}
			dst[y*ds+x] = byte(min(max(v>>7, 0), 255))
		}
	}
}
func Average(src []byte, ss int, dst []byte, ds, w, h int) {
	for y := range h {
		for x := range w {
			dst[y*ds+x] = byte((int(src[y*ss+x]) + int(dst[y*ds+x]) + 1) >> 1)
		}
	}
}
func Error(a, b []int16, n int) (sum uint64) {
	for i := range n {
		d := int64(a[i]) - int64(b[i])
		sum += uint64(d * d)
	}
	return
}

func Quantize(coeff, iscan, qcoeff, dqcoeff []int16, n, zbin, round, quant, shift, dequant int, fp bool) (eob int) {
	for i := range n {
		c := int(coeff[i])
		a := max(c, -c)
		q := 0
		if (fp && a+round >= dequant) || (!fp && a >= zbin) {
			v := min(max(a+round, -32768), 32767)
			if fp {
				q = v * quant >> 16
			} else {
				q = (((v * quant >> 16) + v) * shift) >> 16
			}
			if c < 0 {
				q = -q
			}
		}
		if qcoeff != nil {
			qcoeff[i] = int16(q)
		}
		dqcoeff[i] = int16(q * dequant)
		if q != 0 {
			eob = max(eob, int(iscan[i]))
		}
	}
	return
}

func Gather(src []byte, ss int, pred []byte, ps int, out []int16, w int) {
	for y := range w {
		for x := range w {
			off := ((y/4)*(w/4)+x/4)*16 + (y%4)*4 + x%4
			out[off] = int16(src[y*ss+x]) - int16(pred[y*ps+x])
		}
	}
}
func Temporal(src []byte, ss int, pred []byte, ps, w, strength, weight int, acc, count []uint32) {
	r := 0
	if strength > 0 {
		r = 1 << (strength - 1)
	}
	for y := range w {
		for x := range w {
			d := int(src[y*ss+x]) - int(pred[y*ps+x])
			m := (16 - min((d*d*3+r)>>uint(strength), 16)) * weight
			i := y*w + x
			count[i] += uint32(m)
			acc[i] += uint32(m) * uint32(pred[y*ps+x])
		}
	}
}
func VP8Quant(c, r, q, d, qout, dout []int16) {
	for i := range 16 {
		v := int(c[i])
		a := max(v, -v)
		t := (a + int(r[i])) * int(q[i]) >> 16
		if v < 0 {
			t = -t
		}
		qout[i] = int16(t)
		dout[i] = int16(t * int(d[i]))
	}
}
func AddConstant(pred []byte, ps int, dst []byte, ds, w, h, delta int) {
	for y := range h {
		for x := range w {
			dst[y*ds+x] = byte(min(max(int(pred[y*ps+x])+delta, 0), 255))
		}
	}
}
func TMPredict(dst []byte, ds int, above, left []byte, top byte, w int) {
	for y := range w {
		for x := range w {
			dst[y*ds+x] = byte(min(max(int(above[x])+int(left[y])-int(top), 0), 255))
		}
	}
}
