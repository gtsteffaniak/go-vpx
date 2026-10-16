//go:build !goexperiment.simd || (!arm64 && !amd64) || purego

package pixel

func LoopFilter(_ []byte, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int) bool { return false }
