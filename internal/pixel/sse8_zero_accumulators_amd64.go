//go:build goexperiment.simd && amd64 && !purego

package pixel

import "simd/archsimd"

const sse8VEXZeroAccumulators = true

// Distinct objects prevent a legacy XMM copy of a shared zero value.
var sse8ZeroAccumulatorLo [4]int32
var sse8ZeroAccumulatorHi [4]int32

func sse8ZeroAccumulators() (archsimd.Int32x4, archsimd.Int32x4) {
	return archsimd.LoadInt32x4Array(&sse8ZeroAccumulatorLo), archsimd.LoadInt32x4Array(&sse8ZeroAccumulatorHi)
}
