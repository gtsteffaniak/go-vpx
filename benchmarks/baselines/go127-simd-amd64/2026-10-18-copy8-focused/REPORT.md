# PR43: focused x64 SSE8 controls

Exact `768f2fc2aba7bbd82981494473cbadc026a95570`, base `a3fb408d1323066cf9ed8470c7b69a8ff9b7e520`. Native INTEL(R) XEON(R) PLATINUM 8573C, Go1.27.2 Linux/amd64, GOAMD64v1, SIMD enabled, PGOoff,GOMAXPROCS1. Workflow39091858125/artifact12122992075. ARM-only copy8 helper is inactive onAMD; SSE8source is unchanged. **All40 measured cases are controls, with no newAMD speedup or whole-codec claim.** HOLD MERGE pending full346/C/hosted and root review.

All720timed samples and120 excluded one-shot samples are independently parsed; alltimed bytes/allocations arezero. Every source tree/hash matches pinnedGit (head245/base241/old147). All45 artifactfiles match original44-entrySHA256SUMS plus separately hashed manifest. Actual AVX2/Enabled/new ordinary+race dispatch checks pass. Balanced six order permutations and allcase inventories verified.

All36 public control median parent/current ratios range.987038–1.008654;21medians slower,2cases slower in all6pairs. Worst paired ratio.841438 is retained. This control drift is not attributed to ARM-only production changes.

## All36 live B.Loop controls

| Fixture | Parent/head | Minimum paired ratio | Pairs slower | Old ASM/head |
|---|---:|---:|---:|---:|
| pointer8x8/pitch1312-1320/pattern0 | 1.000818× | 0.985344× | 2/6 | 1.113048× |
| pointer8x8/pitch1312-1320/pattern1 | 1.000126× | 0.992671× | 4/6 | 1.122952× |
| pointer8x8/pitch1312-1320/pattern2 | 0.999181× | 0.989368× | 2/6 | 1.113924× |
| pointer8x8/pitch32-37/pattern0 | 0.997268× | 0.988083× | 4/6 | 1.114055× |
| pointer8x8/pitch32-37/pattern1 | 1.000952× | 0.982717× | 3/6 | 1.108177× |
| pointer8x8/pitch32-37/pattern2 | 0.998601× | 0.984716× | 3/6 | 1.114551× |
| pointer8x8/pitch8-8/pattern0 | 0.998538× | 0.990518× | 3/6 | 1.120862× |
| pointer8x8/pitch8-8/pattern1 | 1.001906× | 0.995833× | 3/6 | 1.123134× |
| pointer8x8/pitch8-8/pattern2 | 1.002292× | 0.980740× | 3/6 | 1.120074× |
| slice8x16/pitch1312-1320/pattern0 | 1.003448× | 0.993174× | 2/6 | 97.620690× |
| slice8x16/pitch1312-1320/pattern1 | 1.001723× | 0.990365× | 1/6 | 97.518952× |
| slice8x16/pitch1312-1320/pattern2 | 1.008654× | 1.006224× | 0/6 | 98.234683× |
| slice8x16/pitch32-37/pattern0 | 1.008322× | 0.996571× | 2/6 | 98.058252× |
| slice8x16/pitch32-37/pattern1 | 1.006572× | 0.996540× | 1/6 | 97.855413× |
| slice8x16/pitch32-37/pattern2 | 1.006218× | 0.975790× | 1/6 | 97.996546× |
| slice8x16/pitch8-8/pattern0 | 1.000344× | 0.994513× | 3/6 | 98.139855× |
| slice8x16/pitch8-8/pattern1 | 0.998624× | 0.990404× | 3/6 | 98.176815× |
| slice8x16/pitch8-8/pattern2 | 1.005866× | 0.997262× | 1/6 | 97.584541× |
| slice8x4/pitch1312-1320/pattern0 | 0.990453× | 0.980109× | 5/6 | 1.257380× |
| slice8x4/pitch1312-1320/pattern1 | 0.996978× | 0.978819× | 3/6 | 1.261726× |
| slice8x4/pitch1312-1320/pattern2 | 0.993459× | 0.985712× | 3/6 | 1.254544× |
| slice8x4/pitch32-37/pattern0 | 0.999684× | 0.844703× | 4/6 | 1.260546× |
| slice8x4/pitch32-37/pattern1 | 0.996850× | 0.841438× | 3/6 | 1.255055× |
| slice8x4/pitch32-37/pattern2 | 0.994219× | 0.972663× | 5/6 | 1.262993× |
| slice8x4/pitch8-8/pattern0 | 0.987038× | 0.905697× | 5/6 | 1.244728× |
| slice8x4/pitch8-8/pattern1 | 0.992024× | 0.970158× | 6/6 | 1.253549× |
| slice8x4/pitch8-8/pattern2 | 0.995906× | 0.913823× | 6/6 | 1.255574× |
| slice8x8/pitch1312-1320/pattern0 | 0.996095× | 0.989529× | 5/6 | 1.243835× |
| slice8x8/pitch1312-1320/pattern1 | 0.995900× | 0.984378× | 5/6 | 1.231936× |
| slice8x8/pitch1312-1320/pattern2 | 0.999846× | 0.991327× | 2/6 | 1.234492× |
| slice8x8/pitch32-37/pattern0 | 0.999538× | 0.985343× | 4/6 | 1.233580× |
| slice8x8/pitch32-37/pattern1 | 0.999485× | 0.974743× | 3/6 | 1.243755× |
| slice8x8/pitch32-37/pattern2 | 0.999383× | 0.926743× | 3/6 | 1.242089× |
| slice8x8/pitch8-8/pattern0 | 1.001436× | 0.996731× | 3/6 | 1.233523× |
| slice8x8/pitch8-8/pattern1 | 0.998975× | 0.984215× | 3/6 | 1.237319× |
| slice8x8/pitch8-8/pattern2 | 1.001852× | 0.990021× | 3/6 | 1.244211× |

