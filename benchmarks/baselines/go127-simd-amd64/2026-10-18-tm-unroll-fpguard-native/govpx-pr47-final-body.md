ARM true-motion prediction now calls one native Go SIMD implementation directly and unrolls disjoint row stores. The checked bool entry delegates to it; aliases and malformed windows retain sequential fallback and panic/partial-write behavior. The amd64 changes carry SSE8 row pointers and replace the FP16 overlap guard’s temporary array/nested loops with ten explicit interval checks. Arithmetic, optional outputs, fallback and upper-vector cleanup remain intact. No dependencies, handwritten assembly or README changes.

Native ARM evidence uses Apple M4 Max, official Go 1.27.2, GOMAXPROCS=1 and PGO off, versus the frozen #46 left-packet winner:

- TM8 improves 1.272–1.289× and TM16 1.409–1.471×, every target pair faster. These incremental gains are separate from earlier retained ≥1.5× wins.
- All nine TM16 medians beat old NEON by 1.031–1.085×, every pair faster. TM8’s median advantage is only 0.35–1.22%, with seven slower pairs retained.
- Whole VP8 720p takes +0.039% time versus parent with mixed pairs and **+9.457% versus old assembly**. The ARM codec gap remains open. All 22 semantic fields and zero allocations match across 18 runs; separate profiles are preserved. Later private TM4, loop-filter and DC experiments are not in this PR.

The completed full native x64 job uses EPYC 9V74, six balanced rounds and the same Go/GOMAXPROCS/PGO settings:

- The actual selected FP16 guard’s 12 disjoint component cases improve **1.853–2.641×** against the frozen parent guard; all 72 pairs exceed 1.8446×. These qualify as retained component wins. Ten overlap-decline controls remain diagnostic, and all alternatives/results are preserved.
- The 31 eligible public FP quantization fixtures, including existing class/no-class callers, improve 1.104–1.204× on medians; one slower pair remains. These public gains are below 1.5× and are reported separately from guard-only results.
- All nine public pointer-SSE fixtures improve 1.1888–1.1909× against parent, every one of 54 pairs faster.

| Full native workload | Candidate time vs parent | Faster parent pairs | Candidate time vs old assembly |
|---|---:|---:|---:|
| VP8 720p | +0.648664% | 4/6 | **+1.919956%, all six slower** |
| VP9 320p | −0.755836% | 5/6 | −44.943077%, all six faster |
| VP9 720p | −1.220876% | 5/6 | −55.654945%, all six faster |

All 498 three-way kernels × six rounds (8,964 samples) allocate zero. All 22 public/profile semantic fields match. All 458 controls are retained, including **240 slower medians and 42 cases slower in every pair**. Copy16 is unchanged; its documented self-copy regressions remain. Old Copy16 is compiler-generated Go, not handwritten assembly. No cause is assigned to unchanged-source timing drift.

Allocation vectors remain explicit: parent VP9 720p round one has 76/60 allocations versus 75/60 elsewhere; the old VP9 profile has 76/120 versus parent/current 75/120. Current whole-codec samples retain their expected allocations. These exceptions prevent a blanket allocation-parity claim.

A separate focused EPYC 7763 job measures pointer fixtures 1.271–1.336× faster than parent, all 54 pairs faster, minimum 1.124512×; medians beat old reference 1.353–1.386×. All 36 public medians beat old reference, but **all 54 slice8x16 pairs regress against parent**, with median time increases about 1.92–2.85%. Actual measured-binary GNU disassembly confirms 37→28 instructions in the pointer hot loop. This codegen observation does not establish the cause of every timing change. All 720 measured samples, 120 excluded oracles and four separate legacy B.N diagnostics remain preserved.

Pristine optimized libvpx v1.16.0 is pinned to `1024874c5919305883187e2953de8fcb4c3d7fa6`:

| Separate runner | VP8 fixed CPU −8 Go/C time | Adaptive VP8 Go/C time | VP9 320p / 720p Go/C time |
|---|---:|---:|---:|
| EPYC 9V74 | 2.122099× | 1.567871× | 2.653299× / 2.392773× |
| Apple M4 Max | 2.088134× | 1.577638× | 2.490255× / 2.482958× |

Ratios above one mean libvpx is faster. Fixed CPU-used −8 with min-q=2/max-q=56 is not fixed quantization. Fixed VP8 on both runners, and adaptive VP8 on this ARM run, match measured output statistics and separately decoded PSNR/SSIM with 58 quality frames. PSNR/SSIM were measured **only in separate quality passes**; each side’s quality-pass output statistics match every timed output. X64 adaptive VP8 and both VP9 comparisons retain unequal output/quality work. Matching statistics do not establish bitstream identity. CPU jobs and research stages are not pooled. The ARM C build’s different cache prefix is documented; compiles are excluded from throughput denominators.

Integrated full-suite, race/checkptr, scalar/purego/build-selection, PGO and pre-commit gates pass. Native evidence verifies 86 source hashes, 31 common fixtures, 338 excluded three-way oracles and 66 separate guard cases/396 timings/66 excluded oracles. All raw samples, rejected experiments, source/codegen proofs, mismatches and #46 evidence remain outside README files. Rosetta establishes correctness only. Root and the CI subagent independently audited the full native evidence. All four exact-head CI checks are green; performance acceptance retains the qualifying component win and the smaller assembly-beating kernels with every measured regression reported.
