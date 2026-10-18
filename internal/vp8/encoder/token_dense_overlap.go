package encoder

// Alias handling preserves libvpx v1.16.0-compatible scalar token semantics.

// Token prepacking is available for a future batched consumer. The default
// token walker stays scalar until consuming these fields improves its total cost.

import "unsafe"

// Overlap is declined before preloading so aliased input/output retains scalar semantics.
func denseTokenOverlap(a unsafe.Pointer, an uintptr, b unsafe.Pointer, bn uintptr) bool {
	if an == 0 || bn == 0 {
		return false
	}
	ap, bp := uintptr(a), uintptr(b)
	if an > ^uintptr(0)-ap || bn > ^uintptr(0)-bp {
		return true
	}
	return ap < bp+bn && bp < ap+an
}
