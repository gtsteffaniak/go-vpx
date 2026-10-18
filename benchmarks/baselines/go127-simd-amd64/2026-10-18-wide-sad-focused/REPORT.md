# PR44: focused x64 SSE8 controls

Exact `46b89f5ab4c6c469755323a61c77fe3a766b4081`, base `bba909597f1a68e4ceeb94dcca0229c765f10d84`, old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Native AMD EPYC 9V45 96-Core Processor, Go 1.27.2 Linux/amd64, GOAMD64=v1, SIMD enabled, PGO off, GOMAXPROCS=1, six balanced rounds at 80 ms. Workflow 39094713024, artifact 12123814047.

All 40 fixtures are unchanged SSE8 controls. No SAD fixture is timed in this focused job. Full 412-case native/C/hosted acceptance and root review remain pending; HOLD MERGE.

All 720 timed samples and 120 excluded one-shot samples are retained; all timed bytes and allocations are zero. Exact pinned source trees match head/base/old counts 250/245/147. All 45 raw files match their saved manifest, including the original 44-entry SHA256SUMS. Public outputs, canaries, sinks, AVX2/Enabled dispatch and ordinary/race checks pass.

The 36 public parent/head median ratios range 0.952417–1.006949×. Thirty median controls are slower; six are slower in every paired round. No case is faster in every round. The minimum paired ratio is 0.915712×. Pointer controls take 2.376–4.996% more time than the parent, with four of nine slower in all six pairs. These regressions remain acceptance-relevant; unchanged helper source is not proof of unchanged performance. New code growth or layout is a possibility, not an established cause.

## All 36 public B.Loop controls

| Fixture | Parent/head | Minimum paired ratio | Pairs slower | Old ASM/head |
|---|---:|---:|---:|---:|
| pointer8x8/pitch1312-1320/pattern0 | 0.952417× | 0.942712× | 5/6 | 1.109206× |
| pointer8x8/pitch1312-1320/pattern1 | 0.962708× | 0.950042× | 6/6 | 1.112735× |
| pointer8x8/pitch1312-1320/pattern2 | 0.958803× | 0.947025× | 6/6 | 1.109323× |
| pointer8x8/pitch32-37/pattern0 | 0.964412× | 0.924812× | 6/6 | 1.110203× |
| pointer8x8/pitch32-37/pattern1 | 0.962700× | 0.945745× | 5/6 | 1.111684× |
| pointer8x8/pitch32-37/pattern2 | 0.963177× | 0.951820× | 5/6 | 1.110041× |
| pointer8x8/pitch8-8/pattern0 | 0.967599× | 0.960839× | 6/6 | 1.127552× |
| pointer8x8/pitch8-8/pattern1 | 0.975394× | 0.951392× | 4/6 | 1.126376× |
| pointer8x8/pitch8-8/pattern2 | 0.976794× | 0.955399× | 4/6 | 1.119050× |
| slice8x16/pitch1312-1320/pattern0 | 0.993407× | 0.982650× | 3/6 | 1.644311× |
| slice8x16/pitch1312-1320/pattern1 | 0.995389× | 0.980819× | 4/6 | 1.639913× |
| slice8x16/pitch1312-1320/pattern2 | 0.990781× | 0.982595× | 4/6 | 1.647983× |
| slice8x16/pitch32-37/pattern0 | 1.002300× | 0.996795× | 2/6 | 1.584454× |
| slice8x16/pitch32-37/pattern1 | 0.999484× | 0.995773× | 1/6 | 1.576714× |
| slice8x16/pitch32-37/pattern2 | 1.001490× | 0.985048× | 3/6 | 1.568212× |
| slice8x16/pitch8-8/pattern0 | 0.998740× | 0.978991× | 3/6 | 1.576868× |
| slice8x16/pitch8-8/pattern1 | 0.994465× | 0.976838× | 5/6 | 1.577656× |
| slice8x16/pitch8-8/pattern2 | 1.000803× | 0.989784× | 3/6 | 1.580003× |
| slice8x4/pitch1312-1320/pattern0 | 1.001928× | 0.987382× | 3/6 | 1.379030× |
| slice8x4/pitch1312-1320/pattern1 | 0.996898× | 0.982166× | 3/6 | 1.375722× |
| slice8x4/pitch1312-1320/pattern2 | 0.999033× | 0.977778× | 3/6 | 1.375658× |
| slice8x4/pitch32-37/pattern0 | 0.999046× | 0.915712× | 3/6 | 1.376842× |
| slice8x4/pitch32-37/pattern1 | 1.006949× | 0.992591× | 3/6 | 1.370430× |
| slice8x4/pitch32-37/pattern2 | 0.995225× | 0.974797× | 4/6 | 1.385252× |
| slice8x4/pitch8-8/pattern0 | 0.994797× | 0.964315× | 4/6 | 1.383562× |
| slice8x4/pitch8-8/pattern1 | 1.006066× | 0.984213× | 2/6 | 1.376969× |
| slice8x4/pitch8-8/pattern2 | 0.988119× | 0.973904× | 5/6 | 1.377037× |
| slice8x8/pitch1312-1320/pattern0 | 0.987430× | 0.964518× | 6/6 | 1.337353× |
| slice8x8/pitch1312-1320/pattern1 | 0.986389× | 0.942503× | 4/6 | 1.348127× |
| slice8x8/pitch1312-1320/pattern2 | 0.982401× | 0.948417× | 5/6 | 1.312142× |
| slice8x8/pitch32-37/pattern0 | 0.990757× | 0.975008× | 5/6 | 1.331892× |
| slice8x8/pitch32-37/pattern1 | 0.978325× | 0.964238× | 6/6 | 1.321348× |
| slice8x8/pitch32-37/pattern2 | 0.986579× | 0.973163× | 5/6 | 1.342818× |
| slice8x8/pitch8-8/pattern0 | 0.986459× | 0.980543× | 5/6 | 1.338863× |
| slice8x8/pitch8-8/pattern1 | 0.985397× | 0.978404× | 5/6 | 1.334862× |
| slice8x8/pitch8-8/pattern2 | 0.988507× | 0.976459× | 5/6 | 1.358590× |

## Four legacy B.N diagnostics

| Fixture | Parent/head | Pairs slower | Old ASM/head |
|---|---:|---:|---:|
| BenchmarkSSE8x16 | 1.007904× | 1/6 | 1.590866× |
| BenchmarkSSE8x4 | 1.009175× | 2/6 | 1.319590× |
| BenchmarkSSE8x8 | 0.996370× | 2/6 | 1.315805× |
| BenchmarkSSE8x8PtrFast | 0.996934× | 3/6 | 1.171795× |

Legacy functions discard their SSE results and lack the live fixture sink. They are separate diagnostics, not pooled with live public acceptance. Old public SSE8 uses actual SSE2 for 4/8/pointer and AVX2 for 16 rows. CPU/protocol-specific ratios establish no broad codec or new PR44 SSE8 gain.

## Preservation

focused-raw.tar.gz preserves every raw byte, including original SHA256SUMS. audit.json retains every sample, ratio, excluded sample, allocation, regression, source hash and dispatch result. root-proof.json is the independent parent audit. raw-sha256.json and sha256.json cover archive members and outer files. No local Go work or repository mutation occurred during this audit.


The complete frozen evidence package is inside `raw-evidence.tar.gz` (6 original files). It preserves all nested raw archives, manifests, reports and auditors byte for byte. Extract to another directory to inspect it; `SHA256SUMS` covers the published outer files.
