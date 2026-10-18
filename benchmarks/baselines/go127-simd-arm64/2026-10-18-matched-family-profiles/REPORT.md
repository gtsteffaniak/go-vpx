# Matched ARM VP8 1200-input profile diagnostic

Exact old build `691cd0512c4829d9ee40d483205dd223f9eb028d` and parent `f8dd8c7472e91c1f147511b04ceec716aca95aa7`, Apple M4 Max, official Go 1.27.2 SIMD, PGO off, one thread/GOMAXPROCS=1. Existing frozen CLIs were hash-checked and executed without rebuilding. See commands.json and environment-and-scope.json. The parent is the measured pre-Copy8 source, rather than the newly staged Copy8 successor.

Two serial profiles used identical 1280×720 VP8 realtime CPU-used=8, 1200 inputs, 30fps, 2500kbps, encode-only controls. Allocation sampling precedes the measured encode pass; CPU profiling includes only that pass. This is diagnostic sampling, not a repeated native benchmark or speed acceptance experiment. Positive CPU-used=8 is adaptive; this run nevertheless matched all recorded workload fields.

## Work parity

All **22 semantic fields and separate zero allocations** match old691, fresh f8 and root's prior f8 1200-input profile exactly. Each encoded 1186 and dropped 14 inputs, output 12,523,129 bytes, keyframe 133,700 bytes, with identical quantizer statistics/histogram. Quality was skipped in every run. Equality of these recorded fields does not establish bitstream identity. Raw JSON and parity.json retain the complete evidence.

## Flat-only operation families

Each CPU sample is assigned exactly once to its nearest recognized operation ancestor; its leaf remains its flat attribution. This attributes inlined utility costs to operations, without adding overlapping cumulative percentages. Explicit classification rules, every leaf contribution and unclassified residual are in family-audit.json. Full pprof flat and cumulative views are separate, with nodecount=0/nodefraction=0.

| Family | old691 sampled seconds | fresh f8 | prior f8 |
|---|---:|---:|---:|
| Loop filter | 1.79 | 2.26 | 2.34 |
| SSE/variance/SAD | 0.36 | 0.60 | 0.66 |
| DC pair inverse/add | 0.64 | 0.17 | 0.28 |
| Prepared tokens/counting/bool writer | 1.29 | 1.20 | 1.07 |
| Forward transforms | 0.35 | 0.39 | 0.47 |
| Quantization | 0.22 | 0.26 | 0.19 |
| Copies/memmove | 0.47 | 0.46 | 0.31 |
| All CPU samples | 8.04 | 8.47 | 8.58 |

The families indicate where further investigation has potential, rather than explaining the net elapsed-time gap causally. One 10ms-sampled capture is noisy; the same f8 build varies substantially by function between the two profiles. Frame-work equality does not remove sampling, layout or profiling effects.

## Most credible remaining family target

Horizontal luma macro plus internal filtering is the clearest remaining diagnostic gap. Old active `mbLoopFilterEdgeH16NEON` + `loopFilterHorizontalEdgesYSharedNEON` account for **0.33s**; current `loopFilterMacroInternalFusedHorizontal` and its inlined helpers account for **0.86s fresh / 0.69s prior**. By comparison, vertical macro+internal is old **1.05s**, fresh **0.79s**, prior **1.07s**: it is hot, but these samples do not establish it as a stable loser. `orientation-and-plane.json` preserves the exact operation-root rule and flat leaves.

The existing optimized frozen-binary disassembly identifies a real instruction-level difference: old horizontal NEON uses UABD per unsigned-byte absolute difference, while current `lfAbsByte` lowers to UMAX+UMIN+SUB. Current generated horizontal source calls that helper for threshold/HEV/neighbor conditions, and extracts two mask halves before branch decisions. The compiler's current archsimd API has no UABD/absolute-difference method. This establishes an avoidable cost relative to old machine code but does not yet supply an available Go API replacement. A three-op equivalent saturated-difference formulation or nearby constant/branch rewrite is not automatically an improvement. Do not re-run rejected shift-constant, ReduceMax, full-tile or partial-C/D cache designs on this evidence alone. No new candidate is proposed or implemented here.

Plane scoring is secondary: old column SSE (`sseBlock16xNDotProd`) is **0.19s**, current row-major `SSEPlane16` including inlined arithmetic **0.37s fresh / 0.30s prior**. The old active DOTPROD arithmetic differs from current signed widening/multiply; prior unsigned-square/SSE16 probes are already archived below 1.5×. These samples justify keeping the operation visible, not repeating those experiments without a distinct dataflow hypothesis.

The prepared token writer remains a large current cost, but the matched old writer is at least as sampled-heavy (0.73s leaf vs 0.55/0.45 current). This profile does not make it a new assembly replacement deficit. DC-pair additions and residual gathering have substantial sampled improvements; preserve their accepted kernel wins.

No production files, branches, index, README or dependencies were changed. No new Go build or performance-acceptance benchmark was run. pprof was invoked through the root project wrapper with an agent-specific ID. Saved sources are narrowly scoped to the relevant exact parent/old implementation. SHA256SUMS covers all report/raw/source/driver files except itself; existing binary hashes are separately pinned rather than copying binaries.
