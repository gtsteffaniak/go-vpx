# ARM48 diagnostic CPU profiles

Two existing 600-input VP8 runs compare the frozen ARM48 candidate and pinned old691 binaries, with identical public options:1280×720,fps30,2500kbps,realtimeCPU8,threads1,GOMAXPROCS1,CGO1,Go1.27.2archsimd,PGOoff. Both saved command records exit zero, all22 semantic fields match and both allocation samples are zero. Quality is skipped; no decoded-quality or bitstream-identity claim is made.

**Every timing value in these runs is excluded from performance acceptance.** These profiles support sampled-function liveness and within-run ranking only. No ratio of sample totals, function flat/cumulative samples, profile duration or ns/frame is used as a causal speedup estimate. The separate balanced300-input acceptance window and C comparisons are not pooled with these diagnostics.

The exact executable hashes/sizes match the integrated source/build proof. The candidate binds to published31123 through the preserved all1937Go source binding; old691 remains the pinned old implementation. No executable blobs or full snapshots are included. Raw profiles, original top reports, JSON/stderr, drivers, command timestamps/exit codes and source/build manifests are preserved. Existing profiles were decoded in pure Python; no Go, pprof tool, build, test or timing process was launched.

## candidate: highest flat-sample ranks

1. `github.com/thesyncim/govpx/internal/pixel.loopFilterMacroInternalFusedHorizontal`
2. `github.com/thesyncim/govpx/internal/pixel.loopFilterMacroInternalFusedVertical`
3. `github.com/thesyncim/govpx/internal/vp8/encoder.writePreparedCoefficientTokenRecords`
4. `github.com/thesyncim/govpx.(*VP8Encoder).selectFastInterFrameModeDecisionDenoise`
5. `github.com/thesyncim/govpx/internal/vp8/encoder.countBlockCoefficientTokensAndRecords`
6. `github.com/thesyncim/govpx/internal/pixel.QuantizeVP8Batch`
7. `github.com/thesyncim/govpx.(*VP8Encoder).buildReconstructingInterFrameCoefficientsWithSegmentation`
8. `github.com/thesyncim/govpx.buildPredictedMacroblockCoefficientsWork`

Observed target liveness (absence is not proof of no calls):

- `TMPredict4Native`: not observed in this saved sample.
- `DCPairInt32Native`: observed in sampled stacks.
- `loopFilterMacroInternalFusedHorizontal`: observed in sampled stacks.
- `loopFilterMacroInternalFusedVertical`: observed in sampled stacks.
- `LoopFilterMacroHorizontal16Unchecked`: observed in sampled stacks.
- `LoopFilterMacroVertical16Unchecked`: observed in sampled stacks.

## old691: highest flat-sample ranks

1. `github.com/thesyncim/govpx/internal/vp8/dsp.loopFilterVerticalEdgesYSharedNEON`
2. `github.com/thesyncim/govpx.(*VP8Encoder).selectFastInterFrameModeDecisionDenoise`
3. `github.com/thesyncim/govpx/internal/vp8/encoder.writePreparedCoefficientTokenRecords`
4. `runtime.memmove`
5. `github.com/thesyncim/govpx.buildPredictedMacroblockCoefficientsWork`
6. `github.com/thesyncim/govpx/internal/vp8/encoder.countBlockCoefficientTokensAndRecords`
7. `github.com/thesyncim/govpx/internal/vp8/dsp.sixTapPredict8x8NEON`
8. `github.com/thesyncim/govpx/internal/vp8/dsp.loopFilterHorizontalEdgesYSharedNEON`

Observed target liveness (absence is not proof of no calls):

- `TMPredict4Native`: not observed in this saved sample.
- `DCPairInt32Native`: not observed in this saved sample.
- `loopFilterMacroInternalFusedHorizontal`: not observed in this saved sample.
- `loopFilterMacroInternalFusedVertical`: not observed in this saved sample.
- `LoopFilterMacroHorizontal16Unchecked`: not observed in this saved sample.
- `LoopFilterMacroVertical16Unchecked`: not observed in this saved sample.

Raw sample values and complete function inventories remain in audit.json for review, without cross-profile quotient calculations. Sampling, inlining, corpus selection and short runs limit causal interpretation. All raw-member and outer checksums verify; no ROOT/README/Git changes occurred.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
