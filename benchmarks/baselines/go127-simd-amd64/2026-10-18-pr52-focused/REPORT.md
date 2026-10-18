# PR52 required focused native audit

Exact head5c5cf3977052368746325c4033a1909e7eacd26a, base6bda90f7eeee6f407d6ed53e299eab505afa440e, old691cd0512c4829d9ee40d483205dd223f9eb028d. Run39131584057/job117687687085 SUCCESS. Native AMD EPYC7763, Go1.27.2 SIMD PGOoff GOMAX1, actual AVX2true/AVX512false. No local build or timing rerun.

SSE40/720timings+120excluded and DC36/648+108 excluded allzeroallocation/oraclesPASS. Public SSE36 vsbase med0.994262–1.003662 (15medianlosses,101slowerpairs); DC36med0.997578–1.004836 (9losses,87slowerpairs). All216oldpaired comparisons faster in each publicfamily. Old slice8x16 records55–63x; this is old-build/native scope, not a new ARMfix speedup. All samples retained.

Optional actualcapability permits Stats18/216timings, not FP48/576. All132 excludedcomponent preflight rows retained. Statsw8h1 med11.635ns bothrefs,3/6faster and3/6slower. Stats18 med0.998558–1.003499,2medianlosses/37slowerpairs, noall6loss. No new1.5xclaim. Actual native GNU capture binds measuredbinarySHA/buildinfo/protocol/ISA/order/source; all nine owner/caller instruction shapes match afterPC/RIP normalization, Pairs8at0x558a20/frame16/239instructions and BLoopat0x5a7480/355instructions identicaladdresses. This is actual capture, not reconstruction; do not infer causality from layout alone.

All recorded sourceblobs and entirepixelcapture source maps Gitverified. Five changedGo paths independently included; recordedCI source inventory is not fullrepository coverage. DSP ARMfixture is sourcebound separately from AMDexecution. Selected/off/race gates pass; AVX512-only correctness appropriately skips absenthardware. Binaries omitted: metadataSHA cross-bound, not independently rehashed absentexecutables. UnchangedAMDproduction means these remain controls, with every loss preserved.


Publication: `raw-evidence.tar.gz` preserves all 5 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
