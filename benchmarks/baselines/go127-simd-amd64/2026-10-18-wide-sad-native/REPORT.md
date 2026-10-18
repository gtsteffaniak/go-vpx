# PR44: native wide SAD and ARM fill performance

Exact tested head `46b89f5ab4c6c469755323a61c77fe3a766b4081`, base `bba909597f1a68e4ceeb94dcca0229c765f10d84`, old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Full native job: INTEL(R) XEON(R) PLATINUM 8573C, go version go1.27.2 linux/amd64. Six balanced rounds, PGO off, GOMAXPROCS=1, GOAMD64=v1. All comparisons below stay within their named CPU and job.

The 45 new 32/64-wide SAD cases measure the changed AMD AVX2 path. The nine new narrow SAD cases, 12 ARM-only fill cases and unchanged SSE8 helpers are controls on AMD. VP9 wide-SAD callers are actual candidates; the measured VP8 narrow-block workload is a control. Source identity does not establish unchanged timing or explain a control regression.

## Native correctness and allocation checks

All 412 cases have six samples per variant and zero timed bytes/allocations: True. All 48 recorded source SHA256 values match exact Git blobs: True. All 22 public/profile non-timing semantic fields match in every dataset: True. Current/parent allocation arrays match exactly: False.

Native validation records 177 top-level PASS events plus 6 required panic-subprobe PASS events, giving 183 targeted PASS events; 548 PASS lines include all other subtests. Active AVX2 SAD ordinary/race selection and actual CPU-off fallback checks: True. There are exactly two expected SAD selected-test skips under actual AVX/AVX2-off and two expected ARM-only token skips. AVX512 is ACTIVE on this Intel8573C host: inherited plane selection passes twice (ordinary and race/checkptr); its source is unchanged in PR44, so these plane timings remain controls. No unexpected skip or failure is present. Every raw validation line is retained.

| Dataset | Allocation samples exactly equal across all variants? | Exceptions |
|---|---|---|
| VP9 320x180 | True | [] |
| VP8 720p | True | [] |
| VP9 720p | False | [{"variant": "previous", "round": 6, "allocations_per_input": 1.2666666666666666}] |
| VP8 720p profile | True | [] |
| VP9 720p profile | True | [] |

Allocation exception: parent VP9 720p round 6 records 1.2666666667 allocations/input (76 allocations/60), while current and old reference are 1.25 in all six rounds and the other five parent rounds are 1.25. This prevents claiming exact current/parent allocation arrays; current has no allocation increase. Every sample and the parent outlier remain explicit.

## All 45 changed wide SAD fixtures

Ratios are reference time / current time: above 1 is faster. Paired minima use each balanced round independently. Old valid wide windows execute handwritten SSE2 sad16ChunksSSE2; this is an actual assembly comparison.

