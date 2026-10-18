# Latest-main quick native libvpx comparison — 2026-10-18

Runtime source: `715645fd212cfc65295f4f6fb1002dcc60856758` (latest main at dispatch). Harness: `45c38dbc0d49470dd2df1c6b74275948900bc18c`; no Go source, go.mod or go.sum differences from that main. Reference: pristine optimized libvpx v1.16.0, commit `1024874c5919305883187e2953de8fcb4c3d7fa6`.

Preliminary short screen requested explicitly by the user. Four synthetic public encoder workloads on each native host, one discarded warmup and two timing rounds per workload, one Go-first and one C-first. Separate decoded PSNR/SSIM passes check the same input and controls after timing. No full CI, broad performance matrix, profile audit or new research ran.

## Controls and measurement

Go 1.27.2, GOEXPERIMENT=simd, CGO_ENABLED=1, PGO off, GOMAXPROCS=1, encoder threads=1, realtime CBR, 30 fps. VP9: 320×180/120 inputs/500 kbps and 1280×720/60 inputs/2500 kbps, cpu-used=8. VP8: 1280×720/60 inputs/2500 kbps, adaptive cpu-used=8 and fixed-speed cpu-used=-8. Fixed-speed still uses matched CBR quantizer bounds 2–56. Each raw JSON records resolved Go options and actual C parity flags, including denoising, dropping, buffering and keyframe controls. The same deterministic generated input is used for both encoders.

Go timing sums EncodeInto calls; C timing uses the required vpxenc internal cumulative codec encode timer. Subprocess, input/file work, allocation sampling, quality decoding and warmup are excluded from headline timings. Mean input-frame cost includes the same submitted inputs, including dropped frames. Latency quantiles from the C helper are not used.

Ratios below are Go encode time / C encode time, equivalently C throughput / Go throughput. Above 1 means libvpx is faster. Two samples are preliminary: both paired ratios are retained; no significance claim or cross-host pooling is made.

Quality equality means identical measured PSNR/SSIM, output sizes/frame-size statistics, encoded/drop counts and quality-frame counts. This measurement does not establish packet-byte identity. Cases with different metrics remain explicitly labeled; the existing broad quality regression threshold is not used to declare iso-quality equivalence. No tolerance was widened or chosen after seeing results.

ARM reuses an already verified optimized C build and source/binary bindings; its config enables NEON, dot-product, I8MM and runtime CPU detection. Prior linked capability evidence for the exact retained binary/archive hashes reports NEON=1, DOTPROD=1, I8MM=1. VPX_SIMD_CAPS and VPX_SIMD_CAPS_MASK are unset. The x64 run builds the pinned C reference with assembly, SSE2/AVX2 and runtime CPU detection, with configuration and optimized archive objects retained.

## Timing results

| Host | Workload | Go ms/input | C ms/input | Go/C time | Paired Go/C ratios | Quality scope |
|---|---|---:|---:|---:|---|---|
| Apple M4 Max ARM64 | vp9-320p | 1.4627 | 0.6367 | 2.2975× | 2.2872/2.3079 | Different quality; equal controls only |
| Apple M4 Max ARM64 | vp9-720p | 10.4974 | 4.8380 | 2.1698× | 2.1747/2.1651 | Different quality; equal controls only |
| Apple M4 Max ARM64 | vp8-720p-adaptive | 7.2750 | 5.0239 | 1.4481× | 1.4492/1.4470 | Equal measured quality/output statistics |
| Apple M4 Max ARM64 | vp8-720p-fixed | 6.1795 | 3.2540 | 1.8990× | 1.8292/1.9734 | Equal measured quality/output statistics |
| AMD EPYC 7763 x64 VM | vp9-320p | 3.9073 | 1.4312 | 2.7301× | 2.7449/2.7156 | Different quality; equal controls only |
| AMD EPYC 7763 x64 VM | vp9-720p | 25.5915 | 10.5746 | 2.4201× | 2.4462/2.3946 | Different quality; equal controls only |
| AMD EPYC 7763 x64 VM | vp8-720p-adaptive | 17.6362 | 10.8984 | 1.6182× | 1.6353/1.6016 | Different quality; equal controls only |
| AMD EPYC 7763 x64 VM | vp8-720p-fixed | 15.2689 | 7.4987 | 2.0362× | 2.0469/2.0256 | Equal measured quality/output statistics |

