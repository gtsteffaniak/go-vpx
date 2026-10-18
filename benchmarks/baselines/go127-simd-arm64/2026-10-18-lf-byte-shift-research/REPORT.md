# ARM byte-shift constant experiment: stopped after screen

Decision: archive the sole candidate. Public picker deltas are small and mixed;
no stable meaningful component gain was established. No fresh-six, whole-codec,
or artificial helper-only benchmark was run. Root source/PRs/README files remain
unchanged.

## Exact source change

Parent `f8dd8c7472e91c1f147511b04ceec716aca95aa7`. The private candidate changes
only `internal/pixel/goav1_loopfilter_arm64.go`: removes two Int8x16 shift globals
and rewrites lfByteShift1, lfByteShift3 and lfByteHalf as ShiftAllRight(1/3/1).
LE/GT/not/other arithmetic, filterShift7/lfShift16, transform helpers, loads,
stores, dispatch and alias checks stay unchanged. Full snapshot comparison
confirmed this single-file production delta. Exact sources/patch and shared
public fixture are included; complete snapshots and binaries are not.

Go 1.27.2 SDK simdARM64.rules 327/331 lowers ShiftAllRight byte operations to
SSHL/USHL with a signed-byte broadcast; rule 392 folds those constant broadcasts
to MOVI.16B. It does not emit immediate SSHR. Actual optimized LLVM confirms the
expected load replacement. Existing exact-f8 profile provenance/saved baseline
LLVM is preserved separately from both fresh test-binary codegen dumps.

| Actual function | Parent instructions | Candidate instructions | Byte global references parent→candidate | Vector stack load/store parent=candidate |
|---|---:|---:|---:|---:|
| loopFilterBytes | 416 | 410 | 3→0 | 2/1 |
| loopFilterMacroInternalFusedHorizontal | 364 | 354 | 4→0 | 0/0 |
| loopFilterMacroInternalFusedVertical | 698 | 680 | 8→0 | 0/19 |

All calls and frame adjustments remain unchanged. Counts include existing STUR
stores (19 in the vertical body); no claim of zero existing spills is made.
SSHL/USHL operation counts are unchanged. Only byte shift-constant setup changes;
16-bit global shift loads remain intact.

## Correctness before timing

Both parent and candidate passed live ARM backend detection, native short suites
for pixel, VP8 DSP/encoder/decoder and VP9 DSP, and race+checkptr2 loop-filter
oracles for pixel, VP8 DSP/decoder and VP9 DSP. Existing independent scalar filter
oracles cover byte arithmetic/exhaustive signed shifts, minimal bounds, aliased
layouts, invalid/fallback geometry and full public outputs. Build/test commands
and complete outputs are preserved under evidence.

## Public screen (native M4 Max)

Official Go 1.27.2, GOEXPERIMENT=simd, PGO off, GOMAXPROCS=1. Three 75ms alternating
parent/candidate rounds (AB,BA,AB) cover all 10 actual public
BenchmarkVP8PickerFilterNative fixtures. Each iteration restores the relevant
source window and calls applyNormalLoopFilterPartialLumaMB. The complete output
buffer must match four separate public edge operations; independent scalar
pixel/DSP oracles ran beforehand. All 60 samples report 0 B/op and 0 allocs/op. No
sample or slower control was discarded. Screen order is preliminary; no
six-balanced-pair confirmation was warranted by these small/mixed changes.

| Public picker case | Parent ns | Candidate ns | Parent/candidate | Candidate time change | Three pair ratios |
|---|---:|---:|---:|---:|---|
| pitch1312/flat | 45.50 | 44.77 | 1.01631× | -1.604% | 0.982044, 1.024614, 1.016306 |
| pitch1312/mixed | 86.15 | 86.89 | 0.99148× | +0.859% | 0.997443, 0.996336, 0.991483 |
| pitch1312/noop | 40.30 | 39.65 | 1.01639× | -1.613% | 1.003279, 1.025445, 1.016546 |
| pitch1312/ramp | 45.32 | 44.61 | 1.01592× | -1.567% | 1.015916, 0.998906, 1.021685 |
| pitch1312/smooth | 86.75 | 85.61 | 1.01332× | -1.314% | 1.017790, 0.997126, 1.008177 |
| pitch32/flat | 45.83 | 44.74 | 1.02436× | -2.378% | 0.937932, 1.008794, 1.024363 |
| pitch32/mixed | 85.74 | 85.43 | 1.00363× | -0.362% | 1.000702, 1.012833, 0.992706 |
| pitch32/noop | 40.12 | 39.91 | 1.00526× | -0.523% | 1.013899, 1.007998, 1.000251 |
| pitch32/ramp | 45.71 | 44.98 | 1.01623× | -1.597% | 1.017119, 1.010056, 1.013975 |
| pitch32/smooth | 85.64 | 86.63 | 0.98857× | +1.156% | 0.982448, 0.995915, 0.989496 |

Median ratios range 0.9886–1.0244×. Candidate times range 2.38% faster to 1.16%
slower, with textured mixed/smooth controls regressing. Fewer instructions and
removed loads alone do not establish a useful public speedup. This report makes
no whole-codec or cross-architecture claim and does not propose default changes.

All six raw logs, every paired ratio/allocation field, complete correctness logs,
optimized parent/candidate LLVM, SDK rules, source/fixture/binary hashes and
original preparation/build/screen/summarization drivers are retained. Drivers
refer to the original private experiment path. SHA256SUMS covers all delivered
files except itself; no binaries/full snapshots are bundled.