| Fixture | Parent/head | Min paired parent/head | Faster pairs | Old SSE2/head | Min paired SSE2/head |
|---|---:|---:|---:|---:|---:|
| 32x16/largepitch/max | 5.003723× | 4.954100× | 6/6 | 2.073224× | 1.993917× |
| 32x16/largepitch/patterned | 5.070159× | 4.992556× | 6/6 | 2.116308× | 1.933254× |
| 32x16/largepitch/zero | 5.069147× | 4.888304× | 6/6 | 2.060465× | 1.990643× |
| 32x16/packed/max | 5.186394× | 4.487737× | 6/6 | 2.010067× | 1.851171× |
| 32x16/packed/patterned | 5.014119× | 4.876970× | 6/6 | 2.045734× | 1.952353× |
| 32x16/packed/zero | 5.037026× | 4.417722× | 6/6 | 2.034294× | 1.855806× |
| 32x16/padded/max | 5.063019× | 4.916157× | 6/6 | 2.041193× | 1.938534× |
| 32x16/padded/patterned | 4.983969× | 4.796529× | 6/6 | 2.002117× | 1.941141× |
| 32x16/padded/zero | 5.062634× | 5.018634× | 6/6 | 2.032238× | 2.003106× |
| 32x32/largepitch/max | 5.931545× | 5.817587× | 6/6 | 2.323219× | 2.197674× |
| 32x32/largepitch/patterned | 6.101630× | 5.790225× | 6/6 | 2.391563× | 2.241884× |
| 32x32/largepitch/zero | 6.000373× | 5.854281× | 6/6 | 2.311497× | 2.207662× |
| 32x32/packed/max | 6.116301× | 6.011257× | 6/6 | 2.401525× | 2.297936× |
| 32x32/packed/patterned | 6.125714× | 5.826830× | 6/6 | 2.405714× | 2.327664× |
| 32x32/packed/zero | 6.126869× | 5.962206× | 6/6 | 2.432733× | 2.326079× |
| 32x32/padded/max | 6.292664× | 5.989641× | 6/6 | 2.452138× | 2.341102× |
| 32x32/padded/patterned | 6.024051× | 5.822454× | 6/6 | 2.346674× | 2.303878× |
| 32x32/padded/zero | 6.112071× | 5.971223× | 6/6 | 2.392823× | 2.349640× |
| 32x64/largepitch/max | 6.665299× | 6.472425× | 6/6 | 2.439975× | 2.367553× |
| 32x64/largepitch/patterned | 6.610221× | 6.305305× | 6/6 | 2.407426× | 2.256273× |
| 32x64/largepitch/zero | 6.744786× | 6.614032× | 6/6 | 2.459219× | 2.269915× |
| 32x64/packed/max | 6.810991× | 6.410760× | 6/6 | 2.518734× | 2.401335× |
| 32x64/packed/patterned | 6.652694× | 6.244131× | 6/6 | 2.457826× | 2.234742× |
| 32x64/packed/zero | 6.780564× | 6.512662× | 6/6 | 2.486938× | 2.271409× |
| 32x64/padded/max | 6.680659× | 6.290873× | 6/6 | 2.466469× | 2.418623× |
| 32x64/padded/patterned | 6.540453× | 6.295069× | 6/6 | 2.432240× | 2.280257× |
| 32x64/padded/zero | 6.731553× | 6.233507× | 6/6 | 2.470444× | 2.395368× |
| 64x32/largepitch/max | 5.797591× | 5.373880× | 6/6 | 2.510413× | 2.396713× |
| 64x32/largepitch/patterned | 5.704342× | 5.518878× | 6/6 | 2.548890× | 2.377562× |
| 64x32/largepitch/zero | 5.813591× | 5.408844× | 6/6 | 2.520298× | 2.390488× |
| 64x32/packed/max | 5.934240× | 5.401114× | 6/6 | 2.540816× | 2.435458× |
| 64x32/packed/patterned | 5.691700× | 5.244987× | 6/6 | 2.452789× | 2.343227× |
| 64x32/packed/zero | 5.698925× | 5.385917× | 6/6 | 2.474215× | 2.347943× |
| 64x32/padded/max | 5.730782× | 5.440457× | 6/6 | 2.450236× | 2.220636× |
| 64x32/padded/patterned | 5.966125× | 5.595316× | 6/6 | 2.497984× | 2.449402× |
| 64x32/padded/zero | 5.933817× | 5.631839× | 6/6 | 2.561197× | 2.409354× |
| 64x64/largepitch/max | 5.855164× | 5.622689× | 6/6 | 3.141809× | 2.996238× |
| 64x64/largepitch/patterned | 5.958693× | 5.448552× | 6/6 | 3.172565× | 2.977735× |
| 64x64/largepitch/zero | 5.829768× | 5.576267× | 6/6 | 3.134035× | 2.831373× |
| 64x64/packed/max | 5.971078× | 5.697264× | 6/6 | 3.125540× | 2.817741× |
| 64x64/packed/patterned | 5.878686× | 5.582444× | 6/6 | 2.998788× | 2.676567× |
| 64x64/packed/zero | 5.909406× | 5.812536× | 6/6 | 2.964801× | 2.628169× |
| 64x64/padded/max | 5.871192× | 5.386905× | 6/6 | 3.183077× | 2.523384× |
| 64x64/padded/patterned | 6.134039× | 5.773762× | 6/6 | 3.272462× | 2.892718× |
| 64x64/padded/zero | 6.154358× | 5.605552× | 6/6 | 3.180829× | 2.748727× |

## All nine narrow SAD controls

| Fixture | Parent/head | Min paired parent/head | Slower pairs | Old SSE2/head |
|---|---:|---:|---:|---:|
| 16x16/largepitch/patterned | 0.998260× | 0.930617× | 4/6 | 1.020884× |
| 16x16/packed/patterned | 0.971262× | 0.939041× | 4/6 | 1.030106× |
| 16x16/padded/patterned | 0.957288× | 0.932615× | 6/6 | 0.995932× |
| 8x16/largepitch/patterned | 0.990351× | 0.971777× | 4/6 | 0.965584× |
| 8x16/packed/patterned | 1.006167× | 0.975564× | 4/6 | 0.980201× |
| 8x16/padded/patterned | 1.007131× | 0.984486× | 3/6 | 0.972123× |
| 8x8/largepitch/patterned | 1.006320× | 0.958169× | 3/6 | 1.119369× |
| 8x8/packed/patterned | 1.015834× | 0.955762× | 3/6 | 1.112916× |
| 8x8/padded/patterned | 1.011300× | 0.963255× | 3/6 | 1.120315× |

