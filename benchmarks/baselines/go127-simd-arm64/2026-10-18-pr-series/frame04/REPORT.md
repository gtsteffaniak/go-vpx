# Reuse immutable chroma classifications instead of repeating Gaussian arithmetic

Baseline `870699ce83f3eb6de44f507ef29f6d6a429e916d` has the same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples; zero allocations. Kernel/component scope only; no whole-codec speedup is claimed.

| Live benchmark scope | Median throughput ratio | Two paired ratios |
|---|---:|---:|
| `BenchmarkFramePR04/skin` | 3.982x | 3.967 / 3.998 |
| `BenchmarkFramePR04/mixed` | 1.645x | 1.650 / 1.641 |
| `BenchmarkFramePR04/gray` | 1.022x | 1.017 / 1.027 |

2 measured cases exceed the requested 15% kernel gain in both pairs. All measured controls and losses are preserved above and in the raw logs. The affected scalar/frozen-parent correctness check passed. The root ML stream and 96-packet VP8 oracle match the unchanged baseline. Source and binary provenance is retained in binding/summary metadata; fixtures may be renamed per PR without changing their measured work.

A 64 KiB immutable chroma classifier is built once at package initialization. All 65,536 byte chroma pairs match the frozen classifier at luma/motion boundaries; nonbyte inputs retain the original arithmetic. The steady-state benchmark excludes one-time initialization and the extra static table footprint.

Full CI, full suites, C comparisons, profiles and broad audits remain deferred under the fixed-c4 codec threshold. README files unchanged.
