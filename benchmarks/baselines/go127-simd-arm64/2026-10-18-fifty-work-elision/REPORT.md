# Fifty short work-elision experiments

Fixed codec baseline: `c4a8bd9f66805fc9b7b4f99e26752cc68505e572`. Final code: `3f51e5c03058c1c4ae30059a4684cfc7f9b73d45` (PR55). Apple M4 Max, Darwin arm64, Go1.27.2, GOEXPERIMENT=simd, CGO1, GOMAXPROCS1, PGOoff. These are preliminary native screens, with competing agent builds and timers stopped during measurement.

There are exactly 50 distinct hypotheses: root IDs 1–25 and Astra Max IDs 26–50. ID 25 compares two zero-quantizer gating alternatives within one hypothesis; both lose their sparse/dense controls and are rejected. The previously verified fused-scoring seed is outside the 50 count. Its source, targeted evidence and original measurements are preserved separately. Final codec results include that seed and both sets of retained changes, so they cannot be attributed to the 50 experiments alone.

## Final codec comparison

Two balanced pairs per canonical workload; entries are baseline/current time ratios, hence throughput speedups. All 22 semantic fields and allocations agree in every pair. Quality is skipped in these encode-only runs; no fresh C, equivalent-quality or assembly victory claim follows from them.

| Workload | Pair 1 | Pair 2 | Median throughput gain | Allocations/frame |
|---|---:|---:|---:|---:|
| VP8 adaptive 720p, 60 frames | 1.003645x | 1.019048x | 1.135% | 0 /0 |
| VP9 320x180, 120 frames | 1.013455x | 1.016174x | 1.481% | 0.616667 /0.616667 |
| VP9 720p, 60 frames | 1.076713x | 1.072717x | 7.471% | 1.25 /1.25 |

The 1.5x fixed-baseline codec threshold is **not met**. Full suites, full matrices, new C comparisons and profile audits remain deferred. The completed x64/libvpx baseline stays in PR54; do not combine its EPYC ratios with these M4 gains. The 1.5x libvpx throughput target remains unmet.

## Component keepers

The final root-only composition received a fresh two-pair 80ms confirmation before merging the VP9 keeper set. Representative medians: local-predictor interior gather 16.657x, predicted 4x4 residual 1.732x, B_PRED residual 1.996x and SSE 1.713x; FDCT constant packets 2.645x and zero packets 1.496x. Single-block zero DCT 1.861x and Walsh 1.504x qualify while Walsh nonzero controls lose 5.8–7.0%; those losses remain in the raw results. Common FDCT controls are 0.990–1.001x. These numbers describe that frozen component composition, not the final whole-codec composition.

Smaller retained hypotheses include DCT conditional-mask subtraction(~3–4% in the initial intended cases), direct transform dispatch(8–13% in initial DCT pairs), and two-block shared quantizer tables(6–14% in intended initial cases). Later combined DCT nonzero controls are 0.983–0.992x, so no final nonzero-DCT gain is claimed. The initial zero/DC variants were refined into the endpoint-filtered constant path; superseded variants remain in the ledger.

Astra's component screens preserve TM bridge/fixed-row, H/DC fixed-row, required-edge-only gathering, redundant threshold removal, narrow rate accumulation, fused 8x8 scoring, uniform residual scoring and exact model endpoint wins. Stable focused confirmation measured redundant-mask removal 4–7%, narrow rate accumulation 13–18%, uniform scoring 1.79–2.43x and model endpoints 1.74–1.81x. Every unsupported-width/common-path loss is retained in the 25-row VP9 ledger; final wrappers select the winning sizes. Original variant measurements are explicitly separate from the final scoped composition.

## Correctness and scope

Final integration passed affected VP8 transform/residual, VP9 score/transform/predictor/model and root live-edge selectors. All 72 VP9 packets across three deterministic streams retain the fixed-baseline SHA256 hashes. The new VP8 checks cover 1050 packet overlap/pattern cases including constant±255, 100 edge/interior gather cases, 400 B_PRED block cases and 28 byte-backed output aliases. Existing full-signed transform/quantizer cases and prior seed aggregate boundaries remain saved. VP9 evidence includes 1440 live-edge backing comparisons, 2240 predictor backing/alias/panic comparisons and 73728 frozen-model states. Native CLI builds and affected AMD64 cross-builds passed; AMD64 performance remains for native CI, with no emulated timing.

All reported component allocations are zero. The root 40ms screens contain material timing variation in some unchanged controls; the raw pairs and subsequent 80ms confirmation are retained, and small effects remain preliminary. No broad suite, race rerun, C rebuild, profile collection or broad evidence audit was performed for this batch. No README, dependency or handwritten-assembly changes.

## Reproduction and raw evidence

`iterations50.tsv` is the combined 50-row decision ledger. `raw.tar.gz` contains the component/codec commands, all raw pairs including losses, targeted logs, fixed compiler/production bindings, per-variant source patches and shared fixtures/drivers. Root shared fixtures evolve in the preserved prepare scripts; each build phase uses identical fixtures on its parent/candidate. Existing binary/source bindings in Astra's proof are retained without a new audit. No executables or full repository trees are packed. Absolute paths in original logs are provenance; the archived directories preserve their file contents.

Run Go commands through `scripts/go.sh` from the Git project, with a distinct CODEX_AGENT_ID and the pinned SDK/experiment environment recorded in `baseline.json`. The source anchor and candidate patches reconstruct individual hypotheses. `whole-screen/summary.json` preserves exact final CLI arguments, all 22 semantic fields, allocation counts and both native pairs. `final-production-binding.json` ties the composition to the code commit above. The final production test fixture has descriptive names; archived initial fixture names remain unchanged so raw benchmark labels can be matched.
