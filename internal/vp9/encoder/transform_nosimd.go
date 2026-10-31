//go:build !goexperiment.simd || (!arm64 && !amd64) || purego

package encoder

import "github.com/gtsteffaniak/go-vpx/internal/vp9/common"

func forwardTransformGo(_ []int16, _ int, _ []int16, _ int, _ common.TxType) bool { return false }

func forwardDCT32Go(_ []int16, _ int, _ []int16, _ bool) bool { return false }
