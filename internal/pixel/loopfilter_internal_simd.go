//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

// LoopFilterInternal16 filters the three four-pixel-spaced internal edges
// of a 16x16 tile. step crosses the edges; outer walks along them. Edges
// retain their original order because each can read the preceding result.
func LoopFilterInternal16(p []byte, step, outer, blimit, limit, thresh int) bool {
	if !Enabled || step <= 0 || outer <= 0 || blimit < 0 || blimit >= 255 || limit < 0 || limit > 255 || thresh < 0 || thresh > 255 {
		return false
	}
	if !((outer == 1 && step >= 16) || (step == 1 && outer >= 8)) {
		return false
	}
	// Three complete eight-row/column windows occupy the exact 16x16 tile.
	_ = p[15*step+15*outer]
	loopFilterBytes(p, 4*step, step, outer, 16, blimit, limit, thresh, 4)
	loopFilterBytes(p, 8*step, step, outer, 16, blimit, limit, thresh, 4)
	loopFilterBytes(p, 12*step, step, outer, 16, blimit, limit, thresh, 4)
	return true
}
