//go:build goexperiment.simd && (arm64 || amd64) && !purego

package dsp

import (
	"github.com/gtsteffaniak/go-vpx/internal/pixel"
	"unsafe"
)

// libvpx v1.16.0 reconintra4x4.c uses avg2/avg3/copies. Four quarter
// samples represent all three exactly, including the single rounding.
func intraDirectionalGo(mode int, dst []byte, ds int, above, left []byte, top byte) bool {
	if !pixel.Enabled {
		return false
	}
	// The direct formulas share their rounded averages instead of materializing
	// four permuted whole-block packets. Keep snapshot/write ordering for aliases.
	if ds >= 4 && len(dst) >= 4 && ds <= (len(dst)-4)/3 && intraBordersDisjoint(dst, above) && intraBordersDisjoint(dst, left) {
		return false
	}
	var border [16]byte
	copy(border[:8], above)
	copy(border[8:12], left)
	border[12] = top
	pixel.Intra4(dst, ds, &border, &intra4Selectors[mode])
	return true
}

func intraBordersDisjoint(dst, src []byte) bool {
	if len(src) == 0 {
		return true
	}
	a, b := uintptr(unsafe.Pointer(&dst[0])), uintptr(unsafe.Pointer(&src[0]))
	return a < b && b-a >= uintptr(len(dst)) || b < a && a-b >= uintptr(len(src))
}
