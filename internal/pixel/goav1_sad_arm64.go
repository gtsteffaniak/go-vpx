// SPDX-License-Identifier: BSD-2-Clause
// Extracted from goav1 at ff46346e73f063ff1ea3a8d7f06ee26cd9baaad4.
// See UPSTREAM.md and ../../third_party/goav1/LICENSE.
//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

func step(p unsafe.Pointer, n int) unsafe.Pointer { return unsafe.Add(p, n) }

// load16 loads the 16 bytes at p. Callers guarantee they are readable.

func load16(p unsafe.Pointer) archsimd.Uint8x16 {
	return archsimd.LoadUint8x16Array((*[16]uint8)(p))
}

func pack2Rows(p unsafe.Pointer, stride int) archsimd.Uint8x16 {
	lo := *(*uint64)(p)
	hi := *(*uint64)(step(p, stride))
	return archsimd.BroadcastUint64x2(lo).SetElem(1, hi).ReshapeToUint8s()
}

func sad8x8Ptr(sp unsafe.Pointer, srcStride int, rp unsafe.Pointer, refStride int) int {
	acc := archsimd.BroadcastUint16x8(0)
	for row := 0; row < 8; row += 2 {
		s := pack2Rows(step(sp, row*srcStride), srcStride)
		r := pack2Rows(step(rp, row*refStride), refStride)
		acc = acc.Add(widen16(absDiffU8x16(s, r)))
	}
	return int(acc.ReduceSum())
}

func sad16ColsDualPtr(sp unsafe.Pointer, srcStride int, rp unsafe.Pointer, refStride int, h int) int {
	even := archsimd.BroadcastUint16x8(0)
	odd := archsimd.BroadcastUint16x8(0)
	for row := 0; row < h; row += 2 {
		even = absAcc16(even, load16(step(sp, row*srcStride)), load16(step(rp, row*refStride)))
		odd = absAcc16(odd, load16(step(sp, (row+1)*srcStride)), load16(step(rp, (row+1)*refStride)))
	}
	return sumU16(even) + sumU16(odd)
}

func sad16ColsX4Ptr(sp, p0, p1, p2, p3 unsafe.Pointer, srcStride, refStride int, h int) (int, int, int, int) {
	a0 := archsimd.BroadcastUint16x8(0)
	a1 := archsimd.BroadcastUint16x8(0)
	a2 := archsimd.BroadcastUint16x8(0)
	a3 := archsimd.BroadcastUint16x8(0)
	for row := 0; row < h; row++ {
		o := row * refStride
		s := load16(step(sp, row*srcStride))
		a0 = absAcc16(a0, s, load16(step(p0, o)))
		a1 = absAcc16(a1, s, load16(step(p1, o)))
		a2 = absAcc16(a2, s, load16(step(p2, o)))
		a3 = absAcc16(a3, s, load16(step(p3, o)))
	}
	return sumU16(a0), sumU16(a1), sumU16(a2), sumU16(a3)
}

func absDiffU8x16(a, b archsimd.Uint8x16) archsimd.Uint8x16 {
	return a.SubSaturated(b).Or(b.SubSaturated(a))
}

func widen16(absd archsimd.Uint8x16) archsimd.Uint16x8 {
	return absd.ExtendLo8ToUint16().Add(absd.HiToLo().ExtendLo8ToUint16())
}

func absAcc16(acc archsimd.Uint16x8, s, r archsimd.Uint8x16) archsimd.Uint16x8 {
	return acc.Add(widen16(absDiffU8x16(s, r)))
}

func sumU16(acc archsimd.Uint16x8) int {
	return int(acc.ExtendLo4ToUint32().Add(acc.HiToLo().ExtendLo4ToUint32()).ReduceSum())
}

func sad16LimitPtr(sp unsafe.Pointer, ss int, rp unsafe.Pointer, rs, h, limit int) int {
	var sum archsimd.Uint16x8
	for y := range h {
		a := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(sp, y*ss)))
		b := archsimd.LoadUint8x16Array((*[16]byte)(unsafe.Add(rp, y*rs)))
		sum = sum.Add(absPairs(a, b))
		if value := sumPairs(sum); value > limit {
			return value
		}
	}
	return sumPairs(sum)
}
