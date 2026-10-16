package dsp

import (
	"github.com/thesyncim/govpx/internal/pixel"
	"unsafe"
)

// Pure-Go fallback for the SAD primitives on architectures without a SIMD
// port. Mirrors libvpx v1.16.0 vpx_dsp/sad.c semantics.

func sadBlock16x16(src []byte, srcStride int, ref []byte, refStride int) int {
	return pixel.SAD16(src, srcStride, ref, refStride, 16)
}

// SAD16x16PtrFast is the pointer-form fallback used when callers have
// already validated bounds. See sad_dispatch_arm64.go for the contract.
func SAD16x16PtrFast(src *byte, srcStride int, ref *byte, refStride int) int {
	return pixel.SAD16(unsafe.Slice(src, 15*srcStride+16), srcStride, unsafe.Slice(ref, 15*refStride+16), refStride, 16)
}

// SAD16x16LimitPtrFast is the limited pointer-form fallback. See
// sad_dispatch_arm64.go for the contract.
func SAD16x16LimitPtrFast(src *byte, srcStride int, ref *byte, refStride int, limit int) int {
	return sadBlockLimitScalarPtr(src, srcStride, ref, refStride, 16, 16, limit)
}

// SAD16x16x4PtrFast mirrors libvpx's vpx_sad16x16x4d interface.
func SAD16x16x4PtrFast(src *byte, srcStride int, ref0 *byte, ref1 *byte, ref2 *byte, ref3 *byte, refStride int, out *[4]uint32) {
	pixel.SAD16x4(unsafe.Slice(src, 15*srcStride+16), srcStride, unsafe.Slice(ref0, 15*refStride+16), unsafe.Slice(ref1, 15*refStride+16), unsafe.Slice(ref2, 15*refStride+16), unsafe.Slice(ref3, 15*refStride+16), refStride, 16, out)
}

// SAD16x16x4LimitPtrFast is the pointer-form x4 limited fallback.
func SAD16x16x4LimitPtrFast(src *byte, srcStride int, ref0 *byte, ref1 *byte, ref2 *byte, ref3 *byte, refStride int, limits *[4]int32, out *[4]uint32) {
	out[0] = uint32(sadBlockLimitScalarPtr(src, srcStride, ref0, refStride, 16, 16, int(limits[0])))
	out[1] = uint32(sadBlockLimitScalarPtr(src, srcStride, ref1, refStride, 16, 16, int(limits[1])))
	out[2] = uint32(sadBlockLimitScalarPtr(src, srcStride, ref2, refStride, 16, 16, int(limits[2])))
	out[3] = uint32(sadBlockLimitScalarPtr(src, srcStride, ref3, refStride, 16, 16, int(limits[3])))
}

func sadBlock16x16Limit(src []byte, srcStride int, ref []byte, refStride int, limit int) int {
	return sadBlockLimitScalar(src, srcStride, ref, refStride, 16, 16, limit)
}

func sadBlockScalarPtr(src *byte, srcStride int, ref *byte, refStride int, width, height int) int {
	srcBase := unsafe.Pointer(src)
	refBase := unsafe.Pointer(ref)
	sad := 0
	for y := 0; y < height; y++ {
		srcRow := unsafe.Add(srcBase, y*srcStride)
		refRow := unsafe.Add(refBase, y*refStride)
		for x := 0; x < width; x++ {
			a := int(*(*byte)(unsafe.Add(srcRow, x)))
			b := int(*(*byte)(unsafe.Add(refRow, x)))
			diff := a - b
			mask := diff >> signShift
			sad += (diff ^ mask) - mask
		}
	}
	return sad
}

func sadBlockLimitScalarPtr(src *byte, srcStride int, ref *byte, refStride int, width, height, limit int) int {
	if pixel.Enabled && width == 16 && height <= 64 {
		return pixel.SAD16Limit(unsafe.Slice(src, (height-1)*srcStride+16), srcStride, unsafe.Slice(ref, (height-1)*refStride+16), refStride, height, limit)
	}

	if pixel.Enabled {
		sum := 0
		for y := range height {
			a := unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(src), y*srcStride)), width)
			b := unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(ref), y*refStride)), width)
			sum += pixel.SAD(a, width, b, width, width, 1)
			if sum > limit {
				return sum
			}
		}
		return sum
	}

	srcBase := unsafe.Pointer(src)
	refBase := unsafe.Pointer(ref)
	sad := 0
	for y := 0; y < height; y++ {
		srcRow := unsafe.Add(srcBase, y*srcStride)
		refRow := unsafe.Add(refBase, y*refStride)
		for x := 0; x < width; x++ {
			a := int(*(*byte)(unsafe.Add(srcRow, x)))
			b := int(*(*byte)(unsafe.Add(refRow, x)))
			diff := a - b
			mask := diff >> signShift
			sad += (diff ^ mask) - mask
		}
		if sad > limit {
			return sad
		}
	}
	return sad
}

func sadBlock16x8(src []byte, srcStride int, ref []byte, refStride int) int {
	return pixel.SAD16(src, srcStride, ref, refStride, 8)
}

func sadBlock8x16(src []byte, srcStride int, ref []byte, refStride int) int {
	return pixel.SAD(src, srcStride, ref, refStride, 8, 16)
}

func sadBlock8x8(src []byte, srcStride int, ref []byte, refStride int) int {
	return pixel.SAD(src, srcStride, ref, refStride, 8, 8)
}

func sadBlock4x4(src []byte, srcStride int, ref []byte, refStride int) int {
	return pixel.SAD(src, srcStride, ref, refStride, 4, 4)
}

func sadBlockScalar(src []byte, srcStride int, ref []byte, refStride int, width, height int) int {
	if width <= 0 || height <= 0 {
		return 0
	}
	_ = src[(height-1)*srcStride+(width-1)]
	_ = ref[(height-1)*refStride+(width-1)]
	srcBase := unsafe.Pointer(&src[0])
	refBase := unsafe.Pointer(&ref[0])
	sad := 0
	for y := 0; y < height; y++ {
		srcRow := unsafe.Add(srcBase, y*srcStride)
		refRow := unsafe.Add(refBase, y*refStride)
		for x := 0; x < width; x++ {
			a := int(*(*byte)(unsafe.Add(srcRow, x)))
			b := int(*(*byte)(unsafe.Add(refRow, x)))
			diff := a - b
			mask := diff >> signShift
			sad += (diff ^ mask) - mask
		}
	}
	return sad
}

func sadBlockLimitScalar(src []byte, srcStride int, ref []byte, refStride int, width, height, limit int) int {
	if pixel.Enabled {
		return sadBlockLimitScalarPtr(unsafe.SliceData(src), srcStride, unsafe.SliceData(ref), refStride, width, height, limit)
	}

	if width <= 0 || height <= 0 {
		return 0
	}
	_ = src[(height-1)*srcStride+(width-1)]
	_ = ref[(height-1)*refStride+(width-1)]
	srcBase := unsafe.Pointer(&src[0])
	refBase := unsafe.Pointer(&ref[0])
	sad := 0
	for y := 0; y < height; y++ {
		srcRow := unsafe.Add(srcBase, y*srcStride)
		refRow := unsafe.Add(refBase, y*refStride)
		for x := 0; x < width; x++ {
			a := int(*(*byte)(unsafe.Add(srcRow, x)))
			b := int(*(*byte)(unsafe.Add(refRow, x)))
			diff := a - b
			mask := diff >> signShift
			sad += (diff ^ mask) - mask
		}
		if sad > limit {
			return sad
		}
	}
	return sad
}
