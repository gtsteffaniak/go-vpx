# PR series 01: fixed VP8 keyframe mode costs

Baseline: `870699ce83f3eb6de44f507ef29f6d6a429e916d`, tree-equivalent to main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2, GOEXPERIMENT=simd, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples, zero allocations.

The live mode-cost functions now index costs prepared once from VP8's immutable keyframe probabilities. Dynamic inter probabilities and invalid-mode behavior retain the existing tree walk. B-mode / Y-mode / UV-mode throughput ratios: **2.679x / 2.861x / 2.949x**. This qualifies under the requested 15% kernel threshold; no end-to-end speedup is claimed.

Independent direct tree-cost comparisons cover all 1000 B-mode contexts, Y/UV leaves and invalid tokens. Affected picker frozen-reference checks, the 96-packet VP8 stream and root ML stream passed. Raw timings and the focused correctness log are retained here. Full CI, suites, new C comparisons and broad audits remain deferred under the fixed-c4 codec gate. README unchanged.
