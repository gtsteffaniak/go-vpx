# PR50 focused native x64 audit

Exact head `10caaa67ed53f8c3c1248f20070e5a4117245852`, base `cb9bbdcd185cee2026d178e23b51bb6939335d75`, old Go/assembly `691cd0512c4829d9ee40d483205dd223f9eb028d`. PR run 39122180914/job 117659831714 **SUCCESS**, completed 11:15:07UTC. This audit covers that PR run only; the separate manual archived49 run is excluded.

**AMD EPYC 7763**, official Go 1.27.2 SIMD, GOAMD64 v1, PGOoff, GOMAXPROCS 1. Actual SDK AVX2=true, Enabled=true, AVX512=false. All 95 artifact-member hashes and exact immutable Git source hashes verify (head 308/base 303/old 147 source files). Shared fixtures match head bytes; base keeps its original inplace API and old alone receives the reference-only adapter. Required selected/off ordinary and race/checkptr tests pass. No selected512 correctness/performance claim is possible: FP32 timed comparisons were explicitly omitted on this CPU.

## Public paired DC: native result

All 36 cases ×six balanced three-way rounds = 648 timed samples and 108 excluded actual one-shots. Every timed/one-shot DC sample is 0B/0alloc; all pre/post output/full-backing/input-canary checks pass. **All 36 case medians and all216 head/reference pairs beat both base and old.** Ratios below are reference/head, so larger means faster current.

| Actual public scope | Cases | Base/head medians | Old/head medians | All pairs ≥1.5×base |
|---|---:|---:|---:|---:|
| Generic pair, disjoint, function-value fixture |9|1.635746–1.668449|1.304031–1.337385|9/9|
| Generic pair, exact alias +same reset work |9|1.463630–1.492329|1.192976–1.222867|0/9|
| Generic pair, disjoint, direct-call control |9|1.709194–1.749923|1.354339–1.393037|9/9|
| New inplace API, exact alias +same reset work |9|1.858986–1.871601|1.502224–1.521623|9/9|

The actual new inplace API has all 54 pairs≥1.823735×base. Its old median ratios all exceed1.5, but only 8/9 cases keep every old pair≥1.5; the lowest old pair is1.472077. Across all 36 cases, minimum pair is1.422340 versus base and1.164574 versus old. Thus 27 cases qualify with both median and every pair≥1.5×base; no universal 1.5 floor is claimed. Generic aliased cases remain useful smaller wins. Inplace fixtures include identical four-row input restoration in each implementation, so these are live public-call measurements with shared reset work, not isolated arithmetic costs.

Actual native GNU code shows old inplace adapter has no live function symbol and its common B.Loop calls old original generic API directly—no extra adapter call. Generic function-value fixture uses the same indirect-call shape on all three refs. Current direct fixture calls the actual new public inplace API. Current generic/inplace owners have 291/186 disassembled instructions, no YMM/ZMM or vector stack accesses. This native code capture is retained; no inference from emulated timing is used.

## Preserved SSE controls

All 40 cases ×18 samples = 720 timings, plus 120 excluded actual one-shots. All timed cases and all public one-shots have zero allocations. Four legacy B.N one-shot setup allocations are recorded separately and do not enter acceptance. Original 36 public medians range 0.996461–1.012598×base; 15 medians and 96 paired samples are slower. Legacy 4 medians range 0.994123–1.044779×base, with 2 medians/13 pairs slower. Every slower pair remains in compact/raw evidence. These source controls do not establish a new SSE algorithm gain.

All 36 public SSE medians/all216 pairs beat old. The very large 8x16 ratio 54–61× is explicitly an old **public-path** comparison: actual old AVX2 kernel mixes legacy MOVQ/PUNPCKLDQ with live YMM operations and VZEROUPPER only at return. It must not be presented as 54–61× pure arithmetic improvement over a clean SIMD kernel. Other public SSE paths are roughly 1.31–1.39×old in this run.

## Optional Stats and unresolved regression

Actual selected512 is unavailable. All 132 excluded optional one-shots pass on both refs; only 18 Stats cases ×12 samples = 216 timings were permitted, all 0B/0alloc. No FP32 timing or selected512 acceptance is claimed from fallback preflights.

Stats w8h1 regresses **17.035→18.91ns (+11.006751%)**, all 6 pairs slower. w8h7 is +1.816969%, also all 6 slower. w8h64 is +3.280387%median, with mixed 3/6 signs. These remain unresolved native regressions; source identity is not a reason to dismiss them as noise.

Exact-Git, same-option Linux reconstruction finds all five active Stats bodies instruction-identical. B.Loop body is identical except one cold error-string relocation. Reconstruction Pairs8/dispatcher/Stats entries shift from mod64=0 to 32; B.Loop entry shifts 32→0. This is a **specific layout hypothesis, not causality proof**. Actual-run pixel binaries were deleted and Stats disassembly was not captured, so their native addresses/alignment are unknown. Reconstructed binaries have different SHA hashes and their addresses are never substituted for measured ones. Both raw addresses/alignment and complete differences are retained. No reconstruction binary was executed.

This is component evidence on one EPYC 7763 run, not a whole-codec/other-CPU claim. No ratios are pooled across runs. All allocations, correctness/ISA logs,648 DC +720 SSE +216 Stats timings,108+120+132 preflights, exact source/fixture/codegen metadata and every regression/outlier are preserved. ROOT/source/PR were read-only; no rerun or dispatch occurred.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
