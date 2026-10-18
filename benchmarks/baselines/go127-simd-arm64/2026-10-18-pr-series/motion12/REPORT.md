# Fill single-block split labels without scanning 16 blocks

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkMotionPR12/fill0` | 1.017x | 0.996 / 1.038 |
| `BenchmarkMotionPR12/fill1` | 0.887x | 0.994 / 0.781 |
| `BenchmarkMotionPR12/fill2` | 1.151x | 1.201 / 1.100 |
| `BenchmarkMotionPR12/fill3` | 2.171x | 2.137 / 2.204 |
| `BenchmarkMotionPR12/score4x4` | 1.004x | 1.004 / 1.004 |

1 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

The selected direct block fill applies only to partition 3. Other partition fills and score traversal keep the original algorithm. The earlier unscoped change lost 25% on fill1 and was narrowed; that original screen is preserved in unscoped-control. The final short screen still has an 11% fill1 control loss despite unchanged arithmetic, and retains the 2.17x selected fill3 gain under the kernel retention rule.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
