# PRIVATE reversible unrolled FP16 guard production trial

UNMEASURED. This alternative is not an accepted optimization and makes no native x64 speed/codec claim. NativeCI must decide its ROI. Strongbase46a08f8e6eb03f6ee4549d975384bd6c952a81fa remains unchanged in the repository. No root/Git/PR edits. Both original125-fileguard and32-filecomponent packages remain immutable.

## Overlay and precise bindings

Four proposed files: quantize_fp_amd64.go (only unrolled aliasguard changes), new componenttest, new native diagnostic script and isolated additions to fullnative script. The original prefix containing all16-lane/narrow8 arithmetic/store/selection/validation and upper-bit cleanup is byte-identical. AllARM source is unchanged; actual liveARM quantizer normalized200-instruction body matches46 exactly.

The66 component names retain their meaning after selecting the trialguard:

| Name | Actual timed callee | Meaningful instructions |
| --- | --- | ---: |
| parent | fp16GuardParentProbe |83, exact frozen46 productionguard |
| single | fp16GuardSingleProbe |75, exact frozen singlecompare |
| unrolled | fp16WindowsDisjoint |88, ACTUAL trialproductionguard |

fp16GuardUnrolledProbe remains oracle-only, and actual selected guard is checked against it and an independent ordered-interval oracle. Linux/Darwin normalized machine bodies match the earlier productionprototypes exactly, changing only functionnames/addresses. Allnoninline costs146/135/185 (frozenunrolled185). Testhelper source is exact copiedlogic. Linux timedloop disassembly shows directstaticcalls; selection occurs outside b.Loop. The unrolled helper has no local80B frame/array/loop/panicBounds call but retains3ABI argument-home stores. Original/single arrayzeroing uses5legacyMOVUPS memory stores; unrolled is scalarGPR-only. No AVX-transition cause or performance benefit is inferred from that fact.

## Correctness completed

OfficialGo1.27.2, rootproject wrapper, private exact46archive,PGOoff,uniqueagent,-p2/GOMAX2. Ordinary/checkptr2 pixel fullFP tests PASS selectedAVX2 and actualAVX/AVX2off. This includes all ten pairwise alias types at all61 eligible sizes; frozen public EOB/panic/fullbacking checks; nil/independentlyshort/length-capacity controls; fullsigned/minimalwindow/invalidtable/disabled state. Component22case oracle verifies frozenparent/single/frozenunrolled/actualunrolled equality and unchanged backing. The extra exhaustiveoffset test is privatevalidation-only, not proposed as anotherrootfile.

Encoder state/scalar tests and existing36 public one-shot fixtures pass ordinary/checkptr2 selected/off. All72ordinary one-shots0B/0alloc;72checked samples match frozenparent allocationvectors exactly (classes cases16B/1alloc instrumentation, othercases0). Purego/nosimd encoder tests pass. Darwin benchmark strictly skips evenoptin: no Rosetta throughput samples. Linux was crossbuilt, not executed. All26 exact validationcommands and rawoutput preserved, including SDK/compiler diagnostics.

## Isolated native harness proposal

Existing498 names perproductionvariant,8964sixroundsamples,338oldexcludedpreflights,31commonreferenceoverlays and91oldrawfiles remain intact. Removing only the isolated currentcomponent hook and two SHA additions restores the exact46 fullscript.77 olduniqueSHApaths are preserved; add componentfixture and diagnosticscript for79. Their exactbytes/hashes and productionguard SHA are in scope.json; the componentlog also records currenthead,CPU/SDK/cache,binary/prod/test/script SHA.

One current-only pixelbinary is compiledPGOoff after unchanged fulltimings/profiles, and existingcleanup removes it. One new fp16-guard-component.txt =>92rawfiles. It contains22geometries×3guards=66names:12disjoint n64/256/1024×4optionalnil policies, and10declinecontrols.66 excluded1x preflights precede396 samples across allsixbalanced guardorders at80ms. Independentoracles/backingchecks remain outside timing; result sink and0alloc checks staylive. Declinecontrols never serve as accepted-speed denominators. No perf-threshold flakiness: only inventory/correctness/allocation defects fail. Existing45mintimeout remains; estimatedextra~45sec+cachedpixelcompile.

NativeLinux x64 and explicitoptin are required. No new test is copied into old/parentreferences. Production guard is exacttrialSHA-fenced. The actualpublic36/codec threevariant comparisons remain the production ROI evidence; component ratios are diagnostic separately. Root may revert/change trialguard after native evidence; current staticcounts do not choose a winner.

## Preservation

Alloverlay bytes, exact commands/logs, LLVM bodies, invariants, bindings and immutable source hashes are frozen here. Executables remain at hashed privatepaths listed in binary-inventory.json; only thosebinaryblobs are omitted from this lean package. Prior packages /private/tmp/govpx-fp16-offset-guard-ready and /private/tmp/govpx-fp16-guard-component-ready are untouched. No further variants/timers initiated.


Publication: `raw-evidence.tar.gz` preserves all 54 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
