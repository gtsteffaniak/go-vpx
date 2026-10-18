# Check VP8 macroblock EOB skippability with word masks

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkR3Components/Skippable/pos-1` | 4.679x | 4.650 / 4.707 |
| `BenchmarkR3Components/Skippable/pos0` | 1.004x | 1.007 / 1.001 |
| `BenchmarkR3Components/Skippable/pos24` | 4.444x | 4.321 / 4.567 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Exact three uint64 loads stay within the EOB[25] array. Mask 0xfe per byte is equivalent to EOB>=2 for every uint8 value; UV/Y2 require zero, and BPred intentionally ignores block24 like the original. Nil returns false.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
