# PR39 ARM Go/libvpx comparison

Fresh native Apple M4 Max / Mac16,6, Go1.27.2 archsimd, PGOoff, GOMAXPROCS=1, one encoder thread. Go head968b518bc20d599606d5dccf03532bebcc4f3b4f. Pristine optimized libvpx v1.16.0 / 1024874c5919305883187e2953de8fcb4c3d7fa6. All1,250 current C source files revalidated against the pinned commit, unchanged reference binary522cae89 from the prior archived ARM build. Actual linked libvpx arm_cpu_caps() returns7 (NEON/DOTPROD/I8MM) both inside and outside the sandbox, with no forced capability environment. Build macros/object inventory and physical CPU sysctls also confirm these capabilities.

Six timing rounds alternate Go-first/reference-first three times; one unmeasured warmup per case, separate decoded PSNR/SSIM pass. Go timing sums EncodeInto calls; C timing is vpxenc cumulative vpx_codec_encode time. Wall/subprocess overhead and repeated-mean latency fields are preserved but excluded from performance ratios. All raw timing/quality JSON, stderr, controls and provenance remain available. No competing local Go builds/tests/profiles ran during measurements.

| Case | Go time / libvpx time | Measured quality/output scope |
|---|---:|---|
| vp9-320p | 2.512025× | Different PSNR, SSIM and payload; unequal output/quality work. |
| vp9-720p | 2.481818× | Different PSNR, SSIM and payload; unequal output/quality work. |
| vp8-720p-adaptive | 1.581560× | PSNR, SSIM, payload and frame counts match. |
| vp8-720p-fixed | 2.104120× | PSNR, SSIM, payload and frame counts match. |

VP8 is1280x720/60inputs/2500kbps; fixed CPU-used=-8 and adaptive CPU-used=8 remain distinct workloads, with actual matching C flags recorded in every report. VP9 is320x180/120inputs/500kbps and1280x720/60inputs/2500kbps, CPU8. These are matched same-input/controls comparisons; matching measured statistics does not assert bitstream identity. VP9 ratios are standings for unequal output/quality work, not equivalent-work speedups. The adaptive VP8 quality/output equality applies only to this exact run; prior AMD adaptive drift/differences remain preserved separately.

All six Go allocation samples are zero for both VP8 cases; VP9 existing allocations are.6166667 and1.25 per input. Both sides’ quality-pass output statistics match ALL six timing outputs in all four cases. All performance ratios above are TIME ratios: greater than1 means Go is slower. No cross-CPU/job comparison or additive optimization-gain claim follows. The fixed VP8 Go/C gap remains open; this benchmark does not turn the separate Go/old-assembly gap into a C-equivalent comparison.

Reference build scripts/configuration are reused unchanged and labeled as provenance. Every actual timing/quality report and linked-runtime-capability probe is new for this head. The Go files are byte-identical to the locally validated residual candidate after the main rebase; no production source or README was changed for this measurement.
