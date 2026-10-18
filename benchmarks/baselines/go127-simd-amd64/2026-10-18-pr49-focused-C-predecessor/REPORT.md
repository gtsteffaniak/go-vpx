# PR49 failed5c37 predecessor: focused/C/hosted evidence

Exact head `5c37b7d742c636ab6b7a0ddeb72e0339b5c09a00`; base `c8bd67df9a00c102163354158c1ee46b0f1527a7`. This is preserved failed-predecessor evidence, not four green checks or merge acceptance. The repaired head is `cb9bbdcd185cee2026d178e23b51bb6939335d75`; it requires new exact-head checks.

The focused EPYC9V74 SDK exposes AVX512/AVX2/Enabled true. Selected512 ordinary and race oracles both PASS; corresponding512-off tests PASS. The fifth full-AVX-off invocation disables AVX/AVX2 but leaves raw SDK AVX512 true, causing the strict off assertion to fail. Effective production Enabled/AVX2 are false. Sixth validation and optional pixel timing phase never execute. No selected512 performance evidence exists here. The next repair explicitly disables cpu.avx512f in all four full-AVX-off shell commands while retaining strict assertions.

The original SSE8 phase completed720 timings+120 excluded one-shot samples, all zero allocations, with exact303/294/147 source hashes and actual pointer loops28→28. All36 fixtures are unchanged controls. Parent/current medians .980723–1.015052×;21 median regressions, two all-six regressions, no all-six win; minimum paired ratio .908322×. The sustained regressions are slice8x8 pitch1312-1320 pattern0 (+.943396% time) and pitch32-37 pattern1 (+1.485343%). All paired data and old-reference comparisons are preserved; no causality or broad codec claim.

Separate optimized-C job on EPYC9V74 succeeds. Fixed-speed VP8 CPU−8/minq2/maxq56 Go/C time2.249773445× matches each timed output and separate58-frame decoded PSNR/SSIM. Fixed speed is not fixed quantizer. Adaptive1.618616915×, VP9320 2.677437040× and7202.489726207× retain unequal work/quality. Each side's separate quality-output statistics match all timings. All32 helper controls/timers/order/allocation checks pass; PSNR/SSIM are separate-pass measurements, not per-timing measurements. Ratios above one mean libvpx faster; jobs are not pooled.

The exhaustive receipt records the state at capture in scope.json; it does not change failed full/focused status. All134 raw files, including empty logs and original checksums, plus complete audits/console/receipts/source hashes are preserved and decompression-verified. Full-native partial evidence is owned separately by ci_fix. No executables, new tests, builds or encode/timing reruns were used to preserve this package.


Publication: `raw-evidence.tar.gz` preserves all 5 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
