# Bounded SIMD residual subtraction: exact-head native audit

Tested `968b518bc20d599606d5dccf03532bebcc4f3b4f` against main `8ee816b14603fd34dccbcb629f021cf5ee57e396`, old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`. Native CPU **EPYC 9V45**, Go 1.27.2, PGO off, six balanced orders. All three checks pass. All 330 kernels have six samples per variant and zero bytes/allocations; fourteen source hashes match. There are 123 targeted checks / 268 PASS lines, no failures/skips, including all new subtraction tests selected and race/checkptr.

All 48 changed public fixtures improve versus the parent SIMD in all 288 pairs: median **1.939–6.626x**, minimum pair **1.771x**. Every pair exceeds 1.5x. The old AMD reference uses **SCALAR FALLBACK FOR ALL 48 CASES**, because its fast selector always returns false,false; old-reference ratios 1.697–19.379x do not measure active x64 assembly. [Pinned old selector](https://github.com/thesyncim/govpx/blob/691cd0512c4829d9ee40d483205dd223f9eb028d/internal/vp9/dsp/subtract_amd64.go). ARM reference behavior differs: contiguous output uses NEON, padded output scalar.

Whole VP9 320 takes **+.463843% more time versus parent** (3/6 faster), while VP9 720 improves **-2.504218%** (5/6). They beat the old reference build 45.4753%/52.8351% on this host. **VP8 is an unchanged-source control: +3.574231% versus parent** (1/6 faster), +1.887788% versus old build. No VP8 benefit or cross-CPU gap-closure claim.

All 22 semantic fields match in each public/profile dataset. Retain old-reference VP9 720 round4 allocation 1.266667 versus modal 1.25; current,parent and profile allocation samples match their controls. All 167 median regressions and eight all-six-slower subpixel controls remain, along with all 282 prior kernel cases. No subtraction fixture regresses.

| Fixture | Parent ns | Current ns | Parent/current | Minimum pair | Old SCALAR/current |
|---|---:|---:|---:|---:|---:|
| 16x16/paddedfalse/pattern0 | 82.7900 | 22.3450 | 3.7051x | 3.4993x | 8.2972x |
| 16x16/paddedfalse/pattern1 | 82.9150 | 22.2350 | 3.7290x | 3.5647x | 8.4124x |
| 16x16/paddedfalse/pattern2 | 82.8100 | 22.1650 | 3.7361x | 3.4586x | 8.4006x |
| 16x16/paddedtrue/pattern0 | 83.1950 | 22.0700 | 3.7696x | 3.4267x | 8.4028x |
| 16x16/paddedtrue/pattern1 | 83.0650 | 22.1700 | 3.7467x | 3.4460x | 8.3897x |
| 16x16/paddedtrue/pattern2 | 82.8950 | 22.1550 | 3.7416x | 3.4032x | 8.4450x |
| 16x8/paddedfalse/pattern0 | 44.2800 | 16.6500 | 2.6595x | 2.4292x | 5.4523x |
| 16x8/paddedfalse/pattern1 | 44.3850 | 16.6250 | 2.6698x | 2.4267x | 5.4328x |
| 16x8/paddedfalse/pattern2 | 44.3800 | 16.6650 | 2.6631x | 2.5878x | 5.4047x |
| 16x8/paddedtrue/pattern0 | 44.3550 | 16.6250 | 2.6680x | 2.5861x | 5.4063x |
| 16x8/paddedtrue/pattern1 | 44.2600 | 16.6800 | 2.6535x | 2.4959x | 5.4134x |
| 16x8/paddedtrue/pattern2 | 44.2750 | 16.6300 | 2.6624x | 2.4319x | 5.3743x |
| 32x16/paddedfalse/pattern0 | 128.8000 | 28.0400 | 4.5934x | 4.1190x | 12.0596x |
| 32x16/paddedfalse/pattern1 | 128.9000 | 28.0150 | 4.6011x | 4.2161x | 12.0667x |
| 32x16/paddedfalse/pattern2 | 128.5000 | 28.2600 | 4.5471x | 4.2159x | 12.0913x |
| 32x16/paddedtrue/pattern0 | 128.3000 | 27.9750 | 4.5862x | 4.3279x | 12.2109x |
| 32x16/paddedtrue/pattern1 | 128.0500 | 28.1200 | 4.5537x | 4.3039x | 12.0928x |
| 32x16/paddedtrue/pattern2 | 128.2500 | 28.0750 | 4.5681x | 4.3296x | 12.1692x |
| 32x32/paddedfalse/pattern0 | 259.4500 | 45.6600 | 5.6822x | 5.3942x | 14.6167x |
| 32x32/paddedfalse/pattern1 | 258.3000 | 45.1750 | 5.7178x | 5.2531x | 14.6718x |
| 32x32/paddedfalse/pattern2 | 251.5500 | 46.3550 | 5.4266x | 5.2089x | 14.4148x |
| 32x32/paddedtrue/pattern0 | 246.8500 | 44.6550 | 5.5279x | 5.2994x | 14.9804x |
| 32x32/paddedtrue/pattern1 | 251.4000 | 44.8100 | 5.6104x | 5.3269x | 14.8773x |
| 32x32/paddedtrue/pattern2 | 253.0000 | 46.3850 | 5.4543x | 4.9530x | 14.3139x |
| 4x4/paddedfalse/pattern0 | 20.1700 | 10.4000 | 1.9394x | 1.8368x | 1.7413x |
| 4x4/paddedfalse/pattern1 | 20.5750 | 10.4150 | 1.9755x | 1.9141x | 1.7484x |
| 4x4/paddedfalse/pattern2 | 20.5350 | 10.4900 | 1.9576x | 1.9303x | 1.7173x |
| 4x4/paddedtrue/pattern0 | 20.4450 | 10.4200 | 1.9621x | 1.8636x | 1.7423x |
| 4x4/paddedtrue/pattern1 | 20.5000 | 10.5300 | 1.9468x | 1.7711x | 1.6966x |
| 4x4/paddedtrue/pattern2 | 20.5050 | 10.4600 | 1.9603x | 1.8818x | 1.7299x |
| 64x32/paddedfalse/pattern0 | 436.3500 | 72.6200 | 6.0087x | 4.7932x | 17.6191x |
| 64x32/paddedfalse/pattern1 | 441.3500 | 71.7100 | 6.1547x | 5.8272x | 17.9333x |
| 64x32/paddedfalse/pattern2 | 444.0000 | 71.7550 | 6.1877x | 5.7865x | 17.9082x |
| 64x32/paddedtrue/pattern0 | 444.0000 | 73.4950 | 6.0412x | 5.6360x | 17.3753x |
| 64x32/paddedtrue/pattern1 | 442.9500 | 72.8300 | 6.0820x | 5.6233x | 17.5614x |
| 64x32/paddedtrue/pattern2 | 444.8500 | 73.0350 | 6.0909x | 5.6807x | 17.3889x |
| 64x64/paddedfalse/pattern0 | 876.0500 | 134.3000 | 6.5231x | 6.1663x | 18.9241x |
| 64x64/paddedfalse/pattern1 | 875.6000 | 132.1500 | 6.6258x | 6.3551x | 19.3795x |
| 64x64/paddedfalse/pattern2 | 873.3500 | 132.7000 | 6.5814x | 6.3648x | 19.2200x |
| 64x64/paddedtrue/pattern0 | 874.4000 | 139.7500 | 6.2569x | 6.0200x | 18.0537x |
| 64x64/paddedtrue/pattern1 | 878.6000 | 140.1500 | 6.2690x | 5.9566x | 18.1484x |
| 64x64/paddedtrue/pattern2 | 879.2000 | 140.3500 | 6.2643x | 6.0261x | 18.4681x |
| 8x8/paddedfalse/pattern0 | 32.9400 | 11.5950 | 2.8409x | 2.7246x | 4.4377x |
| 8x8/paddedfalse/pattern1 | 32.6450 | 11.5800 | 2.8191x | 2.6914x | 4.4296x |
| 8x8/paddedfalse/pattern2 | 32.6450 | 11.6450 | 2.8033x | 2.6128x | 4.4380x |
| 8x8/paddedtrue/pattern0 | 32.6500 | 11.6650 | 2.7990x | 2.6173x | 4.5761x |
| 8x8/paddedtrue/pattern1 | 32.3050 | 11.6800 | 2.7658x | 2.6703x | 4.5107x |
| 8x8/paddedtrue/pattern2 | 32.9150 | 11.6650 | 2.8217x | 2.6450x | 4.5538x |

The eight all-six-slower controls are unchanged subtraction-independent kernels; preserve every observation:

- `vp9/dsp/BenchmarkVP9SubpelPointerNative/16x32/phase0-5/paddedtrue`: +1.797484% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/16x32/phase3-0/paddedfalse`: +1.966944% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/16x32/phase3-0/paddedtrue`: +3.091361% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/16x8/phase3-0/paddedfalse`: +5.084164% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/16x8/phase4-4/paddedfalse`: +6.033810% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/16x8/phase4-4/paddedtrue`: +5.616510% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/32x16/phase3-0/paddedfalse`: +3.337010% parent.
- `vp9/dsp/BenchmarkVP9SubpelPointerNative/32x16/phase3-0/paddedtrue`: +4.059502% parent.

M4 Max evidence is separate, measured against prototype 6a whose tree equals submitted main 8ee. All 48 components gain 1.7127–3.0183xparent, all 288 pairs>=1.6500x. Larger contiguous32/64-column kernels beat active NEON; padded-output public gains are against scalar fallback. Whole VP9 improves 1.3368%/2.9800% versus parent (all six), but remains 8.8426%/35.2237% slower than old build. VP8 control+.4937% remains. All 22 semantic fields and allocation samples match. No cross-host/additive attribution.

Optimized C is a separate EPYC 7763 job with pinned libvpx v1.16.0/1024874c5919305883187e2953de8fcb4c3d7fa6, ASM/SSE2/AVX2/runtime dispatch enabled. Fixed-speed VP8 Go takes **2.362811x libvpx time** (+136.2811%), matching measured PSNR/SSIM/payload/frame counts without proving bitstream identity; this is a VP8 control for this change. Adaptive+77.6720%,VP9+190.5841/+200.7965% compare unequal quality/output work. Every quality-pass output statistic matches timed outputs; all reference timers are `vpxenc-stats`. Never combine native9V45 and C7763 absolute timings or different source checkpoints.

Merged after root review and fresh exact-head/all-three checks as `c2d415fffead020022bd7bb6f40a7bd665795f44`, parent 8ee. Tested and merged main tree both `ddef7aa5fcbfd98955f2c0bbcf1e94ef43972c54`. Remote main verified; no dependent PR existed, then only exact remote codex/residual-simd was deleted. Local checkout/index/branches untouched. Full proof in audit; root-proof.json preserves the independent 160-raw/14-Git-SHA review byte for byte.

All 86 native plus 74 C raw members retain every byte/whitespace; member and outer SHA256 manifests verify. Audit byte-identical to compact JSON. No local Go builds/tests/profiles, source or README changes during audit/remote merge.
