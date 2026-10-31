//go:build govpx_phase_stats

package govpx

import "github.com/gtsteffaniak/govpx/internal/vp9/encoder"

func (e *VP9Encoder) vp9PhaseAttachChoosePartitioningStats(
	args *encoder.ChoosePartitioningArgs, stats *encoder.ChoosePartitioningStats,
) {
	if args != nil {
		args.Stats = stats
	}
}
