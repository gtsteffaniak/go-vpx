# VP8 minimum-cost full-pel search elimination

Preliminary native component screen against `715645fd212cfc65295f4f6fb1002dcc60856758` (all 50 prior PRs merged). Production change: `vp8_encoder_inter_motion_search.go`, `selectInterFrameFullPixelMotionVectorWithSearchStartAndProbsCostTables`, lines 216–226. The root checkout was not edited.

The old dispatch only stops at SAD+MV-SAD cost zero. The ordinary MV-SAD table has a positive floor at zero delta, so an exact match to the integer-pel predictor still evaluates entire NSTEP/DIAMOND/HEX neighborhoods. This patch returns immediately when the initial walk cost equals that global floor. It removes those SAD evaluations and their repeated candidate bounds/cost work.

## Exactness argument and source

Pinned local libvpx v1.16.0 `vp8/encoder/onyx_if.c:1747`, `cal_mvsadcosts`, assigns 300 to each zero-delta component, and larger nonzero logarithmic entries. `vp8/encoder/mcomp.c:54`, `mvsad_err_cost`, combines the two nonnegative entries. The local `internal/vp8/encoder/rd_cost.go:192` table construction and `FullPelMVSADCost16FromDeltas` preserve that floor for every q index. SAD is nonnegative. NSTEP, DIAMOND, final refinement and HEX only move on a strictly lower total cost; therefore no candidate can move when the initial cost equals the floor.

HEX returns the unchanged SAD+MV-SAD walk cost. NSTEP/DIAMOND retain their original variance+MV-error return calculation. Existing zero-cost handling remains first. Exhaustive and invalid-method fallback behavior, first-pass search, and non-nil diagnostic stats paths remain untouched. Fractional refinement still runs exactly as before.

## Correctness

- Frozen original715 dispatch compared for MV, cost, and full `selectFast` result over q=-1..128, three search methods, varying step/final-refine controls, flat/exact/shifted/fractional/partial/dense/improved-start inputs.
- Independent source-derived lower-bound invariant checked for all 128 q tables and all 256 absolute deltas; nil/custom MV-error-cost and exhaustive/invalid-method controls included.
- Existing affected full-pel topology, border/partial, improved-start and final-refine tests pass.
- `TestVP8Algorithm2Stream`: all 96 VP8 packet/flag hashes unchanged. `TestMLCacheStreamNative`: all 72 VP9 root-ML packet hashes unchanged.
- Zero allocations for tested full-pel and complete fast-motion component calls, including controls.

`candidate-correctness.log`, `candidate-final-refine.log`, `base-focused.log`, and `correctness-commands.json` preserve actual commands and passing evidence. The frozen oracle is independent of the new dispatch branch; it reuses unchanged search primitives. No new C comparison, full suite, profile, or matrix was run.

## Native result

Apple M4 Max / Darwin arm64 / Go1.27.2 / GOEXPERIMENT=simd / PGO off / CGO1 / GOMAXPROCS2. Two short balanced 40 ms pairs, A-B-B-A order, serialized with the parent. This is a preliminary component result, not an end-to-end codec claim. Exact-predictor full-pel gains are 1.77–2.65x (flat up to2.66x). HEX full-pel plus fractional refinement improves 1.20–1.24x on exact/flat inputs. Nonmatching controls remain near parity; every raw loss is retained below.

