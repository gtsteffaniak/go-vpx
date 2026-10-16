//go:build goexperiment.simd && amd64 && !purego

package pixel

import "simd/archsimd"

// Sixteen bilinear pixels occupy one AVX2 packet; pairwise dot products
// accumulate signed differences and squared error without narrowing.
func subpelStats16(src []byte, ss int, ref []byte, rs, w, h int, fx, fy [2]int16) (int, uint32) {
	_ = src[h*ss+w]
	_ = ref[(h-1)*rs+w-1]
	defer archsimd.ClearAVXUpperBits()
	x0, x1, y0, y1 := archsimd.BroadcastInt16x16(fx[0]), archsimd.BroadcastInt16x16(fx[1]), archsimd.BroadcastInt16x16(fy[0]), archsimd.BroadcastInt16x16(fy[1])
	round, ones := archsimd.BroadcastInt16x16(64), archsimd.BroadcastInt16x16(1)
	var sum, sse archsimd.Int32x8
	for x := 0; x < w; x += 16 {
		a := archsimd.LoadUint8x16(src[x:]).ExtendToUint16().BitsToInt16().Mul(x0).Add(archsimd.LoadUint8x16(src[x+1:]).ExtendToUint16().BitsToInt16().Mul(x1)).Add(round).ShiftAllRight(7)
		so, ro := ss+x, x
		for range h {
			b := archsimd.LoadUint8x16(src[so:]).ExtendToUint16().BitsToInt16().Mul(x0).Add(archsimd.LoadUint8x16(src[so+1:]).ExtendToUint16().BitsToInt16().Mul(x1)).Add(round).ShiftAllRight(7)
			v := a.Mul(y0).Add(b.Mul(y1)).Add(round).ShiftAllRight(7)
			a = b
			d := v.Sub(archsimd.LoadUint8x16(ref[ro:]).ExtendToUint16().BitsToInt16())
			sum = sum.Add(d.DotProductPairs(ones))
			sse = sse.Add(d.DotProductPairs(d))
			so += ss
			ro += rs
		}
	}
	return int(sum32(sum.GetLo().Add(sum.GetHi()))), uint32(sum32(sse.GetLo().Add(sse.GetHi())))
}
