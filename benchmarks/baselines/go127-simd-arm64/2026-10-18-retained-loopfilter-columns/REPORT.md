# ARM fused vertical input-column retention

The isolated production change keeps the first sixteen input columns in their final transposed form across the sequential VP8 macro and internal edges. This removes compiler-generated stack traffic while preserving the corrected-neighbor carry and edge ordering. The helper, generated ARM source and generator are synchronized; horizontal source is unchanged. No assembly, dependency or README changes are involved.

The fixed source is based on bc4a12a2ea74f6219d26acb5646677016eb6423f. Root integration and broad validation ran on merged PR51 base 6bda90f7eeee6f407d6ed53e299eab505afa440e; the five integrated files match final-changed-source.json exactly. All native measurements below are Apple M4 Max / Darwin arm64 / Go 1.27.2 / GOEXPERIMENT=simd / CGO=1 / PGO off / GOMAXPROCS=1, in explicitly serialized quiet windows.

## Causal code-generation result

Xcode LLVM disassembly of the actual executables shows 733→658 instructions and a 592→320-byte stack frame in loopFilterMacroInternalFusedVertical. Exactly 75 instructions disappear: 17 vector spill stores, 17 corresponding reloads, 34 stack-address additions and 7 vector moves. The transpose shuffle count is unchanged. The compiler already eliminated repeated overlapping source-level row loads in the parent; this report does not claim 32 fewer input loads. codegen-summary.json contains the counts and full LLVM captures are included.

The old Go assembly at 691cd05 keeps paired input columns live across dependent edges. The new Go SIMD representation applies that dataflow principle to the fused 16-row implementation. This is the implemented extra-work removal. Inspection of the Go picker against pinned libvpx picklpf.c found the same trial order and copy/search scope, with no demonstrated extra trial count. Remaining observable ARM instruction differences include three-operation byte absolute difference instead of NEON UABD and expanded correction arithmetic; Go 1.27.2 archsimd does not expose all corresponding instructions. These are leads, not measured explanations for the distinct x64 Go/libvpx codec gap.

## Correctness and memory contract

The existing public guards require disjoint rows of at least 20 bytes and a complete 20×16 window. Initial16-byte loads and the unchanged final row+4 loads fit that window. Retained columns 8..11 and 12..15 are consumed only at the original edge boundaries, before earlier edges can modify those incoming columns; corrected carry values remain authoritative. Arithmetic, masks, HEV decisions, rounding, changed flag and final stores to columns 1..17 are unchanged.

The new direct oracle checks all 16 columns using indexed scalar source bytes over 1,152 pitch/alignment/pattern combinations, clipped spans and full backing immutability. Existing scalar sequential-edge/full-backing correctness and the benchmark pre/post checks pass. Pixel/DSP/decoder tests, targeted race/checkptr, root loop-filter tests and root broad validation pass. Root receipts include the full Go suite, purego suites, generator replay, trace compilation, refreshed PGO, pre-commit, no-assembly and all six build-selection modes. The PGO training results are correctness/build validation, excluded from throughput claims here.

A live pinned-libvpx packet oracle executes with its exact flags in live-oracle-binding.json and final-live-oracle.log: 30 inputs, 19 emitted packets, all packet bytes identical. This is a 64×64 rate-control/drop oracle, not a C performance measurement and not the 720p timing workload. The720p whole runs skip decoded quality and report PSNR/SSIM=0 with quality_frames=0; matching those reported fields is not a decoded-quality proof.

## Whole encoder windows

Each window contains six all-permutation parent/candidate/old691 triples. Workload: VP8, 1280×720, 300 inputs, fps 30, 2500kbps,realtime CPU 8,one thread,encode-only,auto-libvpx=false. The windows remain separate. Each has 22 identical reported semantic fields across 18 samples and zero allocations. Initial candidate source is preserved in a separate local archive whose hash is recorded in reference-binding.json; final source/executables are bound by final-binding.json.