The97.519–98.235×oldASM/current16-row control ratios belong only to this Intel8573C/publicB.Loop call protocol. They do not establish a broad98×codec or newPR43 gain. AVX/SSE transition interaction is a hypothesis; this audit establishes no cause. Oldreference uses actualSSE2 for4/8/pointer andAVX2 for16rows.

## Four separate legacy B.N diagnostics

| Fixture | Parent/head | Pairs slower | Old ASM/head |
|---|---:|---:|---:|
| BenchmarkSSE8x16 | 0.996302× | 5/6 | 104.326923× |
| BenchmarkSSE8x4 | 0.979732× | 5/6 | 0.063522× |
| BenchmarkSSE8x8 | 0.996621× | 4/6 | 1.249642× |
| BenchmarkSSE8x8PtrFast | 0.997052× | 4/6 | 1.118239× |

Original legacy functions discard their returned SSE values and lack the new observable fixture sink. They remain separate weak diagnostics, including unusual old ratios, and are not pooled with live B.Loop acceptance or used to infer an optimization cause. The livefixture scalar/canary/sink checks verify actual outputs before/after each sample.

## Preservation

focused-raw.tar.gz preserves all45files byte for byte, including originalSHA256SUMS. Independent root proof is byte-identical in root-proof.json; audit.json preserves every sample/ratio, one-shot output/allocation, controlregression, sourcehash and feature/test result. raw-sha256.json andsha256.json cover members andouter files. No Go job, source/README/checkout/index/branch change occurred during audit.

## Detailed evidence

`raw-evidence.tar.gz` preserves the complete original frozen evidence package byte for byte, including all audits, root proofs, source hashes, disassembly, original manifests and nested raw benchmark archives. Extract into a separate directory to inspect and verify its original manifest.
