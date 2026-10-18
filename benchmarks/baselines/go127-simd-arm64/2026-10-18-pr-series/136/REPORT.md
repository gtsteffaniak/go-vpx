# Skip probability-search setup for unchanged seeds

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3SearchSeed/modelfalse/alignedfalse` | 1.005x | 0.981 / 1.029 |
| `BenchmarkRound3SearchSeed/modelfalse/alignedtrue` | 1.732x | 1.738 / 1.725 |
| `BenchmarkRound3SearchSeed/modeltrue/alignedfalse` | 1.007x | 1.000 / 1.013 |
| `BenchmarkRound3SearchSeed/modeltrue/alignedtrue` | 6.414x | 6.470 / 6.359 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
