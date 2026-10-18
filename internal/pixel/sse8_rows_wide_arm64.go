//go:build goexperiment.simd && arm64 && !purego

package pixel

import "unsafe"

const sse8WideAvailable = false

func sse8RowsWide(sp unsafe.Pointer, ss int, rp unsafe.Pointer, rs, h int) uint32 { return 0 }
