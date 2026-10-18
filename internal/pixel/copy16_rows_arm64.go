//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

// Copy16Rows follows libvpx v1.16.0 vp8/common/reconinter.c row order.
// A whole row is loaded before its store, preserving per-row memmove aliases.
func Copy16Rows(src []byte, ss int, dst []byte, ds int) bool {
	if !Enabled || (uint(ss)|uint(ds))>>30 != 0 || 15*ss+16 > len(src) || 15*ds+16 > len(dst) {
		return false
	}
	sp, dp := unsafe.Pointer(unsafe.SliceData(src)), unsafe.Pointer(unsafe.SliceData(dst))
	for row := 0; row < 16; row++ {
		v := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(sp, row*ss)))
		v.StoreArray((*[16]byte)(unsafe.Add(dp, row*ds)))
	}
	return true
}
