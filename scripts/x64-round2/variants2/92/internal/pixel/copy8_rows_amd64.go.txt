//go:build goexperiment.simd && amd64 && !purego

package pixel

import "unsafe"

// Copy8Rows preserves libvpx v1.16.0 vp8/common/reconinter.c row order.
// Exact eight-byte row dataflow is adapted from goav1 recon_gosimd.go.
// Each full row is loaded before its store, preserving per-row memmove aliases.
func Copy8Rows(src []byte, ss int, dst []byte, ds int) bool {
	if !Enabled || (uint(ss)|uint(ds))>>30 != 0 || 7*ss+8 > len(src) || 7*ds+8 > len(dst) {
		return false
	}
	sp, dp := unsafe.Pointer(unsafe.SliceData(src)), unsafe.Pointer(unsafe.SliceData(dst))
	if sp != dp || ss != ds {
		for row := range 8 {
			v := *(*uint64)(unsafe.Add(sp, row*ss))
			*(*uint64)(unsafe.Add(dp, row*ds)) = v
		}
	}
	return true
}
