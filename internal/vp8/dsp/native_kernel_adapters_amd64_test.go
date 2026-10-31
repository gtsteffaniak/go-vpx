//go:build amd64 && !purego

package dsp

// Adapters preserve the original seeded assembly differential cases while
// testing the new production Go dispatch against the independent scalar port.
import (
	"github.com/gtsteffaniak/govpx/internal/pixel"
	"unsafe"
)

func sadBlock16x16Native(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 15*ss+16), unsafe.Slice(ref, 15*rs+16)
	return int32(pixel.SAD(a, ss, b, rs, 16, 16))
}
func sadBlock16x16Reference(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 15*ss+16), unsafe.Slice(ref, 15*rs+16)
	return int32(scalarSAD(a, ss, b, rs, 16, 16))
}
func sadBlock16x8Native(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 7*ss+16), unsafe.Slice(ref, 7*rs+16)
	return int32(pixel.SAD(a, ss, b, rs, 16, 8))
}
func sadBlock16x8Reference(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 7*ss+16), unsafe.Slice(ref, 7*rs+16)
	return int32(scalarSAD(a, ss, b, rs, 16, 8))
}
func sadBlock8x16Native(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 15*ss+8), unsafe.Slice(ref, 15*rs+8)
	return int32(pixel.SAD(a, ss, b, rs, 8, 16))
}
func sadBlock8x16Reference(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 15*ss+8), unsafe.Slice(ref, 15*rs+8)
	return int32(scalarSAD(a, ss, b, rs, 8, 16))
}
func sadBlock8x8Native(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 7*ss+8), unsafe.Slice(ref, 7*rs+8)
	return int32(pixel.SAD(a, ss, b, rs, 8, 8))
}
func sadBlock8x8Reference(src *byte, ss int, ref *byte, rs int) int32 {
	a, b := unsafe.Slice(src, 7*ss+8), unsafe.Slice(ref, 7*rs+8)
	return int32(scalarSAD(a, ss, b, rs, 8, 8))
}
func sadBlock16x16x4Native(src *byte, ss int, r0, r1, r2, r3 *byte, rs int, out *[4]uint32) {
	SAD16x16x4PtrFast(src, ss, r0, r1, r2, r3, rs, out)
}
func sixTapPredict16x16Native(dst *byte, ds int, src *byte, ss int, hf, vf *[6]int16, tmp *[21 * 16]byte) {
	a, b := unsafe.Slice(src, 20*ss+21), unsafe.Slice(dst, 15*ds+16)
	pixel.Filter(a, ss, tmp[:], 16, 16, 21, 1, hf[:])
	pixel.Filter(tmp[:], 16, b, ds, 16, 16, 16, vf[:])
}
func loopFilterEdgeH16Native(src *byte, ss int, b, l, t byte) {
	a := unsafe.Slice(src, 7*ss+16)
	loopFilterHorizontalEdgeDispatch(a, ss, b, l, t, 2)
}
func mbLoopFilterEdgeH16Native(src *byte, ss int, b, l, t byte) {
	a := unsafe.Slice(src, 7*ss+16)
	mbLoopFilterHorizontalEdgeDispatch(a, ss, b, l, t, 2)
}
func varianceBlock16xNNative(src *byte, ss int, ref *byte, rs, h int, sum *int32, sse *uint32) {
	a, b := unsafe.Slice(src, (h-1)*ss+16), unsafe.Slice(ref, (h-1)*rs+16)
	su, se := pixel.Stats(a, ss, b, rs, 16, h)
	*sum, *sse = int32(su), se
}
func varianceBlock8x16Native(src *byte, ss int, ref *byte, rs int, sum *int32, sse *uint32) {
	varianceBlock8xNNative(src, ss, ref, rs, 16, sum, sse)
}
func varianceBlock8xNNative(src *byte, ss int, ref *byte, rs, h int, sum *int32, sse *uint32) {
	a, b := unsafe.Slice(src, (h-1)*ss+8), unsafe.Slice(ref, (h-1)*rs+8)
	su, se := pixel.Stats(a, ss, b, rs, 8, h)
	*sum, *sse = int32(su), se
}
func varianceBlock16x16Native(src *byte, ss int, ref *byte, rs int, sum *int32, sse *uint32) {
	varianceBlock16xNNative(src, ss, ref, rs, 16, sum, sse)
}
func varFilterBlock2DBilinearFirstPass16Native(src *byte, ss int, dst *uint16, h int, f0, f1 uint64) {
	pixel.BilinearFirst(unsafe.Slice(src, (h-1)*ss+17), ss, unsafe.Slice(dst, 16*h), 16, h, [2]int16{int16(f0), int16(f1)})
}

