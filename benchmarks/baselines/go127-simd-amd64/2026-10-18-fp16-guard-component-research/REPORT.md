# Private native FP16 guard component proposal

Base46a08f8e6eb03f6ee4549d975384bd6c952a81fa. No repository changes, PR writes, native timings or speed claims. Production guard stays unchanged. Candidates are diagnostic test helpers copied exactly, changing only their names.

The AMD-only pixel test fixture provides22 representative inputs: n64/256/1024 crossed with full, qnil, classesnil and bothnil (12 disjoint cases); every one of10 window pairs has a separate n256 partial-overlap decline control. All inputs have visible allocated windows. Read-only overlaps remain declines. An independent ordered-start/end interval oracle verifies expected results, and all input backing bytes are compared after every test/benchmark. Global result sink prevents unused output. Setup/oracles/full-backing checks occur outside b.Loop; the timed switch is outside each loop and calls a statically named actualparent/frozensingle/frozenunrolled guard. Decline controls are diagnostic, not acceptance denominators.

Ordinary and checkptr2 Rosetta actualAVX2 on/off tests PASS, including existing FP16 tests and22 new independent oracle subcases perrun. The benchmark strictly SKIPS Darwin even when opted in, so no emulated benchmark samples were produced. Linux and Darwin optimized helper bodies are identical to their earlier production prototypes after only address/renamed-symbol normalization; trailing INT3 alignment is excluded. Instruction counts83/75/88; Go inline costs146/135/185, allnoninline. Raw Linux loop disassembly proves directstaticcalls. This establishes comparability, not performance.

## Proposed native harness

Overlay has three files: testfixture; new diagnostic shell script; existing full benchmark script with only an appended isolated diagnostic hook before binarycleanup. Removing that hook restores the exact46 script. All498 old timed names,8964 six-round three-variant samples,31 reference overlays and existing profiles/gates are unchanged. The extra component is current-only, adds one pixel test binary (cleaned by existing wildcard), and one raw artifact fp16-guard-component.txt (91->92). No candidate helpers are copied into assembly/parent reference trees.

The new script requires native Linux x64, explicit timeroptin, strong unchanged productionquant source SHA, selected SIMD backend and GOMAXPROCS1. It records CPU,SDK,cache,head,source/test/binary hashes in the same diagnosticlog. It performs one oracle test phase and66 excluded one-shots, then six balanced permutations ofparent/single/unrolled at80ms. This yields66 names and396 timed samples. Strict inventory and0B/0alloc validation fail on correctness/harness defects; no speed threshold fails CI. All raw results retained. One-shot and timed results are separately tagged. The six permutations give every guard each orderposition twice.

The existing45-minute timeout is unchanged: expected additional timerwork is approximately45seconds plus a cached pixelcompile; actual nativefirst-run overhead remains to be observed. This is a componentdiagnostic; it cannot establish wholecodec benefit. NativeLinux correctness and allocation inventory must pass before judging ratios. Nativeplatform checks rely on the genuine x64 hostedrunner contract; uname alone cannot disprove arbitrary Linux emulation. No local AMDthroughput claims are made.

## Evidence and limits

sh-n PASS on both proposed scripts. Exact existing-script-statement preservation checked. Rootwrapper/private snapshot builds use officialGo1.27.2,PGOoff,GOMAXPROCS2,-p2,projectcache anduniqueagent. Candidate sourcebody/compilerproofs and all rawbuild/test/skip logs are included. Linux binary was crossbuilt, not run. Executables remain private at hashed paths in binary-inventory.json and are omitted from the lean publicationpackage. Scope.json inventories every addition. Root must approve integration aftercurrentPR46 remainsfixed/accepted.


Publication: `raw-evidence.tar.gz` preserves all 33 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
