package dsp

import "github.com/thesyncim/govpx/internal/pixel"

// Pure-Go fallback for VP8 intra-prediction primitives. Mirrors
// libvpx v1.16.0 vp8/common/reconintra.c and vpx_dsp/intrapred.c
// semantics.

func intraDCPredict16x16(dst []byte, dstStride int, above []byte, left []byte, upAvailable bool, leftAvailable bool) {
	intraDCPredictScalar(dst, dstStride, above, left, 16, upAvailable, leftAvailable)
}

func intraDCPredict8x8(dst []byte, dstStride int, above []byte, left []byte, upAvailable bool, leftAvailable bool) {
	intraDCPredictScalar(dst, dstStride, above, left, 8, upAvailable, leftAvailable)
}

func intraVerticalPredict16x16(dst []byte, dstStride int, above []byte) {
	intraVerticalPredictScalar(dst, dstStride, above, 16)
}

func intraVerticalPredict8x8(dst []byte, dstStride int, above []byte) {
	intraVerticalPredictScalar(dst, dstStride, above, 8)
}

func intraHorizontalPredict16x16(dst []byte, dstStride int, left []byte) {
	intraHorizontalPredictScalar(dst, dstStride, left, 16)
}

func intraHorizontalPredict8x8(dst []byte, dstStride int, left []byte) {
	intraHorizontalPredictScalar(dst, dstStride, left, 8)
}

func intraTMPredict16x16(dst []byte, dstStride int, above []byte, left []byte, topLeft byte) {
	if pixel.Enabled {
		pixel.TMPredict(dst, dstStride, above, left, topLeft, 16)
	} else {
		intraTMPredictScalar(dst, dstStride, above, left, topLeft, 16)
	}
}

func intraTMPredict8x8(dst []byte, dstStride int, above []byte, left []byte, topLeft byte) {
	if pixel.Enabled {
		pixel.TMPredict(dst, dstStride, above, left, topLeft, 8)
	} else {
		intraTMPredictScalar(dst, dstStride, above, left, topLeft, 8)
	}
}