The defensible exact measured-quality VP8 fixed-speed gap is **1.90× on ARM and 2.04× on x64**. ARM adaptive VP8 also matches measured quality, with a **1.45×** gap. The x64 adaptive row is not an iso-quality gap: Go has substantially lower PSNR/SSIM and a larger output. VP9 is close in measured quality, but its timing ratios are equal-controls observations rather than exact iso-quality claims.

## Quality and output checks

Deltas are Go minus C. Output and quality statistics from every separate quality pass equal their respective timing runs. Encoded and dropped-frame counts agree between Go and C in all eight rows.

| Host | Workload | Go/C PSNR (dB) | ΔPSNR | Go/C SSIM | ΔSSIM | Δbytes | Encoded/dropped |
|---|---|---|---:|---|---:|---:|---|
| arm64 | vp9-320p | 28.91097432/28.96958290 | -0.05860857 | 0.966373405/0.967011486 | -0.000638081 | -19 | 117/3 |
| arm64 | vp9-720p | 27.51643367/27.54611663 | -0.02968295 | 0.949347795/0.949095623 | +0.000252172 | -405 | 48/12 |
| arm64 | vp8-720p-adaptive | 27.98395081/27.98395081 | +0.00000000 | 0.908503948/0.908503948 | +0.000000000 | +0 | 58/2 |
| arm64 | vp8-720p-fixed | 27.15761087/27.15761087 | +0.00000000 | 0.925927329/0.925927329 | +0.000000000 | +0 | 58/2 |
| amd64 | vp9-320p | 28.91097432/28.96958290 | -0.05860857 | 0.966373405/0.967011486 | -0.000638081 | -19 | 117/3 |
| amd64 | vp9-720p | 27.51643367/27.54611663 | -0.02968295 | 0.949347795/0.949095623 | +0.000252172 | -405 | 48/12 |
| amd64 | vp8-720p-adaptive | 27.98395081/29.17663906 | -1.19268824 | 0.908503948/0.933088086 | -0.024584138 | +7752 | 58/2 |
| amd64 | vp8-720p-fixed | 27.15761087/27.15761087 | +0.00000000 | 0.925927329/0.925927329 | +0.000000000 | +0 | 58/2 |

## Evidence and limits

- `arm64-summary.json` and `amd64-summary.json` contain full precision paired timings, all actual C flags, quality, output and allocation observations. The raw archives retain command manifests, discarded warmups, both timed rounds, separate quality runs, stderr, compiler/configuration and object-list evidence. Binary and source bindings are preserved.
- Native x64 run: [39153782179](https://github.com/thesyncim/govpx/actions/runs/39153782179), C-only manual dispatch. Only libvpx-reference executed; the whole performance matrix and FP quick job were skipped. The lightweight SIMD selection check passed.
- The x64 environment log retains an old static phrase saying “six rounds”; the executed driver summary and 16-command manifest establish four warmups, eight measured invocations (two per case) and four quality invocations. That stale text is not the execution count.
- x64 VP9 720p allocation observations were 75/60 and 76/60 allocations per input; both are retained. VP8 observations are zero. This is a report of the current build, not an allocation regression verdict.
- ARM raw evidence includes one separately labeled VP9 320p noise=0 diagnostic quality probe. It reproduced the same quality/output delta, was excluded from headline timings, and was not pooled with canonical noise=4 cases.
- Timings are host/workload specific and preliminary. In particular ARM fixed-speed paired ratios span 1.83–1.97×. No CPU-to-CPU speed claim, confidence interval, BD-rate result, decoder performance result or general real-video corpus conclusion is established.
- README and production Go source were unchanged. No full CI ran.
