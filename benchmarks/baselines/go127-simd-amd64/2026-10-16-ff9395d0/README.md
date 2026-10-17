# Native amd64 byte-filter checkpoint, 2026-10-16

Assembly baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
SIMD checkpoint: `ff9395d0fc93b9d08d26d533e86c6732d2449577`.
[Successful native workflow](https://github.com/thesyncim/govpx/actions/runs/38955379654).

The hosted VM exposes an AMD EPYC 7763 with AVX2 and no AVX512. Both
variants use official Go 1.27.2, `GOEXPERIMENT=simd`, `GOAMD64=v1`,
`GOMAXPROCS=1`, and PGO disabled. Six serial pairs alternate
assembly/SIMD and SIMD/assembly order. Kernel samples run for 150 ms and
report zero allocations. Native SIMD selection is asserted in `dispatch.txt`.

The earlier `b9badbe3` run used the same CPU model, but a different hosted
VM. These are separate within-run baseline comparisons, not a controlled
same-host comparison of individual optimizations. The intervening `6a7b0345`
run used a different CPU model, EPYC 9V45.

| Public encoder workload | Assembly median | SIMD median | Time change | Paired time changes |
| --- | ---: | ---: | ---: | ---: |
| VP9 320x180 | 6.996 ms/frame | 4.489 ms/frame | -35.8% | -36.2% to -35.1% |
| VP8 720p | 23.187 ms/frame | 28.039 ms/frame | +20.9% | +20.7% to +21.7% |
| VP9 720p | 60.011 ms/frame | 34.547 ms/frame | -42.4% | -43.2% to -41.8% |

The 320x180 workload uses 120 inputs and 500 kbps; both 720p workloads
use 60 inputs and 2,500 kbps. All use 30 fps, realtime CPU-used 8,
one thread, encode-only mode, and no external encoder comparison. VP8
remains slower than assembly in every pair. At `b9badbe3`, its within-run
gap was +411.7%; the changed host and multiple code changes prevent
isolating how much the transition repair or byte-domain filter contributes.

All twelve samples per workload agree on output byte counts, encoded/dropped
frame counts, quantizer statistics/histograms, keyframe/interframe sizes,
and allocation counts. No allocation outliers occur in this run. These
reports retain byte counts, not bitstream hashes; exact bytes are checked
separately by the codec oracle gates.

| Workload | Output bytes | Encoded / dropped | Allocations per input frame |
| --- | ---: | ---: | ---: |
| VP9 320x180 | 246,893 | 117 / 3 | 0.616667 |
| VP8 720p | 685,957 | 58 / 2 | 0 |
| VP9 720p | 642,654 | 48 / 12 | 1.25 |

Separate 120-frame profiles run after all timing pairs. In the SIMD VP8
profile, loop filtering accounts for 23.7% cumulative time, including
22.2% in the byte-wide path; `pixelStatsRows16` accounts for 7.6%,
`forward4Go` for 5.3%, and `Gather` for 4.7%. Profile timings are excluded
from the benchmark medians.

Run `sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT_DIR`
on native Linux/amd64 to reproduce the harness. The raw text/JSON/profile
files are preserved unchanged from the workflow artifact. `summary.json`
contains medians, paired deltas, and every semantic/allocation comparison.
