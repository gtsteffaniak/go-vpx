# ARM copy16 final confirmation

The checked sixteen-row copy helper qualifies under the user’s component rule on Apple M4 Max. Fresh native measurements show fourteen cases at 1.581088–1.839802× the pinned parent, with all 84 corresponding paired gains at least 1.514699×. Both zero-stride exact-in-place controls also improve in all six pairs but remain below 1.5×. Whole-codec measurements are neutral/mixed and are not claimed as an encoder improvement.

## Exact scope

Parent is `34c4e4a17b231a24a3d3774c9621b7e5b80bac03`; candidate is that private snapshot plus `candidate.patch` (seven files). Oldbuild is `691cd0512c4829d9ee40d483205dd223f9eb028d`. Compiler is Go 1.27.2, darwin/arm64, SIMD enabled, PGO disabled. All final sources were formatted with the pinned SDK gofmt before the final correctness runs and builds, then remained frozen through all measurements. Full Go/assembly/module source hashes and all nine executable hashes are preserved. No root checkout, branch, index or README was changed.

The baseline compiler already implements each constant sixteen-byte copy as one vector load and one vector store. This change removes repeated per-row slice bounds/address work after a complete-window proof. It does not replace a runtime.memmove call in these paths. The sampled runtime.memmove profile cost includes other work and is not attributed to this helper.

## Correctness contract

Only `goexperiment.simd && arm64 && !purego` selects the helper. The complementary stub returns false on AMD64, unsupported architectures, purego and experiment-off selections, leaving the original row-copy loops. The ARM guard rejects negative or >=2^30 strides before multiplication; accepted fifteen-stride endpoints fit safely in int64 and must lie in the visible slice length. Every row performs one complete sixteen-byte load before its store, in original row order. This preserves within-row memmove overlap and arbitrary cross-row alias chains, including repeated or overlapping read/write rows. No terminal pointer is formed. Declines occur before reads/stores; the old fallback preserves short-window, spare-capacity, negative/extreme-stride and partial-write/panic behavior. The eight-byte paths are untouched.

Independent scalar tests compare the entire backing state and panic outcome, including alignment, both overlap directions, cross-row aliases, zero/unequal/overlap pitches, exact minimum spans, one-sided/both short slices, spare capacity and negative/huge strides. Direct pixel tests assert successful handled dispatch for eligible layouts and no-write declines for invalid inputs. Formatted-source native ordinary, race plus checkptr level 2, purego and nosimd tests pass. Public fixture checks scalar output, all source bytes and output canaries before and after every timed sample.

## Fresh native component results

Identical fixtures were copied to all three snapshots. Six balanced three-way order permutations were rotated from the initial screen, and DSP/root execution order alternated by round. Each timed case uses testing.B.Loop, 100 ms, one thread and observable output. All sixteen cases improve in all six paired rounds. Every one of 288 final samples reports 0 B/op and 0 allocs/op. Initial 288 samples are retained separately and are not pooled into the fresh estimates.

| Public case | Parent ns | Candidate ns | Oldbuild ns | Parent / candidate | Minimum paired gain | Faster pairs |
|---|---:|---:|---:|---:|---:|---:|
| dsp/pitch16_16/aliasfalse | 11.1400 | 6.0550 | 10.0700 | 1.839802× | 1.603795× | 6/6 |
| dsp/pitch16_16/aliastrue | 9.9290 | 6.0650 | 9.8590 | 1.637098× | 1.520810× | 6/6 |
| dsp/pitch32_37/aliasfalse | 9.8530 | 6.0285 | 9.8455 | 1.634403× | 1.583280× | 6/6 |
| dsp/pitch1312_1320/aliasfalse | 9.7035 | 6.0915 | 10.1780 | 1.592957× | 1.572282× | 6/6 |
| dsp/pitch0_0/aliasfalse | 9.8565 | 6.0975 | 9.6710 | 1.616482× | 1.567437× | 6/6 |
| dsp/pitch0_0/aliastrue | 9.6645 | 6.8180 | 9.6270 | 1.417498× | 1.350141× | 6/6 |
| dsp/pitch1_7/aliasfalse | 9.8330 | 6.0125 | 9.7545 | 1.635426× | 1.596395× | 6/6 |
| dsp/pitch17_31/aliasfalse | 9.7505 | 6.0475 | 9.7970 | 1.612319× | 1.542912× | 6/6 |
| root/pitch16_16/aliasfalse | 10.1600 | 6.1510 | 10.2100 | 1.651764× | 1.555450× | 6/6 |
| root/pitch16_16/aliastrue | 9.5965 | 5.9765 | 9.8720 | 1.605706× | 1.514699× | 6/6 |
| root/pitch32_37/aliasfalse | 9.8575 | 6.0020 | 10.1025 | 1.642369× | 1.612128× | 6/6 |
| root/pitch1312_1320/aliasfalse | 10.1465 | 5.9750 | 9.7235 | 1.698159× | 1.563492× | 6/6 |
| root/pitch0_0/aliasfalse | 9.7020 | 6.0975 | 9.8750 | 1.591144× | 1.563983× | 6/6 |
| root/pitch0_0/aliastrue | 9.4260 | 7.9085 | 9.5105 | 1.191882× | 1.088449× | 6/6 |
| root/pitch1_7/aliasfalse | 9.8595 | 6.1465 | 9.7785 | 1.604084× | 1.576811× | 6/6 |
| root/pitch17_31/aliasfalse | 9.6810 | 6.1230 | 9.6580 | 1.581088× | 1.577614× | 6/6 |

