//go:build goexperiment.simd && (arm64 || amd64) && !purego

package dsp

import (
	"github.com/gtsteffaniak/govpx/internal/pixel"
	"simd/archsimd"
)

// Arithmetic follows libvpx v1.16.0 vp8/common/idctllm.c.
func inverse4Go(input *[16]int16, pred []byte, ps int, dst []byte, ds int, walsh bool, mb []int16) bool {
	if !pixel.Enabled {
		return false
	}
	in0, in1, in2, in3 := pixel.Load4(input[:]), pixel.Load4(input[4:]), pixel.Load4(input[8:]), pixel.Load4(input[12:])
	var out0, out1, out2, out3 archsimd.Int32x4
	{
		if walsh {
			a, b, c, d := in0.Add(in3), in1.Add(in2), in1.Sub(in2), in0.Sub(in3)
			out0, out1, out2, out3 = a.Add(b), c.Add(d), a.Sub(b), d.Sub(c)
		} else {
			a, b := in0.Add(in2), in0.Sub(in2)
			c := in1.Mul(archsimd.BroadcastInt32x4(sinPI8Sqrt2)).ShiftAllRight(16).Sub(in3.Add(in3.Mul(archsimd.BroadcastInt32x4(cosPI8Sqrt2Minus1)).ShiftAllRight(16)))
			d := in1.Add(in1.Mul(archsimd.BroadcastInt32x4(cosPI8Sqrt2Minus1)).ShiftAllRight(16)).Add(in3.Mul(archsimd.BroadcastInt32x4(sinPI8Sqrt2)).ShiftAllRight(16))
			out0, out1, out2, out3 = a.Add(d), b.Add(c), b.Sub(c), a.Sub(d)
		}
		out0 = out0.ShiftAllLeft(16).ShiftAllRight(16)
		out1 = out1.ShiftAllLeft(16).ShiftAllRight(16)
		out2 = out2.ShiftAllLeft(16).ShiftAllRight(16)
		out3 = out3.ShiftAllLeft(16).ShiftAllRight(16)
	}
	in0, in1, in2, in3 = pixel.Transpose4(out0, out1, out2, out3)
	{
		if walsh {
			a, b, c, d := in0.Add(in3), in1.Add(in2), in1.Sub(in2), in0.Sub(in3)
			out0, out1, out2, out3 = a.Add(b), c.Add(d), a.Sub(b), d.Sub(c)
		} else {
			a, b := in0.Add(in2), in0.Sub(in2)
			c := in1.Mul(archsimd.BroadcastInt32x4(sinPI8Sqrt2)).ShiftAllRight(16).Sub(in3.Add(in3.Mul(archsimd.BroadcastInt32x4(cosPI8Sqrt2Minus1)).ShiftAllRight(16)))
			d := in1.Add(in1.Mul(archsimd.BroadcastInt32x4(cosPI8Sqrt2Minus1)).ShiftAllRight(16)).Add(in3.Mul(archsimd.BroadcastInt32x4(sinPI8Sqrt2)).ShiftAllRight(16))
			out0, out1, out2, out3 = a.Add(d), b.Add(c), b.Sub(c), a.Sub(d)
		}
		bias := int32(4)
		if walsh {
			bias = 3
		}
		out0 = out0.Add(archsimd.BroadcastInt32x4(bias)).ShiftAllRight(3).ShiftAllLeft(16).ShiftAllRight(16)
		out1 = out1.Add(archsimd.BroadcastInt32x4(bias)).ShiftAllRight(3).ShiftAllLeft(16).ShiftAllRight(16)
		out2 = out2.Add(archsimd.BroadcastInt32x4(bias)).ShiftAllRight(3).ShiftAllLeft(16).ShiftAllRight(16)
		out3 = out3.Add(archsimd.BroadcastInt32x4(bias)).ShiftAllRight(3).ShiftAllLeft(16).ShiftAllRight(16)
	}
	out0, out1, out2, out3 = pixel.Transpose4(out0, out1, out2, out3)
	if walsh {
		mb[0] = int16(out0.GetElem(0))
		mb[16] = int16(out0.GetElem(1))
		mb[32] = int16(out0.GetElem(2))
		mb[48] = int16(out0.GetElem(3))
		mb[64] = int16(out1.GetElem(0))
		mb[80] = int16(out1.GetElem(1))
		mb[96] = int16(out1.GetElem(2))
		mb[112] = int16(out1.GetElem(3))
		mb[128] = int16(out2.GetElem(0))
		mb[144] = int16(out2.GetElem(1))
		mb[160] = int16(out2.GetElem(2))
		mb[176] = int16(out2.GetElem(3))
		mb[192] = int16(out3.GetElem(0))
		mb[208] = int16(out3.GetElem(1))
		mb[224] = int16(out3.GetElem(2))
		mb[240] = int16(out3.GetElem(3))
	} else {
		pixel.StorePixels4(dst, out0.Add(pixel.LoadPixels4(pred)))
		pixel.StorePixels4(dst[ds:], out1.Add(pixel.LoadPixels4(pred[ps:])))
		pixel.StorePixels4(dst[2*ds:], out2.Add(pixel.LoadPixels4(pred[2*ps:])))
		pixel.StorePixels4(dst[3*ds:], out3.Add(pixel.LoadPixels4(pred[3*ps:])))
	}
	return true
}
