# ARM Copy8 rows: final component confirmation

Retain the correctness-verified ARM helper and public glue under the user's explicit component rule. A representative packed, disjoint root copy improves **1.556568×** versus its pinned SIMD parent and **1.569303×** versus oldbuild. All six pairs are faster; two individual ratios are below 1.5× (minimum 1.403698×). Padded and other ordinary layouts improve less than 1.5×. Whole-encoder results are small/mixed and do not establish an improvement.

## Provenance and contract

Parent: f8dd8c7472e91c1f147511b04ceec716aca95aa7. Oldbuild: 691cd0512c4829d9ee40d483205dd223f9eb028d. Native Apple M4 Max, macOS26.3.1(a), Go1.27.2 SIMD, PGO off, GOMAXPROCS1/encoder thread1. Exact formatted measured sources are frozen in sources/candidate; overlay-sha256.json pins the seven transplant files. No root checkout, index, branch or README was changed by this experiment.

The eight-byte row dataflow is adapted from goav1 0f40b2bd848d8c5c7bdbaa016f2ced3085f7fc7a internal/av1/encoder/recon_gosimd.go, preserving its pinned source/license under upstream-goav1. No dependency is added. Only SIMD ARM64 selects the helper; AMD64, unsupported architectures, purego and experiment-off use the false stub and existing loops. Eight-by-four and four-byte paths are untouched.

Unsigned stride bounds reject negatives and strides>=2^30 before multiplication. Accepted 7*stride+8 endpoints fit int64 and must lie inside visible slice lengths. Exact source/destination pointer equality with equal strides returns handled only AFTER these checks. Every other accepted row performs one full uint64 load before its store, in original order. Within-row overlap, cross-row alias chains, overlapping/repeated rows, unaligned addresses and different strides remain supported. No terminal pointer is formed. Invalid/short inputs decline before any access or write; original fallback preserves slice-capacity extension and partial-write/panic order.

The compiler already uses one GPR load/store per eight-byte copy. This removes repeated slice checks/address work; it does not replace a memmove call or active handwritten assembly in the copy component. The oldbuild component denominator is compiler-generated Go. Its complete codec separately retains the original assembly.

## Correctness and source layouts

Independent public root/DSP oracles stage each eight-byte row, compare the complete backing allocation and panic outcome, and cover both alias directions, exact aliases, different/zero/overlapping strides, minimal lengths, source-only/destination-only/both short windows, spare capacity and negative/huge strides. Direct ARM tests assert handled/decline behavior, no-write declines, all16×16 source/destination alignments and short exact-alias validation. Final ordinary, race/checkptr2, purego and nosimd tests pass. Parent/candidate/oldbuild one-shot public fixtures all pass before timing.

The initial helper lacked the exact-alias shortcut; its root zero-stride exact-alias control regressed to0.697×. All144 exploratory samples and exact initial sources remain preserved. Parent reports an overlapping read-only source/hash audit during that screen, so it is not acceptance evidence or pooled with final results. An early-return exact-alias layout passed validation but cost83 inline units, over Go's80 budget; it was not timed. Its source and compiler diagnostic are preserved. The equivalent single-return/range-loop final layout costs74 and inlines. No further designs were tried.

## Fresh-six public component measurements

Six balanced three-way permutations,100ms testing.B.Loop per case, root/DSP order alternating by round, exact same fixture source in all variants, full scalar/canary checks and observable outputs. All288 final samples report0 B/op and0 allocs/op; every case improves in all six pairs. Exact-alias no-op controls improve2.531–2.848× but are NOT the acceptance basis. Only the packed disjoint root case establishes the new>=1.5× ordinary-copy gain. No broad ARM/x64 claim is inferred.

