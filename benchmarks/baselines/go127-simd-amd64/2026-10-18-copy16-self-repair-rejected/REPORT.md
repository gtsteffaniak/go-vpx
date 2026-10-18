# Rejected PRIVATE AMD Copy16 self-noop repairs

Base46a08f8e. Root47 keeps Copy16productionUNCHANGED to protect its actualdisjointnativewins. These variants have no nativeperformanceevidence/acceptance claim.46native trueexact-self/zero-stridecontrols regressed(root.5096×,DSP.4416×parent); that evidence remainsvalid andretainedbyparent. We didnot trade knownordinarywins for an unmeasurednoopshortcut.

Allproposals place `if ss==ds &&sp==dp{returntrue}` AFTER originalEnabled/unsignedstride/fullvisibleLENguards. It is correctonlyatthatexactpoint: everyrowselfcopy is a semanticnoop; equalbase withunequalstride, shiftedaliasandinvalidshort-visible/capacitywindows mustkeeporiginalorderedmemmove/fallback/partialpanic. Newunitproof tests exactself strides0/1/7/16/31/1312, independentlyshortlen/spare/clippedcaps,nil/extremestrides/disabled. Existingfullbacking/allalias/root/DSPfrozenpanicoracles retained.

| Private variant | Inlinecost | Outcome |
| --- | ---: | --- |
| Naiveoriginalcarry+no-op |90 |reject: losesoriginal80inline,addsABIcall |
| Rangeconditionalcarry/chainedstore |79 |inlinebuthotloop10instructionsvs46carry6;rejectrisk |
| First15+separatefinalrow |87 |noninline |
| Indexedrange |74 |inlinebutadds2IMULperrow;rejectaddresswork |
| Originalcarry/chainedstore |84 |noninline |
| Sharedwhileconditioncarry |86 |noninline |
| Finalparent-proposedremaining15/returnfromloop |83 |noninline |

Naive and79rangevariants eachcompleted9Darwin/Linuxbuilds plusordinary/checkptr2 pixel/root/DSP selected/actualCPUoff fullstate/panic/alias tests, allPASS. Eachhas64public8fixtureone-shots(2packages×2modes×2CPU×8),allordinary0alloc. No emulatedspeedclaim.79range keepsDSPhelper inlined(cost121publicwrapper) androotcost113, but extra per-rowcmp/conditional/jump materiallychangeshotwork. Lastreturn83alsofails80budget. Otherboundedvariants are compile/codegenprobesonly, notclaimedtestedproductioncandidates.

Allrawlogs,sourcevariants,Go-mdiagnostics,actualcallerdisassemblyandbinaryhash/privatepathsretained. Originalnaive/rangecandidate snapshots stayprivate. Package is intentionallyNOT an integrationoverlay. No root/Git/dependency changes. Futurecarefulworkmustrepairself-noop withoutdestroying ordinaryinline/hotloopproperties,thenremeasureallnativecontrols.


Publication: `raw-evidence.tar.gz` preserves all 100 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
