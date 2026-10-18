# ARM TM left-packet SIMD evidence

Original Git base `ac259fce071f3ac5fc45603c6bd7e7926bdfcf73`; strongest timed parent is the already accepted first saturated-byte TM winner, not the original Git implementation. Old assembly reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Apple M4 Max; official Go 1.27.2 SIMD, PGO off, CGO_ENABLED=1 for every build, GOMAXPROCS=1 while measuring.

The only production change versus the strongest parent is the ARM helper. After the existing complete-window and above/output-disjoint checks, a conservative ordered-span test permits hoisting left only when left/output are disjoint. All aliases retain the original per-row left-capture code. The generic hook and other-architecture false stub are byte-identical.

TM16 precomputes saturated positive/negative left-minus-top packets, uses two TBL lookups per row, and stores exactly 16 bytes. TM8 uses exact 8-byte above/left loads, repeats above into both halves and emits two rows per packet; the low-row store precedes the high-row store even when destination rows overlap. Row indices remain within 0..w−1. Pointers advance only when another row/pair remains. Stride below 2^30 and validated visible windows prove all offsets fit int64. The two saturations exactly implement clamp(above+left−top,0,255), since at most one bias is nonzero.

Fresh acceptance preserves 486 samples:27 actual public DSP cases  × 6 balanced three-way rounds  × 100 ms. Every preflight/post-loop independent oracle passed, all samples 0 B/op and 0 allocs/op. TM8 all 9 median gains are 1.652685–1.676489 ×  the strongest parent; all 54 paired gains are 1.607442–1.770860 × . Every TM8 case/pair exceeds the user’s 1.5 ×  component-retention rule. TM16 all 9 median gains are 1.423349–1.480927 ×  and all 54 pairs 1.340609–1.509195 × : faster, with no 1.5 ×  claim. The separate 243-sample 75 ms screen remains separate and is not pooled.

The 9 TM4 controls are unchanged in source but retain all layout/noise regressions: median parent/current ratios 0.987915–1.005314, all paired ratios 0.937269–1.040523, mixed signs. No control is omitted. Old NEON remains faster: old/current median ratios 0.787752–0.796191 for 8 and 0.716083–0.763002 for 16. No assembly-beating claim.

Whole-codec acceptance is a separate 18-run six-balanced-three-way window: VP 8 1280 × 720,300 inputs, realtime cpu-used=8,30 fps,2500 kbps,1 thread, encode-only. Parent/current/old median ns/input are 7,310,921 /7,334,380.5 /6,696,427. Candidate is+0.320883% versus strongest parent,2/6 pairs faster; the paired changes are+2.556143,+1.281731,+0.624472,+0.417170,−1.099354,−0.437765%. This is mixed/flat, with no whole-codec gain claim. Same-run current/old gap is+9.526775%, all 6 pairs slower 7.624040–12.902251%. Ratio-of-medians and median paired ratio are recorded separately. Do not sum gains or subtract gaps across experiment windows.

All 22 semantic fields and allocations match all 18 whole JSONs:298 encoded/2 dropped,3,120,743 output bytes, identical quantizer statistics,0 allocations per input. Quality was explicitly skipped; matching statistics do not prove bitstream identity. The 3 separate 600-input diagnostic profiles likewise match all 22 fields+0 alloc:595 encoded/5 dropped,6,310,889 bytes. Raw profiles and flat/cumulative top 50 views are retained; sampled profiles are not timing acceptance denominators or causal proof.

Ordinary, race/checkptr=2, purego and nosimd tests pass. Inherited full-state oracles cover min spans/clipped capacities, all alignments, input/output aliases, overlapping rows, extreme/negative strides, unsupported shapes, and original panic/partial-write state. Added packet-lane tests use distinct left values, all 256 byte shifts, six top boundaries and five pitches including repeated/overlapping output rows. All 27 public one-shot oracles pass on all 3 fresh builds with 0 allocation. Fixture bytes are identical.

Optimized helper has 160 total instructions including retained alias paths, no frame, calls, bounds panics or vector stack traffic. Packet loops contain 2 TBL, saturating add/sub and exact stores, without per-row scalar broadcasts. AMD pixel/public bodies are normalized instruction-identical to the frozen winner. These codegen diagnostics are not performance claims.

Full overlay contains 3 production and 7 test files relative to original Git base. Incremental patch changes only the helper and adds the packet-lane test relative to the frozen winner. The original accepted winner package remains unchanged; selected strongest source and proofs are included for reproducibility. No root files, README, module dependencies, assembly or toolchain were changed. Executables are omitted from this lean publication and recorded with private path, bytes and SHA256; all full private snapshots/binaries remain unchanged.

## Every fresh case

