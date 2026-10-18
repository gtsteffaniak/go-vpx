# Batch empty VP8 chroma token records and counts

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkR4Components/Count/zero` | 1.659x | 1.691 / 1.626 |
| `BenchmarkR4Components/Count/lumaOnly` | 1.130x | 1.138 / 1.123 |
| `BenchmarkR4Components/Count/denseSmall` | 0.989x | 0.985 / 0.993 |
| `BenchmarkR4Components/Count/denseMixed` | 0.988x | 0.984 / 0.992 |

1 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

Eligibility reads only validated full context/EOB arrays. Fast path requires zero UV contexts, all eight EOBs zero and eight available records. Payload32B/header24B/counter8B are mutually disjoint and disjoint from the complete coefficient/above/left objects by the existing ordered-difference alias proof. It emits identical eight records, increments the same counter modulo int and omits only zero stores; every decline uses the original sequential loop. ARM-only true path; other builds use a false stub. Full macro partial-error/capacity/context state check passes.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
