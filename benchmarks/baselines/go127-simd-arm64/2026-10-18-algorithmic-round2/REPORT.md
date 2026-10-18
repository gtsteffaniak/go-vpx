# Algorithmic work elimination, second 50 iterations

Production head: `7cc9a0f4dace308315bdf19121f28ee912cfbbc5`. Batch parent: latest main `b57aac772af29e85f4b5b34ae33600cc84cd6a0f`. Fixed codec gate: `c4a8bd9f66805fc9b7b4f99e26752cc68505e572`.

50 concrete experiments: 35 retained (17 ARM VP8, 18 x64/component changes), 15 rejected. `iterations.tsv` includes every decision and control loss. All screens are preliminary. Go 1.27.2, SIMD enabled, PGO off, CGO enabled, GOMAXPROCS=1. No dependencies or handwritten assembly added; README files unchanged.

## Retained work reductions

VP8 reuses frame-invariant mode costs, immutable predictor edges and already-computed residual statistics; avoids a full predictor copy/rescan; reconstructs the RD winner once; writes selected reconstruction directly; reuses DCT storage; scores Y without discarded UV predictions; and initializes only selected quantizers, copying identical segment tables.

VP9 consumes aggregate TX4 scores without materialized q/dq/scan/EOB work. Uniform nonzero TX8/16 residuals can avoid residual/transform passes; scratch aliases retain the original write/read ordering. SSE-only consumers, prepared border rows, contiguous plane reductions, long normalized/half-pel strips, high-limit SAD finishing, projection, copy and exact identity-filter exits retain their measured scopes. Low-limit SAD, short strips and 16-wide SSE retain their original arithmetic routes.

| Native component scope | Representative throughput ratios versus parent |
|---|---:|
| M4 Max VP8 fast B_PRED picker | 1.30–1.36x |
| M4 Max VP8 fast keyframe / full RD picker | 1.12–1.15x / about 1.07x |
| M4 Max combined segment quantizer setup | 4.75x |
| M4 Max TX4 dense scoring, sizes 16/64 | 2.13–2.18x |
| M4 Max flat8 identity / prepared border statistics | 1.59–1.87x / 1.71–5.22x |
| EPYC7763 combined TX4 dense scoring, sizes 16/64 | 2.10–2.20x |
| EPYC7763 combined constant TX4/8/16 scoring | 1.81–3.83x |
| EPYC7763 plane SSE / SSE-only blocks / vector variance | 1.25–3.72x / 1.14–3.36x / 1.34–4.13x |

Do not treat these as codec ratios. Preserve control losses: the final x64 screen includes 40 median losses, including a short subpel control at 0.838x, TX8 dense at 0.933x, and a TX16 zero case at 0.936x. Earlier experiments also retain their larger losses and failed ISA attempt. ID86 emitted an unsupported AVX512 instruction on EPYC7763 and was rejected; that code is absent from production.

## Short codec results

All scopes use the root synthetic realtime cpu8 workload, threads 1, encode-only, automatic C comparison disabled. 320 means 320x180/120frames/500kbps; 720 means 1280x720/60frames/2500kbps. VP8 uses adaptive rate control. All 22 nonperformance output fields and allocation rates match.

| Host / reference | VP8 720 throughput | VP9 320 throughput | VP9 720 throughput |
|---|---:|---:|---:|
| M4 Max / batch parent, two paired screens | +0.024% | −1.90% | +0.030% |
| M4 Max / fixed c4, two paired screens | +0.43% | +2.52% | +8.96% |
| EPYC7763 / batch parent, one pair per VP9 scope | not sampled | +0.532% | +0.493% |

One additional compact M4 320p repeat preserved all four ratios: 0.96419,0.99784,0.99048,0.97126; their median is 0.98087. The user directed that the component wins remain and this short variation be treated as noise, with no further probing. ARM VP8 remains flat; the previously identified ARM gap is unresolved.

The fixed-c4 raw runner retains its generic `b57` reference label; `whole-fixed-c4/binding.json` points to the reused immutable c4 binary. Reference interpretation follows the scope, not that generic label. The 1.5x fixed-baseline codec threshold is not met. Full suites, full matrices, new C comparisons, profile audits and broad evidence audits remain deferred. Reuse the completed [PR54 libvpx baseline](https://github.com/thesyncim/govpx/pull/54); do not extrapolate its Go/C ratios across these different hosts.

## Correctness and provenance

Focused scalar/frozen-reference checks passed, including quantizer boundaries, edge clamping, alias ordering, pointer windows, wrapped subpel coefficients, SIMD selection and identity-filter boundaries. The combined x64 run passed 28 ordinary tests per reference, 576 zero-allocation timing samples and both 22-field codec comparisons. All 96 VP8 packets/frame flags match the parent; all 72 VP9 packets retain the fixed-c4 stream hashes. Affected amd64 cross-builds and generator replay passed. No emulated timing was used.

Native runs: [first cohort, EPYC9V45 with AVX512](https://github.com/thesyncim/govpx/actions/runs/39146736111), [second cohort, EPYC7763 AVX2](https://github.com/thesyncim/govpx/actions/runs/39147718292), [combined actual production head, EPYC7763 AVX2](https://github.com/thesyncim/govpx/actions/runs/39149274752), [affected CI](https://github.com/thesyncim/govpx/actions/runs/39149251823).

`raw.tar.gz` contains the compact ledger inputs, independent source patches, fixtures, raw samples, failed attempts, commands and source/binary bindings. The second cohort's 94/100 baseline fixture incorrectly expected a wrapped Y-coordinate to panic; this was corrected and both final ordinary checks passed. The original failed fixture is retained in the source overlay archive. `x64-source-overlays.tar.gz` and the two archive-backed manifests reproduce the independent overlays without adding experiment packages to normal Go discovery. `final-scope.json` identifies the tested production tree; later evidence-only commits reuse these unchanged checks.
