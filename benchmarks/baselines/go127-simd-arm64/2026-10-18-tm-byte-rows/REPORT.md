# ARM TM8/16 saturated-byte SIMD evidence

Exact parent: `ac259fce071f3ac5fc45603c6bd7e7926bdfcf73`. Old assembly reference: `691cd0512c4829d9ee40d483205dd223f9eb028d`. Native Apple M4 Max; official Go 1.27.2 SIMD, PGO off. The accepted fresh runs use CGO_ENABLED=1 for all three builds and GOMAXPROCS=1.

The ARM helper loads the disjoint above packet once and preserves each per-row left read. Saturating byte addition/subtraction implements clamp(above + left − top, 0, 255): only one signed-delta bias is nonzero. Complete visible windows and nonnegative strides below 2^30 are validated before writes; above/output overlap declines. Left/output and read-only input aliases retain their sequential row dependencies. Invalid shapes fall back before stores, preserving original panic and partial-write behavior. AMD, purego and experiment-off builds retain their original implementations.

Fresh component acceptance contains 486 samples: 27 actual public DSP cases × six balanced three-way rounds, 100 ms per case. Every sample passed the independent preflight/post-loop oracle and reported 0 B/op, 0 allocs/op. All 18 TM8/16 target cases and all 108 paired comparisons exceed 1.5× parent. TM8 median gains are 2.105063–2.169680×; TM16 gains are 4.668782–4.733505×. The user explicitly retains qualifying component wins even when whole-codec time is flat.

Old NEON remains faster: old/current median time ratios are 0.475318–0.479051 for TM8 and 0.498718–0.524785 for TM16. This package makes no assembly-beating claim.

The nine TM4 controls retain all regressions and pair signs below. The generic ARM entry frame grew from 16 to 112 bytes and adds two width comparisons; TM4 never executes the helper. TM4 median changes range from −0.130% to +1.511%, with mixed paired signs. No control is discarded.

Whole-codec acceptance is a separate 18-run, six balanced three-way VP8 window: 1280×720, 300 inputs, realtime cpu-used=8, 30 fps, 2500 kbps, one thread, encode-only. Parent/current/old ns per input medians are 7,251,048 / 7,265,317 / 6,621,017. Current is +0.196785% versus parent, only two of six pairs faster: flat/mixed, with no whole-codec gain claimed. The same-run current/old gap is +9.731133%, all six pairs slower. Ratio-of-medians and median-of-paired-ratios are recorded separately in compact.json.

All 22 semantic fields and allocation values match every whole JSON: 298 encoded / 2 dropped, 3,120,743 output bytes, identical quantizer statistics, 0 allocations per input. Quality measurement was explicitly skipped; neither matching payload statistics nor this audit proves bitstream identity.

Three separate 600-input diagnostic profiles use the same CLI binaries and options except frame count and CPU-profile output. All 22 semantic fields and allocations match: 595 encoded / 5 dropped, 6,310,889 bytes, 0 allocations. Their raw profiles and flat/cumulative top-50 views are retained. These sampled profiles are not timing denominators or causal proof of a performance gain.

Correctness includes frozen original SIMD row-capture semantics, distinct original pixel scalar and public DSP disabled semantics, full backing canaries, clipped capacities, min spans, alignments, exact/partial/cross-row aliases, repeated rows, extreme strides, unsupported widths, short inputs, and panic/partial-state equivalence. Native ordinary, race/checkptr=2, purego, experiment-off, AMD selected Rosetta and actual AVX/AVX2-off checks pass. Rosetta results establish correctness only. The helper has 94 optimized ARM instructions, no calls, no stack frame or vector stack traffic; AMD public bodies are normalized instruction-identical. Instruction counts are diagnostic evidence, not timing results.

Limitations and failed evidence are preserved: the initial 243-sample screen used CGO=1 for the candidate and CGO=0 for references and is provisional, never pooled with fresh acceptance. The initial bounds-checked-left helper source and codegen are retained. An initial test incorrectly reused pixel scalar semantics for the disabled public DSP fallback; its failing log/source and the corrected row-capture oracle are retained. Production code was unchanged by that test correction.

The source proof verifies 784 parent production/module Git blobs. Only the three production overlay files differ; modules and every README remain unchanged. The overlay contains all six test files as well. The identical public benchmark fixture was overlaid on parent and old references; it calls the actual public DSP functions and consumes every output pixel.

Compiled executables are deliberately omitted from this lean publication. OMITTED-BINARIES.json records every retained private executable path, byte size and SHA-256. Full private snapshots and binaries remain unchanged. The raw archive contains selected exact reference sources, all measured/failed variants, validation/build/codegen logs, all raw benchmark/whole/profile outputs, drivers and source proofs. It contains no complete source snapshots.

## All fresh public cases

