# Final ARM50 integrated evidence audit

Independent audit PASS. Apple M4 Max, official Go 1.27.2 SIMD, CGO 1, PGOoff, one runtime thread. No new performance timers were run by this audit. ROOT and all earlier snapshots remain unchanged. Every bound Go source set and executable hash was verified. Binaries are omitted with exact private path, SHA-256 and byte-size receipts. All slower samples and outliers remain present.

The separate generated-LF confirmation contains 360 samples: 20 public axis fixtures and 10 actual decoder picker fixtures, six alternating pairs at 100 ms. All samples pass their pre/post output oracle with zero allocations. Horizontal controls have median speedups 0.987686–1.065874×, 38/60 pairs faster and 22 slower. Vertical medians 1.004791–1.080118×, 50/60 faster, 9 slower, 1 tie. Actual picker medians 1.030306–1.068141×, 53/60 faster and 7 slower; worst paired ratio 0.786860 remains included. Ratios compare to the exact strong SAT parent of that window, not old assembly.

The final EOB confirmation contains 180 samples from 15 actual prepared-token fixtures with independent public output/coder-state/full-buffer checks before and after B.Loop, all zero allocation. The all-zero-block stream using default probabilities is 1.593549×LF-only, all 6 pairs at least 1.568383×. The all-zero p230 median is 1.482457×, not 1.5; minimum pair 1.465379. Sparse default/p230 medians 1.234578/1.232360×. Dense and p128 controls are mixed, with median regressions down to 0.987559× and worst dense pair 0.911353. p0 denotes default probabilities. This component result must not be generalized to sparse/dense token streams or whole-codec throughput. Mandatory accepted DC component wins remain retained in the separately frozen saturation stage; none of their earlier results are pooled with this final window.

The matched whole-codec window is 32 runs: SAT parent, generated LF-only, LF+EOB candidate and old691, eight balanced rotations/reverse orders. Input is VP8 realtime 1280×720, 300 frames, fps 30, 2500kbps, cpu-used 8, one thread. All 22 non-timing/non-allocation fields match across every run and allocations are zero. Quality is skipped, so matching zero-valued quality fields establish neither measured PSNR/SSIM nor bitstream identity.

| Variant | Median ms/input |
|---|---:|
| Strong SAT parent |7.060871|
| LF-only |7.0044455|
| LF+EOB |7.0234435|
| old691 |6.6984715|

LF-only is 0.799129% faster than SAT (8/8 pairs faster). LF+EOB is 0.530069% faster than SAT (8/8), but 0.271228% slower than LF-only (5/8 slower). The combined candidate remains 4.851435% slower than old691 (5/8 slower); the three pairs faster than old are retained along with the five losses. LF-only remains 4.567818% slower than old. These are ratios of medians; all individual paired ratios and median paired percentages are separately recorded. No gap is derived by subtracting results from different windows, no cross-CPU inference is made, and the whole ARM assembly gap remains open.

The LF360 source map predates the final architecture-specific Stats type-alias repair. The final whole32 and EOB180 use the separately bound alias-rebased source. Exact source maps/deltas preserve that distinction; LF machine-source bodies are identical between those snapshots, but the full source sets are not pooled or called byte-identical.

Generator history is explicit: generator.log and generator-attempt1.log preserve the first column-assignment assertion failure; generator-attempt2.log is empty. An independent new scratch replay of the unmodified final generator exits 0 and reproduces all 449 copied/replayed Go files exactly, including ARM/AMD macro bodies, old rowstore and new rowvalues helper. No frozen snapshot was regenerated in place. Source inputs, expected hashes and reproduction instructions are included.

The saved full race/checkptr=2 run fails eight zero-allocation assertions. All eight were independently rerun on the exact final strong SAT parent under the same flags, packages and selectors and reproduce with the same values: DSP BilinearPredict16x16, SixTapPredict16x16 and InverseWalsh each 1; decoder ReconstructInterFrameGrid 3, PredictIntra 16, ReconstructBPredIntraMacroblock 16, TransformMacroblockTokens 1 and ParseStateHeaderWithReaderAndProbs 1. These are observed instrumentation baseline limitations, not a passing full instrumented suite. Changed-path scoped race/checkptr, purego/nosimd and ordinary complete package tests pass. No threshold or test was weakened.

Complete earlier packet/saturation/rejected-values/EOB-review/Stats-alias histories remain in their immutable packages, referenced by exact outer hashes with their reports/manifests copied here. This final lean package contains all final360+180+32 raw outputs, every candidate production delta, generator/proofs/failures, build commands, fixtures, baseline-failure rerun and binary/source receipts. It does not copy entire source snapshots or executables.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