## Whole-codec measurements

| Workload | Current time vs parent | Median paired time change | Faster pairs | Current time vs old reference |
|---|---:|---:|---:|---:|
| VP9 320x180 | -1.701597% | -1.835431% | 5/6 | -44.161957% |
| VP8 720p | -2.287169% | +0.147410% | 3/6 | -28.722667% |
| VP9 720p | -4.505876% | -4.173051% | 6/6 | -54.397432% |

Percentages compare candidate time with reference time. A percentage relative to old time is not the inverse percentage for old speed relative to candidate. Ratio of medians and median paired change differ; both are retained. Independent profiles establish sampled liveness rather than additive kernel or codec gains.

## Every unchanged SSE8 public control

| Fixture | Parent/head | Min paired parent/head | Slower pairs | Old reference/head |
|---|---:|---:|---:|---:|
| pointer8x8/pitch1312-1320/pattern0 | 1.001282× | 0.975148× | 3/6 | 1.111562× |
| pointer8x8/pitch1312-1320/pattern1 | 0.996392× | 0.971013× | 5/6 | 1.109162× |
| pointer8x8/pitch1312-1320/pattern2 | 1.004702× | 0.976690× | 3/6 | 1.110858× |
| pointer8x8/pitch32-37/pattern0 | 1.007268× | 0.970785× | 4/6 | 1.142343× |
| pointer8x8/pitch32-37/pattern1 | 0.981198× | 0.946164× | 3/6 | 1.088628× |
| pointer8x8/pitch32-37/pattern2 | 0.999568× | 0.957083× | 1/6 | 1.119749× |
| pointer8x8/pitch8-8/pattern0 | 1.005204× | 0.988143× | 2/6 | 1.107321× |
| pointer8x8/pitch8-8/pattern1 | 1.018242× | 0.974598× | 3/6 | 1.089876× |
| pointer8x8/pitch8-8/pattern2 | 1.000561× | 0.967361× | 4/6 | 1.120241× |
| slice8x16/pitch1312-1320/pattern0 | 0.973966× | 0.958238× | 5/6 | 96.248086× |
| slice8x16/pitch1312-1320/pattern1 | 0.985910× | 0.943524× | 4/6 | 96.001523× |
| slice8x16/pitch1312-1320/pattern2 | 0.988799× | 0.943326× | 3/6 | 98.802626× |
| slice8x16/pitch32-37/pattern0 | 0.997295× | 0.955588× | 4/6 | 96.522411× |
| slice8x16/pitch32-37/pattern1 | 0.985254× | 0.974164× | 6/6 | 99.107489× |
| slice8x16/pitch32-37/pattern2 | 0.984674× | 0.948012× | 4/6 | 96.321839× |
| slice8x16/pitch8-8/pattern0 | 0.987038× | 0.943946× | 4/6 | 95.043843× |
| slice8x16/pitch8-8/pattern1 | 0.976282× | 0.947917× | 4/6 | 97.016067× |
| slice8x16/pitch8-8/pattern2 | 0.971648× | 0.940609× | 5/6 | 96.360153× |
| slice8x4/pitch1312-1320/pattern0 | 0.975457× | 0.951324× | 6/6 | 1.185274× |
| slice8x4/pitch1312-1320/pattern1 | 0.974454× | 0.950310× | 6/6 | 1.213676× |
| slice8x4/pitch1312-1320/pattern2 | 0.948834× | 0.919060× | 6/6 | 1.167010× |
| slice8x4/pitch32-37/pattern0 | 0.995150× | 0.928796× | 3/6 | 1.212167× |
| slice8x4/pitch32-37/pattern1 | 0.964044× | 0.917559× | 5/6 | 1.224740× |
| slice8x4/pitch32-37/pattern2 | 0.960874× | 0.935097× | 5/6 | 1.190803× |
| slice8x4/pitch8-8/pattern0 | 0.972873× | 0.940795× | 5/6 | 1.175697× |
| slice8x4/pitch8-8/pattern1 | 0.977988× | 0.911293× | 5/6 | 1.189079× |
| slice8x4/pitch8-8/pattern2 | 0.982483× | 0.924675× | 4/6 | 1.180373× |
| slice8x8/pitch1312-1320/pattern0 | 0.972475× | 0.951366× | 5/6 | 1.261943× |
| slice8x8/pitch1312-1320/pattern1 | 0.992040× | 0.935578× | 3/6 | 1.259879× |
| slice8x8/pitch1312-1320/pattern2 | 0.988803× | 0.939014× | 3/6 | 1.230576× |
| slice8x8/pitch32-37/pattern0 | 0.976047× | 0.920057× | 6/6 | 1.231224× |
| slice8x8/pitch32-37/pattern1 | 0.994731× | 0.952828× | 3/6 | 1.244058× |
| slice8x8/pitch32-37/pattern2 | 0.986087× | 0.947345× | 4/6 | 1.238501× |
| slice8x8/pitch8-8/pattern0 | 0.965462× | 0.932370× | 6/6 | 1.245298× |
| slice8x8/pitch8-8/pattern1 | 0.969365× | 0.938176× | 4/6 | 1.231630× |
| slice8x8/pitch8-8/pattern2 | 0.968068× | 0.934258× | 4/6 | 1.246668× |

