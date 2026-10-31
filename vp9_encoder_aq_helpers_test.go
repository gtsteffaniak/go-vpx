package govpx

import vp9dec "github.com/gtsteffaniak/govpx/internal/vp9/decoder"

func vp9SegmentCountsForGrid(grid []vp9dec.NeighborMi) [vp9dec.MaxSegments]int {
	var counts [vp9dec.MaxSegments]int
	for _, mi := range grid {
		if mi.SegmentID < vp9dec.MaxSegments {
			counts[mi.SegmentID]++
		}
	}
	return counts
}
