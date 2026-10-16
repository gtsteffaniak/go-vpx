# Native amd64 SIMD checkpoint, 2026-10-16

Assembly baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
SIMD checkpoint: `6a7b0345d7e47644cb130d7e9220a70d25409d24`.
[Successful native workflow](https://github.com/thesyncim/govpx/actions/runs/38953587137).

The hosted VM exposes an AMD EPYC 9V45 with AVX2 and AVX512. Both variants
use official Go 1.27.2, `GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`,
and PGO disabled. Six serial pairs alternate assembly/SIMD and SIMD/assembly
order. Kernel benchmarks use 150 ms per sample and report zero allocations.
Native SIMD selection is asserted in `dispatch.txt`.

This CPU differs from the EPYC 7763 used for the preceding `b9badbe3` run.
These results cannot isolate the AVX/SSE transition repair on that older
CPU. Compare each variant against its baseline on the same runner.

| Public encoder workload | Assembly median | SIMD median | Time change | Paired time changes |
| --- | ---: | ---: | ---: | ---: |
| VP9 320x180 | 4.587 ms/frame | 2.656 ms/frame | -42.1% | -43.3% to -39.7% |
| VP8 720p | 12.271 ms/frame | 20.485 ms/frame | +66.9% | +64.0% to +73.3% |
| VP9 720p | 39.922 ms/frame | 21.715 ms/frame | -45.6% | -47.2% to -44.9% |

The 320x180 workload uses 120 input frames and 500 kbps; both 720p workloads
use 60 input frames and 2,500 kbps. All use 30 fps, realtime CPU-used 8,
one thread, encode-only mode, and no external encoder comparison. The
remaining VP8 regression is material; this is not overall performance parity
with the removed assembly.

All twelve samples per workload agree on output byte counts, encoded/dropped
frame counts, and quantizer statistics/histograms. These reports contain
byte counts, not bitstream hashes; exact bytes are tested separately by the
codec oracle gates.

| Workload | Output bytes | Encoded / dropped | Allocations per input frame |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667, all samples |
| VP8 720p | 685,957 | 58 / 2 | 0, all samples |
| VP9 720p | 642,654 | 48 / 12 | 1.25 median, both variants |

One VP9 720p assembly sample records 76 allocations over 60 inputs rather
than the 75 in the other five assembly and all six SIMD samples. Raw values
are retained in the individual JSON reports and `summary.json`.

The raw CPU profiles and flat/cumulative top tables were collected after
all timing pairs, on separate 120-frame 720p runs. In the SIMD VP8 profile,
loop filtering accounts for 33.1% cumulative time, including 26.7% in the
sixteen-position narrow filter; statistics account for 12.7%. Byte packing
is 2.1% flat. Profile timings are excluded from the benchmark medians.

Run `sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR`
on native Linux/amd64 to reproduce the harness. `environment.txt` retains
the compiler and CPU metadata; `*-kernels.txt` retains raw repeated kernel
results; `*-encoder-*.json` retains public timing samples; `*-profile*`
retains the independent profiles and their reports. `summary.json` contains
medians, all paired deltas, allocation samples, and parity checks.
