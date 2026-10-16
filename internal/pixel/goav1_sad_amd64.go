// SPDX-License-Identifier: BSD-2-Clause
// Extracted from goav1 at ff46346e73f063ff1ea3a8d7f06ee26cd9baaad4.
// See UPSTREAM.md and ../../third_party/goav1/LICENSE.
//go:build goexperiment.simd && amd64 && !purego

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
	acc := archsimd.BroadcastUint64x2(0)
	for row := 0; row < 8; row += 2 {
		s := pack2Rows(step(sp, row*srcStride), srcStride)
		r := pack2Rows(step(rp, row*refStride), refStride)
		acc = acc.Add(s.SumOf8AbsDiff(r))
	}
	return sumU64(acc)
}

func sad16ColsDualPtr(sp unsafe.Pointer, srcStride int, rp unsafe.Pointer, refStride int, h int) int {
	acc := archsimd.BroadcastUint64x2(0)
	for row := 0; row < h; row++ {
		acc = acc.Add(load16(step(sp, row*srcStride)).SumOf8AbsDiff(load16(step(rp, row*refStride))))
	}
	return sumU64(acc)
}

func sad16ColsX4Ptr(sp, p0, p1, p2, p3 unsafe.Pointer, srcStride, refStride int, h int) (int, int, int, int) {
	a0 := archsimd.BroadcastUint64x2(0)
	a1 := archsimd.BroadcastUint64x2(0)
	a2 := archsimd.BroadcastUint64x2(0)
	a3 := archsimd.BroadcastUint64x2(0)
	for row := 0; row < h; row++ {
		o := row * refStride
		s := load16(step(sp, row*srcStride))
		a0 = a0.Add(s.SumOf8AbsDiff(load16(step(p0, o))))
		a1 = a1.Add(s.SumOf8AbsDiff(load16(step(p1, o))))
		a2 = a2.Add(s.SumOf8AbsDiff(load16(step(p2, o))))
		a3 = a3.Add(s.SumOf8AbsDiff(load16(step(p3, o))))
	}
	return sumU64(a0), sumU64(a1), sumU64(a2), sumU64(a3)
}

func sumU64(acc archsimd.Uint64x2) int {
	return int(acc.GetElem(0) + acc.GetElem(1))
}