| Case | Parent ns | Current ns | Old ns | Parent/current | Old/current | Six parent/current paired ratios |
|---|---:|---:|---:|---:|---:|---|
| 4x4/pitch4/flat | 6.915000 | 7.019500 | 2.979500 | 0.985113 | 0.424460 | 1.001172, 0.937669, 0.982583, 1.000434, 0.990397, 0.989437 |
| 4x4/pitch4/mixed | 6.950000 | 6.943000 | 2.988000 | 1.001008 | 0.430362 | 1.001178, 1.003022, 0.994666, 0.983437, 1.051794, 0.997275 |
| 4x4/pitch4/clipped | 6.929500 | 6.920500 | 2.983000 | 1.001300 | 0.431038 | 0.970551, 0.998119, 1.014027, 0.999133, 1.009536, 1.002888 |
| 4x4/pitch37/flat | 6.895500 | 6.928000 | 2.975000 | 0.995309 | 0.429417 | 0.982711, 0.998695, 0.991223, 1.003031, 1.003470, 0.996391 |
| 4x4/pitch37/mixed | 6.931000 | 6.973000 | 2.981000 | 0.993977 | 0.427506 | 1.011129, 0.983231, 0.992825, 1.002018, 0.985899, 1.000577 |
| 4x4/pitch37/clipped | 6.933000 | 6.946500 | 3.003500 | 0.998057 | 0.432376 | 0.997970, 1.004347, 0.971115, 0.997698, 0.998848, 1.002444 |
| 4x4/pitch1312/flat | 6.922000 | 6.933500 | 2.996000 | 0.998341 | 0.432105 | 0.992654, 0.998120, 0.997270, 1.001301, 1.012276, 0.982663 |
| 4x4/pitch1312/mixed | 6.907500 | 6.922000 | 2.983500 | 0.997905 | 0.431017 | 0.983962, 0.995952, 1.004057, 0.992950, 0.998700, 1.001011 |
| 4x4/pitch1312/clipped | 6.940000 | 6.964000 | 2.989500 | 0.996554 | 0.429279 | 0.961931, 0.987638, 0.988361, 1.003034, 0.997561, 1.000577 |
| 8x8/pitch8/flat | 16.965000 | 7.840500 | 3.756000 | 2.163765 | 0.479051 | 2.098403, 2.214221, 2.116853, 2.102102, 2.234151, 2.162606 |
| 8x8/pitch8/mixed | 17.085000 | 7.941000 | 3.774500 | 2.151492 | 0.475318 | 2.113003, 2.201387, 2.117736, 2.125827, 2.205309, 2.160873 |
| 8x8/pitch8/clipped | 16.850000 | 7.905000 | 3.774000 | 2.131562 | 0.477419 | 2.170453, 2.127769, 2.118752, 2.131148, 2.145560, 2.118003 |
| 8x8/pitch37/flat | 16.760000 | 7.924500 | 3.768000 | 2.114960 | 0.475487 | 2.093053, 2.131665, 2.323829, 2.205419, 2.076962, 2.134889 |
| 8x8/pitch37/mixed | 16.680000 | 7.890500 | 3.769000 | 2.113934 | 0.477663 | 2.106938, 2.111559, 2.111154, 2.116696, 2.258428, 2.127741 |
| 8x8/pitch37/clipped | 16.765000 | 7.874500 | 3.755500 | 2.129024 | 0.476919 | 2.150096, 2.130065, 2.082188, 2.121637, 2.223063, 2.176232 |
| 8x8/pitch1312/flat | 16.630000 | 7.900000 | 3.760000 | 2.105063 | 0.475949 | 2.096897, 2.132063, 2.100946, 2.132495, 2.241401, 2.082599 |
| 8x8/pitch1312/mixed | 17.160000 | 7.909000 | 3.773000 | 2.169680 | 0.477051 | 2.218509, 2.142228, 2.239756, 2.152138, 2.189074, 2.098593 |
| 8x8/pitch1312/clipped | 16.820000 | 7.862500 | 3.737500 | 2.139269 | 0.475358 | 2.160185, 2.138446, 2.203454, 2.241750, 2.128966, 2.109159 |
| 16x16/pitch16/flat | 54.925000 | 11.705000 | 5.837500 | 4.692439 | 0.498718 | 4.635897, 4.686106, 4.697774, 4.685957, 4.659322, 4.752348 |
| 16x16/pitch16/mixed | 55.045000 | 11.645000 | 5.832000 | 4.726921 | 0.500816 | 4.691113, 4.682051, 4.772372, 4.615900, 4.718036, 4.766178 |
| 16x16/pitch16/clipped | 55.130000 | 11.680000 | 5.832000 | 4.720034 | 0.499315 | 4.720971, 4.641525, 4.763841, 4.633389, 4.717428, 4.807491 |
| 16x16/pitch37/flat | 55.240000 | 11.670000 | 6.023500 | 4.733505 | 0.516153 | 4.713790, 4.636749, 4.795652, 4.587894, 4.607084, 4.929825 |
| 16x16/pitch37/mixed | 54.985000 | 11.680000 | 6.051000 | 4.707620 | 0.518065 | 4.746528, 4.693231, 4.665532, 4.666667, 4.717707, 4.859772 |
| 16x16/pitch37/clipped | 55.045000 | 11.635000 | 6.038000 | 4.730984 | 0.518951 | 4.677835, 4.723130, 4.814750, 4.705042, 4.717818, 4.759275 |
| 16x16/pitch1312/flat | 55.050000 | 11.640000 | 6.108500 | 4.729381 | 0.524785 | 4.663519, 4.713672, 4.594300, 4.723404, 4.800172, 4.765517 |
| 16x16/pitch1312/mixed | 55.105000 | 11.720000 | 6.123000 | 4.701792 | 0.522440 | 4.742437, 4.761945, 4.623737, 4.783691, 4.631933, 4.700512 |
| 16x16/pitch1312/clipped | 55.185000 | 11.820000 | 6.102000 | 4.668782 | 0.516244 | 4.667243, 4.798801, 4.626904, 4.683587, 4.652614, 4.637815 |


Publication: `raw-evidence.tar.gz` preserves all 16 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
