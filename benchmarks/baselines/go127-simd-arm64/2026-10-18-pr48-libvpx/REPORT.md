# ARM48 integrated Go / optimized libvpx comparison

This measures the exact current ARM48 integrated snapshot over merged PR47 `d2610cd0747edd72f7656068dfd4604d29e45307`, before a PR48 Git commit existed. All 1,937 Go files match the compiled snapshot/build manifest and current ROOT at preflight, with exactly 13 production and six test changes. Final publication must add a Git-to-snapshot binding; this report does not pretend the binary was compiled from a future commit.

Native Apple M4 Max / Mac16,6, official Go 1.27.2 archsimd, CGO=1, PGO off, GOMAXPROCS=1, one encoder thread. Val’s prebuilt CLI SHA256 is `c5c4395d55ece52bfeb2c1d9dd4d19cef90337477117629d0085bf50f853b7ae`; the preserved build command exited zero and uses the explicit `/private/tmp/codex-projects` project cache. No Go or C rebuild ran for this comparison.

Pristine optimized libvpx v1.16.0 / `1024874c5919305883187e2953de8fcb4c3d7fa6` reuses the unchanged binary `522cae89…` and static archive `84e97dc4…`. Fresh verification covers all 1,250 pinned C source files, configuration, optimized object inventory and binary/archive hashes. Fresh linked probes report caps=7 (NEON/DOTPROD/I8MM), inside and outside the sandbox, with no forced capability environment. Oracle-patched encoders were not substituted.

Each case has one discarded warmup, six timing pairs with three Go-first/three C-first orders and rotating workload order, then a separate decoded-quality pass after all timings. All 24 measured reference timers and eight warmup/quality timers use cumulative vpx_codec_encode statistics. Go sums EncodeInto time. Fixture construction, wall/subprocess overhead and repeated-mean latency placeholders are excluded from performance ratios.

| Workload | Go time / libvpx time | Measured output/quality scope |
|---|---:|---|
| vp9-320p | 2.434797079× | Different payload, PSNR and SSIM; unequal work |
| vp9-720p | 2.540569955× | Different payload, PSNR and SSIM; unequal work |
| vp8-720p-adaptive | 1.535884623× | Matching timed output statistics and separately decoded PSNR/SSIM; 58 quality frames |
| vp8-720p-fixed | 2.045833301× | Matching timed output statistics and separately decoded PSNR/SSIM; 58 quality frames |

Ratios above one mean libvpx is faster. Fixed VP8 denotes CPU-used −8 with min-q=2/max-q=56, not a fixed quantizer. Both VP8 workloads freshly match measured Go/C payload, frame/drop counts and separately decoded PSNR/SSIM. Every timed VP8 pair mutually matches output statistics, and each side’s separate quality-pass output statistics match every timed output. PSNR/SSIM were measured only in separate quality passes, not independently in each timed process. Matching statistics do not establish bitstream identity.

VP9 retains different payload and decoded quality work. Each side’s quality output statistics still match all its timed outputs; that does not make the Go/C workloads equivalent. All six Go allocation samples are zero for both VP8 cases, 37/60 per input for VP9 320p and 75/60 for VP9 720p. No allocation exception occurred in this run; all samples and mismatches remain in audit.json and original raw JSON.

This is a separate same-run M4 comparison. It is not pooled with PR47, x64 jobs or other research stages, and no component gain or causal improvement is inferred from comparing those stages. All 66 original reports/stderr/commands/summary files are preserved byte-for-byte, together with provenance, driver and progress logs, source/build manifests and runtime/configuration proof. Every decompressed raw member and outer checksum is verified. Executable blobs are omitted; complete hashes remain recorded. No source, README, Git/index/branch, build or timing rerun occurred during this lane. All jobs stopped after quality END.


Publication: `raw-evidence.tar.gz` preserves all 12 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.

Publication binding: all1,937 compiled Go source blobs match source commit `d05c390826ad56012abf6ae501470729a6aaab88` exactly; see `compiled-source-binding.json`. Frozen inner notes saying final Git binding is pending describe their original pre-publication state. Subsequent evidence-only commits keep those runtime sources.
