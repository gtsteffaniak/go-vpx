# Private x64 SSE8 AVX2 codegen/correctness probe

No native x64 speedup is established. No root checkout, CI head, README,
Git branch or PR changes were made. The variants remain private and require
native benchmark evidence before acceptance.

The baseline common SSE8 packs two exact8-byte rows into a128-bit packet,
computes byte Max-Min and two128-bit square/dot accumulators. The new x64
helper zero-extends the full16-byte packet into256-bit int16 lanes, performs
signed subtraction, and adds one256-bit DotProductPairs result. Differences
are[-255,255], so each product and pair sum is exact. The complete8x8 sum is
at most4,161,600; the bounded64-row slice sum at most33,292,800 fits int32.

SDK1.27.2 has no Uint64x4.SetElem. The first private attempt's compile errors
are retained. The accepted packing form uses BroadcastUint64x4 of exact64-bit
loads, InterleaveLoGrouped (actual SDK spelling), GetLo, then ExtendToUint16.
LLVM21 shows VMOVQ, VPBROADCASTQ, VPUNPCKLQDQ, VEXTRACTI128, VPMOVZXBW,
VPSUBW, VPMADDWD and VPADDD. Broadcasts are register-source, not fused memory
broadcasts. Old Go objdump cannot decode these AVX2 bytes correctly; LLVM is
authoritative. Upper bits clear through ClearAVXUpperBits before return.

## Variants and codegen results

- slice-bool: helper called only after existing Enabled/complete endpoint
  checks AND common h>=8 guard. ARM false stub inlines away. The AMD helper
  declines unsupported height/CPU defensively; old four-row path remains.
- slice-direct: AMD compile-time availability true, ARM false. Common
  h>=8+availability dispatch returns helper uint32 directly; there is no
  runtime bool fallback or its argument preservation. Caller endpoint and
  CPU guards own the helper's preconditions. ARM optimized code has no helper
  call. Exact odd row uses[a,0] and[b,0], reading only8bytes each.
- pointer-unrolled4-bool: four explicit packets pass all correctness checks,
  but compiler hoists16 row loads, spilling/reloading two XMM temporaries with
  legacy MOVUPS amid live YMM data. This variant is blocked by codegen; it is
  retained to show why source unrolling alone is not an improvement.
- pointer-loop4-bool and pointer-loop2pairs-bool: fixed4-iteration loop and
  two-pair body repeated twice, respectively. Both retain the bool API and
  compile with clean VEX-only hot packets/reduction, no hot calls or spills.
- pointer-direct: separate two-pair/two-iteration uint32 API. Unsupported CPU,
  nil or unsupported strides execute EXACT original SSE8(unsafe.Slice(...))
  expression inside pixel. DSP forwards one call. Its wrapper is inlined
  away; saved public caller invokes pixel.SSE8Ptr with four register args,
  without the bool fallback/header construction. Linux/Darwin fast body is
  clean VEX-only; fallback and stack-growth calls are outside the fast body.

Trusted pointer guards require Enabled, actual AVX2, nonnil pointers and
nonnegative strides below2^30. The existing public pointer contract already
requires both complete8x8 windows; no object length is invented. Stride and
row arithmetic fits64-bit. Read-only aliases, zero strides and overlapping
rows preserve results. Common slice geometry/bounds and panic behavior are
unchanged. Purego/experiment-off/other architectures retain original fallback;
final generic stub excludes BOTH arm64 and amd64. Final ARM pointer stub,
purego and experiment-off tests pass.

## Verification and limits

All credible variants pass existing TestSquaredError/TestSSE8ExactRows,
including every byte pair, heights1..65 represented by the existing matrix,
short inputs, asymmetric/zero/overlapping strides, minimal spans, alignments,
extreme differences and immutable backing canaries. Pointer variants also
have independent exact-window/alignment/alias/every-byte-pair tests and all36
public fixture one-shot oracles, zero allocations. Rosetta test logs explicitly
show AVX2=true. Checkptr1 repeats correctness. Actual AVX/AVX2-off runs show
false feature flags, scalar results and preserved fallback behavior.

The direct pointer test compares original and new invalid-geometry failures
in subprocesses: checkptr throws for invalid unsafe.Slice geometry and cannot
be tested with recover. Ordinary and instrumented nil/negative/extreme
fallback behavior agrees. Positive huge strides are not fabricated with tiny
objects: both old/new APIs require the caller-validated complete object.

No emulated speed claim is made from the one-shot ns/op values. Native x64
six-pair representative measurements, including four-row controls and caller
costs, remain required. Both direct variants are plausible candidates only.

All relevant exact sources, drivers, raw logs, codegen and audit fields are
retained; binaries/full source snapshots are excluded. Recreate private
archives from parent-commit.txt, overlay each variant's sources and run the
exact drivers from the root checkout/cache scope. SHA256SUMS covers every
packaged file except itself.
