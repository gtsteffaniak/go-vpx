//go:build goexperiment.simd && arm64 && !purego

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
func high4(v archsimd.Int16x8) archsimd.Int32x4   { return v.HiToLo().ExtendLo4ToInt32() }
func mac(lo, hi archsimd.Int32x4, v, c archsimd.Int16x8) (archsimd.Int32x4, archsimd.Int32x4) {
	return lo.Add(v.MulWidenLo(c)), hi.Add(v.HiToLo().MulWidenLo(c))
}
func packBytes(v archsimd.Int16x8) archsimd.Uint8x16 { return v.SaturateToUint8() }

// Same narrowing layout as goav1's planePackPairI16/planeNarrowWord.
func store8(b []byte, lo, hi archsimd.Int32x4) {
	l := lo.SaturateToInt16().ToBits().ReshapeToUint64s()
	h := hi.SaturateToInt16().ToBits().ReshapeToUint64s()
	storeBytes8(b, l.InterleaveLo(h).ReshapeToUint16s().BitsToInt16().SaturateToUint8())
}
func sadAdd(acc archsimd.Uint32x4, a, b archsimd.Uint8x16) archsimd.Uint32x4 {
	d := a.Max(b).Sub(a.Min(b))
	lo := d.ExtendLo8ToUint16()
	hi := d.HiToLo().ExtendLo8ToUint16()
	return acc.Add(lo.ExtendLo4ToUint32()).Add(lo.HiToLo().ExtendLo4ToUint32()).Add(hi.ExtendLo4ToUint32()).Add(hi.HiToLo().ExtendLo4ToUint32())
}
func sumSquares(v archsimd.Int32x4) uint64 {
	s := v.MulWidenLo(v).Add(v.HiToLo().MulWidenLo(v.HiToLo()))
	return uint64(s.GetElem(0)) + uint64(s.GetElem(1))
}
func Store4(s []int16, v archsimd.Int32x4) {
	v = v.ShiftAllLeft(16).ShiftAllRight(16)
	v.SaturateToInt16().StorePart(s[:4])
}
func StorePixels4(s []byte, v archsimd.Int32x4) {
	v.SaturateToInt16().SaturateToUint8().StorePart(s[:4])
}

func sum32(v archsimd.Int32x4) int32    { return v.ReduceSum() }
func sumU32(v archsimd.Uint32x4) uint32 { return v.ReduceSum() }
func Max4(v archsimd.Int32x4) int32     { return v.ReduceMax() }

// Four int16 outputs occupy precisely the eight bytes written here. This
// store is unaligned-safe on amd64 and arm64 and stays within the Go slice.
func put4(out []int16, v uint64) { _ = out[3]; *(*uint64)(unsafe.Pointer(&out[0])) = v }