The ~95–99× old-reference/head 16-row SSE8 ratios are specific to this Intel8573C live B.Loop call protocol and inherited SSE8 implementation. They are not a new PR44 gain or a codec-wide ratio; this audit establishes no transition or layout cause. Old SSE8 uses actual SSE2 for 4/8/pointer and AVX2 for 16 rows. Old AMD subtraction is scalar for all 48 subtraction fixtures; old copy helpers are compiler-generated Go, so the label old reference does not imply handwritten assembly for every fixture.

All 219 median regressions and 42 all-six-slower kernel controls remain in audit.json with every sample. No control regression or outlier is pruned. Twelve fill control cases and every legacy B.N diagnostic are preserved separately. The independent focused EPYC9V45 job is frozen in /private/tmp/govpx-pr44-focused-evidence-ready/ (preserve the focused bundle alongside the full archive when copying to a repository baseline). Its 36 live SSE8 controls include 30 slower medians and six all-six-slower cases; pointer time is +2.376–4.996%. No SAD was timed there; do not pool focused and full CPUs/jobs or infer a cause from unchanged source.

## ARM fill scope

Fresh M4 Max native 8x8 FillBlock medians improve 1.59083–1.61632× against the already-optimized fixed-store parent; every pair is at least 1.50616×. Actual DSP n16 packed/no-neighbor improves 1.6274×. Other 4x4/16x16 component gains are below 1.5×, and one-row 8-byte controls are slower (~0.955×). These are incremental comparisons against an existing fixed-store winner; gains against different parents are not multiplied. ARM whole VP8 median time is −0.091934% versus the same strong locked parent, mixed 4/6 faster, and candidate time remains +8.670979% versus old reference in that run. Old NEON 8x8/16x16 fill is not unrolled; the earlier hypothesis was withdrawn. Actual DSP n4 performance was not measured. Full local raw reports/proofs and the corrected provenance are committed in the PR.

## Optimized libvpx, separate CPU/job

CPU: AMD EPYC 7763 64-Core Processor. Pinned libvpx v1.16.0 / 1024874c5919305883187e2953de8fcb4c3d7fa6, optimized assembly/SSE2/AVX2 with runtime dispatch. Reference timers use vpxenc-stats. Fixed speed means CPU-used −8 with minq=2/maxq=56, not a fixed quantizer. Matching measured PSNR/SSIM/payload/frame/drop counts does not establish bitstream identity.

| Workload | Go time / libvpx time | Measured quality/payload/counts match? | Quality-output matches timing? |
|---|---:|---|---|
| vp9-320p | 2.895904× | False | {"govpx": true, "libvpx": true} |
| vp9-720p | 2.910286× | False | {"govpx": true, "libvpx": true} |
| vp8-720p-adaptive | 1.806403× | False | {"govpx": true, "libvpx": true} |
| vp8-720p-fixed | 2.369923× | True | {"govpx": true, "libvpx": true} |

libvpx is faster at these ratios. Adaptive VP8 and VP9 rows perform unequal quality/output work and retain that limitation. Every actual encoder option, timer, timing/quality output sample and separate-CPU identity remains in audit.json and raw logs. No cross-job causal attribution or broad libvpx parity claim is made.

## Preservation

native-raw.tar.gz and libvpx-raw.tar.gz preserve every raw byte; raw-sha256.json covers all members. audit.json preserves all 412 kernel datasets, whole-codec outputs, per-round semantic/allocation checks, all regressions, source hashes, feature-selection tests, profiles and C controls. Parent independent C/full proof JSON and scripts are included; focused raw and independent proof are in the linked separate frozen bundle. scope.txt records exact merge receipt when authorized; sha256.json covers every outer file. No source, README, checkout, index or local branch was modified by this audit.


The complete frozen evidence package is inside `raw-evidence.tar.gz` (18 original files). It preserves all nested raw archives, manifests, reports and auditors byte for byte. Extract to another directory to inspect it; `SHA256SUMS` covers the published outer files.
