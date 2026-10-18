# Write only final VP8 coefficient token contexts

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkR3Components/Contexts/is4false` | 3.540x | 3.703 / 3.377 |
| `BenchmarkR3Components/Contexts/is4true` | 3.600x | 3.966 / 3.235 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Fast path requires coefficient, above and left objects mutually disjoint using ordered address differences; all overlaps use the frozen original walk. Terminal EOB entries produce exactly the final bottom-row/top-column contexts; Y2 preservation in is4x4 remains unchanged. Alias/full-backing test passes.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
