# Private AMD SSE8 entry repair

Exact base c638b7c76eac857affa0780820e86c211cdd84d3 (PR42 remains unchanged by this work). Go 1.27.2, GOEXPERIMENT=simd, GOAMD64=v1, PGO off. No performance claim: Rosetta runs and one-shot public fixtures establish correctness only. Native CI must measure whether this repairs the EPYC regressions.

## Source changes

The shared file retains ssePacket. The ARM SSE8 function moves verbatim into an ARM-only file; ARM compile-time-false stubs remain unchanged. The AMD entry now contains the 256-bit pair-row loop directly, after the original complete-window endpoints. It reduces into a scalar, executes inline ClearAVXUpperBits, and returns. There is no SSE8-to-wide-helper call. The narrow h<8 path and unusual-geometry fallback restore exact pre42 ab3 arithmetic and ordinary zero initialization, removing the new global-zero loads. The now-unreferenced AMD zero-global/helper file is deleted. The pointer source is identical to the previously reviewed direct-fence variant. The wide dispatch test now checks the public entry; a new test checks disabled-backend scalar/alias/canary behavior, including odd/boundary heights and unused extreme strides.

Bounds, read-only alias behavior, numeric range, input immutability, endpoint panic order and original CPU-off scalar behavior are preserved. Both bounded slice endpoints execute before wide arithmetic. An internal valid fast path has no remaining language-level panic site. The trusted pointer contract remains complete, valid 8x8 windows; removing defer does not promise cleanup on exceptional unwinding after invalid nonnil memory is supplied. Unsupported CPU/nil/unusual pointer geometry retains the original fallback before any YMM arithmetic.

## Correctness

Selected AVX2 and actual AVX/AVX2-off AMD pixel tests pass under checkptr=2. Coverage includes all byte pairs, exact/minimal windows, zero/overlapping/unequal strides, mutually aliased read-only inputs, bounds/fallback/invalid-pointer diagnostics and backend-disabled scalar checks. All 36 identical public DSP one-shot scalar/canary fixtures pass under checkptr=2 both selected and CPU-off, with zero allocation readings. Fresh native ARM focused tests, purego and nosimd pixel/DSP checks also pass. Complete raw outputs are retained. No benchmark latency from these runs is used as performance evidence.

## Optimized code

Linux and Darwin SSE8 each decode to 373 instructions across all branches, scalar fallback, cold panic sites and alignment padding; pointer remains 149. Wide-reachable CFGs contain zero helper/callee calls, zero legacy vector operations, zero vector stack traffic and exactly one VZEROUPPER after GPR result extraction. The only calls anywhere in the slice function are cold runtime.panicBounds. Four restored legacy XMM zero moves belong exclusively to h<8/unusual branches, are bypassed by the wide path, and match the pre42 initialization design. The slice entry has a 16-byte stack frame and no stack-growth/callee call; pointer's cold fallback and stack-growth calls remain. Instruction counts and CFG diagnostics are not native speed measurements.

## Transplant and reproduction

Apply candidate.patch to the exact base; candidate/ pins changed/new files and parent/ pins previous versions including the deleted AMD zero source. audit.json records source/binary hashes, source identity assertions, both target CFGs and public oracle counts. all-go-module-source-sha256.txt covers the private snapshot. reproduce.sh creates a private snapshot from the exact base, applies the patch and invokes the root project wrapper for correctness and codegen. No root files, branch/index or README was changed by this experiment.
