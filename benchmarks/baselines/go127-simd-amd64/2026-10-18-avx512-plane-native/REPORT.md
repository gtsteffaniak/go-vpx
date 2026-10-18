# PR40: native AVX-512 plane-error candidate

Exact head `34c4e4a17b231a24a3d3774c9621b7e5b80bac03`, parent `c2d415fffead020022bd7bb6f40a7bd665795f44`. All three checks passed. Root independently reviewed the evidence and authorized the guarded merge. Whole-codec improvement is not claimed.

Native CPU: INTEL(R) XEON(R) PLATINUM 8573C; go version go1.27.2 linux/amd64; PGO off, GOMAXPROCS=1. Separate C CPU: AMD EPYC 7763 64-Core Processor. Do not combine or causally compare measurements between jobs.

The full AVX512 F/CD/BW/DQ/VL feature bundle was exposed. TestSSEPlane16AVX512 executed and passed in ordinary and race/checkptr validation, with AVX2=true and AVX512=true both times and no skip. Every changed fixture has an eligible width. The current VP8 profile sampled ssePlane32AVX512 for 0.01s, independently confirming public-path liveness.

## Four changed public plane fixtures

| Fixture | Parent speedup | Minimum paired speedup | Pairs faster | Old-reference speedup |
|---|---:|---:|---:|---:|
| 1088p | 0.977455× | 0.956588× | 2/6 | 1.945351× |
| 128x128 | 1.590938× | 1.572279× | 6/6 | 3.884434× |
| 720p | 1.556973× | 1.488458× | 6/6 | 2.647709× |
| 720pPartial | 1.672263× | 1.650316× | 6/6 | 6.939267× |

128x128 and 720pPartial exceed 1.5× parent in every pair. 720p median exceeds 1.5×, but its minimum pair is 1.488458×. **1088p regresses**: current median takes 2.306549% more time than parent, with four of six pairs slower. No blanket AVX512 or frame-size win is supported.

## Whole-codec results

| Public case | Current time vs parent | Pairs faster | Current time vs old build |
|---|---:|---:|---:|
| VP9 320x180 | +0.757471% | 2/6 | -43.860896% |
| VP8 720p | +0.684300% | 2/6 | -28.375147% |
| VP9 720p | +0.719344% | 2/6 | -52.735243% |

VP8 is the changed-code candidate and is slower in this whole-codec run; VP9 is unchanged control work. All old-build public comparisons are faster on this CPU. These results do not close remaining ARM gaps or establish cross-CPU performance.

## Correctness, allocation and controls

All 330 cases × six samples × three variants have zero bytes and allocations. All 16 recorded source SHA256 values match exact-head Git blobs. Native validation records 123 top-level test executions plus six retained panic probes (129 targeted checks), 274 PASS lines including subtests, with no failures or skips. All 22 semantic fields agree in all three public and two profile datasets. Every allocation sample matches: VP8 0, VP9 320 .616667, VP9 720 1.25, VP9 profile .625. No allocation outlier was discarded.

All 95 median regressions and 3 all-six-slower cases are retained in audit.json. The latter are unchanged controls:

| Control | Current time vs parent |
|---|---:|
| SSE8RowsNative/slice8x4/pitch32-37/pattern1 | +1.096690% |
| SSE8RowsNative/slice8x4/pitch32-37/pattern2 | +0.291375% |
| SSE8RowsNative/slice8x8/pitch1312-1320/pattern0 | +0.313761% |

Remaining old-reference gaps include legacy SSE8x8, 4x4 six-tap, quant batch2 and limited SAD. Source/backend scope matters: all 48 old-reference AMD subtraction cases use a scalar fallback, so their old-build comparisons are not assembly wins. All 326 unchanged kernels are controls. Short profile cumulative samples are not additive; the largest current VP8 sampled bodies include fused vertical filtering, coefficient counting/token writing and DCT batching.

## Optimized upstream libvpx

| Case | Go / libvpx time | Comparison scope |
|---|---:|---|
| vp9-320p | 2.913086× | Unequal measured quality/output work |
| vp9-720p | 2.941364× | Unequal measured quality/output work |
| vp8-720p-adaptive | 1.759965× | Unequal measured quality/output work |
| vp8-720p-fixed | 2.233204× | Measured PSNR/SSIM/payload/frame/drop counts match |

Libvpx is faster. Fixed-speed VP8 matching the measured quantities does not imply identical bitstreams. Every quality-pass output stat matches timed output in this job. All reference timers are vpxenc-stats; pinned libvpx v1.16.0 commit1024874c5919305883187e2953de8fcb4c3d7fa6 uses runtime detection and verified SSE2/AVX2/assembly objects. Preserve separate-host scope.

## Preservation

native-raw.tar.gz contains all 86 native files; libvpx-raw.tar.gz all 74 C files. No raw sample was dropped or rewritten. raw-sha256.json records all 160 member hashes; sha256.json covers the outer files. audit.json carries each fixture, pair, semantic field, allocation sample, regression, profile summary, source hash and workflow/artifact identity.

## Merge proof

Merged as `2905d810c8d075c283f259f7baa6f531c6ec4766`; tested and merged tree both `1f55c9097fee5cc3c58380b1dfee6599a5c16f16`. Remote main and expected parent were verified. PR41 was retargeted to this equal-tree main without rewriting its measured head before deleting only the exact merged PR40 remote source branch. No local checkout, index or branch changed. Full proof and independent root verification are embedded in audit.json; root-proof.json preserves the independent proof bytes.
