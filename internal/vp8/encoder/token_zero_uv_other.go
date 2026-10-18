//go:build !goexperiment.simd || !arm64 || purego || govpx_oracle_trace

package encoder

func zeroUVTokenRecords(_ *TokenContextPlanes, _ *TokenContextPlanes, _ *MacroblockCoefficients, _ *coefficientTokenCounts, _ *InterCoefficientTokenRecords) bool {
	return false
}
