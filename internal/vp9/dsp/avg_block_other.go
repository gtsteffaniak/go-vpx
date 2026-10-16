package dsp

import "github.com/thesyncim/govpx/internal/pixel"

func avg8x8QuadAsm(src []uint8, off, stride int, out *[4]int32) bool {
	if !pixel.Enabled {
		return false
	}
	for i := range 4 {
		base := off + (i/2)*8*stride + (i%2)*8
		out[i] = int32((pixel.Sum(src[base:], stride, 8, 8) + 32) >> 6)
	}
	return true
}
