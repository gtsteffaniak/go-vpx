# PR51 optimized libvpx native audit

Exact headbc4a12a2ea74f6219d26acb5646677016eb6423f, run39126944619/job117673948242 SUCCESS. Actual CPU is AMD EPYC9V74. Only the optimized-libvpx sibling is audited here; the main-native sibling, required focused9V45 and manual jobs are separate and never pooled.

All74raw artifact files and SHA inventory are preserved. Source checks bind50exact immutable GoCLI/driver/workflow files. Native Go SIMD selection executes. Pristine libvpxv1.16.0 is pinned to1024874c5919305883187e2953de8fcb4c3d7fa6, with recorded source-archiveSHA and binary hashes. Actual build configuration enables optimization/runtimeCPUdetection/VP8/VP9 and SSE2/AVX2; corresponding optimized assembly objects and CPU flags are present. No separate linked-C runtime-capability probe is in this artifact, so no specific per-kernel dispatch is invented.

All32helper invocations (4warmups+24timings+4qualitypasses) have exact recorded controls and valid timers. Six rounds alternate reference-first, with rotated case order, one thread/GOMAXPROCS1/PGOoff. Every C timer reports vpxenc-stats, never wall-clock fallback; subprocess/file/setup overhead is excluded from that internal timer. Go allocations/input are0forVP8,37/60forVP9320 and75/60forVP9720 on all32rows, with no exceptions. C allocation counts are not measured by this artifact.

Go/C TIME ratios above1 mean C is FASTER. Fixed-speed VP8720 is11.9420055ms/5.9234415ms=2.016058654x: Go takes101.606% more time. CPU-used−8 is fixed speed, NOT a fixed quantizer (minq2/maxq56 remain). Each timed payload/count/bitrate/keyframe/interframe statistic matches both the separate quality output and the other implementation:723593bytes,58encoded/2dropped; decoded58-framePSNR27.15761087 andSSIM0.925927329 match. These are measured output/quality equalities, not bitstream identity or per-timed-encode PSNR measurements.

Adaptive VP8720 Go/C TIME is1.499925698x, with unequal work: Go685957bytes/PSNR27.9839508/SSIM0.908503948 versus C678205bytes/PSNR29.1766391/SSIM0.933088086;58/2counts match. C has higher measured quality and lower payload in this run. VP9320 TIME2.656758380x andVP9720 TIME2.435401020x retain close but unequal payload/quality results and all117/3and48/12counts. No equivalent-quality performance claim is made for these three cases. There is no Go-beats-C claim.

For EACH side of EACH case, the separate quality-pass output statistics match EVERY timed sample in THIS run. No adaptive timed-payload drift is observed here; drift from prior experiments is not imported. All per-round TIME ratios, output/quality differences, actual libvpx flags, allocation values and control arrays remain in compact.json. Separate quality passes exclude encode-only and time additional decoding; they are never mixed into encode-only timing medians.

Allraw/config/object/source/command/quality/error logs, audit source, recorded binary/archive hashes and exactCLI sources are archived. Executable blobs/full source snapshots are omitted; recorded hashes are not misrepresented as independent rehashes of absent binaries. PriorPR50auditor templates/history are clearly separate from current passed checks. No ROOT/README/Git changes, local builds, encodes, reruns or performance thresholds were introduced by this audit.


Publication: `raw-evidence.tar.gz` preserves all 5 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
