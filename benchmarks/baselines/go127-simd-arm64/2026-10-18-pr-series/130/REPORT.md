# Average pixels directly in the byte domain

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3Average/w16/pad0` | 1.186x | 1.183 / 1.190 |
| `BenchmarkRound3Average/w16/pad7` | 1.186x | 1.187 / 1.184 |
| `BenchmarkRound3Average/w32/pad0` | 1.221x | 1.220 / 1.223 |
| `BenchmarkRound3Average/w32/pad7` | 1.231x | 1.235 / 1.228 |
| `BenchmarkRound3Average/w4/pad0` | 1.056x | 1.052 / 1.060 |
| `BenchmarkRound3Average/w4/pad7` | 1.041x | 1.053 / 1.028 |
| `BenchmarkRound3Average/w64/pad0` | 1.277x | 1.277 / 1.276 |
| `BenchmarkRound3Average/w64/pad7` | 1.240x | 1.254 / 1.226 |
| `BenchmarkRound3Average/w8/pad0` | 1.118x | 1.116 / 1.120 |
| `BenchmarkRound3Average/w8/pad7` | 1.119x | 1.117 / 1.122 |

6 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
