# Repaired PR42: native SSE8 and standalone ARM token kernel

Exact `f8dd8c7472e91c1f147511b04ceec716aca95aa7`, base `cf3117194f0590fa0af35951a32e8fadb6da9ab3`. Root independent focused/full review declares PERFREADY; all four checks pass. Merged through guarded normal GitHub --merge; exact-tree and ancestry proof below. The rejectedc638 head and every original regression remain in their separate committed archive.

Full native CPU AMD EPYC 9V45 96-Core Processor; Go1.27.2/PGOoff/GOMAXPROCS1. Independent focused CPU AMD EPYC 9V74 80-Core Processor. Optimized C runs on a separate AMD EPYC 7763 64-Core Processor. Compare each dataset only within its own CPU/job; no pooling or additive/cross-job attribution.

## Full native: all36 live public SSE8 fixtures

All36 medians improve1.001986–1.778815×parent;27cases improve in all six pairs. The nine trusted-pointer cases qualify at1.688280–1.778815×median with every paired gain>=1.646377×. Smaller slice8x4 gains are mixed and include slower pairs. All36 medians beat actual old SSE2/AVX2 reference kernels.

| Fixture | Parent/head | Minimum paired gain | Pairs faster | Old ASM/head |
|---|---:|---:|---:|---:|
| pointer8x8/pitch1312-1320/pattern0 | 1.778815× | 1.723325× | 6/6 | 1.164465× |
| pointer8x8/pitch1312-1320/pattern1 | 1.766099× | 1.692677× | 6/6 | 1.193984× |
| pointer8x8/pitch1312-1320/pattern2 | 1.755032× | 1.688147× | 6/6 | 1.170249× |
| pointer8x8/pitch32-37/pattern0 | 1.765845× | 1.724501× | 6/6 | 1.186590× |
| pointer8x8/pitch32-37/pattern1 | 1.764828× | 1.686405× | 6/6 | 1.175055× |
| pointer8x8/pitch32-37/pattern2 | 1.768634× | 1.701370× | 6/6 | 1.173724× |
| pointer8x8/pitch8-8/pattern0 | 1.688280× | 1.646377× | 6/6 | 1.187100× |
| pointer8x8/pitch8-8/pattern1 | 1.758574× | 1.658810× | 6/6 | 1.172496× |
| pointer8x8/pitch8-8/pattern2 | 1.725006× | 1.685663× | 6/6 | 1.180255× |
| slice8x16/pitch1312-1320/pattern0 | 1.262102× | 1.211879× | 6/6 | 1.650531× |
| slice8x16/pitch1312-1320/pattern1 | 1.244696× | 1.226350× | 6/6 | 1.654880× |
| slice8x16/pitch1312-1320/pattern2 | 1.239636× | 1.121435× | 6/6 | 1.647787× |
| slice8x16/pitch32-37/pattern0 | 1.234632× | 1.223097× | 6/6 | 1.553584× |
| slice8x16/pitch32-37/pattern1 | 1.252772× | 1.238490× | 6/6 | 1.562609× |
| slice8x16/pitch32-37/pattern2 | 1.228858× | 1.199043× | 6/6 | 1.575295× |
| slice8x16/pitch8-8/pattern0 | 1.248612× | 1.242330× | 6/6 | 1.565878× |
| slice8x16/pitch8-8/pattern1 | 1.251322× | 1.247642× | 6/6 | 1.582814× |
| slice8x16/pitch8-8/pattern2 | 1.262638× | 1.243752× | 6/6 | 1.587115× |
| slice8x4/pitch1312-1320/pattern0 | 1.010543× | 0.889254× | 4/6 | 1.358461× |
| slice8x4/pitch1312-1320/pattern1 | 1.015355× | 0.975171× | 5/6 | 1.378839× |
| slice8x4/pitch1312-1320/pattern2 | 1.009772× | 0.938409× | 4/6 | 1.345670× |
| slice8x4/pitch32-37/pattern0 | 1.020675× | 0.999577× | 5/6 | 1.384810× |
| slice8x4/pitch32-37/pattern1 | 1.015483× | 0.996268× | 4/6 | 1.402781× |
| slice8x4/pitch32-37/pattern2 | 1.001986× | 0.970418× | 4/6 | 1.357546× |
| slice8x4/pitch8-8/pattern0 | 1.005368× | 0.987828× | 4/6 | 1.378381× |
| slice8x4/pitch8-8/pattern1 | 1.002837× | 0.981696× | 3/6 | 1.377154× |
| slice8x4/pitch8-8/pattern2 | 1.011486× | 0.957965× | 4/6 | 1.397893× |
| slice8x8/pitch1312-1320/pattern0 | 1.255556× | 1.115802× | 6/6 | 1.398264× |
| slice8x8/pitch1312-1320/pattern1 | 1.260870× | 1.249248× | 6/6 | 1.415283× |
| slice8x8/pitch1312-1320/pattern2 | 1.256848× | 1.239286× | 6/6 | 1.378490× |
| slice8x8/pitch32-37/pattern0 | 1.258829× | 1.206311× | 6/6 | 1.388199× |
| slice8x8/pitch32-37/pattern1 | 1.249429× | 1.216635× | 6/6 | 1.401458× |
| slice8x8/pitch32-37/pattern2 | 1.246816× | 1.188861× | 6/6 | 1.384042× |
| slice8x8/pitch8-8/pattern0 | 1.249847× | 1.184395× | 6/6 | 1.378845× |
| slice8x8/pitch8-8/pattern1 | 1.241421× | 1.099437× | 6/6 | 1.376938× |
| slice8x8/pitch8-8/pattern2 | 1.265189× | 1.204060× | 6/6 | 1.387673× |

