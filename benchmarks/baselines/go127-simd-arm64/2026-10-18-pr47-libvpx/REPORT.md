This is a separate Apple M4 Max comparison of exact PR #47 b6209532fbc30cf9567e0cdc78a21d61a60fe900 against optimized libvpx v1.16.0, not the later private TM4, loop-filter or DC experiments.

Both VP8 workloads in this ARM run match the measured output statistics and separately decoded PSNR/SSIM. Quality was measured in separate quality passes; the timed processes were checked against those passes' output statistics. PSNR/SSIM were not measured independently in every timed process. Output statistics do not establish bitstream identity. VP9 output/quality differs, so its time ratios describe the specified workload with unequal output work. The fixed CPU-used VP8 case is not a fixed-quantizer test.

The frozen report's "shared cache" wording refers to the same git-common-dir-derived project identity, not an identical filesystem cache prefix. The C lane used the default /var/folders/.../C/codex-projects/99feb44aac13f74a/go-build prefix, while root agents explicitly use /private/tmp/codex-projects/99feb44aac13f74a/go-build. Neither cache path contains .codex/worktrees. See the accompanying exact cache-scope clarification. Compiles are excluded from all performance denominators; no source, timer, binary or raw result was changed for this clarification.


Frozen report follows; the publication notice above clarifies its quality-pass and cache-prefix scope.

# PR47 native ARM Go/libvpx comparison

Exact Go head `b6209532fbc30cf9567e0cdc78a21d61a60fe900`; Apple M4 Max / Mac16,6, official Go 1.27.2 archsimd, CGO=1, PGO off, GOMAXPROCS=1 and one encoder thread. The immutable archived source was built through the shared project cache wrapper. All 1,925 source files stayed unchanged. Binary SHA256: `69c0ce69921e35f060ac49e2c0216b817890763746f934c18c9140325a335788`.

Pristine optimized libvpx v1.16.0 / `1024874c5919305883187e2953de8fcb4c3d7fa6` reused the existing reference binary. All 1,250 current C source files match the pinned commit; binary, static archive, configuration and object inventory match the archived build. Fresh linked runtime probes report NEON/DOTPROD/I8MM (caps=7) inside and outside the sandbox, with no forced capability override. Oracle-patched encoders were not substituted.

Each workload has six timing pairs, three Go-first and three reference-first, with rotating workload order and discarded warmups. Go sums EncodeInto durations; C uses cumulative vpx_codec_encode statistics. All 24 measured reference timers and eight warmup/quality timers use vpxenc-stats. Fixture construction, subprocess overhead, wall time and repeated-mean latency fields are excluded from ratios. Four decoded-quality passes follow all timing rounds.

| Workload | Go time / libvpx time | Measured quality/output scope |
|---|---:|---|
| vp9-320p | 2.490255310× | Different payload, PSNR and SSIM; unequal work |
| vp9-720p | 2.482958391× | Different payload, PSNR and SSIM; unequal work |
| vp8-720p-adaptive | 1.577637600× | Matching PSNR, SSIM, payload and frame/drop counts; 58 quality frames |
| vp8-720p-fixed | 2.088133921× | Matching PSNR, SSIM, payload and frame/drop counts; 58 quality frames |

A TIME ratio greater than one means libvpx is faster. Fixed VP8 means CPU-used −8 with min-q=2/max-q=56, not a fixed quantizer. Both VP8 cases match measured payload, frame/drop counts, PSNR and SSIM on this exact run. Every timed Go/C output mutually matches these statistics, and each side’s separate quality output matches all six timed outputs. Matching statistics do not establish bitstream identity.

Fixed VP8: PSNR 27.157610871537482, SSIM 0.9259273287259608, payload 723,593 bytes, 58 encoded and two dropped frames. Adaptive VP8: PSNR 27.983950811796223, SSIM 0.9085039479857895, payload 685,957 bytes, 58 encoded and two dropped frames. Adaptive equality is observed here and is not assumed for other runs.

VP9 retains different quality and payload work. Its separate quality output matches every timed output on each side; that does not make Go and C outputs identical. All allocation samples are preserved: VP8 zero, VP9 320p 37/60 per input, VP9 720p 75/60 per input, with no exception in these six rounds. Earlier checkpoint allocation exceptions remain in their own archives.

The archive preserves all 66 original reports, stderr files, commands and summary, plus 17 provenance files. All 83 raw member hashes and outer checksums are verified. Every paired ratio, option, output statistic and allocation remains reviewable. These same-run M4 results are not pooled with earlier ARM or x64 jobs and do not imply additional component gains. All timing, quality and build jobs stopped before the next research window.

Preparation context is retained in `provenance/notes.txt`: an initial archive request named nonexistent go.sum and failed before the build; sandbox CPU sysctl access was unavailable, so a fresh native read was used; Go build information omits the explicit PGO-off flag, which is pinned by the actual command and exit-zero receipt. Build VCS metadata is disabled for the immutable archive, with exact Git head/tree, source, archive and binary hashes recorded externally. No timing samples were rerun or discarded. Binary blobs are omitted from this lean archive; their full hashes are recorded. No root source, README, index, checkout or branch was changed.


Publication: `raw-evidence.tar.gz` preserves all 7 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
