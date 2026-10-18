# PR50 optimized libvpx comparison

Exact Go head `10caaa67ed53f8c3c1248f20070e5a4117245852`; successful C job117659832238 in run39122180995. This is the separate native EPYC7763 job, not the full-native or focused runner. Go1.27.2/archsimd, GOMAXPROCS1/encoderthreads1/PGOoff are recorded. Exact Git source hashes cover50 CLI/driver/workflow files; workflow checks out the PR head and the driver builds directly without source overlays. Recorded Go/C executable SHA metadata is preserved; executable bytes are omitted and cannot be independently rehashed from these raw files.

Pristine libvpx v1.16.0 pin `1024874c5919305883187e2953de8fcb4c3d7fa6` is downloaded from that exact archive endpoint, without oracle patches. Archive SHA/build/configuration/object inventory and recorded vpxenc SHA are retained. Actual CPU flags expose SSE2/AVX2; vpx_config/config.mk enable optimized assembly, SSE2/AVX2 and runtime dispatch. No separate linked C capability probe exists in this artifact; the report does not invent one.

| Workload | Go time / libvpx time | Scope |
|---|---:|---|
| VP8 CPU−8 | 2.262491245× | Each timed output and separate58-frame PSNR/SSIM match |
| VP8 adaptiveCPU8 | 1.680179244× | Unequal payload/quality work |
| VP9 320 | 2.868383359× | Unequal payload/quality work |
| VP9 720 | 2.826897097× | Unequal payload/quality work |

**These are Go/C TIME ratios: Go takes more time and libvpx is faster.** Fixed VP8 means fixed CPU speed−8 with minq2/maxq56, not fixed quantizer. Six measured pairs per workload balance three Go-first and three C-first calls with rotating workload order, plus four discarded warmups and four separate decoded-quality passes. Go sums EncodeInto time; C uses cumulative codec-internal vpxenc-stat timing. Wall/subprocess costs and repeated-mean latency placeholders are not performance denominators.

Every one of32 helper controls, orders and timers is checked. Each side's separate quality-pass output statistics match all six timed outputs; PSNR/SSIM are measured ONLY in separate passes, not in every timing process. Fixed VP8 additionally matches mutual measured payload/counts/PSNR/SSIM; that is not bitstream identity. Adaptive/VP9 unequal work and every mismatching field are preserved rather than labelled an equivalent-work speedup. No comparison pools other CPUs/jobs or infers a codec optimization from different measurement windows.

All32 Go allocation samples match their expected values, including warmups and quality: VP8zero, VP9320 37/60, VP9720 75/60 per input. No exception occurred in this job. The separate PR49 C76/60 sample remains valid in its own immutable archive; this run does not revise it.

All74 original raw files, including empty stderr files, are preserved byte-for-byte. Complete results/output/options/allocations, exact source files, reproducible auditors, provenance and CI receipts are archived and decompression/hash verified. A source-only audit initially looked for a configure flag in the build log; it was corrected to the actual config.mk CONFIGURE_ARGS, without changing raw bytes or rerunning encodes. No builds, quality passes, timing reruns, source or README edits occurred during audit/freezing. The full PR50 benchmark remains pending at this captured receipt; this successful C job alone grants no merge approval.


Publication: `raw-evidence.tar.gz` preserves all 7 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
