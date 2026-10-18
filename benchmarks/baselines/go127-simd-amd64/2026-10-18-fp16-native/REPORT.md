# PR45: native FP16 and ARM fused DC evidence

Tested head `ac259fce071f3ac5fc45603c6bd7e7926bdfcf73`, base `2ca5bab1110426401e4f9dbfd6fa0da8f568c717`, old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Full native and focused jobs each use AMD EPYC7763; C uses a separate EPYC9V74 job. Equal model names do not make separate jobs interchangeable. All comparisons remain scoped to their job, compiler, fixture and options.

Native FP n256 qualifies under the user’s component rule: all12 medians are 2.323684–2.765502×parent and all72 paired ratios are at least2.2576×. All12 also beat old reference in every pair (median1.087658–13.261730×; minimum1.085191×). WholeVP9 is mixed/small; these kernel gains do not establish a universal codec or assembly win.

## Correctness, dispatch and allocation controls

All91 native rawfiles and63 exact Git source hashes are verified. All463 cases ×6rounds ×3variants =8,334 timed kernel samples have0B/op and0allocs/op. All22 non-timing semantic fields match across three public workloads and both profile workloads. Preserve81DC and144FP one-shot preflights separately; these225 samples are excluded from timed statistics.

Four TestFP16Selected PASS events prove ordinary/race AVX2 active and actualGODEBUG AVX/AVX2 off. This selected test passes under CPU-off rather than skipping. InheritedAVX512 is inactive on this host and its two feature-test skips are expected. ARM-only token tests and actualCPU-off wideSAD skips are retained. ARM DC helpers are constant-false on AMD; public fallback/alias/partial-panic tests run ordinary/race. The DC supplement state test has two verbose named current PASS events plus three nonverbose reference-preflight completions; do not report five named PASS events.

**Allocation exceptions remain explicit.** Old-reference VP9 720p round4 records76/60=1.2666666667 allocations/input; other17whole samples are75/60=1.25. CURRENT VP9 profile records76/120=0.6333333333, compared with parent/old75/120=0.625, an observed extra allocation. All whole-current VP9 720p samples remain1.25 and all timed kernels remain zeroalloc. This archive does not claim identical profile allocations or dismiss the difference as noise.

## All36 new public FP fixtures

The24 n64/n256 cases are eligible for the changed16-lane AVX2 path;12 n16 cases exercise the prior eight-lane path. Exact-LEN/CAP outputs are aligned32 or offset2; classes use aligned32 or offset1. Actual immutableScanOrder alignment is logged, not forced. Independent scalarEOB/q/dq/classes, full backing/canaries, immutable coefficient/ScanOrder state and observable sinks are checked before and after timing. All3 pinned implementations receive identical fixture bytes.

