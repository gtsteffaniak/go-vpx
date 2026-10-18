//go:build goexperiment.simd && arm64 && !purego

package pixel

import (
	"simd/archsimd"
	"unsafe"
)

var tm4LeftRows = [16]byte{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3}

// TMPredict4Native retains libvpx v1.16.0 vpx_dsp/intrapred.c clipping.
// publicChecked preserves the DSP public entry's original upfront checks.
// Above and left may be hoisted only when both are disjoint from all output.
// Invalid windows and aliases retain the original SIMD row-capture fallback.
func TMPredict4Native(dst []byte, ds int, above, left []byte, top byte, publicChecked bool) {
	if publicChecked {
		_ = above[3]
		_ = left[3]
		_ = dst[3*ds+3]
	}
	if !Enabled {
		TMPredictScalar(dst, ds, above, left, top, 4)
		return
	}
fast:
	for {
		if uint(ds)>>30 != 0 {
			break fast
		}
		span := 3*ds + 4
		if span > len(dst) || len(above) < 4 || len(left) < 4 {
			break fast
		}
		dp, ap, lp := unsafe.Pointer(unsafe.SliceData(dst)), unsafe.Pointer(unsafe.SliceData(above)), unsafe.Pointer(unsafe.SliceData(left))
		d, a, l := uintptr(dp), uintptr(ap), uintptr(lp)
		// A four-byte input overlaps the bounded output span exactly when
		// its start is within [dst-3, dst+span-1]. No pointer addition is used.
		if a-d+3 < uintptr(span)+3 {
			break fast
		}
		if l-d+3 < uintptr(span)+3 {
			break fast
		}
		av := archsimd.BroadcastUint32x4(*(*uint32)(ap)).ReshapeToUint8s()
		lv := archsimd.Uint32x4{}.SetElem(0, *(*uint32)(lp)).ReshapeToUint8s().LookupOrZero(archsimd.LoadUint8x16Array(&tm4LeftRows))
		t := archsimd.BroadcastUint8x16(top)
		pos, neg := lv.SubSaturated(t), t.SubSaturated(lv)
		tile := av.AddSaturated(pos).SubSaturated(neg)
		if ds == 4 {
			tile.StoreArray((*[16]byte)(dp))
		} else {
			rows := tile.ReshapeToUint32s()
			*(*uint32)(dp) = rows.GetElem(0)
			*(*uint32)(unsafe.Add(dp, ds)) = rows.GetElem(1)
			*(*uint32)(unsafe.Add(dp, 2*ds)) = rows.GetElem(2)
			*(*uint32)(unsafe.Add(dp, 3*ds)) = rows.GetElem(3)
		}
		return
	}

	// Exact original SIMD evaluation order for four-pixel rows: capture left
	// once per row, then form the output slice before loading the above packet.
	t := archsimd.BroadcastInt32x4(int32(top))
	for y := 0; y < 4; y++ {
		v := archsimd.BroadcastInt32x4(int32(left[y])).Sub(t)
		StorePixels4(dst[y*ds:], LoadPixels4(above).Add(v))
	}
}
