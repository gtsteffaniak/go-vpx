# Stop first-pass motion search after the earliest exact match

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3FirstPass/dense/16` | 1.000x | 1.009 / 0.992 |
| `BenchmarkRound3FirstPass/dense/8` | 0.999x | 0.996 / 1.001 |
| `BenchmarkRound3FirstPass/flat/16` | 69.034x | 68.586 / 69.482 |
| `BenchmarkRound3FirstPass/flat/8` | 64.182x | 63.049 / 65.316 |
| `BenchmarkRound3FirstPass/same/16` | 1.985x | 1.996 / 1.973 |
| `BenchmarkRound3FirstPass/same/8` | 1.959x | 1.950 / 1.968 |
| `BenchmarkRound3FirstPass/shift/16` | 1.581x | 1.574 / 1.588 |
| `BenchmarkRound3FirstPass/shift/8` | 1.577x | 1.572 / 1.583 |

6 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
