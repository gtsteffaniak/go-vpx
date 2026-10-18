# Hoist coefficient-band cost rows during rate scoring

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3KnownCost/tx0/dense` | 1.024x | 1.022 / 1.027 |
| `BenchmarkRound3KnownCost/tx0/sparse` | 1.023x | 1.020 / 1.027 |
| `BenchmarkRound3KnownCost/tx0/zero` | 0.852x | 0.849 / 0.856 |
| `BenchmarkRound3KnownCost/tx1/dense` | 1.377x | 1.367 / 1.387 |
| `BenchmarkRound3KnownCost/tx1/sparse` | 1.128x | 1.134 / 1.122 |
| `BenchmarkRound3KnownCost/tx1/zero` | 0.851x | 0.846 / 0.857 |
| `BenchmarkRound3KnownCost/tx2/dense` | 1.302x | 1.461 / 1.143 |
| `BenchmarkRound3KnownCost/tx2/sparse` | 1.132x | 1.124 / 1.141 |
| `BenchmarkRound3KnownCost/tx2/zero` | 0.848x | 0.848 / 0.847 |
| `BenchmarkRound3KnownCost/tx3/dense` | 1.551x | 1.549 / 1.554 |
| `BenchmarkRound3KnownCost/tx3/sparse` | 1.146x | 1.139 / 1.153 |
| `BenchmarkRound3KnownCost/tx3/zero` | 0.846x | 0.889 / 0.803 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
