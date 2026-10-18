# Skip unused VP8 denoiser chroma prediction

Preliminary native component screen against `715645fd212cfc65295f4f6fb1002dcc60856758`. Published onto merged main499ef24e as runtime778614b3; production bytes match the measured independent715 overlay. Only test names changed to execute through the existing fast CI prefix. The combined main96 VP8/72 ML streams, frozen denoiser and focused predictor helper checks pass.

## Removed work

`vp8_encoder_denoiser.go`, `applyDenoiserToInterMacroblock`, always reconstructed motion-compensated Y, U and V. Pinned libvpx v1.16.0 `vp8/encoder/denoising.c:622` and the existing Go filter only consume UV predictors when the denoiser is not Y-only, motion magnitude is zero, and Y filtering succeeds.

This change uses `PredictWholeMVInterLumaWithState` when Y-only mode or nonzero motion already guarantees UV will not be read. The helper shares the original clamp, reference-offset and DSP implementations, retains validation of all reference planes before writing Y, and never touches chroma destinations or residual scratch. Zero-motion YUV blocks retain the full predictor. If the helper declines, the original full predictor and partial/split/unprepared-state fallbacks run unchanged.

The only deliberately different private state is unused `mcU`/`mcV` scratch after a successful luma-only prediction. Its only readers are in this function's zero-motion UV filter branch, and that branch reconstructs both planes before every read. A later UV-consumer regression verifies the lifetime boundary.

## Correctness evidence

- 32 frozen715 denoiser contexts: subpel motion, odd integer motion, even full-pel motion, Y-only, zero-motion YUV, rejected moving Y, split fallback, and unprepared-state fallback; each crossed with full/partial edges and separate/aliased source-overlay planes. Each runs three times and then a zero-motion UV consumer: 128 complete frozen comparisons.
- Compared overlay flags, all source/filtered planes, all running-average planes and borders, denoiser state, Y scratch, fallback images, and reconstruction scratch. Unused chroma scratch is intentionally unspecified until the next consumer, when it is compared too.
- New helper matches the existing complete predictor for full/partial images, corners, integer/fractional/extreme MVs and six-tap/bilinear/fullpixel configs; invalid modes, nil/unprepared states, short output, missing chroma and invalid zero-MV coordinates reject. Output row padding remains intact.
- Existing denoiser state/runtime-control tests and prepared-reference predictor tests pass.
- `TestVP8Algorithm2Stream`: all96 VP8 packet/flag hashes unchanged. `TestMLCacheStreamNative`: all72 VP9 rootML packet hashes unchanged.
- Complete denoiser operation is zero-allocation in focused assertion and every native case.

No new C comparison, profile, full suite, or broad matrix was run. Commands and logs are preserved in `correctness-commands.json`, `base-focused.log`, `candidate-correctness.log`, and `decoder-correctness.log`.

## Native component result

Apple M4 Max / Darwin arm64 / Go1.27.2 / GOEXPERIMENT=simd / CGO1 / PGO off / GOMAXPROCS2. Two balanced40ms pairs in A-B-B-A order, serialized with parent timing. This measures the complete `applyDenoiserToInterMacroblock` call with predictor construction, filtering, copies and state writes live, not just the helper. The source remains distinct from the overlay in these timings; frozen tests additionally cover aliasing/partial edges. No end-to-end codec improvement is claimed. Denoising must be enabled to reach this path.

| Case | Baseline ns pair1, pair2 | Candidate ns pair1, pair2 | Speedup pair1, pair2 |
|---|---:|---:|---:|
| subpel | 290.7, 289.7 | 194.1, 195.1 | 1.498x, 1.485x |
| odd-integer | 116.6, 116.1 | 77.03, 76.48 | 1.514x, 1.518x |
| fullpel | 80.31, 79.28 | 77.91, 76.67 | 1.031x, 1.034x |
| y-only | 70.1, 70.51 | 58.25, 58.46 | 1.203x, 1.206x |
| zero-yuv | 159.3, 163.5 | 165.0, 167.5 | 0.965x, 0.976x |
| rejected | 290.4, 293.2 | 196.4, 193.9 | 1.479x, 1.512x |
| split-fallback | 169.9, 172.1 | 175.1, 174.0 | 0.970x, 0.989x |
| state-fallback | 318.4, 316.9 | 319.1, 319.9 | 0.998x, 0.991x |

The moving subpel/odd-integer/rejected cases improve about48–52%; Y-only improves20%. Even full-pel copies improve3%. Zero-motion YUV and fallback controls contain small losses (all retained above); these are preliminary40ms samples. Exact raw observations remain in `native-{1,2}-{base,candidate}.log`, `native-summary.json`, and `native-commands.json`. Source, fixture, binary and pinned-C hashes are in `binding.json`. Fixed c4a8bd9f whole-codec gate remains unmet.

## Whole-codec short screen on merged main

Go1.27.2SIMD/PGOoff/CGO1/GOMAXPROCS1/threads1; VP8 1280×720,60inputs,30fps,2500kbps,realtime/noise4. Two balanced pairs against parent499 and fixed baselinec4. All22 non-performance fields and zero allocations agree. No separate quality/C/profile/full validation was repeated.

| Mode | Reference | Paired speedups | Median |
|---|---|---|---:|
| cpu-8 | fixed-c4 | 1.0476×, 1.0360× | 1.0418× |
| cpu-8 | parent499 | 1.0108×, 1.0216× | 1.0162× |
| cpu8 | fixed-c4 | 1.0099×, 0.9819× | 0.9959× |
| cpu8 | parent499 | 1.0102×, 0.9583× | 0.9843× |

Fixed-speed is preliminarily about1.6% faster than parent; adaptive is noisy and about1.6% lower in this screen, with a−4.2% second pair. The change is retained for its correctness-verified≈1.5× moving-block denoiser component win; no2% codec or libvpx-gap closure claim. Zero-motion/fallback component losses and all codec losses remain visible. The1.5× fixed-baseline whole-codec gate remains unmet.

The `previous-*` files consolidate x64 native feedback and green merge receipts for PR107/108, plus the preserved inconclusive unused-token-table probe. Their Intel component timings are kept separate from earlier EPYC libvpx timings. They do not establish a combined whole-codec speedup. README remains unchanged.
