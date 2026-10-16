//go:build goexperiment.simd && amd64 && !purego

package pixel

import "simd/archsimd"

// Keep feature checks in the compiler's feature domain, as in goav1.
var Enabled = archsimd.X86.AVX2()
