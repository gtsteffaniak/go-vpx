# PR series 02: fast keyframe predictor preparation

Baseline: `870699ce83f3eb6de44f507ef29f6d6a429e916d`, same tree as main `98aa69f178b8f3213c25dc1b1475f57b30f9621e`. Apple M4 Max / Darwin arm64, Go 1.27.2 SIMD, PGO off, CGO enabled, GOMAXPROCS=1. Preliminary two balanced 40 ms samples, zero allocations.

Fast Y/UV pickers prepare invariant predictor edges once per candidate set. The Y scorer predicts only Y, avoiding four discarded UV predictions. Fast Y throughput: **1.949x / 1.883x** for flat/dense source patterns; fast UV: **1.504x / 1.525x**. These are component results qualifying under the requested 15% kernel threshold; no codec speedup is claimed.

Frozen-parent selected modes/rates/distortions and relevant reconstruction state match at all macroblock borders. The omitted Y-scorer UV writes are discarded by the encoder; the 96-packet VP8 oracle and root ML stream passed unchanged. Focused correctness and all raw samples are retained. Full CI and broad validation remain deferred. README unchanged.
