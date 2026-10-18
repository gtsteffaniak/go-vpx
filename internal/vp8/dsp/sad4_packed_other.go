//go:build !goexperiment.simd || !arm64 || purego

package dsp

func sad4RowsGo(_ []byte, _ int, _ []byte, _ int) (int, bool) { return 0, false }
