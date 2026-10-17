# Native amd64 SIMD checkpoint

Baseline: `691cd0512c4829d9ee40d483205dd223f9eb028d`.
SIMD: `bb1a6ef3476345a7bc250dee3abd6d53768b3913`.

[Native GitHub Actions run](https://github.com/thesyncim/govpx/actions/runs/38947952120).
Run `sh scripts/benchmark-amd64.sh 691cd0512c4829d9ee40d483205dd223f9eb028d OUTPUT`
on native Linux/amd64. Compiler, CPU, feature selection and cache environment
are recorded alongside the raw results. Six alternating A/B pairs use one
Go thread and disable PGO. Kernel benchtime is 150ms.

VP9 takes 34.4% less time at 320x180 and 40.1% less at 720p. VP8 at 720p
regresses by 221.9%; subsequent loop-filter work targets that gap. Every
workload agrees on encoded byte counts, frame counts, quantizer statistics
and allocations; exact bitstream oracle checks run separately in CI.
These are checkpoint results, not measurements of later commits.
