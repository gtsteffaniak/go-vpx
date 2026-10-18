# Rejected ARM horizontal edge unrolling

Apple M4 Max, official Go1.27.2, GOEXPERIMENT=simd, PGO off, GOMAXPROCS=1. Exact parent46b89f5ab4c6c469755323a61c77fe3a766b4081. Only the ARM horizontal macro-plus-three-internal-edge loop is unrolled; arithmetic, loads, stores, corrected-neighbor carries and the separate ARM vertical / all AMD source remain unchanged. Generator reproduction is included. This is distinct from earlier pre-gate, byte-shift, abs and packet-cache experiments.

Three diagnostic parent/candidate pairs at75ms use the unchanged actual public VP8 picker fixture: two pitches times five restored input patterns, live output checked against separate original public edge operations and full-buffer equality. Orders AB/BA/AB are a small exploratory screen, not six balanced acceptance rounds. All60 timed samples are zero bytes/allocations. Pixel macro/internal differential, declined-write and zero-correction oracles, decoder loop-filter tests and race/checkptr2 pass. The7200 macro/internal cases compare the independent existing widened-eight-lane network; the public fixture compares separate edge operations. Neither is claimed to be a new standalone C oracle.

LLVM disassembly changes364 to557 entries, no calls or vector stack operations. Final source and binary hashes are recorded; these static counts are not timing evidence.

All10 public medians improve only1.016826–1.034839×; nine cases are faster in all3 pairs. One pitch1312/mixed pair regresses (ratio0.979700). No component reaches1.5×. Given the modest screen signal and larger code, this candidate is rejected; no fresh-six or whole-codec run follows and no production branch is changed. Every timed sample, paired result, source, generator, one-shot exclusion, build/check log and disassembly is retained here. Compiled binaries remain private; their SHA256 values are preserved.

| case | parent ns | candidate ns | speedup | all3 paired parent/current |
|---|---:|---:|---:|---|
| pitch32/flat | 46.080 | 44.780 | 1.029031× | 1.025824, 1.048293, 1.015409 |
| pitch32/ramp | 46.170 | 45.050 | 1.024861× | 1.055050, 1.015316, 1.026228 |
| pitch32/smooth | 89.600 | 86.790 | 1.032377× | 1.041860, 1.008296, 1.044111 |
| pitch32/mixed | 89.320 | 87.610 | 1.019518× | 1.020990, 1.007011, 1.019518 |
| pitch32/noop | 40.600 | 39.550 | 1.026549× | 1.026043, 1.018565, 1.045687 |
| pitch1312/flat | 46.040 | 44.490 | 1.034839× | 1.040452, 1.000220, 1.035514 |
| pitch1312/ramp | 45.760 | 44.790 | 1.021657× | 1.007993, 1.034203, 1.021657 |
| pitch1312/smooth | 89.110 | 87.310 | 1.020616× | 1.025765, 1.008705, 1.016309 |
| pitch1312/mixed | 88.950 | 87.440 | 1.017269× | 1.017269, 0.979700, 1.041308 |
| pitch1312/noop | 40.490 | 39.820 | 1.016826× | 1.013836, 1.014025, 1.029131 |


The complete frozen evidence package is inside `raw-evidence.tar.gz` (36 original files). It preserves all nested raw archives, manifests, reports and auditors byte for byte. Extract to another directory to inspect it; `SHA256SUMS` covers the published outer files.