| Fixture | Parent/head median | Min paired parent/head | Faster pairs | Old reference/head median | Min paired old/head |
|---|---:|---:|---:|---:|---:|
| n16/no-classes/zero/aligned32 | 1.018093× | 1.013897× | 6/6 | 0.736171× | 0.733590× |
| n16/no-classes/zero/offset2-classes1 | 1.017961× | 1.016919× | 6/6 | 0.755232× | 0.753339× |
| n16/no-classes/sparse/aligned32 | 1.017515× | 1.015704× | 6/6 | 0.789372× | 0.787348× |
| n16/no-classes/sparse/offset2-classes1 | 1.018697× | 1.016593× | 6/6 | 0.817480× | 0.815727× |
| n16/no-classes/dense/aligned32 | 1.018854× | 1.016311× | 6/6 | 1.030433× | 1.026082× |
| n16/no-classes/dense/offset2-classes1 | 1.017950× | 1.016316× | 6/6 | 1.035603× | 1.031157× |
| n16/classes/zero/aligned32 | 1.008866× | 1.006774× | 6/6 | 1.935332× | 1.929651× |
| n16/classes/zero/offset2-classes1 | 1.008866× | 1.007817× | 6/6 | 1.935593× | 1.932986× |
| n16/classes/sparse/aligned32 | 1.010169× | 1.008353× | 6/6 | 1.993742× | 1.992703× |
| n16/classes/sparse/offset2-classes1 | 1.008998× | 1.007290× | 6/6 | 1.993349× | 1.990350× |
| n16/classes/dense/aligned32 | 1.008732× | 1.000775× | 6/6 | 2.485860× | 2.464360× |
| n16/classes/dense/offset2-classes1 | 1.009254× | 0.998969× | 5/6 | 2.490485× | 2.470876× |
| n64/no-classes/zero/aligned32 | 1.355886× | 1.352507× | 6/6 | 0.826238× | 0.825114× |
| n64/no-classes/zero/offset2-classes1 | 1.292029× | 1.289054× | 6/6 | 0.760766× | 0.758535× |
| n64/no-classes/sparse/aligned32 | 1.357539× | 1.354130× | 6/6 | 0.846614× | 0.844280× |
| n64/no-classes/sparse/offset2-classes1 | 1.291658× | 1.278484× | 6/6 | 0.775823× | 0.767623× |
| n64/no-classes/dense/aligned32 | 1.356304× | 1.354194× | 6/6 | 1.071630× | 1.068726× |
| n64/no-classes/dense/offset2-classes1 | 1.293612× | 1.289376× | 6/6 | 0.941195× | 0.936957× |
| n64/classes/zero/aligned32 | 1.378830× | 1.376879× | 6/6 | 4.192201× | 4.180739× |
| n64/classes/zero/offset2-classes1 | 1.389210× | 1.382911× | 6/6 | 4.208512× | 4.195830× |
| n64/classes/sparse/aligned32 | 1.381740× | 1.374306× | 6/6 | 4.333395× | 4.318771× |
| n64/classes/sparse/offset2-classes1 | 1.388292× | 1.383038× | 6/6 | 4.354402× | 4.340854× |
| n64/classes/dense/aligned32 | 1.379026× | 1.373035× | 6/6 | 5.187935× | 5.159978× |
| n64/classes/dense/offset2-classes1 | 1.394203× | 1.378529× | 6/6 | 5.189341× | 5.157875× |
| n256/no-classes/zero/aligned32 | 2.765502× | 2.668196× | 6/6 | 1.511354× | 1.508104× |
| n256/no-classes/zero/offset2-classes1 | 2.344023× | 2.263711× | 6/6 | 1.087658× | 1.085191× |
| n256/no-classes/sparse/aligned32 | 2.706801× | 2.673100× | 6/6 | 1.540288× | 1.535403× |
| n256/no-classes/sparse/offset2-classes1 | 2.323684× | 2.269485× | 6/6 | 1.115666× | 1.108707× |
| n256/no-classes/dense/aligned32 | 2.709880× | 2.672783× | 6/6 | 1.664498× | 1.659021× |
| n256/no-classes/dense/offset2-classes1 | 2.325867× | 2.296632× | 6/6 | 1.203975× | 1.202638× |
| n256/classes/zero/aligned32 | 2.512977× | 2.415987× | 6/6 | 9.281765× | 9.170481× |
| n256/classes/zero/offset2-classes1 | 2.356725× | 2.284451× | 6/6 | 8.551695× | 8.543290× |
| n256/classes/sparse/aligned32 | 2.451252× | 2.446076× | 6/6 | 10.104561× | 10.076854× |
| n256/classes/sparse/offset2-classes1 | 2.343816× | 2.257628× | 6/6 | 9.325665× | 9.120471× |
| n256/classes/dense/aligned32 | 2.486150× | 2.447410× | 6/6 | 13.261730× | 12.409584× |
| n256/classes/dense/offset2-classes1 | 2.355455× | 2.290568× | 6/6 | 12.201941× | 11.233963× |

n64 medians improve1.291658–1.394203×parent, below1.5×. Four n16 and five n64 new-case medians are slower than old reference. n16 parent medians improve1.008732–1.018854×, while the existing n16 fused-class controls below regress; differing fixtures remain separate.

## All10 existing FP caller fixtures

| Fixture | Parent/head median | Min paired parent/head | Faster pairs | Old reference/head median | Min paired old/head |
|---|---:|---:|---:|---:|---:|
| BenchmarkVP9FusedFPClassesNative/4x4/zero | 0.981347× | 0.976843× | 0/6 | 1.852404× | 1.843376× |
| BenchmarkVP9FusedFPClassesNative/4x4/sparse | 0.979482× | 0.976960× | 0/6 | 1.911923× | 1.902823× |
| BenchmarkVP9FusedFPClassesNative/4x4/dense | 0.978870× | 0.923313× | 0/6 | 2.392098× | 2.252478× |
| BenchmarkVP9FusedFPClassesNative/8x8/zero | 1.470543× | 1.456160× | 6/6 | 4.333494× | 4.305635× |
| BenchmarkVP9FusedFPClassesNative/8x8/sparse | 1.461347× | 1.442763× | 6/6 | 4.458672× | 4.416272× |
| BenchmarkVP9FusedFPClassesNative/8x8/dense | 1.467903× | 1.459954× | 6/6 | 5.378628× | 5.335622× |
| BenchmarkVP9FusedFPClassesNative/16x16/zero | 2.558007× | 2.485779× | 6/6 | 9.374608× | 9.348123× |
| BenchmarkVP9FusedFPClassesNative/16x16/sparse | 2.555688× | 2.491460× | 6/6 | 10.215348× | 10.174220× |
| BenchmarkVP9FusedFPClassesNative/16x16/dense | 2.552623× | 2.479573× | 6/6 | 12.879736× | 12.393405× |
| BenchmarkVP9FPNoClassesNative | 1.457394× | 1.453374× | 6/6 | 0.938732× | 0.936649× |