| Public case | Parent ns | Candidate ns | Oldbuild ns | Parent/candidate | Minimum paired ratio | Faster pairs |
|---|---:|---:|---:|---:|---:|---:|
| dsp/pitch0_0/aliasfalse | 5.0595 | 3.9585 | 5.1035 | 1.278136× | 1.231052× | 6/6 |
| dsp/pitch0_0/aliastrue | 4.4420 | 1.7550 | 4.8200 | 2.531054× | 2.514105× | 6/6 |
| dsp/pitch1312_1320/aliasfalse | 5.3835 | 3.9685 | 5.0755 | 1.356558× | 1.193610× | 6/6 |
| dsp/pitch17_31/aliasfalse | 5.0395 | 3.7430 | 5.0140 | 1.346380× | 1.316829× | 6/6 |
| dsp/pitch1_7/aliasfalse | 5.1370 | 3.8115 | 5.0285 | 1.347763× | 1.257007× | 6/6 |
| dsp/pitch32_37/aliasfalse | 5.0890 | 3.8565 | 5.1590 | 1.319590× | 1.132193× | 6/6 |
| dsp/pitch8_8/aliasfalse | 5.2240 | 3.7890 | 5.3840 | 1.378728× | 1.314366× | 6/6 |
| dsp/pitch8_8/aliastrue | 4.5170 | 1.7610 | 4.6635 | 2.565020× | 2.513850× | 6/6 |
| root/pitch0_0/aliasfalse | 5.0455 | 3.8680 | 5.0685 | 1.304421× | 1.252753× | 6/6 |
| root/pitch0_0/aliastrue | 4.9675 | 1.7440 | 4.4600 | 2.848337× | 2.536571× | 6/6 |
| root/pitch1312_1320/aliasfalse | 5.0755 | 4.0050 | 5.0375 | 1.267291× | 1.197368× | 6/6 |
| root/pitch17_31/aliasfalse | 5.0320 | 3.6785 | 5.0530 | 1.367949× | 1.312973× | 6/6 |
| root/pitch1_7/aliasfalse | 5.0495 | 3.8460 | 5.1180 | 1.312923× | 1.191774× | 6/6 |
| root/pitch32_37/aliasfalse | 5.1130 | 3.7535 | 5.1785 | 1.362195× | 1.320637× | 6/6 |
| root/pitch8_8/aliasfalse | 5.8060 | 3.7300 | 5.8535 | 1.556568× | 1.403698× | 6/6 |
| root/pitch8_8/aliastrue | 4.6590 | 1.7415 | 4.5650 | 2.675280× | 2.549570× | 6/6 |

## Whole VP8 measured separately

Each length uses six fresh balanced parent/candidate/oldbuild permutations:1280×720,60 or300 inputs,30fps,2500kbps,realtime CPU-used+8,one thread,encode-only,PGO off,no libvpx autodetection. CPU-used+8 is adaptive. All22 recorded semantic/configuration fields match exactly in every variant/round, and allocations separately match at zero (23 non-timing fields including allocations). Quality is skipped in this encode-only workload; matching recorded statistics are not a bitstream identity claim.

| Inputs | Parent ms/input | Candidate ms/input | Oldbuild ms/input | Candidate time vs parent | Faster parent pairs | Candidate time vs oldbuild |
|---|---:|---:|---:|---:|---:|---:|
|60|7.3342505|7.3575705|6.5313160|+0.317960%|4/6|+12.650659%|
|300|7.3290805|7.3451660|6.6777765|+0.219475%|2/6|+9.994187%|

Parent paired time ranges are-1.354488%..+1.629794% at60 and-0.408606%..+0.768934% at300. Small mixed results establish no whole-encoder gain. The complete oldbuild codec remains faster. Its assembly gap is scoped to this same native run; no cross-run subtraction or arithmetic attribution is made. Every slower sample and all36 raw JSON/stderr pairs are retained.

## Evidence and reproduction

The frozen overlay/patch, exact original context/identical fixtures, three source layouts, both compiler inliner reports, optimized selected code/callers, source/executable hashes, all432 exploratory+fresh component samples and whole raw/audits are included. verification.json independently recalculates medians, ratios, allocations and all22 semantic fields from raw files. drivers/reproduce.py reconstructs private snapshots via git archive and invokes the root project wrapper; executing it explicitly starts correctness/build/native timing. Original experiment drivers are retained too. No complete source snapshot or binaries are copied. raw-evidence.tar.gz contains the delivery members with an inner SHA manifest; SHA256SUMS covers every outer file except itself.
