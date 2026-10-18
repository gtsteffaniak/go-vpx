# Reuse existing token costs for large coefficient magnitudes

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. AMD EPYC 7763 / Linux amd64 (AVX2 native), Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkRound3ExtraCost/0` | 0.872x | 0.871 / 0.873 |
| `BenchmarkRound3ExtraCost/1` | 0.878x | 0.877 / 0.879 |
| `BenchmarkRound3ExtraCost/17` | 0.896x | 0.876 / 0.916 |
| `BenchmarkRound3ExtraCost/32767` | 10.419x | 10.494 / 10.344 |
| `BenchmarkRound3ExtraCost/32768` | 10.549x | 10.596 / 10.502 |
| `BenchmarkRound3ExtraCost/32769` | 0.979x | 0.967 / 0.990 |
| `BenchmarkRound3ExtraCost/4095` | 0.877x | 0.879 / 0.874 |
| `BenchmarkRound3ExtraCost/4096` | 12.069x | 10.667 / 13.471 |
| `BenchmarkRound3ExtraCost/8192` | 10.587x | 10.596 / 10.578 |

4 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream passed all 72 VP9 packet hashes. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Independent scalar/frozen-reference fixtures cover the affected output and state; see the focused correctness log.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
