//go:build !goexperiment.simd || !arm64 || purego

package pixel

func Copy16Rows(src []byte, ss int, dst []byte, ds int) bool { return false }