func varFilterBlock2DBilinearSecondPass16Native(src *[17 * 16]uint16, dst *byte, h int, f0, f1 uint64) {
	d := unsafe.Slice(dst, 16*h)
	varFilterBlock2DBilinearSecondPass16(src, d, h, [2]int16{int16(f0), int16(f1)})
}
func sixTapPredict16x16Reference(dst *byte, ds int, src *byte, ss int, hf, vf *[6]int16, tmp *[21 * 16]byte) {
	a, b := unsafe.Slice(src, 20*ss+21), unsafe.Slice(dst, 15*ds+16)
	for y := range 21 {
		for x := range 16 {
			v := 64
			for k := range 6 {
				v += int(a[y*ss+x+k]) * int(hf[k])
			}
			tmp[y*16+x] = ClipPixel(v >> 7)
		}
	}
	for y := range 16 {
		for x := range 16 {
			v := 64
			for k := range 6 {
				v += int(tmp[(y+k)*16+x]) * int(vf[k])
			}
			b[y*ds+x] = ClipPixel(v >> 7)
		}
	}
}
func loopFilterEdgeH16Reference(src *byte, ss int, b, l, t byte) {
	a := unsafe.Slice(src, 7*ss+16)
	loopFilterHorizontalEdgeScalar(a, ss, b, l, t, 2)
}
func mbLoopFilterEdgeH16Reference(src *byte, ss int, b, l, t byte) {
	a := unsafe.Slice(src, 7*ss+16)
	mbLoopFilterHorizontalEdgeScalar(a, ss, b, l, t, 2)
}
func varianceBlock16xNReference(src *byte, ss int, ref *byte, rs, h int, sum *int32, sse *uint32) {
	a, b := unsafe.Slice(src, (h-1)*ss+16), unsafe.Slice(ref, (h-1)*rs+16)
	su, se := varianceBlock16xNScalar(a, ss, b, rs, h)
	*sum, *sse = int32(su), uint32(se)
}
func varianceBlock8x16Reference(src *byte, ss int, ref *byte, rs int, sum *int32, sse *uint32) {
	varianceBlock8xNReference(src, ss, ref, rs, 16, sum, sse)
}
func varianceBlock8xNReference(src *byte, ss int, ref *byte, rs, h int, sum *int32, sse *uint32) {
	a, b := unsafe.Slice(src, (h-1)*ss+8), unsafe.Slice(ref, (h-1)*rs+8)
	su, se := varianceBlock8xNScalar(a, ss, b, rs, h)
	*sum, *sse = int32(su), uint32(se)
}
func varianceBlock16x16Reference(src *byte, ss int, ref *byte, rs int, sum *int32, sse *uint32) {
	varianceBlock16xNReference(src, ss, ref, rs, 16, sum, sse)
}
func varFilterBlock2DBilinearFirstPass16Reference(src *byte, ss int, dst *uint16, h int, f0, f1 uint64) {
	a := unsafe.Slice(src, (h-1)*ss+17)
	d := unsafe.Slice(dst, 16*h)
	for y := range h {
		for x := range 16 {
			d[y*16+x] = uint16((int(a[y*ss+x])*int(uint16(f0)) + int(a[y*ss+x+1])*int(uint16(f1)) + 64) >> 7)
		}
	}
}
func varFilterBlock2DBilinearSecondPass16Reference(src *[17 * 16]uint16, dst *byte, h int, f0, f1 uint64) {
	d := unsafe.Slice(dst, 16*h)
	varFilterBlock2DBilinearSecondPass16Scalar(src, d, h, [2]int16{int16(f0), int16(f1)})
}
