# Compute first-pass VP8 luma variance with fused SIMD stats

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkR4Components/Mean/dim16` | 10.404x | 9.906 / 10.903 |
| `BenchmarkR4Components/Mean/dim17` | 10.758x | 10.035 / 11.481 |
| `BenchmarkR4Components/Mean/dim32` | 11.250x | 11.277 / 11.224 |

3 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Fast path requires a complete positive-stride in-bounds16x16 source block using division bounds. Stats reads a valid16-byte zero reference row repeatedly with stride0, never writes inputs. Sum<=65280 and SSE<=16646400 fit original types, and the original +128 rounded int64 mean-square correction is preserved exactly. Edge/invalid windows retain the original clamped scalar path.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
