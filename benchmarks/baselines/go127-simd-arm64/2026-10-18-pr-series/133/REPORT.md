# Use byte arithmetic for x64 true-motion prediction

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3TM/w16/top0` | 6.870x | 6.896 / 6.843 |
| `BenchmarkRound3TM/w16/top128` | 6.880x | 6.877 / 6.884 |
| `BenchmarkRound3TM/w16/top255` | 6.909x | 6.892 / 6.926 |
| `BenchmarkRound3TM/w32/top0` | 1.019x | 0.999 / 1.039 |
| `BenchmarkRound3TM/w32/top128` | 1.039x | 1.040 / 1.038 |
| `BenchmarkRound3TM/w32/top255` | 1.040x | 1.046 / 1.035 |
| `BenchmarkRound3TM/w4/top0` | 0.986x | 0.988 / 0.983 |
| `BenchmarkRound3TM/w4/top128` | 0.966x | 0.950 / 0.982 |
| `BenchmarkRound3TM/w4/top255` | 0.982x | 0.983 / 0.980 |
| `BenchmarkRound3TM/w8/top0` | 2.824x | 2.788 / 2.860 |
| `BenchmarkRound3TM/w8/top128` | 2.878x | 2.888 / 2.868 |
| `BenchmarkRound3TM/w8/top255` | 2.855x | 2.851 / 2.859 |

6 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