| Case | Parent ns | Current ns | Old ns | Parent/current | Old/current | Six parent/current ratios |
|---|---:|---:|---:|---:|---:|---|
| 4 x 4/pitch 4/flat | 6.978000 | 6.944500 | 2.988000 | 1.004824 | 0.430269 | 1.034041, 1.014895, 0.995109, 1.003747, 0.991070, 0.951597 |
| 4 x 4/pitch 4/mixed | 6.977000 | 6.948500 | 2.985000 | 1.004102 | 0.429589 | 0.997242, 1.011841, 1.007785, 0.976344, 1.002586, 0.993448 |
| 4 x 4/pitch 4/clipped | 6.983500 | 6.976000 | 2.982000 | 1.001075 | 0.427466 | 0.999709, 0.995272, 1.040523, 1.024957, 0.968269, 0.988247 |
| 4 x 4/pitch 37/flat | 6.973500 | 7.007000 | 2.995500 | 0.995219 | 0.427501 | 1.003042, 0.982850, 1.028930, 0.975780, 1.009573, 0.998283 |
| 4 x 4/pitch 37/mixed | 6.979500 | 7.026500 | 3.005500 | 0.993311 | 0.427738 | 0.969292, 1.025090, 0.984075, 1.005495, 1.035245, 0.998006 |
| 4 x 4/pitch 37/clipped | 6.948500 | 7.033500 | 3.005000 | 0.987915 | 0.427241 | 0.937269, 0.991692, 0.975473, 1.005057, 0.980243, 1.000720 |
| 4 x 4/pitch 1312/flat | 6.936000 | 6.974500 | 2.983500 | 0.994480 | 0.427773 | 0.981725, 1.005664, 0.995540, 1.034517, 0.968358, 0.993139 |
| 4 x 4/pitch 1312/mixed | 6.999500 | 6.962500 | 2.990000 | 1.005314 | 0.429443 | 1.007801, 0.993256, 1.012363, 0.995462, 1.008501, 1.007942 |
| 4 x 4/pitch 1312/clipped | 6.936500 | 7.019000 | 2.990000 | 0.988246 | 0.425987 | 1.003931, 1.005362, 1.016796, 0.988034, 1.016244, 0.968400 |
| 8 x 8/pitch 8/flat | 7.881500 | 4.741500 | 3.767500 | 1.662238 | 0.794580 | 1.670659, 1.686208, 1.655005, 1.641188, 1.659574, 1.677923 |
| 8 x 8/pitch 8/mixed | 7.936500 | 4.752000 | 3.783500 | 1.670139 | 0.796191 | 1.639186, 1.683362, 1.665061, 1.669200, 1.669881, 1.677086 |
| 8 x 8/pitch 8/clipped | 7.933500 | 4.740500 | 3.765000 | 1.673558 | 0.794220 | 1.691419, 1.675596, 1.673207, 1.647320, 1.654367, 1.687592 |
| 8 x 8/pitch 37/flat | 7.880000 | 4.768000 | 3.756000 | 1.652685 | 0.787752 | 1.668726, 1.681712, 1.628881, 1.645005, 1.656988, 1.658021 |
| 8 x 8/pitch 37/mixed | 7.924000 | 4.755000 | 3.759000 | 1.666456 | 0.790536 | 1.663694, 1.657029, 1.770860, 1.637411, 1.671813, 1.675955 |
| 8 x 8/pitch 37/clipped | 7.853500 | 4.739500 | 3.757500 | 1.657031 | 0.792805 | 1.654596, 1.686570, 1.699029, 1.607442, 1.634162, 1.654714 |
| 8 x 8/pitch 1312/flat | 7.894500 | 4.736500 | 3.761500 | 1.666737 | 0.794152 | 1.655400, 1.666384, 1.680541, 1.649789, 1.646656, 1.685919 |
| 8 x 8/pitch 1312/mixed | 7.922500 | 4.749000 | 3.754500 | 1.668246 | 0.790587 | 1.687086, 1.686333, 1.661992, 1.642992, 1.631185, 1.690863 |
| 8 x 8/pitch 1312/clipped | 7.936500 | 4.734000 | 3.746000 | 1.676489 | 0.791297 | 1.679574, 1.689017, 1.719220, 1.664140, 1.662060, 1.659770 |
| 16 x 16/pitch 16/flat | 11.730000 | 8.073500 | 5.846000 | 1.452901 | 0.724097 | 1.438627, 1.461279, 1.448780, 1.453597, 1.494428, 1.427878 |
| 16 x 16/pitch 16/mixed | 11.770000 | 8.099000 | 5.830000 | 1.453266 | 0.719842 | 1.457019, 1.466583, 1.445828, 1.437477, 1.462936, 1.492519 |
| 16 x 16/pitch 16/clipped | 11.740000 | 8.145000 | 5.832500 | 1.441375 | 0.716083 | 1.482138, 1.488057, 1.484403, 1.434676, 1.488889, 1.409517 |
| 16 x 16/pitch 37/flat | 11.810000 | 8.086500 | 6.036500 | 1.460459 | 0.746491 | 1.503574, 1.504232, 1.444257, 1.407398, 1.509195, 1.412617 |
| 16 x 16/pitch 37/mixed | 11.680000 | 8.206000 | 6.045000 | 1.423349 | 0.736656 | 1.480869, 1.454097, 1.453104, 1.340609, 1.420002, 1.418173 |
| 16 x 16/pitch 37/clipped | 11.785000 | 8.226000 | 6.043000 | 1.432653 | 0.734622 | 1.501513, 1.459818, 1.467492, 1.409670, 1.374881, 1.387668 |
| 16 x 16/pitch 1312/flat | 11.775000 | 8.124000 | 6.158500 | 1.449409 | 0.758063 | 1.424469, 1.489148, 1.476278, 1.389220, 1.439598, 1.450879 |
| 16 x 16/pitch 1312/mixed | 11.730000 | 8.057500 | 6.134000 | 1.455787 | 0.761278 | 1.458333, 1.482228, 1.425489, 1.424744, 1.485707, 1.411721 |
| 16 x 16/pitch 1312/clipped | 11.860000 | 8.008500 | 6.110500 | 1.480927 | 0.763002 | 1.502778, 1.496479, 1.505821, 1.440238, 1.493523, 1.457640 |


Publication: `raw-evidence.tar.gz` preserves all 18 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