| Window | Parent median ns/frame | Candidate median | Old691 median | Candidate time vs parent | Candidate time vs old |
| --- | ---: | ---: | ---: | --- | --- |
| Initial | 7281296.5 | 7182789.5 | 7021717.0 | -1.352877%; 6/6 pairs faster | +2.293919%; 6/6 pairs slower |
| Final | 7224948.0 | 7122281.0 | 7060167.0 | -1.421007%; 5/6 pairs faster | +0.879781%; 4/6 pairs slower |

The final window includes one +2.831397% paired loss versus parent. The initial candidate is 2.293919% slower than old at medians (all 6pairs slower); final is 0.879781% slower (4/6 pairs slower). Neither number closes the gap or proves a universal effect. All samples, including outliers and slower controls, remain in the raw files and summary JSON.

## Component controls

Each component fixture verifies independent scalar/full backing before and after timing, restores input during B.Loop, and consumes output. Initial: 6 balanced pairs at 75ms/case, 408 rows. Final: 10 balanced pairs at 150ms/case, 680 rows. There are 24 axis cases plus 10 picker cases per executable (34 total); the initial independent review incorrectly called the axis fixture 48 cases, conflating its two variants. Its original receipt is preserved and the final independent audit corrects that count. All component rows report 0B/op and 0allocs/op.

Speedup means parent time divided by candidate time. These gains are below 1.5×; the keeper justification is the causal removal of extra work plus repeated whole-encoder benefit. All final median losses are retained below.

| Case | Initial speedup | Final speedup | Final pairs faster /10 |
| --- | ---: | ---: | ---: |
| Axis/verticalfalse/pitch32/flat | 1.004402 | 0.860257 | 0 |
| Axis/verticalfalse/pitch32/ramp | 1.001544 | 0.871183 | 0 |
| Axis/verticalfalse/pitch32/smooth | 0.996746 | 0.984011 | 2 |
| Axis/verticalfalse/pitch32/mixed | 0.990590 | 0.984865 | 3 |
| Axis/verticalfalse/pitch32/noop | 0.996694 | 0.868065 | 1 |
| Axis/verticalfalse/pitch32/macroOnly | 0.999528 | 0.890470 | 0 |
| Axis/verticalfalse/pitch1312/flat | 1.000517 | 0.865963 | 0 |
| Axis/verticalfalse/pitch1312/ramp | 1.008348 | 0.867935 | 0 |
| Axis/verticalfalse/pitch1312/smooth | 1.004818 | 0.981003 | 3 |
| Axis/verticalfalse/pitch1312/mixed | 1.004978 | 0.986283 | 3 |
| Axis/verticalfalse/pitch1312/noop | 1.008140 | 0.868449 | 0 |
| Axis/verticalfalse/pitch1312/macroOnly | 1.004247 | 0.892243 | 0 |
| Axis/verticaltrue/pitch32/flat | 1.019762 | 1.019181 | 9 |
| Axis/verticaltrue/pitch32/ramp | 1.026581 | 1.020384 | 9 |
| Axis/verticaltrue/pitch32/smooth | 1.002867 | 1.002774 | 7 |
| Axis/verticaltrue/pitch32/mixed | 1.010871 | 0.998039 | 3 |
| Axis/verticaltrue/pitch32/noop | 1.019713 | 1.021915 | 10 |
| Axis/verticaltrue/pitch32/macroOnly | 1.030780 | 1.035386 | 9 |
| Axis/verticaltrue/pitch1312/flat | 1.020653 | 1.018628 | 9 |
| Axis/verticaltrue/pitch1312/ramp | 1.023889 | 1.009518 | 7 |
| Axis/verticaltrue/pitch1312/smooth | 0.994011 | 0.999272 | 5 |
| Axis/verticaltrue/pitch1312/mixed | 1.006103 | 0.999561 | 6 |
| Axis/verticaltrue/pitch1312/noop | 1.025308 | 1.015395 | 9 |
| Axis/verticaltrue/pitch1312/macroOnly | 1.032379 | 1.028494 | 9 |
| Picker/pitch32/flat | 1.028410 | 1.035090 | 10 |
| Picker/pitch32/ramp | 1.035838 | 1.038197 | 10 |
| Picker/pitch32/smooth | 0.983488 | 0.977017 | 2 |
| Picker/pitch32/mixed | 0.984353 | 0.976972 | 1 |
| Picker/pitch32/noop | 1.014055 | 1.014124 | 10 |
| Picker/pitch1312/flat | 1.032783 | 1.031977 | 9 |
| Picker/pitch1312/ramp | 1.038376 | 1.026642 | 9 |
| Picker/pitch1312/smooth | 0.981888 | 0.974255 | 1 |
| Picker/pitch1312/mixed | 0.982425 | 0.980234 | 1 |
| Picker/pitch1312/noop | 1.013384 | 1.024444 | 10 |

