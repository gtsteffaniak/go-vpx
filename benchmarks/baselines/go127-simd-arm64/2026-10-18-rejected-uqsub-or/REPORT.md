# ARM unsigned saturated-absolute difference: rejected after screen

Parent `768f2fc2aba7bbd82981494473cbadc026a95570` (PR43).
Stop/archive: no useful public component gain; no root, README or PR changes.
No fresh-six/whole-codec experiment was warranted.

The sole candidate changes ARM lfAbsByte from unsigned Max−Min to
`a.SubSaturated(b).Or(b.SubSaturated(a))`. An architecture helper split keeps
AMD's original expression. Loads/stores, thresholds, CPU dispatch, aliases,
geometry and fallback contracts remain unchanged. Scoped prior report/patch
review found the rejected clamped signed-saturating-abs experiment, not this
exact unsigned UQSUB/OR expression. This is not a repeat of signed/clamped math.

Pinned SDK lowers unsigned SubSaturatedUint8x16 directly to UQSUB and OR to ORR.
Actual optimized LLVM confirms2UQSUB+OR replacing Max+Min+Sub. Both use3ops;
totals416/364/698 instructions in generic/horizontal/vertical bodies remain
unchanged. Calls unchanged; vector stack stores1/0/19 unchanged. No lower
instruction-count, latency or execution-port speedup is claimed.

Both variants passed all65,536 unsigned byte pairs, live ARM backend detection,
five-package short suites(pixel,VP8 DSP/encoder/decoder,VP9 DSP), and four-package
loop-filter race/checkptr2(pixel,VP8 DSP/decoder,VP9 DSP). Existing fused-filter
oracle covers7,200 geometry/alignment/trial/threshold combinations; aliases,
minimal windows, bounds and sentinel-preserving writes are covered as well.
Complete commands and logs are in evidence; no thresholds/assertions changed.

Native Apple M4 Max; official Go1.27.2; GOEXPERIMENT=simd; PGO off; GOMAXPROCS=1.
Three75ms alternating AB/BA/AB rounds of the actual public picker fixture, with
source window restored each iteration and complete output matched against four
separate public edge operations. All60 samples are0B/op/0allocs. All30 paired
ratios/slower controls retained, including first-parent-flat60.58ns outlier.

| Public case | Parent ns | Candidate ns | Parent/candidate | Candidate time change |
|---|---:|---:|---:|---:|
| pitch1312/flat | 45.33 | 45.38 | 0.99890× | +0.110% |
| pitch1312/mixed | 86.73 | 88.30 | 0.98222× | +1.810% |
| pitch1312/noop | 39.98 | 40.26 | 0.99305× | +0.700% |
| pitch1312/ramp | 45.89 | 45.82 | 1.00153× | -0.153% |
| pitch1312/smooth | 87.44 | 87.08 | 1.00413× | -0.412% |
| pitch32/flat | 45.29 | 45.95 | 0.98564× | +1.457% |
| pitch32/mixed | 86.10 | 87.30 | 0.98625× | +1.394% |
| pitch32/noop | 40.33 | 40.06 | 1.00674× | -0.669% |
| pitch32/ramp | 45.17 | 45.23 | 0.99867× | +0.133% |
| pitch32/smooth | 86.39 | 87.56 | 0.98664× | +1.354% |

Median ratios0.9822–1.0067×: candidate times range0.67%faster to1.81%slower.
Mixed-texture controls regress at both pitches. These near-noise/mixed changes
do not establish useful improvement. No x64 or whole-codec claim. Host released;
no further Go/build/test/timing jobs were started after the screen.

Factual correction to the earlier conversational dual-copy feasibility note:
pinned internal/abi/abi_arm64.go has IntArgRegs=16, FloatArgRegs=16; twelve scalar
argument words do not exceed ARM Go's integer argument-register count. The
incorrect argument is withdrawn. No private dual-copy report file existed and
no production decision adopted it. That caller's actual rejection stands on
intervening UV copies/spatial mutation and different Y-copy sources.

Exact source/helper/test/public fixture, full LLVM, correctness/build logs,
raw6logs/commands, all samples/ratios and environment/source/binary hashes are
preserved. Binaries/full snapshots excluded. Drivers retain original private
experiment paths. SHA256SUMS covers all delivered members except itself.
