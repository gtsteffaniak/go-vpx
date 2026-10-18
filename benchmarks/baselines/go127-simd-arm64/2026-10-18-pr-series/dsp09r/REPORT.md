# Use shared direct directional averages with snapshot alias fallback

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkDSPPRDirectional/mode2` | 1.107x | 1.101 / 1.112 |
| `BenchmarkDSPPRDirectional/mode3` | 1.305x | 1.293 / 1.317 |
| `BenchmarkDSPPRDirectional/mode4` | 1.072x | 1.068 / 1.075 |
| `BenchmarkDSPPRDirectional/mode5` | 1.142x | 1.129 / 1.155 |
| `BenchmarkDSPPRDirectional/mode6` | 1.212x | 1.214 / 1.209 |
| `BenchmarkDSPPRDirectional/mode7` | 1.128x | 1.147 / 1.109 |
| `BenchmarkDSPPRDirectional/mode8` | 1.205x | 1.196 / 1.215 |
| `BenchmarkDSPPRDirectional/mode9` | 1.240x | 1.244 / 1.237 |

4 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