| Component | Baseline ns pair1, pair2 | Candidate ns pair1, pair2 | Speedup pair1, pair2 |
|---|---:|---:|---:|
| fullpel/NSTEP/flat | 245.7, 243.3 | 94.12, 91.55 | 2.610x, 2.658x |
| fullpel/NSTEP/exact | 240.6, 237.3 | 92.74, 91.52 | 2.594x, 2.593x |
| fullpel/NSTEP/shifted | 244.7, 236.3 | 92.36, 91.12 | 2.649x, 2.593x |
| fullpel/NSTEP/partial | 398.0, 354.4 | 214.2, 206.7 | 1.858x, 1.715x |
| fullpel/NSTEP/dense | 338.1, 332.4 | 331.8, 333.3 | 1.019x, 0.997x |
| fullpel/NSTEP/improved | 171.5, 173.2 | 172.0, 171.2 | 0.997x, 1.012x |
| fullpel/DIAMOND/flat | 181.8, 172.6 | 93.37, 91.88 | 1.947x, 1.879x |
| fullpel/DIAMOND/exact | 168.1, 170.1 | 91.88, 91.29 | 1.830x, 1.863x |
| fullpel/DIAMOND/shifted | 174.8, 166.0 | 91.09, 93.84 | 1.919x, 1.769x |
| fullpel/DIAMOND/partial | 329.5, 284.9 | 226.3, 221.6 | 1.456x, 1.286x |
| fullpel/DIAMOND/dense | 236.3, 229.6 | 232.4, 229.2 | 1.017x, 1.002x |
| fullpel/DIAMOND/improved | 140.3, 137.8 | 138.3, 138.4 | 1.014x, 0.996x |
| fullpel/HEX/flat | 175.8, 170.5 | 72.19, 72.2 | 2.435x, 2.361x |
| fullpel/HEX/exact | 162.7, 163.9 | 72.05, 72.27 | 2.258x, 2.268x |
| fullpel/HEX/shifted | 170.3, 163.3 | 72.29, 71.58 | 2.356x, 2.281x |
| fullpel/HEX/partial | 295.4, 329.5 | 195.3, 185.8 | 1.513x, 1.773x |
| fullpel/HEX/dense | 242.1, 241.6 | 238.0, 239.2 | 1.017x, 1.010x |
| fullpel/HEX/improved | 170.7, 176.3 | 169.1, 169.9 | 1.009x, 1.038x |
| fast/hex/flat | 566.3, 573.0 | 464.6, 464.0 | 1.219x, 1.235x |
| fast/hex/exact | 557.6, 567.0 | 464.4, 464.3 | 1.201x, 1.221x |
| fast/hex/dense | 684.1, 693.3 | 681.9, 681.4 | 1.003x, 1.017x |
| fast/hex/improved | 613.9, 614.1 | 612.8, 615.6 | 1.002x, 0.998x |

Exact source, fixture, C source and binary SHA256s are in `binding.json`. Raw samples are `native-{1,2}-{base,candidate}.log`; `native-commands.json` records exits/timing and `native-summary.json` retains paired numbers. Whole-codec impact remains to be measured by the parent on its canonical workload. The fixed c4a8bd9f full-validation gate remains unmet.

## Parent canonical whole-codec short screen

Runtime8d977bda; source SHA matches the measured component candidate exactly. Tests were renamed for the existing fast CI selector; their benchmark body and production source are unchanged. VP8 1280×720/60inputs/2500kbps/30fps/noise4/realtime, Go1.27.2SIMD,threads1,GOMAXPROCS1,CGO1,PGOoff. Two balanced samples against latest-main715 and immutablec4 baseline. All22 output semantics and zero allocations match. Separate quality decoding/C comparison was not repeated.

| Mode | Reference | Paired throughput speedups | Median |
|---|---|---|---:|
| cpu-8 | fixed-c4 | 1.0077×, 1.0380× | 1.0229× |
| cpu-8 | main | 1.0083×, 1.0043× | 1.0063× |
| cpu8 | fixed-c4 | 0.9825×, 1.0374× | 1.0099× |
| cpu8 | main | 0.9972×, 1.0354× | 1.0163× |

Latest-main fixed-speed gain is about0.6% and adaptive about1.6%, both preliminary and below the2% codec target. Retain for the1.77–2.65× exact-match full-pel kernel improvement and20–24% complete-HEX component gain. This does not establish closure of the libvpx gap. The two variants in native-manifest.json are independent overlays against715; they are not a combined-codec claim. README unchanged; full validation remains deferred.
