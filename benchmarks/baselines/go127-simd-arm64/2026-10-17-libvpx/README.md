# Native ARM optimized libvpx comparison

Apple M4 Max, Go 1.27.2, SIMD enabled, PGO off, one Go execution thread and one encoder thread. Codec sources are `211125e5`; the binary additionally includes the reference-order harness committed with this evidence. Pristine libvpx v1.16.0 enables NEON, dot product, I8MM and runtime CPU dispatch. See the source validation manifest, compiler/configuration, object list and binary hashes.

| Workload | Go ms/input | C ms/input | Go time delta | Encoded/dropped, both |
|---|---:|---:|---:|---|
| vp9-320p | 1.644036 | 0.637233 | +158.00% | 117/3 |
| vp9-720p | 12.588807 | 4.842000 | +159.99% | 48/12 |
| vp8-720p-adaptive | 8.471075 | 5.115767 | +65.59% | 58/2 |
| vp8-720p-fixed | 6.996804 | 3.220525 | +117.26% | 58/2 |

Six runs per workload alternate Go-first/reference-first, with discarded warmups and a separate decoded quality pass. The shared driver records the exact same generated inputs and corresponding CBR, quantizer, GOP, buffering, drop, denoising and speed controls. Timings sum Go EncodeInto calls and stock vpxenc internal vpx_codec_encode calls; subprocess/file wall times are recorded separately. Input frame counts are the denominator, including dropped inputs. The helper’s legacy reference latency quantile fields repeat the mean and are not used.

VP8 adaptive `cpu-used=8` and fixed `cpu-used=-8` are distinct workloads. Within each, both codecs have identical measured PSNR, SSIM, payload size and encoded/drop counts across all samples. This alone does not prove bitstream identity. VP9 outputs and measured quality differ slightly:

| Workload | PSNR Go/C | SSIM Go/C | Payload Go/C |
|---|---|---|---|
| vp9-320p | 28.910974/28.969583 | 0.966373405/0.967011486 | 246893/246912 |
| vp9-720p | 27.516434/27.546117 | 0.949347795/0.949095623 | 642654/643059 |

VP9 is a matched-input/control comparison, not an assertion of identical output or coding work. All six paired Go/C ratios exceed zero in each workload; raw samples and quality details are retained in summary.json and individual result files. These results show a substantial remaining gap against optimized C even where the Go implementation beats the previous Go assembly baseline.

Reproduce with scripts/benchmark-libvpx.py and an optimized pristine dual-codec vpxenc; build-arm.sh records this ARM build. Native x64 CI now builds the same immutable upstream version independently and retains corresponding evidence.
