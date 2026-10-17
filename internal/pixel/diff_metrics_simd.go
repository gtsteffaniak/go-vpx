//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import "simd/archsimd"

func diffMetricsSIMD(src []byte, ss int, ref []byte, rs, w, h int) (sad, sum int, sse uint32) {
	var abs archsimd.Uint32x4
	var squares, signed archsimd.Int32x4
	for y := range h {
		a, b := src[y*ss:y*ss+w], ref[y*rs:y*rs+w]
		x := 0
		for ; x+16 <= w; x += 16 {
			av, bv := archsimd.LoadUint8x16(a[x:]), archsimd.LoadUint8x16(b[x:])
			abs = sadAdd(abs, av, bv)
			squares, signed = pixelDiffAcc(squares, signed, av, bv)
		}
		if x+8 <= w {
			av, bv := load8(a[x:]), load8(b[x:])
			abs = sadAdd(abs, av, bv)
			squares, signed = pixelDiffAcc(squares, signed, av, bv)
			x += 8
		}
		for ; x < w; x++ {
			d := int(a[x]) - int(b[x])
			sad += max(d, -d)
			sum += d
			sse += uint32(d * d)
		}
	}
	return sad + int(sumU32(abs)), sum + int(sum32(signed)), sse + uint32(sum32(squares))
}
