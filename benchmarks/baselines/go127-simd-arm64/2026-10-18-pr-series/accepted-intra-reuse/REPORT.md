# Accepted VP8 whole-intra predictor reuse — preliminary short screen

Parent: `7e0e37b7bbbcda65523f5581e7879c7d2d93827c` (latest main when this pass started). Apple M4 Max, Darwin arm64, Go 1.27.2, GOEXPERIMENT=simd, CGO_ENABLED=1, PGO off, GOMAXPROCS=1 for native timings. Root scripts/go.sh supplies the shared project build cache.

Accepted DC/V/H/TM macroblocks already hold the predictor before forward transform and quantization. Add inverse residual directly to those pixels instead of staging decoder tokens and building the same predictor again. Mirror serial and row-worker paths. B_PRED retains its sequential dependency order and drops only unused token staging.

Codec: built-in deterministic 1280x720, 60 frames, 30 fps, 2500 kbit/s, realtime, threads=1, default noise sensitivity, encode-only. Two balanced A/B then B/A pairs; short screens are preliminary.

| Scope | Paired throughput ratios | Aggregate ratio |
|---|---|---|
| VP8 fixed cpu=-8 | 1.04333, 1.03330 | 1.03831x |
| VP8 adaptive cpu=+8 | 1.00175, 0.99185 | 0.99680x |
| Accepted intra DC component | 1.91857, 1.90275 | 1.91066x |
| Accepted intra mixed residual component | 1.32688, 1.24808 | 1.28748x |
| Accepted intra skipped component | 2.79406, 2.99267 | 2.89337x |

All codec non-performance fields match; allocation counts remain zero. Component and codec results are separate. The component parent uses a test-only legacy adapter preserved alongside its source; the whole-codec parent is unchanged main.

Correctness: 768 focused full-plane legacy comparisons across Y/UV modes, borders/interior, quantizers, and mixed Y2/EOB, including coefficient immutability. Existing 96 VP8 packet/flag hashes and 72 VP9 ML packet hashes match. Final integrated tree also passes affected threaded-keyframe, identical-thread-bitstream and B_PRED reconstruction tests. No full suite, profile audit, new C comparison, or full CI.

`arm-short` preserves all samples, commands and binary SHA256 bindings. `manifest.json` binds production files. The short native x64 dispatch measures the same affected VP8 workloads; its result will be retained separately.

Rejected probes are retained: rerouting full encoder loop filtering through fused luma plus chroma was flat/slightly slower after a slow first parent outlier; the common ONE token writer specialization was mostly 1–3% at component level, below retention threshold. Neither is included in production. No README changes.
