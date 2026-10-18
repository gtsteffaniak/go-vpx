# PR45: focused SSE8 controls

Exact head `ac259fce071f3ac5fc45603c6bd7e7926bdfcf73`, base `2ca5bab1110426401e4f9dbfd6fa0da8f568c717`, old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Native AMD EPYC 7763 64-Core Processor, Go1.27.2 Linux/amd64, PGOoff, GOMAXPROCS1, GOAMD64v1; six balanced rounds at80ms. Workflow39098241911 / artifact12125885875.

All40 fixtures are unchanged-source SSE8 controls. No FP or DC kernel is timed in this focused job, and no whole-codec or FP speedup is inferred. All720 timed samples have0B/op and0allocs/op;120 one-shot samples are preserved separately/excluded. Exact pinned source trees verify261head/250parent/147old files, and all45raw files match original44-entrySHA256SUMS plus the separately hashed manifest. AVX2/Enabled and ordinary/race correctness gates pass.

Public36 parent/head median ratios span0.996232–1.070521×;23cases are faster in all6pairs, six medians slower(all slice8x4), none slower in all6. Minimum paired ratio0.960232 is retained. Equal source does not establish the cause of timing drift. All36 public cases beat old reference in every pair; median1.034857–62.632896×. Large16-row ratios are CPU/call-protocol-specific inherited controls, not newPR45 or codec-wide gains; no transition/layout cause is established. Old public4/8/pointer uses actualSSE2 and16rows AVX2.

| Public fixture | Parent/head | Min paired parent/head | Slower pairs | Old/head |
|---|---:|---:|---:|---:|
| pointer8x8/pitch1312-1320/pattern0 | 1.055003× | 1.033358× | 0/6 | 1.057014× |
| pointer8x8/pitch1312-1320/pattern1 | 1.042509× | 0.987503× | 1/6 | 1.043658× |
| pointer8x8/pitch1312-1320/pattern2 | 1.062390× | 1.014043× | 0/6 | 1.067999× |
| pointer8x8/pitch32-37/pattern0 | 1.044831× | 1.034816× | 0/6 | 1.063391× |
| pointer8x8/pitch32-37/pattern1 | 1.056874× | 0.995361× | 1/6 | 1.067939× |
| pointer8x8/pitch32-37/pattern2 | 1.070521× | 1.019884× | 0/6 | 1.057285× |
| pointer8x8/pitch8-8/pattern0 | 1.026479× | 0.960232× | 1/6 | 1.034857× |
| pointer8x8/pitch8-8/pattern1 | 1.060099× | 0.986134× | 1/6 | 1.079413× |
| pointer8x8/pitch8-8/pattern2 | 1.066233× | 1.039474× | 0/6 | 1.074303× |
| slice8x16/pitch1312-1320/pattern0 | 1.022009× | 1.019802× | 0/6 | 62.058277× |
| slice8x16/pitch1312-1320/pattern1 | 1.023270× | 1.019827× | 0/6 | 62.178095× |
| slice8x16/pitch1312-1320/pattern2 | 1.023249× | 1.012225× | 0/6 | 55.412275× |
| slice8x16/pitch32-37/pattern0 | 1.030957× | 1.030000× | 0/6 | 62.632896× |
| slice8x16/pitch32-37/pattern1 | 1.033771× | 1.030549× | 0/6 | 62.601626× |
| slice8x16/pitch32-37/pattern2 | 1.031572× | 1.002433× | 0/6 | 55.876837× |
| slice8x16/pitch8-8/pattern0 | 1.031638× | 1.027194× | 0/6 | 62.189826× |
| slice8x16/pitch8-8/pattern1 | 1.033830× | 1.015207× | 0/6 | 62.197393× |
| slice8x16/pitch8-8/pattern2 | 1.032268× | 1.031037× | 0/6 | 55.637605× |
| slice8x4/pitch1312-1320/pattern0 | 0.999519× | 0.998399× | 5/6 | 1.332051× |
| slice8x4/pitch1312-1320/pattern1 | 1.000962× | 0.996373× | 1/6 | 1.336221× |
| slice8x4/pitch1312-1320/pattern2 | 0.999199× | 0.996273× | 4/6 | 1.331909× |
| slice8x4/pitch32-37/pattern0 | 0.999199× | 0.995846× | 3/6 | 1.331056× |
| slice8x4/pitch32-37/pattern1 | 0.999947× | 0.995344× | 3/6 | 1.334544× |
| slice8x4/pitch32-37/pattern2 | 0.999573× | 0.996910× | 3/6 | 1.332977× |
| slice8x4/pitch8-8/pattern0 | 0.996232× | 0.987458× | 3/6 | 1.335845× |
| slice8x4/pitch8-8/pattern1 | 1.001497× | 1.000000× | 0/6 | 1.334580× |
| slice8x4/pitch8-8/pattern2 | 1.000962× | 0.994587× | 2/6 | 1.333832× |
| slice8x8/pitch1312-1320/pattern0 | 1.028846× | 1.027473× | 0/6 | 1.342033× |
| slice8x8/pitch1312-1320/pattern1 | 1.026978× | 1.025571× | 0/6 | 1.339735× |
| slice8x8/pitch1312-1320/pattern2 | 1.030220× | 1.026484× | 0/6 | 1.342033× |
| slice8x8/pitch32-37/pattern0 | 1.028362× | 1.024635× | 0/6 | 1.340805× |
| slice8x8/pitch32-37/pattern1 | 1.026520× | 1.024635× | 0/6 | 1.340192× |
| slice8x8/pitch32-37/pattern2 | 1.028859× | 1.027447× | 0/6 | 1.343106× |
| slice8x8/pitch8-8/pattern0 | 1.027918× | 1.026508× | 0/6 | 1.342334× |
| slice8x8/pitch8-8/pattern1 | 1.026039× | 1.019056× | 0/6 | 1.339424× |
| slice8x8/pitch8-8/pattern2 | 1.026978× | 1.024680× | 0/6 | 1.340649× |

## Four separate legacy B.N diagnostics

| Fixture | Parent/head | Slower pairs | Old/head |
|---|---:|---:|---:|
| BenchmarkSSE8x16 | 1.020109× | 0/6 | 60.650988× |
| BenchmarkSSE8x4 | 1.001115× | 3/6 | 1.391471× |
| BenchmarkSSE8x8 | 1.000387× | 2/6 | 1.343935× |
| BenchmarkSSE8x8PtrFast | 0.959762× | 6/6 | 1.094785× |

Legacy diagnostics discard returned SSE and lack the public live sink; they are not pooled with B.Loop acceptance. Every sample/allocation/output/source hash and slower case remains in audit.json/raw. Independent root proof is byte-identical in root-proof.json. All45raw bytes including originalSHA256SUMS are preserved; inner/outer SHA256 verified. Final PR45 acceptance and merge receipt live in the separate full native/C bundle. No local Go or repository mutation occurred during audit.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
