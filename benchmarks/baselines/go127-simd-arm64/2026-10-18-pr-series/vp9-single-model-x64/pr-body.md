The VP9 non-RD picker ran the small-block RD model before the large-block model, then overwrote every small-model result whenever the large model was valid. Run the small model only as the fallback, in both filter sweep and ordinary candidate scoring. Prediction/variance and large-model evaluation remain unchanged.

Preliminary M4 Max/Go1.27.2 SIMD single-thread realtime CPU8,720p60 throughput: paired1.03734×/1.00740× (median+2.24%). All output fields match and codec allocations remain unchanged at1.25/frame. The320p slow-parent outlier is retained without a gain claim. Raw results, source and binary bindings are under benchmarks/baselines/go127-simd-arm64/2026-10-18-vp9-single-rd-model.

Validation: prepared-predictor differential oracle plus unchanged96 VP8 and72 VP9 ML packet hashes. Native x64 affected-codec results are recorded below. No full CI, C comparison, profile audit or README changes.

### Native x64 codec feedback

[Quick run 39161954456](https://github.com/thesyncim/govpx/actions/runs/39161954456) compared this head with integrated costgate baseline e9d4d08b on AMD EPYC 7763 (actual AVX2 enabled, AVX512 unavailable). Two balanced VP9 720p/60-input CPU8 pairs improved throughput by **0.58% and 0.79%** (encode time −0.57%/−0.78%). All 22 output fields matched; codec allocations stayed unchanged at **1.25/frame (75/60)**. The existing frame-context setup benchmark was an unchanged control (.997–1.006×, zero allocations), not an optimization claim. Prepared-predictor, cost-transition, 96 VP8 and 72 ML checks passed on both variants. These two pairs remain preliminary; no C, profile, or full CI was added.
