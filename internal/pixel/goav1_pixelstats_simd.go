// SPDX-License-Identifier: BSD-2-Clause
// Extracted from goav1 at ff46346e73f063ff1ea3a8d7f06ee26cd9baaad4.
// See UPSTREAM.md and ../../third_party/goav1/LICENSE.
//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import (
	"encoding/binary"
	"simd/archsimd"
)

func pixelDiffAcc(sse, sum archsimd.Int32x4, a, b archsimd.Uint8x16) (archsimd.Int32x4, archsimd.Int32x4) {
	lo := a.ExtendLo8ToUint16().BitsToInt16().Sub(b.ExtendLo8ToUint16().BitsToInt16())
	hi := widenHi8(a).BitsToInt16().Sub(widenHi8(b).BitsToInt16())
	d0 := lo.ExtendLo4ToInt32()
	d1 := hiInt16(lo).ExtendLo4ToInt32()
	d2 := hi.ExtendLo4ToInt32()
	d3 := hiInt16(hi).ExtendLo4ToInt32()
	sse = sse.Add(d0.Mul(d0)).Add(d1.Mul(d1)).Add(d2.Mul(d2)).Add(d3.Mul(d3))
	sum = sum.Add(d0).Add(d1).Add(d2).Add(d3)
	return sse, sum
}

func reducePixelAcc(sse, sum archsimd.Int32x4) (total, signed int64) {
	total = int64(sse.GetElem(0)) + int64(sse.GetElem(1)) + int64(sse.GetElem(2)) + int64(sse.GetElem(3))
	signed = int64(sum.GetElem(0)) + int64(sum.GetElem(1)) + int64(sum.GetElem(2)) + int64(sum.GetElem(3))
	return total, signed
}

func pixelStatsRows16(src []byte, srcStride int, ref []byte, refStride int, w, h int) (sse uint32, sum int32) {
	sseAcc := archsimd.BroadcastInt32x4(0)
	sumAcc := archsimd.BroadcastInt32x4(0)
	for r := range h {
		srow := src[r*srcStride : r*srcStride+w]
		rrow := ref[r*refStride : r*refStride+w]
		for c := 0; c < w; c += 16 {
			sseAcc, sumAcc = pixelDiffAcc(sseAcc, sumAcc, archsimd.LoadUint8x16(srow[c:c+16]), archsimd.LoadUint8x16(rrow[c:c+16]))
		}
	}
	total, signed := reducePixelAcc(sseAcc, sumAcc)
	return uint32(total), int32(signed)
}

func pixelStatsPairs8(src []byte, srcStride int, ref []byte, refStride int, h int) (sse uint32, sum int32) {
	sseAcc := archsimd.BroadcastInt32x4(0)
	sumAcc := archsimd.BroadcastInt32x4(0)
	r := 0
	for ; r+2 <= h; r += 2 {
		sseAcc, sumAcc = pixelDiffAcc(sseAcc, sumAcc, packRows8(src, r*srcStride, srcStride), packRows8(ref, r*refStride, refStride))
	}
	total, signed := reducePixelAcc(sseAcc, sumAcc)
	if r < h {
		t, sg := pixelStatsPureGo(src[r*srcStride:], srcStride, ref[r*refStride:], refStride, 8, 1)
		total += int64(t)
		signed += int64(sg)
	}
	return uint32(total), int32(signed)
}

func packRows8(p []byte, off, stride int) archsimd.Uint8x16 {
	lo := binary.LittleEndian.Uint64(p[off:])
	hi := binary.LittleEndian.Uint64(p[off+stride:])
	return archsimd.BroadcastUint64x2(lo).SetElem(1, hi).ReshapeToUint8s()
}

func pixelStatsSIMD(src []byte, srcStride int, ref []byte, refStride int, w, h int) (sse uint32, sum int32) {
	switch {
	case w == 8:
		return pixelStatsPairs8(src, srcStride, ref, refStride, h)
	case w >= 16 && w&15 == 0:
		return pixelStatsRows16(src, srcStride, ref, refStride, w, h)
	}
	return pixelStatsPureGo(src, srcStride, ref, refStride, w, h)
}

// pixelStatsPureGo is govpx's glue for the extracted goav1 fallback contract.
func pixelStatsPureGo(src []byte, ss int, ref []byte, rs, w, h int) (uint32, int32) {
	sum, sse := StatsScalar(src, ss, ref, rs, w, h)
	return sse, int32(sum)
}
