package govpx

import "github.com/gtsteffaniak/go-vpx/internal/testutil"

func appendVP9I420(out []byte, img Image) []byte {
	return testutil.AppendI420Planes(out, img.Width, img.Height,
		img.Y, img.YStride,
		img.U, img.UStride,
		img.V, img.VStride)
}
