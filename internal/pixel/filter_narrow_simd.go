//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import (
	"encoding/binary"
	"simd/archsimd"
	"unsafe"
)

// Standard VP9 filters have at most 128 total negative weight, and every
// coefficient is in [-128,128]. Accumulating negative products first keeps
// the signed int16 sum exact. Subsequent positive products may saturate only
// above 32767; those results already clip to byte 255 after rounding by 7.
// This ordering therefore permits eight int16 lanes without premature
// saturation changing a result. Arbitrary filters retain the int32 kernel.
func filterNarrow8(src []byte, ss int, dst []byte, ds, w, h, step int, taps []int16) bool {
	if len(taps) != 8 || w <= 0 || (w != 4 && w%8 != 0) || h <= 0 || ss < 0 || ds < 0 || step < 0 {
		return false
	}
	mode := narrow8Mode(taps)
	if mode == 0 {
		return false
	}
	_ = src[(h-1)*ss+w-1+7*step]
	_ = dst[(h-1)*ds+w-1]
	sp := unsafe.Pointer(unsafe.SliceData(src))
	if step == ss && filterVerticalNarrow8(src, ss, dst, ds, w, h, taps, mode) {
		return true
	}
	c0, c1, c2, c3 := archsimd.BroadcastInt16x8(taps[0]), archsimd.BroadcastInt16x8(taps[1]), archsimd.BroadcastInt16x8(taps[2]), archsimd.BroadcastInt16x8(taps[3])
	c4, c5, c6, c7 := archsimd.BroadcastInt16x8(taps[4]), archsimd.BroadcastInt16x8(taps[5]), archsimd.BroadcastInt16x8(taps[6]), archsimd.BroadcastInt16x8(taps[7])
	bias := archsimd.BroadcastInt16x8(64)
	n := 8
	if w == 4 {
		n = 4
	}
	for y := range h {
		for x := 0; x < w; x += n {
			p := unsafe.Add(sp, y*ss+x)
			var p0, p1, p2, p3, p4, p5, p6, p7 archsimd.Int16x8
			if step == 1 && len(src)-y*ss-x >= 16 {
				a := load16(p)
				z := archsimd.Uint8x16{}
				p0 = widen8(a)
				p1 = widen8(z.ConcatShiftBytesRight(a, 1))
				p2 = widen8(z.ConcatShiftBytesRight(a, 2))
				p3 = widen8(z.ConcatShiftBytesRight(a, 3))
				p4 = widen8(z.ConcatShiftBytesRight(a, 4))
				p5 = widen8(z.ConcatShiftBytesRight(a, 5))
				p6 = widen8(z.ConcatShiftBytesRight(a, 6))
				p7 = widen8(z.ConcatShiftBytesRight(a, 7))
			} else {
				p0, p1, p2, p3 = widen8(loadNPtr(p, n)), widen8(loadNPtr(unsafe.Add(p, step), n)), widen8(loadNPtr(unsafe.Add(p, 2*step), n)), widen8(loadNPtr(unsafe.Add(p, 3*step), n))
				p4, p5, p6, p7 = widen8(loadNPtr(unsafe.Add(p, 4*step), n)), widen8(loadNPtr(unsafe.Add(p, 5*step), n)), widen8(loadNPtr(unsafe.Add(p, 6*step), n)), widen8(loadNPtr(unsafe.Add(p, 7*step), n))
			}

			acc := bias
			if mode == 1 {
				acc = mulAdd16(acc, p0, c0)
				acc = mulAdd16(acc, p2, c2)
				acc = mulAdd16(acc, p5, c5)
				acc = mulAdd16(acc, p7, c7)
				acc = acc.AddSaturated(p1.Mul(c1)).AddSaturated(p3.Mul(c3)).AddSaturated(p4.Mul(c4)).AddSaturated(p6.Mul(c6))
			} else {
				acc = mulAdd16(acc, p0, c0)
				acc = mulAdd16(acc, p1, c1)
				acc = mulAdd16(acc, p6, c6)
				acc = mulAdd16(acc, p7, c7)
				acc = acc.AddSaturated(p2.Mul(c2)).AddSaturated(p3.Mul(c3)).AddSaturated(p4.Mul(c4)).AddSaturated(p5.Mul(c5))
			}
			packed := packBytes(filterShift7(acc))
			if n == 4 {
				binary.LittleEndian.PutUint32(dst[y*ds+x:], packed.ReshapeToUint32s().GetElem(0))
			} else {
				storeBytes8(dst[y*ds+x:], packed)
			}
		}
	}
	return true
}

func narrow8Mode(taps []int16) int {
	regular, smooth := true, true
	negative := 0
	for k, v := range taps {
		if v < -128 || v > 128 {
			return 0
		}
		if v < 0 {
			negative -= int(v)
		}
		switch k {
		case 0, 2, 5, 7:
			regular = regular && v <= 0
		default:
			regular = regular && v >= 0
		}
		switch k {
		case 0, 1, 6, 7:
			smooth = smooth && v <= 0
		default:
			smooth = smooth && v >= 0
		}
	}
	if negative > 128 || !regular && !smooth {
		return 0
	}
	if regular {
		return 1
	}
	return 2
}
