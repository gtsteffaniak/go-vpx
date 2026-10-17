//go:build goexperiment.simd && arm64 && !purego

package pixel

import "simd/archsimd"

// Complete sixteen-byte rows feed two NEON int16 packets. The horizontal
// pair stays in registers for reuse by the following output row.
func subpelStats16(src []byte, ss int, ref []byte, rs, w, h int, fx, fy [2]int16) (int, uint32) {
	_ = src[h*ss+w]
	_ = ref[(h-1)*rs+w-1]
	x0, x1, y0, y1 := archsimd.BroadcastInt16x8(fx[0]), archsimd.BroadcastInt16x8(fx[1]), archsimd.BroadcastInt16x8(fy[0]), archsimd.BroadcastInt16x8(fy[1])
	round := archsimd.BroadcastInt16x8(64)
	var sum, sse archsimd.Int32x4
	for x := 0; x < w; x += 16 {
		r0, r1 := archsimd.LoadUint8x16(src[x:]), archsimd.LoadUint8x16(src[x+1:])
		a0 := filterShift7(mulAdd16(widen8(r0).Mul(x0).Add(round), widen8(r1), x1))
		a1 := filterShift7(mulAdd16(widen8(highBytes(r0)).Mul(x0).Add(round), widen8(highBytes(r1)), x1))
		so, ro := ss+x, x
		for range h {
			r0, r1 = archsimd.LoadUint8x16(src[so:]), archsimd.LoadUint8x16(src[so+1:])
			b0 := filterShift7(mulAdd16(widen8(r0).Mul(x0).Add(round), widen8(r1), x1))
			b1 := filterShift7(mulAdd16(widen8(highBytes(r0)).Mul(x0).Add(round), widen8(highBytes(r1)), x1))
			v0 := filterShift7(mulAdd16(a0.Mul(y0).Add(round), b0, y1))
			v1 := filterShift7(mulAdd16(a1.Mul(y0).Add(round), b1, y1))
			a0, a1 = b0, b1
			r := archsimd.LoadUint8x16(ref[ro:])
			d0, d1 := v0.Sub(widen8(r)), v1.Sub(widen8(highBytes(r)))
			sum = sum.Add(sumDiff16(d0.Add(d1)))
			sse = sse.Add(square8(d0)).Add(square8(d1))
			so += ss
			ro += rs
		}
	}
	return int(sum32(sum)), uint32(sum32(sse))
}
