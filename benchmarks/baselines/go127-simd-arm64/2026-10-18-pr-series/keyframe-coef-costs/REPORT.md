# Reuse fixed keyframe coefficient costs in VP8 RD scoring

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkVP8KeyframeCoefCosts/0` | 1.143x | 1.136 / 1.149 |
| `BenchmarkVP8KeyframeCoefCosts/1` | 3.468x | 3.483 / 3.452 |
| `BenchmarkVP8KeyframeCoefCosts/2` | 4.933x | 4.902 / 4.965 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The fixed keyframe coefficient probability pointer uses one prepared table. Dynamic inter probabilities still use their existing table or original scoring. The checked table scorer preserves invalid type/context/skipDC/token/eob behavior; direct independent coefficient scoring agrees for zero, dense and extreme coefficients at every valid and invalid boundary.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
