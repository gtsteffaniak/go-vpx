# Failed PR49 native job: completed measurements, parser failure

Exact head `41596b4e0ad1087e4afe4020825d03cb38e2c4ce`, base `c8bd67df9a00c102163354158c1ee46b0f1527a7`, old Go/assembly reference `691cd0512c4829d9ee40d483205dd223f9eb028d`.
Workflow run39115190163/job117639458187 FAILED at10:06:42UTC after42m26. The final optional Python parser rejected a verbose name-only preflight header. This package does not label the GitHub job green. All raw files remain unchanged.

Hardware: **AMD EPYC 7763 64-Core Processor**. Real SDK AVX2=true, AVX512=false. Selected512 oracle skipped6 times. The optional log explicitly declares512 unavailable and omits all FP32 timings. No selected512 correctness/performance acceptance is possible from this runner.

The independent data audit verifies93raw hashes,116exact failed-head source hashes,33common fixtures and one reference-only adapter,534cases×18=9612main samples,446excluded ordinary one-shots,396guard samples/66excluded, and216Stats samples/132optional excluded preflights. Every timed micro sample is0B/0alloc. All22codec semantic fields and all allocation samples agree, including profile rows. All303median regressions and112all-six-slower cases are preserved in full-proof.json.

## Actual public Stats18, paired parent/current

One deterministic pattern per width/height, not multiple patterns: bytes(i*37+11) versus bytes(i*19+53), strides width+13/width+17 and starts3/5. Thus w32h1 is one case with source/reference strides45/49.

| Case | Parent ns | Current ns | Parent/current | Paired min–max | All6 slower |
|---|---:|---:|---:|---:|---:|
| w32/h1 | 15.925 | 19.12 | 0.832897 | 0.824780–0.844080 | True |
| w32/h4096 | 24414 | 23298 | 1.047945 | 1.043800–1.053208 | False |
| w32/h4097 | 24458 | 23306 | 1.049430 | 1.043829–1.056060 | False |
| w32/h64 | 393.6 | 374.9 | 1.049880 | 1.048164–1.051042 | False |
| w32/h7 | 51.615 | 49.89 | 1.034576 | 1.016993–1.039694 | False |
| w32/h8 | 57.435 | 55.42 | 1.036359 | 1.023890–1.043668 | False |
| w64/h1 | 19.685 | 19.485 | 1.010264 | 0.985937–1.034950 | False |
| w64/h4096 | 40128 | 38962 | 1.029927 | 1.024469–1.040902 | False |
| w64/h4097 | 40206 | 38916 | 1.033135 | 1.025826–1.037658 | False |
| w64/h64 | 633.65 | 613.2 | 1.033350 | 1.031958–1.037170 | False |
| w64/h7 | 77.965 | 76.16 | 1.023700 | 1.022053–1.026561 | False |
| w64/h8 | 87.435 | 85.22 | 1.025992 | 1.017105–1.028098 | False |
| w8/h1 | 17.05 | 15.985 | 1.066625 | 1.059813–1.072146 | False |
| w8/h4096 | 8677.5 | 8669.5 | 1.000923 | 0.998846–1.021438 | False |
| w8/h4097 | 8677.5 | 8672 | 1.000634 | 0.995188–1.002652 | False |
| w8/h64 | 147.6 | 143.9 | 1.025712 | 0.966510–1.050907 | False |
| w8/h7 | 29.905 | 28.405 | 1.052808 | 1.034107–1.057082 | False |
| w8/h8 | 26.56 | 25.6 | 1.037500 | 1.035895–1.042578 | False |

w32h1 is a reproducible20.063% time regression (15.925→19.12ns, all6pairs slower). Long width32/64 cases gain roughly3–5%; width8short cases gain up to6.66%. No1.5× result is claimed for this Stats trial.

## Whole codecs on this same native runner

| Workload | Old Go+ASM ms | Parent ms | Current ms | Current time vs parent | Current time vs old | All6 faster than old |
|---|---:|---:|---:|---:|---:|---:|
| VP9 320x180 | 6.997566 | 4.107698 | 4.119705 | +0.292305% | -41.126605% | True |
| VP8 720p | 23.381346 | 18.513858 | 18.652705 | +0.749965% | -20.223988% | True |
| VP9 720p | 60.519698 | 28.995812 | 28.998161 | +0.008103% | -52.084756% | True |

Old means actual pinned691cd051 Go/handwritten assembly built with the same official Go1.27.2/PGOoff/GOMAX1. It is not the external libvpx C oracle. Do not pool these x64 results with separate local ARM whole-codec windows. All per-round output, quantizer, drop/count fields and allocation samples remain in raw/full proof.

The two-line parser proposal accepts only verbose name-only preflight headers; inventory, positive timed values and allocation checks remain unchanged. Replay on the complete real immutable optional log passes132excluded+216timed without running a new benchmark. The original42minute job remains failed; a repaired exact-head CI run is still required.

No executables or full repository snapshot are included. Raw.pgo files are profiling data and remain intact. Repaired root production source is deliberately excluded: included116sources are the exact failed41596 commit.


Publication: `raw-evidence.tar.gz` preserves all 237 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
