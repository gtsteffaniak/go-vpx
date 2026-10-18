//go:build goexperiment.simd && (arm64 || amd64) && !purego

package encoder

import (
	"simd/archsimd"
	"unsafe"

	"github.com/thesyncim/govpx/internal/pixel"
	"github.com/thesyncim/govpx/internal/vp9/common"
)

// blockYrdScore specializes the libvpx v1.16.0 vp9_pickmode.c:781-853
// transform/quantize/SATD/error walk to its aggregate-only consumer. Neither
// qcoeff/dqcoeff nor a scan position escapes BlockYrd. A zero sum(abs(q)) is
// exactly EOB==0; the DC-only EOB==1 rate is the same sum.
func blockYrdScore(src []int16, bw, bh int, tx common.TxSize,
	r, q, d [2]int16, sse uint64, coeff []int16, nonzero bool,
) (BlockYrdResult, bool) {
	if !pixel.Enabled || tx < common.Tx8x8 || tx > common.Tx16x16 ||
		d[0] < 3 || d[0] > 2048 || d[1] < 3 || d[1] > 2048 {
		return BlockYrdResult{}, false
	}
	n := 4 << uint(tx)
	units := (bw / n) * (bh / n)
	res := BlockYrdResult{Valid: true, Skippable: true}
	known := sse < BlockYrdUnknownSSE
	if known {
		res.SSE = int64(sse << 4)
	}
	// Subtraction already reports this fact. Zero Hadamard input produces
	// only zero coefficients, and these quantizer rounds are below dequant.
	if !nonzero {
		if known {
			res.Dist = res.SSE
		} else {
			res.Rate = units << VP9ProbCostShift
		}
		return res, true
	}
	var rate int
	var dist int64
	for y := 0; y < bh; y += n {
		for x := 0; x < bw; x += n {
			hadamardBlockYrd(src[y*bw+x:], bw, coeff, n)
			rtx, etx := blockYrdQuantScore(coeff, n*n, r, q, d)
			rate += rtx
			dist += int64(etx >> 2)
		}
	}
	res.Skippable = rate == 0
	if res.Skippable && known {
		res.Dist = res.SSE
		return res, true
	}
	res.Rate = rate<<(2+VP9ProbCostShift) + units<<VP9ProbCostShift
	res.Dist = dist
	return res, true
}

// The caller supplies n=64/256 complete coefficients and the exact BlockYrd
// tables r=floor(factor*d/128), q=floor(65536/d), with 3<=d<=2048.
// Thus q>=0 and qabs*d<=32767: dequantization cannot wrap. Sign restoration
// can be omitted because both abs(q) and (coeff-dq)^2 are sign invariant.
// For a=abs(coeff), clipped=min(a+r,32767), the reciprocal truncation gives
// -r<=a-qabs*d<1+1.5*d<=3073. Consequently every squared error fits int32,
// each four-lane accumulator is below 2^31, and the 256-coefficient total
// is below 2^32. Larger or unusual tables retain the general path.
func blockYrdQuantScore(coeff []int16, n int, r, q, d [2]int16) (int, uint32) {
	cp := unsafe.Pointer(unsafe.SliceData(coeff))
	acR := archsimd.BroadcastUint16x8(uint16(r[1]))
	acQ, acD := archsimd.BroadcastInt16x8(q[1]), archsimd.BroadcastInt16x8(d[1])
	vr := acR.SetElem(0, uint16(r[0]))
	vq, vd := acQ.SetElem(0, q[0]), acD.SetElem(0, d[0])
	max16 := archsimd.BroadcastUint16x8(32767)
	rate, err := archsimd.Uint32x4{}, archsimd.Uint32x4{}
	for i := 0; i < n; i += 8 {
		c := archsimd.LoadInt16x8Array((*[8]int16)(unsafe.Add(cp, i*2)))
		mag := c.Abs().ToBits()
		a := mag.Add(vr).Min(max16).BitsToInt16()
		t := blockYrdMulHigh16(a, vq).Masked(a.GreaterEqual(vd))
		rate = rate.Add(blockYrdSumPairs16(t))
		diff := mag.Sub(t.Mul(vd).ToBits()).BitsToInt16()
		err = err.Add(blockYrdSquarePairs16(diff))
		vr, vq, vd = acR, acQ, acD
	}
	return int(blockYrdSum32(rate)), blockYrdSum32(err)
}
