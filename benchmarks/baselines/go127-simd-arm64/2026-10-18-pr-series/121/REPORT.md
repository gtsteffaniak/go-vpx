# Skip empty VP8 coefficient EOB scans with word tests

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkR4Components/EOB/pos-1` | 3.375x | 3.218 / 3.531 |
| `BenchmarkR4Components/EOB/pos0` | 0.709x | 0.706 / 0.712 |
| `BenchmarkR4Components/EOB/pos15` | 0.992x | 0.985 / 0.999 |

1 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

For valid fixed 16-coefficient arrays, four exact word reads prove only the all-zero identity. Every nonzero block retains the original full zigzag scan; raster halves are not mistaken for scan halves. skipDC outside0/1 keeps the original fallback. No writes or coefficient/order changes. Early-nonzero control loss remains explicit.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
