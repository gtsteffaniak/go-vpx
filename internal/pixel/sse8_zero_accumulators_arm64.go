//go:build goexperiment.simd && arm64 && !purego

package pixel

import "simd/archsimd"

const sse8VEXZeroAccumulators = false

func sse8ZeroAccumulators() (archsimd.Int32x4, archsimd.Int32x4) {
	return archsimd.Int32x4{}, archsimd.Int32x4{}
}
