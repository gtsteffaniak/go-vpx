//go:build goexperiment.simd && arm64 && !purego

// Fixed shifts preserve the libvpx v1.16.0 VP8 transform arithmetic.
package encoder

import "simd/archsimd"

var vp8ShiftLeft2 = archsimd.BroadcastInt32x4(2)

func vp8Left2(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftLeft2) }

var vp8ShiftLeft3 = archsimd.BroadcastInt32x4(3)

func vp8Left3(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftLeft3) }

var vp8ShiftRight3 = archsimd.BroadcastInt32x4(-3)

func vp8Right3(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftRight3) }

var vp8ShiftRight4 = archsimd.BroadcastInt32x4(-4)

func vp8Right4(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftRight4) }

var vp8ShiftRight12 = archsimd.BroadcastInt32x4(-12)

func vp8Right12(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftRight12) }

var vp8ShiftRight16 = archsimd.BroadcastInt32x4(-16)

func vp8Right16(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftRight16) }

var vp8ShiftSign = archsimd.BroadcastInt32x4(-31)

func vp8Sign(v archsimd.Int32x4) archsimd.Int32x4 { return v.Shift(vp8ShiftSign) }

func vp8Zero4(v archsimd.Int32x4) bool { return v.ToBits().ReduceMax() == 0 }
