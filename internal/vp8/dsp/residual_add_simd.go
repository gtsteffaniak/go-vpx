//go:build goexperiment.simd && (arm64 || amd64) && !purego

package dsp

import (
	"github.com/thesyncim/govpx/internal/pixel"
	"unsafe"
)

func addResidual4Go(dst []byte, ds int, residual *[16]int16) bool {
	if !pixel.Enabled || ds < 4 || ds > (len(dst)-4)/3 || len(dst) < 4 {
		return false
	}
	// Scalar updates must stay ordered if the residual aliases the byte destination.
	a, b := uintptr(unsafe.Pointer(&dst[0])), uintptr(unsafe.Pointer(residual))
	if a <= b && b-a < uintptr(len(dst)) || b < a && a-b < 32 {
		return false
	}
	for y := 0; y < 4; y++ {
		pixel.StorePixels4(dst[y*ds:], pixel.LoadPixels4(dst[y*ds:]).Add(pixel.Load4(residual[y*4:])))
	}
	return true
}
