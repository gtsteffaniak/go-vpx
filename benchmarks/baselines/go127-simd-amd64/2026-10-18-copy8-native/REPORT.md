# PR43: ARM Copy8 component and x64 controls

Exact `768f2fc2aba7bbd82981494473cbadc026a95570`, base `a3fb408d1323066cf9ed8470c7b69a8ff9b7e520`. All four checks pass. Independent root full/focused/C reviews pass; explicit PERFREADY and guarded ordinary merge are complete.

## Retained ARM component

Native M4Max/Go1.27.2/PGOoff representative packed disjoint root copy improves1.556568×parent/1.569303×old compiler-generated Go copy. All6 pairs improve, but two are below1.5× (minimum1.403698×). Exact-alias no-op gains are excluded as acceptance evidence. Other ordinary layouts improve less. The old component is GoGPRcopy, not handwritten assembly.

Whole ARM VP8 is mixed: candidate time+0.317960%parent at60inputs(4/6faster) and+0.219475%at300(2/6). Candidate takes+12.650659%/+9.994187%time relative to the old assembly-containing codec. These percentages use old-build time as denominator; they are not inverse old-build speed percentages. All22 recorded semantic/configuration fields plus separate zeroallocation rates match. Quality is skipped in this encode-only workload; no bitstreamidentity claim.

The checked8byte row helper preserves visible-window validation before exact-alias handling, row ordering/overlap and original fallback partial-write/panic behavior. All seven measured source/fixture files are preserved. AMD/purego/experimentoff uses the false stub and original loops.

## Full native x64 controls

Native CPU AMD EPYC 7763 64-Core Processor; Go1.27.2/PGOoff/GOMAXPROCS1. All346 component cases and allpublic/profile timings are unchanged-path AMD controls. No newAMD SIMD or causal speedup is claimed. ActualAVX2 executes; inheritedAVX512 is unavailable and correctly skips twice. ARM-token-only SIMD also skips twice with its fallback checks passing.

| New Copy8 control | Time vs parent | Pairs faster | Time vs old Go copy |
|---|---:|---:|---:|
| pitch0_0/aliasfalse | +0.242365% | 2/6 | -8.088889% |
| pitch0_0/aliastrue | -0.040167% | 4/6 | -3.260447% |
| pitch1312_1320/aliasfalse | +0.048146% | 2/6 | -7.603379% |
| pitch17_31/aliasfalse | -0.130265% | 5/6 | -5.271829% |
| pitch1_7/aliasfalse | -0.241313% | 2/6 | -8.174145% |
| pitch32_37/aliasfalse | +0.242014% | 1/6 | -8.200355% |
| pitch8_8/aliasfalse | +0.240154% | 2/6 | -7.450111% |
| pitch8_8/aliastrue | +0.145619% | 2/6 | -3.184466% |

## Whole x64 controls

Time percentage is100×(median(current)/median(reference)−1), using six samples each. Median paired changes are shown separately; no sample was dropped.

| Case | Ratio-of-medians time vs parent | Median paired change | Pairs faster | Time vs old build |
|---|---:|---:|---:|---:|
| VP9 320x180 | +0.261022% | -0.151485% | 4/6 | -40.390821% |
| VP8 720p | +2.368013% | +1.387650% | 1/6 | -18.754408% |
| VP9 720p | +2.718161% | +1.590911% | 1/6 | -48.465024% |

These full-run control regressions remain visible. They are not attributed to the ARM-only Copy8 source, and no x64 VP9 or codec-gap closure is claimed.

## Correctness, allocation exception and all regressions

All346cases×6samples×3variants have0B/op and0allocs/op. All41 CI-recorded source hashes match exact Git blobs; source mapping uses longest complete relative path so the root and DSP Copy8 fixture basenames cannot collide. The pinned driver overlays exactly24 identical test-only fixtures onto each reference, with productionfiles pinned.

Native validation:143top-levelPASS plus6retainedpanicprobes=149passing targeted checks;294totalPASS including subtests. Root and DSP Copy8 state tests pass ordinary+race/checkptr (4executions), as do existingCopy16state tests(4). Only four expected architectural skips occur:ARM-only token SIMD twice and unavailable inheritedAVX512 twice. No unexpected skip or failure.

**All22semantic fields match** in all3public and2profile datasets. Allocations are assessed separately: current and parentVP9720 are exactly1.25 in all6rounds, while old-reference round3 is1.2666666667 (76allocations/60inputs versus75/60in other runs). That one old-build allocation outlier is retained, not treated as a current regression or hidden. All other public/profile allocation samples match:VP8zero,VP9320 .616667,VP9720profile .625.