## Horizontal control investigation

The final shared H/V benchmark has large H losses on cheap inputs. The H leaf has exactly the same 364 instruction lines, bytes and address 0x1001b9310 in original parent, initial candidate and final candidate. Final vs parent shared caller has identical 660 instructions after relocation normalization, identical public guard logic and H/memmove targets, but a different address:parent 0x1001f7990, initial 0x1001f2b70, final 0x1001f7860. Its timed H call is at 0x1001f82fc /0x1001f34dc /0x1001f81cc. The timed input-reset memmove target remains 0x100086b60. Actual disassemblies and normalized diffs are included.

A separate frozen diagnostic compares a 504-instruction H-only owner, with no V call and identical fixtures in the two source snapshots. It retains the scalar/full-backing checks and reset semantics; the compiler can specialize its H-only copy/guards, so its absolute timings are not directly interchangeable with the shared fixture. Six balanced 100ms pairs cover 144 rows. Seven of 12 medians remain losses ; 38 of 72 paired measurements are slower. Every result has zero allocations.

| Independent H case | Speedup | Pairs faster /6 |
| --- | ---: | ---: |
| pitch32/flat | 1.001173 | 4 |
| pitch32/ramp | 0.989813 | 1 |
| pitch32/smooth | 0.987233 | 1 |
| pitch32/mixed | 1.009315 | 4 |
| pitch32/noop | 0.995448 | 4 |
| pitch32/macroOnly | 0.993213 | 1 |
| pitch1312/flat | 1.013075 | 3 |
| pitch1312/ramp | 0.975531 | 3 |
| pitch1312/smooth | 1.006441 | 3 |
| pitch1312/mixed | 0.988949 | 2 |
| pitch1312/noop | 0.995685 | 2 |
| pitch1312/macroOnly | 1.014493 | 6 |

Six all-permutation triples using the untouched original binaries reproduce the shared flat discrepancy:initial candidate 1.001539×/1.002474× at pitches 32/1312 (3/6 pairs faster each), final 0.864084×/0.864862× (all 12 paired measurements slower). These 36 rows are preserved with the 144 independent rows. This establishes executable/caller-dependent behavior, not added H arithmetic or a uniquely proven hardware/layout mechanism. The original bad controls have not been discarded, fixed, or averaged away. There is no universal horizontal-performance claim.

## Reproduction and scope

measure.py/final-measure.py/horizontal-measure.py, commands.json files, source maps, executable hashes and initial/final/review receipts bind every window. Independent 1x preflight logs are correctness-only and excluded from performance summaries. Source snapshots and binaries stay in their recorded local paths; this publication bundle intentionally excludes executable binaries and full source trees. Changed sources are captured with .txt suffixes to avoid creating accidental Go packages. The initial full source archive remains local and is hash-bound.

This ARM fix does not explain the separate native x64 fixed-work Go/libvpx 2.016× time ratio or claim an x64 speedup. The old Go assembly comparison, the C packet oracle and remote x64 C timing are distinct references. A separate causal review of the x64 feature-first quantizer guard follows; it is not part of these five source files.


Publication: `raw-evidence.tar.gz` preserves all 343 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
