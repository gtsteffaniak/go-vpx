# Native FP guard diagnostic: independent audit PASS

Head e2fffb9f649b4a757c285671899eddf624e3013a; parent5c5cf3977052368746325c4033a1909e7eacd26a. Native run39133345380/job117692814590 succeeded. Prior39132874998 failed resolving a branch name before builds/timing; its failed log is retained.

AMD EPYC7763, official Go1.27.2, Linux/amd64, CGO1, PGOoff, GOMAXPROCS1. ActualAVX2/Enabled=true, AVX512=false. No selected512 result is claimed.

All420raw files and producer manifest revalidated. All Git source maps, sole compiled quantizer guard-reorder delta, source/binary receipts, selected/race/checkptr/off tests,168excluded1x oracles,1008zero-allocation component samples, sixbalancedAB/BAorders,12whole samples,48excluded diagnostic profiles and two excluded dispatch traces independently pass. All22semantic fields match. Every slower component/pair and allocation sample remains in raw/summary/strict-proof.

Component medians span0.986109–1.013681×;30/84medians are slower. Whole320×180 VP9/120inputs median time −0.044985%, fourof sixpairs faster; paired changes +0.230942%,−0.080256%,−0.079747%,−0.346790%,+0.088472%,−0.466894%. All12timed allocation rates0.616666667; all48profile rates are preserved separately.

Native GNU code verifies unchanged0xa8 SIMD frame. For n64 on this non512CPU both guards execute two rejection instructions (featurecompare/branch versus sizecompare/branch). For n256 the size-first candidate executes six guard instructions versus two in parent. Thus the prior hypothetical six-to-two non512 claim does not describe this reverse experiment. Excluded trace sees106998n16entries and459144eligible n64calls; zero n256calls/selected512. The459144avoided feature queries did not establish a throughput win.

Merged24profiles/variant report quantizeFPBlockSIMD flat0.15s→0.16s and cumulative0.23s→0.23s. Sampling is diagnostic; near-flat throughput does not establish the historical+2.264%regression cause. No production acceptance or root edits performed.

This lean package omits executable binaries/full snapshots. Exact binary hashes, sizes/address maps, build information, machine code and source receipts remain.


Publication: `raw-evidence.tar.gz` preserves all 431 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
