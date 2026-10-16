//go:build goexperiment.simd && arm64 && !purego

package pixel

// ASIMD is mandatory on the arm64 platforms supported by Go.
const Enabled = true

// Native 649237be timing selects the general vertical internal-filter body
// on x64; ARM retains its measured exact-width store specialization.
const internalVerticalGeneric = false
