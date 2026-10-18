VP8 temporal denoising now filters directly from the raw luma source into its running average and lets accepted coefficient generation consume that average. This removes the raw-to-overlay staging copy and the accepted-average-to-overlay copy. Existing SIMD arithmetic, aliased partial-edge fallback, spatial-filter ordering and serial/threaded dependencies remain intact; caller-owned source stays read-only.

Preliminary Apple M4 Max / Go 1.27.2 SIMD single-thread 720p60 realtime codec throughput: fixed cpu=-8 **1.0441×** (1.0404×/1.0478× pairs); adaptive cpu=+8 **1.0555×** (1.0955×/1.0156×, variation retained). Zero codec allocations and identical output fields. Luma-only component improves 1.22×; other measured controls and raw source/binary bindings are preserved in the experiment report.

Validation: average-only versus scalar decisions/full-average/source immutability; frozen denoiser signal/state/spatial comparisons; 32 pinned noise-enabled packets covering complete/partial frames and Threads=1/2; 96 normal VP8 packet/flag hashes;72 VP9 ML hashes and allocation checks. Native x64 feedback and the fixture correction are recorded below. No full CI, C comparison, profile audit or README changes.

### Fixed-speed stream oracle correction

The original tiny-frame fixture used adaptive `cpu-used=+8`, whose real wall-clock feedback invalidated fixed packet hashes. The oracle-only [diagnostic 39162342059](https://github.com/thesyncim/govpx/actions/runs/39162342059) produced the same isolated native stream on baseline and candidate, while both failed the old ARM pin. The corrected fixture sets `-8` before encoder construction, asserts effective speed 8 after each frame, and uses pins captured first from the immutable baseline. All four fixed-speed stream hashes now match on both native references. The original adaptive failure evidence remains preserved; denoiser production code was unchanged by this fixture repair.

### Native x64 results and accepted tradeoff

[Quick codec run 39162855642](https://github.com/thesyncim/govpx/actions/runs/39162855642) used AMD EPYC 7763 (actual AVX2 enabled, AVX512 unavailable), VP8 720p/60 inputs, 2500 kbps, one thread, PGO off and two balanced AB/BA rounds per CPU setting. **All four codec samples were slower:**

| CPU-used | Round | Encode-time change | Throughput change |
|---|---:|---:|---:|
| −8 | 1 | +0.683% | −0.679% |
| −8 | 2 | +0.715% | −0.710% |
| +8 | 1 | +0.144% | −0.144% |
| +8 | 2 | +1.911% | −1.875% |

All 22 output/semantic fields matched and all codec allocation samples were zero. Denoiser full-pel, odd-integer, state-fallback and Y-only components improved **1.212–1.222×**; zero-YUV improved 1.122×. The rejected-denoiser component regressed to **0.9564×** (paired 0.9581×/0.9546×), with other controls preserved. This is an explicitly accepted small tradeoff for the ARM codec gains and x64 component gains; it is **not an x64 codec gain**. No further validation, C, profile or full-matrix run was added.
