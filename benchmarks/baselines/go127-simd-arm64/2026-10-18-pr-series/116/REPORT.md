# Evaluate VP8 coefficient token costs without path loops

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkR3Components/DirectCost/token1` | 1.508x | 1.511 / 1.505 |
| `BenchmarkR3Components/DirectCost/token4` | 2.015x | 2.072 / 1.958 |
| `BenchmarkR3Components/DirectCost/token10` | 2.356x | 2.269 / 2.442 |

3 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Token validity is checked before fixed-path dispatch. The direct expressions use the original CoefTree edge order and same probability indices; elision and sentinel behavior are unchanged for all block/band/context combinations.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
