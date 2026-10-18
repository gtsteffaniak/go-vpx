# ARM TM direct-entry unroll evidence

Retention attribution correction: the user mandated keeping gains at or above 1.5× and separately requested assembly-beating performance. Keeping the smaller gains against the optimized parent here is engineering judgment consistent with that goal. The historical raw report remains unchanged; see the [scalar-contract supplement](../2026-10-18-tm-native-unroll-scalar-contract/REPORT.md) for corrected source and codegen identity evidence.

Original Git base `ac259fce071f3ac5fc45603c6bd7e7926bdfcf73`; strongest acceptance parent is the exact frozen left-packet TM winner, not the original Git implementation. Old assembly reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Apple M4 Max; official Go 1.27.2 SIMD, PGO off, CGO_ENABLED=1 for all six DSP/CLI builds, GOMAXPROCS=1 while measuring.

The final candidate removes the bool fallback-liveness layer from ARM public TM8/16 calls, restores the DSP wrappers’ inlining and unrolls four row pairs / sixteen rows. The actual timed TM8/16 callers make one direct call to TMPredictNative. Original two-table-lookup saturated-byte arithmetic is retained. No packed-store, above-bias-swap or Float64 bit-store variant is combined with this candidate.

TMPredictBytes now contains only its original eligibility/window/above-disjoint guards, then delegates and returns true. Declines retain the no-read/no-write contract. Native guards, original stateful alias path and original SIMD fallback retain their semantics; the native source is byte-identical to the measured unroll prototype. This removes duplicate arithmetic. The original DSP CPU-disabled scalar path remains exact, including its upfront checks and row-capture behavior, which differs from the pixel scalar per-pixel fallback.

Fresh acceptance contains 486 samples: 27 actual public DSP cases × six balanced three-way rounds, 100 ms per case. All independent preflight/post-loop full-state oracles pass and every sample reports 0 B/op, 0 allocs/op. Against the strongest left-packet parent, TM8 median gains are 1.271666–1.288978×, all 54 paired gains 1.254242–1.397494×. TM16 median gains are 1.408664–1.470906×, all 54 paired gains 1.343549–1.510027×. Every target pair is faster, but no 1.5× median claim is made against this strong baseline. Gains are measured directly, never multiplied from prior screens.

Against same-run old NEON, all nine TM16 medians improve 1.030930–1.084735× and all 54 pairs are faster, 1.013380–1.162186×. TM8 old/current median ratios are only 1.003495–1.012231×: a small 0.35–1.22% median advantage with seven of 54 pairs slower (paired range 0.988886–1.064709×). TM8 is near the old-reference timing floor; no universal or every-pair assembly-beating claim. The user explicitly authorized retaining repeatable assembly-beating component gains below the 1.5× threshold against an already optimized parent.

All nine TM4 controls are preserved: median parent/current ratios 0.993219–1.031803× and paired ratios 0.955907–1.047496×, with mixed signs. The full table and every raw sample remain below/in compact.json.

Whole-codec validation is a separate 18-run six-balanced-three-way VP8 window: 1280×720, 300 inputs, realtime cpu-used=8, 30 fps, 2500 kbps, one thread, encode-only. Parent/current/old median ns/input are 7,274,262.5 / 7,277,098.5 / 6,648,388.5. Candidate is +0.038987% versus parent, only one of six pairs faster; paired changes are +0.130904, −1.237900, +0.946466, +0.072949, +0.425459, +0.154407%. This is effectively flat/mixed, not a whole-codec win. The same-run current/old gap remains +9.456577%, all six pairs slower 3.081939–10.630667%. Ratio-of-medians and median paired ratio are recorded separately. Do not compare or subtract gaps across windows or hardware.

All 22 semantic fields and allocations match every whole JSON: 298 encoded / 2 dropped, 3,120,743 bytes, identical quantizer statistics and zero allocations per input. Three separate 600-input diagnostic profiles also match all 22 fields plus zero allocation: 595 encoded / 5 dropped, 6,310,889 bytes. Raw profiles and six flat/cumulative top-50 views are retained. Quality was explicitly skipped; matching these statistics does not prove bitstream identity. Profiles are diagnostic samples, not timing acceptance denominators.

