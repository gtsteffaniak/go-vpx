# Repaired PR42: independent native focused SSE8 results

Exact head `f8dd8c7472e91c1f147511b04ceec716aca95aa7`, baseline `cf3117194f0590fa0af35951a32e8fadb6da9ab3`, old ASM `691cd0512c4829d9ee40d483205dd223f9eb028d`. Native AMD EPYC 9V74 80-Core Processor; Go1.27.2 Linux/amd64, GOAMD64v1, SIMD on, PGO off, GOMAXPROCS1. Workflow39089608484/artifact12121679836. This separate focused job makes no whole-codec claim. Full338/C/hosted and root review remain required before merge.

All 720 timed samples are independently reparsed: 40cases × three variants × six balanced order permutations. Every timed allocation is zero. All 120 one-shot oracle samples are preserved and excluded from medians. Identical scalar-oracle/canary fixture copied to both references; fixture bytes match exact head. Source trees and 622 recorded source SHA256 values (head241/base234/old147) match Git blobs. All raw hashes and selected/race tests pass, with actual AVX2=true/Enabled=true.

All36 live public medians improve 1.010991–1.482594× versus parent;31/36 improve in all six pairs. No median or all-pair1.5× claim is supported. Slice8x4 has mixed pairs and minimum0.869623×; every slower sample is retained. Stronger8x8/8x16/pointer classes improve all6pairs. All36 medians also beat their actual old ASM kernels (SSE2 for8x4/8x8/pointer8x8, AVX2 for8x16). No cross-CPU attribution or cause for the historical151ns anomaly is claimed.

## All36 live B.Loop fixtures

| Fixture | Parent/head | Minimum paired gain | Pairs faster | Old ASM/head |
|---|---:|---:|---:|---:|
| pointer8x8/pitch1312-1320/pattern0 | 1.480234× | 1.477688× | 6/6 | 1.160322× |
| pointer8x8/pitch1312-1320/pattern1 | 1.481943× | 1.432011× | 6/6 | 1.160574× |
| pointer8x8/pitch1312-1320/pattern2 | 1.480939× | 1.477638× | 6/6 | 1.158484× |
| pointer8x8/pitch32-37/pattern0 | 1.480451× | 1.369362× | 6/6 | 1.159833× |
| pointer8x8/pitch32-37/pattern1 | 1.479151× | 1.399889× | 6/6 | 1.159181× |
| pointer8x8/pitch32-37/pattern2 | 1.482594× | 1.474144× | 6/6 | 1.160498× |
| pointer8x8/pitch8-8/pattern0 | 1.480369× | 1.477306× | 6/6 | 1.160416× |
| pointer8x8/pitch8-8/pattern1 | 1.480668× | 1.473500× | 6/6 | 1.161248× |
| pointer8x8/pitch8-8/pattern2 | 1.482024× | 1.440992× | 6/6 | 1.160138× |
| slice8x16/pitch1312-1320/pattern0 | 1.195384× | 1.186561× | 6/6 | 1.707123× |
| slice8x16/pitch1312-1320/pattern1 | 1.196258× | 1.042273× | 6/6 | 1.709395× |
| slice8x16/pitch1312-1320/pattern2 | 1.196258× | 1.141230× | 6/6 | 1.707404× |
| slice8x16/pitch32-37/pattern0 | 1.195860× | 1.184834× | 6/6 | 1.705812× |
| slice8x16/pitch32-37/pattern1 | 1.197133× | 1.194909× | 6/6 | 1.707288× |
| slice8x16/pitch32-37/pattern2 | 1.194986× | 1.147441× | 6/6 | 1.706725× |
| slice8x16/pitch8-8/pattern0 | 1.196258× | 1.194268× | 6/6 | 1.705812× |
| slice8x16/pitch8-8/pattern1 | 1.195782× | 1.194909× | 6/6 | 1.705929× |
| slice8x16/pitch8-8/pattern2 | 1.194831× | 1.193959× | 6/6 | 1.704573× |
| slice8x4/pitch1312-1320/pattern0 | 1.027966× | 1.012011× | 6/6 | 1.308890× |
| slice8x4/pitch1312-1320/pattern1 | 1.010991× | 1.005491× | 6/6 | 1.308581× |
| slice8x4/pitch1312-1320/pattern2 | 1.014865× | 1.000705× | 6/6 | 1.309215× |
| slice8x4/pitch32-37/pattern0 | 1.021169× | 0.869623× | 5/6 | 1.305999× |
| slice8x4/pitch32-37/pattern1 | 1.012932× | 0.999441× | 5/6 | 1.306087× |
| slice8x4/pitch32-37/pattern2 | 1.015996× | 0.992207× | 5/6 | 1.308787× |
| slice8x4/pitch8-8/pattern0 | 1.019722× | 0.990600× | 5/6 | 1.306780× |
| slice8x4/pitch8-8/pattern1 | 1.022034× | 1.017734× | 6/6 | 1.308131× |
| slice8x4/pitch8-8/pattern2 | 1.021263× | 0.982184× | 5/6 | 1.309090× |
| slice8x8/pitch1312-1320/pattern0 | 1.083730× | 1.079548× | 6/6 | 1.270408× |
| slice8x8/pitch1312-1320/pattern1 | 1.070680× | 1.058224× | 6/6 | 1.255111× |
| slice8x8/pitch1312-1320/pattern2 | 1.071797× | 1.069151× | 6/6 | 1.256720× |
| slice8x8/pitch32-37/pattern0 | 1.079812× | 1.052956× | 6/6 | 1.265816× |
| slice8x8/pitch32-37/pattern1 | 1.079176× | 1.072332× | 6/6 | 1.267494× |
| slice8x8/pitch32-37/pattern2 | 1.078226× | 1.048449× | 6/6 | 1.264815× |
| slice8x8/pitch8-8/pattern0 | 1.070745× | 1.057822× | 6/6 | 1.253708× |
| slice8x8/pitch8-8/pattern1 | 1.067307× | 1.059528× | 6/6 | 1.250140× |
| slice8x8/pitch8-8/pattern2 | 1.073272× | 1.061773× | 6/6 | 1.256171× |

## Separate legacy B.N diagnostics

| Fixture | Parent/head | Pairs faster | Old ASM/head |
|---|---:|---:|---:|
| BenchmarkSSE8x16 | 1.193464× | 6/6 | 1.865795× |
| BenchmarkSSE8x4 | 1.067990× | 6/6 | 1.340709× |
| BenchmarkSSE8x8 | 1.105839× | 6/6 | 1.327236× |
| BenchmarkSSE8x8PtrFast | 1.459284× | 6/6 | 1.208731× |

Legacy B.N diagnostics retain their original protocol and are not pooled with B.Loop acceptance. Independent root verification is preserved byte for byte in root-proof.json. focused-raw.tar.gz holds all45 artifact files, including its original SHA256SUMS (44 entries); raw-sha256.json separately hashes all45. Every tar member and outer hash was verified. No source, README, local checkout/index/branch or Go job changed during auditing.
