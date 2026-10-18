# Consume denoised luma directly from the accepted average

The temporal denoiser previously copied raw Y into an overlay, filtered into the running-average plane, and copied the final average back into that overlay before coefficients read it. The new average-only entrypoint reads the raw signal without changing it and lets accepted coefficients consume the running average directly. This removes both signal copies. Spatial-filter output uses the same average-backed route. The original aliased partial-edge path and unsupported-layout fallback remain available.

The existing SIMD arithmetic is unchanged. The shared complete filter only gains a final signal-copy flag; its regular API retains the previous behavior. Encoder glue records the selected Y backing plane and supplies its stride to accepted coefficient generation in both serial and row-threaded builders. The no-filter spatial path likewise uses its existing average rather than copying it back to overlay.

Preliminary native Apple M4 Max / Darwin ARM64, Go 1.27.2, SIMD enabled, PGO off, CGO enabled, GOMAXPROCS=1. Baseline is exact PR110 runtime `2bbc6431`, equivalent to main `88fe29f77914a38e96a23adb28bd0275b827ce87`.

| Whole codec workload | Paired speedups | Median |
|---|---|---|
| VP8 1280x720, 60 frames, fixed cpu -8 | 1.04039x, 1.04780x | 1.04409x |
| VP8 1280x720, 60 frames, adaptive cpu +8 | 1.09545x, 1.01561x | 1.05553x |

Both codec workloads have identical semantic fields and zero allocations/frame. The adaptive spread is preserved; the short screens do not establish a precise long-run percentage. Existing denoiser component benchmarks improve 1.017–1.220x, with the Y-only case at 1.220x. Raw AB/BA timing rows, all controls, commands and binary/source hashes are under `screen/` and `source-binding.json`. No x64 performance claim is made.

Focused correctness passed:

- 560 average-only complete-filter comparisons against the scalar implementation, including rejected/corrected cases, distinct strides, motion and increase settings, and unchanged raw signal bytes.
- Existing frozen denoiser consumer/state checks for all eight motion/fallback patterns, complete/partial edges, aliased/separate signal storage, spatial filtering, repeated calls and following UV consumers. Comparisons still check complete source, average, state and scratch bytes; the selected Y signal is checked at the consumer boundary because unused overlay Y bytes are deliberately no longer written.
- 32 noisy-stream packet hashes captured from the immutable baseline before candidate compilation: complete and partial dimensions, Threads=1 and Threads=2, eight frames each.
- Existing 96-packet VP8 and 72-packet ML stream hashes, zero-allocation denoiser checks, and focused fixed-N threading determinism.

The average is separately owned from the caller's raw source and the overlay. The complete-kernel alias guards remain intact. A coefficient builder consumes its current macroblock before later macroblocks can change its average through spatial filtering; the existing row dependency synchronization enforces this order for threaded builders. Partial-edge staging retains the original aliased signal path. No README, new dependencies, handwritten assembly, full CI, full suite, profile or new C comparison was added.
