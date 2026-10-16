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
	if pixel.Max4(a.Abs().Max(b.Abs()).Max(c.Abs()).Max(d.Abs())) > 255 {
		return false
	}

	one := archsimd.BroadcastInt32x4(1)
	zero := archsimd.Int32x4{}
	var v0, v1, v2, v3 archsimd.Int32x4
	if walsh {
		x := a.Add(c).ShiftAllLeft(2)
		y := b.Add(d).ShiftAllLeft(2)
		z := b.Sub(d).ShiftAllLeft(2)
		w := a.Sub(c).ShiftAllLeft(2)
		v0 = x.Add(y).Add(x.NotEqual(zero).ToInt32x4().And(one))
		v1 = w.Add(z)
		v2 = w.Sub(z)
		v3 = x.Sub(y)
	} else {
		x := a.Add(d).ShiftAllLeft(3)
		y := b.Add(c).ShiftAllLeft(3)
		z := b.Sub(c).ShiftAllLeft(3)
		w := a.Sub(d).ShiftAllLeft(3)
		v0 = x.Add(y)
		v2 = x.Sub(y)
		v1 = z.Mul(archsimd.BroadcastInt32x4(2217)).Add(w.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(14500)).ShiftAllRight(12)
		v3 = w.Mul(archsimd.BroadcastInt32x4(2217)).Sub(z.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(7500)).ShiftAllRight(12)
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
		v0 = v0.Sub(v0.ShiftAllRight(31)).Add(three).ShiftAllRight(3)
		v1 = v1.Sub(v1.ShiftAllRight(31)).Add(three).ShiftAllRight(3)
		v2 = v2.Sub(v2.ShiftAllRight(31)).Add(three).ShiftAllRight(3)
		v3 = v3.Sub(v3.ShiftAllRight(31)).Add(three).ShiftAllRight(3)
	} else {
		x := a.Add(d)
		y := b.Add(c)
		z := b.Sub(c)
		w := a.Sub(d)
		seven := archsimd.BroadcastInt32x4(7)
		v0 = x.Add(y).Add(seven).ShiftAllRight(4)
		v2 = x.Sub(y).Add(seven).ShiftAllRight(4)
		v1 = z.Mul(archsimd.BroadcastInt32x4(2217)).Add(w.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(12000)).ShiftAllRight(16).Add(w.NotEqual(zero).ToInt32x4().And(one))
		v3 = w.Mul(archsimd.BroadcastInt32x4(2217)).Sub(z.Mul(archsimd.BroadcastInt32x4(5352))).Add(archsimd.BroadcastInt32x4(51000)).ShiftAllRight(16)
	}
	pixel.Store4(output[:], v0)
	pixel.Store4(output[4:], v1)
	pixel.Store4(output[8:], v2)
	pixel.Store4(output[12:], v3)
	return true
}
