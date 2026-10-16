package encoder

import "github.com/thesyncim/govpx/internal/pixel"

func hadamard8x8Into(src []int16, stride int, coeff []int16) {
	if !hadamardGo(src, stride, coeff, 8) {
		hadamard8x8Scalar(src, stride, coeff)
	}
}

func hadamard16x16Into(src []int16, stride int, coeff []int16) {
	if !hadamardGo(src, stride, coeff, 16) {
		hadamard16x16Scalar(src, stride, coeff)
	}
}

func satdAbsSum(coeff []int16, n int) int {
	if pixel.Enabled {
		return pixel.AbsSum(coeff, n)
	}
	return satdAbsSumScalar(coeff, n)
}
