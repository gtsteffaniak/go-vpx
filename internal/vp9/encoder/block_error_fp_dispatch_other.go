package encoder

import "github.com/gtsteffaniak/govpx/internal/pixel"

func blockErrorFPDispatch(coeff, dqcoeff []int16, n int) uint64 {
	if pixel.Enabled {
		return uint64(pixel.Error(coeff, dqcoeff, n))
	}

	return blockErrorFPScalar(coeff, dqcoeff, n)
}
