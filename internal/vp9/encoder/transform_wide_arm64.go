//go:build goexperiment.simd && arm64 && !purego

package encoder

func forwardDCTWideGo(input []int16, stride int, output []int16, n int) bool { return false }