The original B.N SSE8x8 fixture is a separate diagnostic:1.204022×parent/1.305477×old ASM, all six parent pairs faster. It is not pooled with live B.Loop acceptance.

## Independent focused job

[Focused workflow](https://github.com/thesyncim/govpx/actions/runs/39089608484) / [artifact12121679836](https://github.com/thesyncim/govpx/actions/runs/39089608484/artifacts/12121679836). All720 timed samples and120 excluded one-shot samples remain preserved. All36 public medians improve1.010991–1.482594×parent;31cases improve in every pair. Pointer medians1.479151–1.482594×are below1.5×onthisCPU; slice8x4 is mixed. All36 medians beat oldASM as well. Four legacy B.N diagnostics stay separate. Source hashes match all622 archived Git blobs, and actual AVX2 dispatch/new ordinary/race tests pass.

Frozen focused bundle is `govpx-pr42-repaired-focused-evidence-ready/`, provided alongside this bundle; its manifest and independent audit/proof are referenced inside audit.json. It contains all45 focused artifact files (including original44-entrySHA256SUMS), every36+4case ratio, root proof and outer/raw hashes. These45 files are additional to the160 fullnative/C files below.

## Whole encoder results

| Public case | Time vs parent | Pairs faster | Time vs old build |
|---|---:|---:|---:|
| VP9 320x180 | +0.337444% | 3/6 | -46.267794% |
| VP8 720p | -0.536194% | 3/6 | -5.067986% |
| VP9 720p | -1.530390% | 4/6 | -54.025554% |

VP8 may execute the changed SSE8 paths; VP9 is unchanged control work. Paired signs remain mixed, so no stable whole-encoder win or universal assembly-gap closure is claimed.

## Validation, allocations and regressions

All338 cases × six samples × three variants have zero bytes/allocations. All36 CI-recorded source hashes match exact Git blobs. Native logs record141top-levelPASS plus six retained panic probes (147passing targeted checks),292PASS lines including subtests. All four new SSE8 tests pass ordinary and race/checkptr, with AVX2=true/Enabled=true. The inherited AVX512 test executes and passes twice on9V45. The ARM-only token SIMD test correctly skips twice with ARM SIMD required; its scalar, stream and unsupported-path availability tests pass twice. No unexpected skip/failure is hidden.

All22semantic fields agree in three public and two profile datasets. Every allocation sample matches, including nonzero controls: VP8zero,VP9320 .616667,VP9720 1.25,VP9profile .625. No allocation/timing outlier is removed.

All 117 median regressions and 8 all-six-slower unchanged SAD controls remain in audit.json, along with every slower SSE8 pair. All-six-slower controls:

| Control | Time vs parent |
|---|---:|
| SAD16x16LimitEarly | +7.150313% |
| SADLimitSliceNative/pitch1312/extremefalse/early | +6.893653% |
| SADLimitSliceNative/pitch1312/extremetrue/early | +5.916748% |
| SADLimitSliceNative/pitch1312/extremetrue/middle | +2.833438% |
| SADLimitSliceNative/pitch32/extremefalse/bound | +3.801023% |
| SADLimitSliceNative/pitch32/extremefalse/early | +6.853108% |
| SADLimitSliceNative/pitch32/extremefalse/wide | +3.773585% |
| SADLimitSliceNative/pitch32/extremetrue/middle | +4.260870% |

## Standalone ARM token component

The retained M4Max/Go1.27.2/PGOoff standalone token record-base kernel improves2.583–2.594×versus the scalar LUT/store, every pair>2.53×and zero allocations. Default encoder integration is deliberately unwired because both integration trials regress; the component is not credited to whole-codec results. ARM SSE8 restores its original body. Root correctness, race/checkptr, PGO, noasm and build-selection evidence remains in the submitted reports.

## Optimized libvpx

| Case | Go/libvpx time | Scope |
|---|---:|---|
| vp9-320p | 2.932213× | Unequal measured quality/output work |
| vp9-720p | 2.948447× | Unequal measured quality/output work |
| vp8-720p-adaptive | 1.714338× | Unequal measured quality/output work |
| vp8-720p-fixed | 2.193076× | Matched measured PSNR/SSIM/payload/frame/drop counts |

Libvpx remains faster. Fixed speed means CPU-used=−8 with minq2/maxq56, not fixed quantizer. Independent C proof confirms PSNR27.1576108715, SSIM.9259273287 and58qualityframes. Fixed matching measured quantities does not establish bitstream identity. Every quality-pass output statistic matches timing in this job. Assembly/SSE2/AVX2/runtime detection and vpxenc-stats timers are verified for pinned libvpxv1.16.0/1024874c5919305883187e2953de8fcb4c3d7fa6. No cross-job attribution.

## Preservation

The two raw archives retain all86native+74Cfiles byte for byte. Raw member and outer hashes are verified; independent root full-proof bytes are preserved separately and both independent proofs are in audit.json. Complete allocations, semantic comparisons, pair metrics, regressions, profiles, runtime selection and source/CI identities remain reviewable. The historical old-reference harness label assembly does not mean every unrelated function is assembly: all48AMD subtraction fixtures were scalar fallback. Old SSE8 fixture kernels above were actual SSE2/AVX2. No source/README/localcheckout/index/branch was changed by this audit.

## Merge proof

Merged `a3fb408d1323066cf9ed8470c7b69a8ff9b7e520` through guarded normal GitHub --merge. Tested and merged tree both `af1d126948659f4fe436a6fbb93f326f5a57c91b`; parents exactly `[cf3117194f0590fa0af35951a32e8fadb6da9ab3, f8dd8c7472e91c1f147511b04ceec716aca95aa7]`. Remote main was verified, and measured head ancestry is preserved. The source remote remains exactf8 for the root successor; no checkout/index/local branch was changed. Full proof, four passing premerge checks and byte-recoverable independent C proof are embedded in audit.json.
