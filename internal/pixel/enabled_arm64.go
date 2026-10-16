//go:build goexperiment.simd && arm64 && !purego

package pixel

// ASIMD is mandatory on the arm64 platforms supported by Go.
const Enabled = true
