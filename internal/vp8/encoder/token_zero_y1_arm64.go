//go:build goexperiment.simd && arm64 && !purego && !govpx_oracle_trace

// Empty Y1 token semantics follow libvpx v1.16.0 vp8/encoder/tokenize.c.
package encoder

import (
	"encoding/binary"
	"github.com/gtsteffaniak/govpx/internal/pixel"
	"github.com/gtsteffaniak/govpx/internal/vp8/tables"
	"simd/archsimd"
	"unsafe"
)

// zeroY1TokenRecords handles a complete empty luma walk with zero contexts.
// The macroblock caller has already checked pointers and consumed Y2.
// Decline before mutation for insufficient capacity or any read/write overlap.
func zeroY1TokenRecords(above, left *TokenContextPlanes, coeffs *MacroblockCoefficients, counts *coefficientTokenCounts, records *InterCoefficientTokenRecords, blockType, skipDC int) bool {
	if !pixel.Enabled || binary.LittleEndian.Uint32(above.Y1[:])|binary.LittleEndian.Uint32(left.Y1[:]) != 0 {
		return false
	}
	if cap(records.Records)-len(records.Records) < 16 {
		return false
	}
	eobs := archsimd.LoadUint8x16Array((*[16]byte)(coeffs.EOB[:16]))
	if eobs.ReduceMax() > byte(skipDC) {
		return false
	}
	index := len(records.Records)
	// These are the only modified windows: the payload, slice header and EOB count.
	payload := unsafe.Pointer(unsafe.Add(unsafe.Pointer(unsafe.SliceData(records.Records)), uintptr(index)*4))
	header := unsafe.Pointer(records)
	counter := unsafe.Pointer(&(*counts)[blockType][skipDC][0][tables.DCTEOBToken])
	cp, ap, lp := unsafe.Pointer(coeffs), unsafe.Pointer(above), unsafe.Pointer(left)
	if denseTokenOverlap(payload, 64, header, 24) || denseTokenOverlap(payload, 64, counter, unsafe.Sizeof(int(0))) || denseTokenOverlap(header, 24, counter, unsafe.Sizeof(int(0))) {
		return false
	}
	for _, w := range [3]struct {
		p unsafe.Pointer
		n uintptr
	}{{payload, 64}, {header, 24}, {counter, unsafe.Sizeof(int(0))}} {
		if denseTokenOverlap(w.p, w.n, cp, unsafe.Sizeof(*coeffs)) || denseTokenOverlap(w.p, w.n, ap, unsafe.Sizeof(*above)) || denseTokenOverlap(w.p, w.n, lp, unsafe.Sizeof(*left)) {
			return false
		}
	}
	value := uint32(tables.DCTEOBToken | (((blockType*tables.CoefBands+skipDC)*tables.PrevCoefContexts)*coefficientProbRowSize)<<coefficientTokenRecordProbOffsetShift)
	packet := archsimd.BroadcastUint32x4(value)
	records.Records = records.Records[:index+16]
	for j := 0; j < 16; j += 4 {
		packet.StoreArray((*[4]uint32)(unsafe.Add(payload, uintptr(j)*4)))
	}
	(*counts)[blockType][skipDC][0][tables.DCTEOBToken] += 16
	return true
}