Oldbuild copy operations are compiler-generated Go vector copies, not handwritten assembly. Candidate copy operations are faster than that oldbuild control in all sixteen medians; no broad AMD64 claim is made. Exact in-place copies were measured explicitly: original code can skip each equal-address row, while the candidate performs its ordered loads/stores. Their measured gains and the lower zero-stride controls are retained.

## Whole VP8 results

Each length uses six fresh balanced parent/candidate/oldbuild permutations at 1280×720, realtime, CPU-used +8, one thread, 30 fps and 2500 kbps, encode-only, PGO off. Positive CPU-used is adaptive. All 22 recorded non-timing semantic fields plus allocation rate (23 fields total) match across every sample and all three implementations at each length. This includes payload size, frame/drop counts, quantizer and quality statistics; equal statistics do not prove bitstream identity.

| Inputs | Parent ms/frame | Candidate ms/frame | Oldbuild ms/frame | Candidate vs parent | Faster pairs | Paired change range |
|---|---:|---:|---:|---:|---:|---:|
| 60 | 7.539709 | 7.509548 | 6.698419 | -0.400035% | 4/6 | -1.841978% to +0.671964% |
| 300 | 7.389298 | 7.368984 | 6.765731 | -0.274911% | 4/6 | -1.865736% to +1.280537% |

Both lengths have mixed paired signs, so neither establishes a stable whole-encoder gain. The oldbuild remains faster for the complete codec despite the candidate’s faster copy component. All whole-codec raw JSON and stderr are preserved, including every slower candidate sample. Results are scoped to this shared M4 Max desktop; sampled profiles, instruction counts and cross-run differences are not used as causal timing claims.

## Evidence and reproduction

`sources/candidate/` contains the seven transplant files; parent/oldbuild contain relevant original context and identical public fixtures. `candidate.patch` contains all seven files relative to the clean pinned parent. For measurement, its two public fixture files are also overlaid identically into parent and oldbuild; their production sources remain pinned. Initial/final kernel raw logs, all 36 whole-codec JSON reports and stderr, validation logs, complete sample JSON, build/benchmark command lists, codegen and source/executable hashes are retained. `drivers/reproduce.py` reconstructs clean snapshots and runs the formatted candidate protocol through the root project wrapper. Running it explicitly starts native builds/tests/timing; packaging did not run new benchmarks.

`verification.json` records independent raw-file/sample/semantic/source checks. `SHA256SUMS` covers every other package member. The package is frozen after verification.

Root integration preserves all four measured production files and both timed benchmark fixtures exactly; see root-proof.json. The direct ARM-only test file gained an additional1,024 full-backing cases spanning every source/destination byte alignment0..15, two stridepairs and disjoint/aliased minimum windows after timing. These supplemental tests pass under race/checkptr2. The original measured test remains inside the frozen raw archive; current test and supplemental validation logs are separately hashed in root-validation.tar.gz. No benchmark body or measured production source changed. Root PGO refresh, noasm/SIMD-selection/pre-commit gates, and x64 Rosetta checkptr public fallback validation pass.

The submitted native x64 harness retains all330 preceding cases and adds8 DSP copy controls (338total), records23 source hashes, and checks current root/DSP full-state copy oracles ordinarily and under race/checkptr. The new helper is ARM-only; x64/purego/experiment-off stubs inline tofalse. New x64 copy timings therefore measure the existing Go compiler vector-copy path and cannot establish an ARM kernel speedup. README files are unchanged.
