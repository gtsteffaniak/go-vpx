# Native amd64 batched-DCT checkpoint, 2026-10-16

Assembly baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
SIMD checkpoint: `3c94d2cc89db12f67b148676332383ba90248a6b`.
[Successful native workflow](https://github.com/thesyncim/govpx/actions/runs/38956724630).

Both variants run on the same hosted EPYC 7763 VM (AVX2, no AVX512), with
official Go 1.27.2, `GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`,
and PGO disabled. Six pairs alternate assembly/SIMD and SIMD/assembly
order; kernel samples run for 150 ms. All kernels report zero allocations.
`dispatch.txt` asserts that the runtime selected native SIMD.

Only `internal/vp8/encoder/dct_batch_native_test.go` is copied onto the
assembly baseline, providing the same residual-domain benchmark fixture.
Baseline production source remains unchanged at the pinned commit. The
25-block DCT medians are 217.25 ns for assembly and 216.15 ns for SIMD
(-0.51%); paired changes range from -0.87% to +0.55%, effectively parity
within these samples. Both implementations receive identical coefficients
in [-255,255]. This benchmark is separate from the production encoder runs.

| Public encoder workload | Assembly median | SIMD median | Time change | Paired time changes |
| --- | ---: | ---: | ---: | ---: |
| VP9 320x180 | 6.966 ms/frame | 4.445 ms/frame | -36.2% | -36.3% to -36.1% |
| VP8 720p | 23.166 ms/frame | 27.470 ms/frame | +18.6% | +17.4% to +34.7% |
| VP9 720p | 59.909 ms/frame | 34.636 ms/frame | -42.2% | -42.8% to -41.1% |

The fourth VP8 SIMD sample takes 30.982 ms/frame; the other five take
27.265–28.024 ms/frame. All samples are retained. VP8 remains slower than
assembly in every pair. Earlier EPYC 7763 runs used separate hosted VMs;
those comparisons do not isolate individual changes on the same host.

All workloads use realtime CPU-used 8, one thread, 30 fps, encode-only mode,
and no external encoder comparison. VP9 320x180 uses 120 inputs at 500 kbps;
the 720p workloads use 60 inputs at 2,500 kbps.

| Workload | Output bytes | Encoded / dropped | Allocations per input frame |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Every byte count, encoded/dropped count, quantizer statistic/histogram,
keyframe/interframe size, and allocation count matches across all twelve
samples per workload. There are no allocation outliers. Exact bitstreams
are tested separately by codec oracles; these reports retain byte counts.

Independent 120-frame profiles run after all timing pairs. In the SIMD VP8
profile, the byte-wide loop filter accounts for 20.8% cumulative CPU time,
`pixelStatsRows16` for 13.1%, `Gather` for 6.1%, and the batched DCT for
2.8%. Profile timings are excluded from the benchmark medians.

Check out the recorded SIMD checkpoint and run
`sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR`
on native Linux/amd64 to reproduce the harness. Raw workflow text, JSON,
profiles, and metadata are preserved unchanged. `summary.json` contains
all medians, samples, paired changes, and semantic/allocation checks.
