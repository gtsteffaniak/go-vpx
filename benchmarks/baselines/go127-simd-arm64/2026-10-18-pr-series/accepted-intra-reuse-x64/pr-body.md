Accepted VP8 DC/V/H/TM macroblocks now add inverse residuals to their existing predictor, avoiding a second predictor build and decoder-token staging. Serial and row-worker paths match; B_PRED keeps its sequential reconstruction.

Preliminary Apple M4 Max / Go 1.27.2 SIMD, 720p60 realtime single-thread codec screen: fixed cpu=-8 **1.038× throughput** (paired 1.043×/1.033×), adaptive cpu=+8 **0.997×**. Accepted intra DC/skip components improve 1.91×/2.89×; mixed residual 1.29×. All losses and raw bindings are retained under `benchmarks/baselines/go127-simd-arm64/2026-10-18-pr-series/accepted-intra-reuse`.

Validation: focused full-plane legacy oracle, 96 VP8 packet/flag hashes, 72 VP9 ML packet hashes, affected threaded keyframe/bitstream and B_PRED checks; zero allocations. Adds optional affected-codec selection to the existing short native x64 workflow so VP8 fixed/adaptive uses two balanced 720p pairs. No full CI, C comparison, profile audit or README changes.


### Native x64 codec feedback

[Quick native run 39160666370](https://github.com/thesyncim/govpx/actions/runs/39160666370) compared this exact head with 7e0e37b7 on an AMD EPYC 7763 (AVX2 enabled, AVX512 unavailable), using VP8 720p, 60 inputs, 2500 kbps, one thread, PGO off and two balanced AB/BA rounds. Fixed `cpu-used=-8` improved throughput by **2.24% and 2.13%** (2.19% median); adaptive `+8` was mixed at +3.92% and −0.11%, so that result remains preliminary. All four pairs retained identical 22 output/semantic fields and zero allocations. The unrelated subpel control was flat/mixed. No C, profile, or full-matrix run was added.
