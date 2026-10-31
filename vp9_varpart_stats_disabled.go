//go:build !govpx_phase_stats

package govpx

import "github.com/gtsteffaniak/govpx/internal/vp9/encoder"

func (e *VP9Encoder) vp9PhaseAttachChoosePartitioningStats(
	*encoder.ChoosePartitioningArgs, *encoder.ChoosePartitioningStats,
) {
}
