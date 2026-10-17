//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import (
	"encoding/binary"
	"simd/archsimd"
	"unsafe"
)

// Load4 reads exactly four int16 samples, including minimally sized rows.
func Load4(s []int16) archsimd.Int32x4 {
	_ = s[3]
	return archsimd.Uint64x2{}.SetElem(0, *(*uint64)(unsafe.Pointer(unsafe.SliceData(s)))).ReshapeToUint16s().BitsToInt16().ExtendLo4ToInt32()
}

// Store4 preserves Go's wrapping int16 conversion, without saturating it.
// Transpose4 follows goav1's in-register 4x4 int32 transpose.
func Transpose4(a, b, c, d archsimd.Int32x4) (archsimd.Int32x4, archsimd.Int32x4, archsimd.Int32x4, archsimd.Int32x4) {
	ab0 := a.InterleaveLo(b).ToBits().ReshapeToUint64s()
	ab1 := a.InterleaveHi(b).ToBits().ReshapeToUint64s()
	cd0 := c.InterleaveLo(d).ToBits().ReshapeToUint64s()
	cd1 := c.InterleaveHi(d).ToBits().ReshapeToUint64s()
	return ab0.InterleaveLo(cd0).ReshapeToUint32s().BitsToInt32(), ab0.InterleaveHi(cd0).ReshapeToUint32s().BitsToInt32(), ab1.InterleaveLo(cd1).ReshapeToUint32s().BitsToInt32(), ab1.InterleaveHi(cd1).ReshapeToUint32s().BitsToInt32()
}
func LoadPixels4(s []byte) archsimd.Int32x4 {
	return archsimd.Uint32x4{}.SetElem(0, binary.LittleEndian.Uint32(s)).ReshapeToUint8s().ExtendLo8ToUint16().BitsToInt16().ExtendLo4ToInt32()
}

func loadU16(s []uint16) archsimd.Int32x4 {
	_ = s[3]
	return archsimd.Uint64x2{}.SetElem(0, *(*uint64)(unsafe.Pointer(unsafe.SliceData(s)))).ReshapeToUint16s().ExtendLo4ToUint32().BitsToInt32()
}
func storeU16(s []uint16, v archsimd.Int32x4) {
	_ = s[3]
	Store4(unsafe.Slice((*int16)(unsafe.Pointer(unsafe.SliceData(s))), 4), v)
}
