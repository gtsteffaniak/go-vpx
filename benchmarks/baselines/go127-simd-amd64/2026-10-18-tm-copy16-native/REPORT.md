# PR46: native Copy16 and ARM TM evidence

Tested head `46a08f8e6eb03f6ee4549d975384bd6c952a81fa`, base `78ce74d39806d74ecf91bc3d082899a74d7bc85f`, old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. This full native job used **INTEL(R) XEON(R) PLATINUM 8573C**, Go1.27.2 Linux/amd64, PGOoff, GOMAXPROCS1, six balanced150ms rounds. Workflow [39101922555](https://github.com/thesyncim/govpx/actions/runs/39101922555).

All77 exact Git source hashes and31 identical common fixtures match. All498×6×3=8964 timed kernel samples have0B/op and0allocs/op. All338 one-shot correctness/allocation samples are preserved and excluded:81DC+144FP+81TM+32Copy16. All22 semantic fields across three public workloads and two profiles match; EVERY allocation sample is exactly equal across candidate,parent,old reference on this run. Public allocations: VP8whole/profile0; VP9320whole37/60; VP9720whole75/60; VP9profile75/120. No allocation outlier occurred on this run; prior exceptions are retained in their earlier archives.

Actual AVX512 is active and TestSSEPlane16AVX512 passes ordinary+race twice. FP probes pass active twice/off twice. Expected skips are ARM-only TM arithmetic4, ARM-only standalone token2, CPU-off SADselected2. NewCopy16 probes each pass4; root/DSP partialpanic8, inherited publicCopy16state4, TMcommonstate7. DCstate has two named verbose PASS events plus three nonverbose reference-preflight completions, not five named PASS lines. Original logs remain untouched.

## Actual Copy16 candidates

The pinned old Copy16 reference is compiler-generated Go, not handwritten assembly. Root and DSP calls are timed separately in every round. Seven of eight cases in each caller are faster in all six parent pairs. All six disjoint root cases exceed1.5× in every pair (minimum1.666667×); the pitched self-alias root case has median1.547868× but minimum1.432177×, so it does not meet an all-pairs1.5× claim. All seven winning DSP cases exceed1.5× in every pair.

**Both zero-stride self-alias cases regress in all six pairs.** Root parent/candidate0.509592× and DSP0.441570× (current time+96.236%/+126.465% versus parent). Exact self-copy semantics are unchanged; the branch is preserved for a follow-up fast no-op repair. Other wins do not erase these results. No transition/layout cause is established.

| Caller / geometry | Parent/candidate | Min paired parent/candidate | Faster parent pairs | Old compiler-Go/candidate | Min paired old/candidate |
|---|---:|---:|---:|---:|---:|
| root pitch0_0/aliasfalse | 1.704121× | 1.666667× | 6/6 | 1.517885× | 1.490061× |
| root pitch0_0/aliastrue | 0.509592× | 0.468595× | 0/6 | 0.445648× | 0.436580× |
| root pitch1312_1320/aliasfalse | 1.707488× | 1.670741× | 6/6 | 1.526911× | 1.491979× |
| root pitch16_16/aliasfalse | 1.775304× | 1.729772× | 6/6 | 1.588259× | 1.548311× |
| root pitch16_16/aliastrue | 1.547868× | 1.432177× | 6/6 | 1.375302× | 1.325599× |
| root pitch17_31/aliasfalse | 1.733307× | 1.708949× | 6/6 | 1.545239× | 1.520623× |
| root pitch1_7/aliasfalse | 1.718603× | 1.684615× | 6/6 | 1.534929× | 1.502308× |
| root pitch32_37/aliasfalse | 1.692843× | 1.681260× | 6/6 | 1.510251× | 1.490783× |
| DSP pitch0_0/aliasfalse | 2.408741× | 2.399657× | 6/6 | 2.699292× | 2.686001× |
| DSP pitch0_0/aliastrue | 0.441570× | 0.437807× | 0/6 | 0.518523× | 0.508868× |
| DSP pitch1312_1320/aliasfalse | 2.411295× | 2.389175× | 6/6 | 2.700430× | 2.681489× |
| DSP pitch16_16/aliasfalse | 2.398945× | 2.379209× | 6/6 | 2.692756× | 2.674475× |
| DSP pitch16_16/aliastrue | 2.080442× | 2.064170× | 6/6 | 2.465459× | 2.399658× |
| DSP pitch17_31/aliasfalse | 2.282577× | 2.274322× | 6/6 | 2.561182× | 2.548928× |
| DSP pitch1_7/aliasfalse | 2.398700× | 2.343938× | 6/6 | 2.689266× | 2.642805× |
| DSP pitch32_37/aliasfalse | 2.443129× | 2.389999× | 6/6 | 2.736055× | 2.686239× |

## Whole codecs and controls

| Public workload | Candidate time vs parent | Faster parent pairs | Candidate time vs old | Faster old pairs |
|---|---:|---:|---:|---:|
| VP9 320x180 | +0.142172% | 2/6 | -44.320866% | 6/6 |
| VP8 720p | -0.041056% | 2/6 | -28.414515% | 6/6 |
| VP9 720p | +0.663959% | 3/6 | -54.869765% | 6/6 |

VP8 is a candidate whole-codec result: effectively flat median with only2/6 faster pairs. VP9 is an unchanged-source control; its small increases are retained. No broad whole-codec improvement or cross-CPU explanation is inferred. Among482 other kernels (including27ARM-TM-disabled controls),319 medians increase and41 are slower in every pair; every sample and every regression is retained in audit.json. All498 kernels, both self-copy regressions and all unchanged control regressions remain reviewable.

The separate focused SSE8 job used EPYC7763, not this full Intel8573C job. All720timed+120excluded samples/source/dispatch checks pass and all timed allocations are0. Its20/36 public medians increase, five cases are all-six slower; all9pointer medians increase3.471–6.534%, four are all-six slower. Legacy B.Npointer is+6.072% all-six slower. No unchanged-source regression is dismissed as noise or assigned a cause. Full focused details/raw are in the companion bundle.

## Separate optimized libvpx reference

Optimized C runner **AMD EPYC 9V74 80-Core Processor** uses pinnedlibvpxv1.16.0 / `1024874c5919305883187e2953de8fcb4c3d7fa6`, assembly/SSE2/AVX2/runtime dispatch enabled. All24reference timers arevpxenc-stats. Allquality-output statistics match their corresponding timed outputs on every run.

| Workload | Go/libvpx TIME ratio | Scope |
|---|---:|---|
| vp8-720p-fixed | 2.282511644× | Matched measured PSNR/SSIM/payload/counts and58qualityframes |
| vp8-720p-adaptive | 1.661755936× | Different output/quality work |
| vp9-320p | 2.664478208× | Different output/quality work |
| vp9-720p | 2.576795879× | Different output/quality work |

A TIME ratio above one means libvpx is faster. Fixed VP8 denotesCPU-used−8 withmin-q2/max-q56, not a fixed quantizer. Matching measured output statistics/quality does not establish bitstream identity. Adaptive VP8/VP9 use different output and quality work; their ratios cannot support equivalent-work claims. Separate jobs are never pooled, even when CPU model names happen to match.

## Evidence and acceptance scope

Raw archives preserve91native+74C files (165); the companion focused archive preserves45files including its originalSHA256SUMS. Every decompressed raw member matches its savedSHA256 manifest and original downloaded bytes. Independent root proofs/scripts, separate independent audits, source77SHA proof, common31 inventory, all338preflights, allregressions and allallocation samples are included. Root parser initially expected pitch1312_1312; actual frozen fixture/raw are1312_1320. This analysis-only selector correction is recorded separately; raw files were not changed.

ARM TM gains are established only by committed nativeM4 reports, not the27AMDcontrol timings. PR46 retains qualifying component gains under explicit user policy while documenting zero-stride self-copy and unchanged-control regressions. Readiness/merge require explicit rootPERFREADY, published final PRbody, fresh exacthead/base and allfourCI successes. Normal guarded GitHub merge completed at2026-10-18T06:32:50Z: `2df3eaf5fb547db07bd970cd65bf504ecb287590`. Parents are `78ce74d39806d74ecf91bc3d082899a74d7bc85f` and `46a08f8e6eb03f6ee4549d975384bd6c952a81fa`; merged/main tree exactly equals tested tree`3863fe98bbb0f1ab2de14a49b8121e78b9358087`. Explicit rootPERFREADY accepted documented self-copy/control regressions while retaining disjoint kernel wins. Allfour fresh checks wereSUCCESS and published finalbody was byte-verified beforemerge. Source branch was preserved for successor coordination. No production/README/index/checkout/localbranch mutation occurred during this audit.


Publication: `raw-evidence.tar.gz` preserves all 28 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
