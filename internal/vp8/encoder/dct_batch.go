package encoder

// Whole-MB batched forward 4x4 DCT entry point. Mirrors libvpx v1.16.0
// vp8/encoder/encodemb.c vp8_transform_mb / vp8_transform_mbuv.
//
// ForwardDCT4x4Batch processes count consecutive 4x4 int16 blocks. Native
// kernels transform two blocks across eight lanes, with per-block handling
// for the remainder and scalar handling outside the VP8 residual range.
// Input and output may be the same slice. Results match the scalar DCT.
func ForwardDCT4x4Batch(input []int16, output []int16, count int) {
	forwardDCT4x4BatchSIMD(input, output, count)
}