All 134 median kernel regressions and 2 all-six-slower controls are preserved in audit.json. The latter are:

| Unchanged control | Median time vs parent | Largest paired time change |
|---|---:|---:|
| GatherPackedPredictionNative/w8/pitch32/extremetrue | +0.168776% | +0.759494% |
| SubtractRowsNative/16x8/paddedfalse/pattern1 | +0.496983% | +25.505498% |

The subtract control round5+25.505498%timing outlier is retained. Its historical AMD old-reference path is scalar; the harness assembly label does not make that subtraction comparison assembly. Short independent profiles establish samples/hotspots only, with no additive cumulative attribution.

## Separate focused controls

[Focused workflow39091858125](https://github.com/thesyncim/govpx/actions/runs/39091858125), Intel8573C, preserves720timed+120excludedone-shot samples and all45artifactfiles. All36SSE8 and4legacy cases are controls. Parent/current public median ratios.987038–1.008654, with21median and2all-six slower controls. OldASM/current16-row97.519–98.235×ratios are specific to this CPU/public call protocol; no broad98×codec or newPR43gain is supported. Legacy discarded-return B.N diagnostics remain separate. Frozen sibling `govpx-pr43-focused-evidence-ready/` contains all40ratios/raw/hashes/rootproof and is referenced inside audit.json.

## Optimized libvpx

Separate Cjob CPU AMD EPYC 7763 64-Core Processor; same CPUmodel does not make this the same VM/run as fullnative. Reference build verifies assembly/SSE2/AVX2/runtime detection and vpxenc-stats timers, pinnedlibvpxv1.16.0/1024874c5919305883187e2953de8fcb4c3d7fa6.

| Case | Go/libvpx time | Scope |
|---|---:|---|
| vp9-320p | 2.932200× | Unequal measured quality/output work |
| vp9-720p | 2.982942× | Unequal measured quality/output work |
| vp8-720p-adaptive | 1.805886× | Unequal measured quality/output work |
| vp8-720p-fixed | 2.372684× | Matching measured PSNR/SSIM/payload/frame/drop counts |

Libvpx remains faster. Fixed means CPU-used−8 with minq2/maxq56, not fixed quantizer. Matching measured quantities does not prove identical bitstreams. Every quality-pass outputstat matches timing; adaptive/VP9 remain unequal-work diagnostics. Independent Cproof is embedded byte-recoverably with its hash in audit.json.

## Preservation

native-raw.tar.gz andlibvpx-raw.tar.gz preserve86+74=160rawfiles byte for byte. Focused45files are preserved separately, totaling205rawfiles. Member/outer hashes are verified. Every timing/allocation exception, semantic field, selected test/dispatch result, source identity, profile summary and regression remains in the audit/raw data. No README, source, checkout, index, local branch or Go job is changed by this read-only audit.

## Code generation and independent proofs

The three original codegen files are preserved byte for byte: [normalized proof](amd64-normalized-codegen-proof.json), [base denoiser LLVM](base-denoiser-llvm.txt), [head denoiser LLVM](head-denoiser-llvm.txt). Root compares6522symbols (2799project symbols): all VP9/pixel kernels match within normalization, while one VP8denoiser function has equivalent register/commutative-CMP choices. RIP-relative data displacements/comments are ignored; this is cross-compiled codegen evidence, not native timing, binary identity or data-layout proof. No causal attribution is made. Independent full root proof, with the explicit old-reference allocation outlier, is preserved in root-proof.json; the independent C proof is directly preserved in root-C-proof.json. Focused proof/raws remain in the frozen sibling bundle and referenced inside audit.json.

## Merge proof

Merged `bba909597f1a68e4ceeb94dcca0229c765f10d84` through guarded normal GitHub --merge. Tested and merged tree both `03f8f4e4d8bdf2a29d58dcd5ea36f977577ffb26`; parents exactly `[a3fb408d1323066cf9ed8470c7b69a8ff9b7e520, 768f2fc2aba7bbd82981494473cbadc026a95570]`. Head ancestry and remote main verified. Source remote remains exact768f for the root successor; no checkout/index/local branch was changed. Complete proof and four passing premerge checks are embedded in audit.json.

## Detailed evidence

`raw-evidence.tar.gz` preserves the complete original frozen evidence package byte for byte, including all audits, root proofs, source hashes, disassembly, original manifests and nested raw benchmark archives. Extract into a separate directory to inspect and verify its original manifest.