Existing4x4 fused-class controls take about1.90–2.16% more time than parent, all6pairs slower; retain all samples. Existing no-classes median improves1.457394×parent but current time remains+6.526632%versus old reference. Old fused-class rows measure composite public quantization plus class handling; do not label every old-reference operation as pure handwritten assembly.

## Whole-codec measurements

| Workload | Current time vs parent | Median paired time change | Faster parent pairs | Current time vs old |
|---|---:|---:|---:|---:|
| VP9320 | +0.489540% | +0.362033% | 1/6 | -40.760481% |
| VP8720 | -2.171776% | -2.279946% | 6/6 | -21.554460% |
| VP9720 | -0.578321% | -0.669192% | 4/6 | -52.647330% |

VP9 320 is+0.489540%parent (five of six slower); VP9 720 is−0.578321% (four of six faster). Retain this mixed outcome without claiming a broad VP9 win. MeasuredVP8 is an unchanged-AMD control: its−2.171776%parent/all6faster result is not attributed to FP or inactiveARM DC helpers. Candidate time versus old reference is−21.554460%VP8,−40.760481%VP9320 and−52.647330%VP9720 in this job. Ratio of medians and median paired differences are separate statistics.

## AMD controls and focused SSE8

The generic control category has417 cases, excluding all36 new FP and all10 existing FP rows shown separately. Preserve all192 control median regressions and61 all-six-slower controls in root-full-proof.json/audit.json. New15DC supplement cases and old12DC cases are AMD controls; none is dropped. Source identity does not prove unchanged timing or establish a cause.

The independent focused EPYC7763 job has720 timed samples and120 excluded one-shot oracles, all timed bytes/allocations zero. Exact source trees verify261head/250parent/147old files. SSE8 helper and fixture source is unchanged. Public36 control medians span0.996232–1.070521×parent;23cases are faster in all6pairs, six have slower medians, none is slower in all6; minimum paired ratio0.960232 is retained. All36 public cases beat old reference in every pair, median1.034857–62.632896×. The large16-row ratio is CPU/call-protocol-specific inherited behavior, not a newFP gain or broad codec claim; no transition/layout cause is established. Four legacyB.N diagnostics discard return values and stay separate from liveB.Loop acceptance.

## Optimized libvpx, separate EPYC9V74 job

Pinned libvpxv1.16.0/1024874c5919305883187e2953de8fcb4c3d7fa6 uses optimized assembly, SSE2, AVX2 and runtime dispatch. All24 measured C timers are vpxenc-stats and each case has three Go-first/three C-first runs. A ratio above1 below means libvpx is faster. FixedVP8 means CPU-used−8,minq2/maxq56, not a fixed quantizer.

| Case | Go time / libvpx time | Measured quality/payload/count match? | Quality output matches every timed run? |
|---|---:|---|---|
| vp8-720p-fixed | 2.269763× | True | True |
| vp8-720p-adaptive | 1.653070× | False | True |
| vp9-320p | 2.679722× | False | True |
| vp9-720p | 2.552145× | False | True |

FixedVP8 measuredPSNR27.1576108715/SSIM0.9259273287,723593outputbytes,58encoded/2dropped/58qualityframes match between implementations and all timed output-control fields. This does not prove bitstream identity. AdaptiveVP8 and VP9 perform unequal quality/output work; these ratios do not establish equivalent-quality codec gaps or improvement versus a different job. Fixed Go allocation samples are all0; VP9320 all0.6166667 andVP9720 all1.25 in this C job, separate from native-profile exceptions.

## Preservation and scope

All91native+74C rawfiles (165total), source/dispatch/oracle/profile/control logs, all225excluded native preflights and every allocation exception are preserved byte for byte. The separate focused bundle preserves45rawfiles including originalSHA256SUMS; keep both bundles together (210rawfiles total). Root independent full proof JSON/script, independent C/focused audits and scripts, source SHA manifests and final CI/merge receipts are included. No source, README, checkout, index or local branch is changed by this audit. ARMDC performance evidence is committed separately and is not pooled with AMD measurements.


Publication: `raw-evidence.tar.gz` preserves all 25 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
