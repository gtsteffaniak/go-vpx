//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"encoding/binary"
	"simd/archsimd"
	"unsafe"
)

func load8(b []byte) archsimd.Uint8x16 {
	return archsimd.Uint64x2{}.SetElem(0, binary.LittleEndian.Uint64(b)).ReshapeToUint8s()
}
func storeBytes8(b []byte, v archsimd.Uint8x16) {
	binary.LittleEndian.PutUint64(b, v.ReshapeToUint64s().GetElem(0))
}

func widen8(v archsimd.Uint8x16) archsimd.Int16x8 { return v.ExtendLo8ToUint16().BitsToInt16() }
func high4(v archsimd.Int16x8) archsimd.Int32x4 {
	return v.ToBits().ReshapeToUint32s().BitsToInt32().PermuteScalars(2, 3, 0, 1).ToBits().ReshapeToUint16s().BitsToInt16().ExtendLo4ToInt32()
}
func mac(lo, hi archsimd.Int32x4, v, c archsimd.Int16x8) (archsimd.Int32x4, archsimd.Int32x4) {
	c32 := c.ExtendLo4ToInt32()
	return lo.Add(v.ExtendLo4ToInt32().Mul(c32)), hi.Add(high4(v).Mul(c32))
}
func packBytes(v archsimd.Int16x8) archsimd.Uint8x16 {
	return v.Max(archsimd.Int16x8{}).Min(archsimd.BroadcastInt16x8(255)).ToBits().ReshapeToUint8s().PermuteOrZero(archsimd.LoadInt8x16Array(&packByteIndices))
}
func store8(b []byte, lo, hi archsimd.Int32x4) {
	storeBytes8(b, packBytes(lo.SaturateToInt16Concat(hi)))
}
func sadAdd(acc archsimd.Uint32x4, a, b archsimd.Uint8x16) archsimd.Uint32x4 {
	return acc.Add(a.SumOf8AbsDiff(b).ReshapeToUint32s())
}
func sumSquares(v archsimd.Int32x4) uint64 {
	a := v.MulWidenEven(v)
	v = v.PermuteScalars(1, 0, 3, 2)
	s := a.Add(v.MulWidenEven(v))
	return uint64(s.GetElem(0)) + uint64(s.GetElem(1))
}

var packByteIndices = [16]int8{0, 2, 4, 6, 8, 10, 12, 14, -128, -128, -128, -128, -128, -128, -128, -128}

var packInt16Indices = [16]int8{0, 1, 4, 5, 8, 9, 12, 13, -128, -128, -128, -128, -128, -128, -128, -128}

func Store4(s []int16, v archsimd.Int32x4) {
	put4(s, v.ToBits().ReshapeToUint8s().PermuteOrZero(archsimd.LoadInt8x16Array(&packInt16Indices)).ReshapeToUint64s().GetElem(0))
}
func StorePixels4(s []byte, v archsimd.Int32x4) {
	binary.LittleEndian.PutUint32(s, packBytes(v.SaturateToInt16Concat(v)).ReshapeToUint32s().GetElem(0))
}
func sum32(v archsimd.Int32x4) int32 {
	v = v.Add(v.PermuteScalars(2, 3, 0, 1))
	return v.Add(v.PermuteScalars(1, 0, 3, 2)).GetElem(0)
}
func sumU32(v archsimd.Uint32x4) uint32 { return uint32(sum32(v.BitsToInt32())) }
func Max4(v archsimd.Int32x4) int32 {
	v = v.Max(v.PermuteScalars(2, 3, 0, 1))
	return v.Max(v.PermuteScalars(1, 0, 3, 2)).GetElem(0)
}

// Four int16 outputs occupy precisely the eight bytes written here. This
// store is unaligned-safe on amd64 and arm64 and stays within the Go slice.
func put4(out []int16, v uint64) { _ = out[3]; *(*uint64)(unsafe.Pointer(&out[0])) = v }

func highBytes(v archsimd.Uint8x16) archsimd.Uint8x16 {
	return v.ReshapeToUint32s().PermuteScalars(2, 3, 0, 1).ReshapeToUint8s()
}
