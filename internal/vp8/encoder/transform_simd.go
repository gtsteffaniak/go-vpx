//go:build goexperiment.simd && (arm64 || amd64) && !purego

// Arithmetic follows libvpx v1.16.0.
package encoder

import (
	"github.com/thesyncim/govpx/internal/pixel"
	"simd/archsimd"
)

func forward4Go(input []int16, stride int, output *[16]int16, walsh bool) bool {
	if !pixel.Enabled {
		return false
	}
	a, b, c, d := pixel.Transpose4(pixel.Load4(input), pixel.Load4(input[stride:]), pixel.Load4(input[2*stride:]), pixel.Load4(input[3*stride:]))
	// DCT products use the pixel residual bound. Walsh uses only int32
	// sums/shifts and is exact for all int16 DC coefficients (including Y2).
	if !walsh && pixel.Max4(a.Abs().Max(b.Abs()).Max(c.Abs()).Max(d.Abs())) > 255 {
		return false
	}

	one := archsimd.BroadcastInt32x4(1)
	zero := archsimd.Int32x4{}
	var v0, v1, v2, v3 archsimd.Int32x4
	if walsh {
		x := vp8Left2(a.Add(c))
		y := vp8Left2(b.Add(d))
		z := vp8Left2(b.Sub(d))
		w := vp8Left2(a.Sub(c))
		v0 = x.Add(y).Add(x.NotEqual(zero).ToInt32x4().And(one))
		v1 = w.Add(z)
		v2 = w.Sub(z)
		v3 = x.Sub(y)
	} else {
		x := vp8Left3(a.Add(d))
		y := vp8Left3(b.Add(c))
		z := vp8Left3(b.Sub(c))
		w := vp8Left3(a.Sub(d))
		v0 = x.Add(y)
		v2 = x.Sub(y)
		v1 = vp8Right12(z.Mul(archsimd.BroadcastInt32x4(2217)).Add(w.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(14500)))
		v3 = vp8Right12(w.Mul(archsimd.BroadcastInt32x4(2217)).Sub(z.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(7500)))
	}
	a, b, c, d = pixel.Transpose4(v0, v1, v2, v3)
	if walsh {
		x := a.Add(c)
		y := b.Add(d)
		z := b.Sub(d)
		w := a.Sub(c)
		v0 = x.Add(y)
		v1 = w.Add(z)
		v2 = w.Sub(z)
		v3 = x.Sub(y)
		three := archsimd.BroadcastInt32x4(3)
		v0 = vp8Right3(v0.Sub(vp8Sign(v0)).Add(three))
		v1 = vp8Right3(v1.Sub(vp8Sign(v1)).Add(three))
		v2 = vp8Right3(v2.Sub(vp8Sign(v2)).Add(three))
		v3 = vp8Right3(v3.Sub(vp8Sign(v3)).Add(three))
	} else {
		x := a.Add(d)
		y := b.Add(c)
		z := b.Sub(c)
		w := a.Sub(d)
		seven := archsimd.BroadcastInt32x4(7)
		v0 = vp8Right4(x.Add(y).Add(seven))
		v2 = vp8Right4(x.Sub(y).Add(seven))
		v1 = vp8Right16(z.Mul(archsimd.BroadcastInt32x4(2217)).Add(w.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(12000))).Add(w.NotEqual(zero).ToInt32x4().And(one))
		v3 = vp8Right16(w.Mul(archsimd.BroadcastInt32x4(2217)).Sub(z.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(51000)))
	}
	pixel.Store4(output[:], v0)
	pixel.Store4(output[4:], v1)
	pixel.Store4(output[8:], v2)
	pixel.Store4(output[12:], v3)
	return true
}
