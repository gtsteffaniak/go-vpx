//go:build !goexperiment.simd || (!arm64 && !amd64) || purego

package dsp

func addResidual4Go(_ []byte, _ int, _ *[16]int16) bool { return false }