Correctness includes full native pixel/DSP packages, exhaustive public byte arithmetic (511 signed deltas × 256 left values × both widths), frozen scalar/original SIMD oracles, min spans/clipped capacities, independent alignments, input/output and cross-row aliases, overlapping rows, negative/extreme strides, unsupported shapes and exact panic/partial-write state. Ordinary, race/checkptr=2, purego, nosimd and actual AMD selected/off tests pass. AMD pixel/public codegen is normalized-identical to the strong parent. Rosetta establishes correctness only.

ARM native code is 364 instructions with a 16-byte frame, no hot calls or vector stack spills. The bool refactor preserves normalized native and actual caller bodies from the measured prototype. Fresh ordinary builds also match preflight codegen; only two post-loop checksum-sink relocation immediates differ, and both targets were verified against symbol tables. The actual timed region is identical. Compiler -m records inner/public wrapper costs 65/72 and confirms inlining. These diagnostics are not performance measurements.

All earlier research is preserved, including every slower control and pair: packed-store, bias-swap, void-dispatch, direct-glue and initial unroll screens total 1,215 samples, separate from final fresh acceptance. Float64 raw-bit stores passed 99,456 full-state cases including signaling/quiet NaNs but removed no instructions, so no timing was run. History contains exact changed sources, drivers, build/oracle/codegen logs, screen commands/timestamps and raw outputs. The original frozen TM packages are referenced by manifest SHA and remain unchanged.

Buildinfo tooling initially expected an explicit PGO-off field which Go 1.27.2 omits; verification now checks the recorded build commands. The initial parser source/note is preserved. No codec correctness or performance gate was weakened. The codegen comparison likewise preserves its initial relocation-only mismatch note and exact follow-up proof.

No root source, index, branch, README, dependency, assembly or toolchain was changed. The complete original-Git overlay contains six production and eight test files; incremental.patch against the strongest left-packet winner contains exactly five production changes and one added public-arithmetic test. Executables are omitted with exact private paths, byte sizes and SHA256; full private snapshots/binaries remain unchanged. The raw archive contains selected sources rather than complete snapshots.

## Every fresh public case

