# ARM49 Go / pristine optimized libvpx

Exact KEEP158 prebuilt CLI `62f174c3a7876a3a93be53150c42c7c2a96c39a17b76db1ccfbcdb697bda531f` binds all1946 compiled Go files to code commit `c10edbcd6753743331f0ad92b55e94737f747c85`. A future evidence-only PR head must preserve those bytes. Fresh pristine libvpx v1.16.0 pin `1024874c5919305883187e2953de8fcb4c3d7fa6` source1250/config/object inventory/binary522cae89…/static84e97dc4… and native caps7 NEON/DOTPROD/I8MM are verified. Apple M4 Max/Mac16,6; official Go1.27.2 CGO1 archsimd PGOoff GOMAXPROCS1/threads1. Build concurrency GOMAX2 is recorded separately; no build occurs in this comparison.

Four discarded warmups,24balanced Go/C timing pairs and four separate decoded-quality passes use the historical same-input protocol. C timing is cumulative codec-internal vpxenc statistics; subprocess/wall and repeated-mean legacy latency placeholders are not performance denominators. All32 helper reports and commands/stderr are retained unchanged.

| Workload | Go time / libvpx time | Quality/output scope |
|---|---:|---|
| vp9-320p | 2.485216339× | Different payload and decoded quality; unequal work |
| vp9-720p | 2.508104154× | Different payload and decoded quality; unequal work |
| vp8-720p-adaptive | 1.497734717× | Matched timed output stats and separate58-frame decoded quality |
| vp8-720p-fixed | 2.035751548× | Matched timed output stats and separate58-frame decoded quality |

Ratios above one mean libvpx is faster. Fixed VP8 denotes CPU-used−8 with minq2/maxq56, not fixed quantizer. Both VP8 cases freshly match measured output/counts and separate PSNR/SSIM; every timed output matches its side’s quality-pass output statistics. PSNR/SSIM are measured ONLY in separate quality passes, not independently in each timed encode. Matching aggregate statistics are not bitstream identity.

VP9 retains unequal payload and reconstruction quality despite equal encoded/drop counts. The audit preserves every mismatching field and paired sample. All32 Go allocation samples are0 for VP8,37/60 for VP9 320 and75/60 for VP9 720; no allocation exception occurred. Raw reports preserve every other timed field. No pooling with PR48/x64 or other research stages, and no causal optimization claim follows from comparing separate windows.

The plan documents that the current bench cannot request constant Q or guaranteed equal quality. Noise0 and Q-ladder diagnostics remain proposals and were not executed. Source/binary checks and separate ROOT code-binding close the earlier pending-commit preparation scope; the original proof text remains immutable. Root independent proof was generated before the code commit and therefore retains its historical pending-publication wording; publication-code-binding.json supersedes only that identity status.

No executable/full snapshot is included. Complete source/binary/config hashes, drivers, raw stderr including empty files, preparation history, source proof, independent audits and execution attestation are preserved. Audit/freezing performed no builds, warmups, encodes, quality or timing reruns. All decompressed member hashes and outer checksums verify.


Publication: `raw-evidence.tar.gz` preserves all 11 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
