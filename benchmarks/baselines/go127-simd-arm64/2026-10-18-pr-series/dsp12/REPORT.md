# Avoid unused signed-sum work in four-wide SSE reductions

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkDSPPRSSE4/h4/stride4` | 1.259x | 1.262 / 1.256 |
| `BenchmarkDSPPRSSE4/h4/stride24` | 1.255x | 1.244 / 1.266 |
| `BenchmarkDSPPRSSE4/h8/stride4` | 1.223x | 1.227 / 1.219 |
| `BenchmarkDSPPRSSE4/h8/stride24` | 1.232x | 1.225 / 1.238 |

4 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
