package dsp

// Scalar fast paths for the unscaled (step 16) 8-tap convolve, used when the
// SIMD kernels are not built in. They compute exactly what convolveHoriz,
// convolveVert and their averaging variants compute; the gains come from picking the filter once per
// block, unrolling the taps, slicing each row up front so the compiler can
// drop bounds checks, and walking the vertical pass row by row instead of
// column by column.

func convolveTap8(s []byte, f *[8]int16) uint8 {
	_ = s[7]
	sum := int32(s[0])*int32(f[0]) + int32(s[1])*int32(f[1]) +
		int32(s[2])*int32(f[2]) + int32(s[3])*int32(f[3]) +
		int32(s[4])*int32(f[4]) + int32(s[5])*int32(f[5]) +
		int32(s[6])*int32(f[6]) + int32(s[7])*int32(f[7])
	return clipPixel((sum + 64) >> 7)
}

// convolveHorizUnscaled: src row y starts at srcStart+y*srcStride, where
// srcStart already includes the -3 tap offset and the integer x position.
func convolveHorizUnscaled(src []byte, srcStride int, dst []byte, dstStride int,
	f *[8]int16, w, h, srcStart int, avg bool) {
	for y := 0; y < h; y++ {
		row := src[srcStart+y*srcStride : srcStart+y*srcStride+w+7]
		out := dst[y*dstStride : y*dstStride+w]
		if avg {
			for x := range out {
				out[x] = uint8((int(out[x]) + int(convolveTap8(row[x:x+8], f)) + 1) >> 1)
			}
			continue
		}
		for x := range out {
			out[x] = convolveTap8(row[x:x+8], f)
		}
	}
}

// convolveVertUnscaled: src row 0 of the 8-row window for output row y is
// srcStart+y*srcStride.
func convolveVertUnscaled(src []byte, srcStride int, dst []byte, dstStride int,
	f *[8]int16, w, h, srcStart int, avg bool) {
	f0, f1, f2, f3 := int32(f[0]), int32(f[1]), int32(f[2]), int32(f[3])
	f4, f5, f6, f7 := int32(f[4]), int32(f[5]), int32(f[6]), int32(f[7])
	for y := 0; y < h; y++ {
		b := srcStart + y*srcStride
		r0 := src[b : b+w]
		r1 := src[b+srcStride : b+srcStride+w]
		r2 := src[b+2*srcStride : b+2*srcStride+w]
		r3 := src[b+3*srcStride : b+3*srcStride+w]
		r4 := src[b+4*srcStride : b+4*srcStride+w]
		r5 := src[b+5*srcStride : b+5*srcStride+w]
		r6 := src[b+6*srcStride : b+6*srcStride+w]
		r7 := src[b+7*srcStride : b+7*srcStride+w]
		out := dst[y*dstStride : y*dstStride+w]
		for x := range out {
			sum := int32(r0[x])*f0 + int32(r1[x])*f1 + int32(r2[x])*f2 + int32(r3[x])*f3 +
				int32(r4[x])*f4 + int32(r5[x])*f5 + int32(r6[x])*f6 + int32(r7[x])*f7
			p := clipPixel((sum + 64) >> 7)
			if avg {
				p = uint8((int(out[x]) + int(p) + 1) >> 1)
			}
			out[x] = p
		}
	}
}