| Case | Parent ns | Current ns | Old ns | Parent/current | Old/current | Six parent/current ratios |
|---|---:|---:|---:|---:|---:|---|
| 4x4/pitch4/flat | 6.945500 | 6.941500 | 3.007000 | 1.000576 | 0.433192 | 0.955907, 1.000870, 1.007538, 1.015704, 0.997700, 1.028954 |
| 4x4/pitch4/mixed | 6.957000 | 7.004500 | 2.980500 | 0.993219 | 0.425512 | 1.005098, 1.005389, 0.968124, 1.005462, 1.000288, 0.983465 |
| 4x4/pitch4/clipped | 6.993000 | 6.941500 | 2.987000 | 1.007419 | 0.430310 | 0.994084, 1.018561, 0.984437, 1.006926, 1.007910, 1.017677 |
| 4x4/pitch37/flat | 6.998500 | 6.951000 | 2.978500 | 1.006834 | 0.428499 | 0.979304, 1.002611, 1.010248, 0.973567, 1.029780, 1.031650 |
| 4x4/pitch37/mixed | 6.988500 | 6.935000 | 2.977500 | 1.007714 | 0.429344 | 1.004971, 1.008115, 1.015506, 1.004027, 1.033223, 1.007053 |
| 4x4/pitch37/clipped | 6.993000 | 6.955000 | 2.977500 | 1.005464 | 0.428109 | 1.011331, 1.008632, 1.001005, 1.001150, 1.034230, 1.015047 |
| 4x4/pitch1312/flat | 6.974000 | 6.929000 | 2.983000 | 1.006494 | 0.430509 | 1.004966, 1.010114, 1.000577, 1.003606, 1.029153, 0.990791 |
| 4x4/pitch1312/mixed | 7.003500 | 6.941500 | 2.983000 | 1.008932 | 0.429734 | 0.998981, 1.017518, 1.001000, 1.003297, 1.013424, 0.999281 |
| 4x4/pitch1312/clipped | 7.170000 | 6.949000 | 2.980500 | 1.031803 | 0.428911 | 1.005146, 1.039873, 0.992961, 1.027770, 1.047496, 1.036882 |
| 8x8/pitch8/flat | 4.736500 | 3.719500 | 3.748000 | 1.273424 | 1.007662 | 1.289424, 1.315599, 1.257842, 1.266060, 1.320968, 1.274805 |
| 8x8/pitch8/mixed | 4.771500 | 3.715500 | 3.737500 | 1.284215 | 1.005921 | 1.292888, 1.287022, 1.275464, 1.344071, 1.311542, 1.286211 |
| 8x8/pitch8/clipped | 4.795000 | 3.720000 | 3.733000 | 1.288978 | 1.003495 | 1.382328, 1.286940, 1.256946, 1.300027, 1.310502, 1.288943 |
| 8x8/pitch37/flat | 4.777500 | 3.716500 | 3.750000 | 1.285484 | 1.009014 | 1.397494, 1.283715, 1.278988, 1.277466, 1.282510, 1.285253 |
| 8x8/pitch37/mixed | 4.739500 | 3.727000 | 3.766500 | 1.271666 | 1.010598 | 1.286422, 1.305356, 1.254242, 1.274940, 1.293136, 1.256085 |
| 8x8/pitch37/clipped | 4.739500 | 3.722000 | 3.751500 | 1.273375 | 1.007926 | 1.278362, 1.379378, 1.254463, 1.271141, 1.301156, 1.257987 |
| 8x8/pitch1312/flat | 4.751500 | 3.719000 | 3.763000 | 1.277628 | 1.011831 | 1.260372, 1.275390, 1.285868, 1.273631, 1.301130, 1.347309 |
| 8x8/pitch1312/mixed | 4.745500 | 3.720000 | 3.765500 | 1.275672 | 1.012231 | 1.346615, 1.260417, 1.272458, 1.279226, 1.284918, 1.268621 |
| 8x8/pitch1312/clipped | 4.747000 | 3.723500 | 3.769000 | 1.274876 | 1.012220 | 1.315673, 1.260997, 1.262667, 1.315949, 1.286251, 1.276121 |
| 16x16/pitch16/flat | 8.096000 | 5.658000 | 5.833000 | 1.430894 | 1.030930 | 1.426056, 1.447185, 1.428671, 1.365751, 1.465426, 1.409297 |
| 16x16/pitch16/mixed | 8.090000 | 5.623000 | 5.841000 | 1.438734 | 1.038769 | 1.405960, 1.448871, 1.446441, 1.431798, 1.510027, 1.422108 |
| 16x16/pitch16/clipped | 8.057500 | 5.651500 | 5.844500 | 1.425728 | 1.034150 | 1.394197, 1.408388, 1.453416, 1.398622, 1.488286, 1.459919 |
| 16x16/pitch37/flat | 8.304000 | 5.645500 | 6.034000 | 1.470906 | 1.068816 | 1.394886, 1.412651, 1.472060, 1.471578, 1.424496, 1.465243 |
| 16x16/pitch37/mixed | 8.167000 | 5.688000 | 6.036000 | 1.435830 | 1.061181 | 1.428673, 1.445286, 1.412665, 1.433491, 1.414044, 1.456010 |
| 16x16/pitch37/clipped | 8.097000 | 5.722500 | 6.062000 | 1.414941 | 1.059327 | 1.371793, 1.476719, 1.379667, 1.361926, 1.430726, 1.478307 |
| 16x16/pitch1312/flat | 8.104000 | 5.647000 | 6.125500 | 1.435098 | 1.084735 | 1.388998, 1.433227, 1.394159, 1.438562, 1.503894, 1.439007 |
| 16x16/pitch1312/mixed | 8.139000 | 5.674000 | 6.112500 | 1.434438 | 1.077282 | 1.343549, 1.428522, 1.432820, 1.465971, 1.458795, 1.417773 |
| 16x16/pitch1312/clipped | 8.016000 | 5.690500 | 6.113500 | 1.408664 | 1.074334 | 1.424726, 1.380201, 1.376596, 1.404463, 1.478831, 1.460802 |


Publication: `raw-evidence.tar.gz` preserves all 23 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
