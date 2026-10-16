//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import "simd/archsimd"

func ssePacket(lo, hi archsimd.Int32x4, a, b archsimd.Uint8x16) (archsimd.Int32x4, archsimd.Int32x4) {
	dl := widen8(a).Sub(widen8(b))
	dh := widenHi8(a).BitsToInt16().Sub(widenHi8(b).BitsToInt16())
	return lo.Add(square8(dl)), hi.Add(square8(dh))
}

// SSE16 computes squared error over sixteen columns, without the signed
// sum needed by variance. Only the requested rows and columns are read;
// accumulation wraps modulo2^32, as in Stats.
func SSE16(src []byte, ss int, ref []byte, rs, h int) uint32 {
	if !Enabled {
		return SSE16Scalar(src, ss, ref, rs, h)
	}
	var lo, hi archsimd.Int32x4
	for y := range h {
		lo, hi = ssePacket(lo, hi, archsimd.LoadUint8x16(src[y*ss:]), archsimd.LoadUint8x16(ref[y*rs:]))
	}
	return uint32(sum32(lo.Add(hi)))
}

// SSE8 packs pairs of exact eight-byte rows into one SIMD packet. An odd
// final row reads eight bytes, including for a minimally sized allocation.
func SSE8(src []byte, ss int, ref []byte, rs, h int) uint32 {
	if !Enabled {
		return SSE8Scalar(src, ss, ref, rs, h)
	}
	var lo, hi archsimd.Int32x4
	y := 0
	for ; y+2 <= h; y += 2 {
		lo, hi = ssePacket(lo, hi, packRows8(src, y*ss, ss), packRows8(ref, y*rs, rs))
	}
	if y < h {
		lo, hi = ssePacket(lo, hi, load8(src[y*ss:]), load8(ref[y*rs:]))
	}
	return uint32(sum32(lo.Add(hi)))
}
