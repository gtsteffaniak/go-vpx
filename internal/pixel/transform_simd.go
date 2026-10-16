//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import "simd/archsimd"

// Load4 reads exactly four int16 samples, including minimally sized rows.
func Load4(s []int16) archsimd.Int32x4 {
	v, _ := archsimd.LoadInt16x8Part(s[:4])
	return v.ExtendLo4ToInt32()
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
	v, _ := archsimd.LoadUint8x16Part(s[:4])
	return v.ExtendLo8ToUint16().BitsToInt16().ExtendLo4ToInt32()
}
