# Rejected ARM batched FDCT fixed-window probe

The safe fixed-array probe is rejected. On Apple M4 Max, Go 1.27.2, SIMD enabled and PGO disabled, its public operation improved by only 1.023–1.123× over the pinned parent. Every median remains below the user’s 1.5× component acceptance threshold, and every median remains slower than the pinned old NEON implementation. No whole-codec performance claim is made.

Parent: `968b518bc20d599606d5dccf03532bebcc4f3b4f`. Old assembly: `691cd0512c4829d9ee40d483205dd223f9eb028d`. Candidate: that exact parent plus `candidate.patch`, uncommitted and isolated outside the repository. The root source, index, branches and README files were not modified.

## Change and contract

The candidate keeps the existing full-count endpoint checks and arithmetic. Within each two-block iteration, checked `[32]int16` views replace repeated dynamic row slices; exact `[8]int16` array loads/stores replace slice APIs. It introduces no unsafe operations. It retains the per-pair residual magnitude guard, full-signed-range scalar fallback, sequential two-block processing, exact in-place behavior and odd-block caller handling. No coefficient, shift, transpose or range threshold changed.

The ordinary optimized batch body changes from 372 to 288 instructions and from 18 to 6 cold `panicBounds` call sites. These are diagnostics, not performance acceptance evidence. Removing that bookkeeping did not remove the dominant batch throughput limit. The remaining old-NEON gap is measured here; the earlier profile’s 0.15 s Go cumulative versus 0.13 s NEON sample totals are only liveness evidence and are not used as timing ratios.

## Native public measurements

Twelve identical `ForwardDCT4x4Batch` fixture cases were copied into parent, candidate and old-assembly snapshots. Each case computes an independent scalar expected result before timing, verifies complete output, input immutability and full output canaries, uses `testing.B.Loop`, and consumes output after timing. Patterns stay inside the common ±255 residual contract. All six permutations of parent/candidate/assembly order were used, one 100 ms benchmark sample per case per round, one thread, no competing agent builds or timing. All 216 samples report 0 B/op and 0 allocs/op.

| Public case | Parent ns | Candidate ns | Old NEON ns | Parent / candidate | Faster parent pairs |
|---|---:|---:|---:|---:|---:|
| blocks2/mixed | 8.8235 | 7.8585 | 7.1340 | 1.1228× | 6/6 |
| blocks2/flat | 8.8185 | 7.8680 | 7.0575 | 1.1208× | 6/6 |
| blocks2/extreme | 8.8370 | 7.8765 | 7.1510 | 1.1219× | 6/6 |
| blocks8/mixed | 32.0750 | 31.2200 | 28.1000 | 1.0274× | 6/6 |
| blocks8/flat | 31.9900 | 31.0350 | 27.8500 | 1.0308× | 6/6 |
| blocks8/extreme | 32.2250 | 31.3850 | 27.9300 | 1.0268× | 6/6 |
| blocks16/mixed | 63.9900 | 62.4000 | 55.5600 | 1.0255× | 6/6 |
| blocks16/flat | 63.6650 | 61.8250 | 54.9500 | 1.0298× | 6/6 |
| blocks16/extreme | 63.8250 | 61.9250 | 54.8000 | 1.0307× | 6/6 |
| blocks25/mixed | 102.8000 | 100.4000 | 86.3900 | 1.0239× | 6/6 |
| blocks25/flat | 101.7000 | 99.3650 | 85.6100 | 1.0235× | 5/6 |
| blocks25/extreme | 102.9000 | 99.5700 | 87.7000 | 1.0334× | 5/6 |

All twelve candidate medians improve on the parent. Ten cases improve in all six paired rounds; blocks25/flat and blocks25/extreme improve in five of six. Those slower samples remain in the raw evidence. No speedups were summed with earlier optimizations, and no cross-CPU comparison is used.

## Correctness and evidence

Parent and candidate pass the existing independent scalar batch tests, including random full-int16 inputs and exact in-place operation at counts 0/1/2/3/8/16/25. Candidate also passes those tests with race instrumentation and checkptr level 2. The exact copied public fixture passes its scalar oracle at one iteration and in every measured sample on all three implementations. Full-int16 behavior of the old assembly is outside that common residual fixture and was not asserted.

`raw/` contains all 18 timing logs (12 rows each); `validation/` contains the three one-iteration fixture checks, ordinary parent/candidate correctness logs, candidate race/checkptr result and all three build commands/logs. `codegen/` contains ordinary optimized parent/candidate disassemblies. `sources/` contains only relevant exact source/test files, not entire snapshots. `reproduce.py` reconstructs pinned snapshots from Git and overlays these files before building/running the identical protocol through the root project wrapper. It does not alter the checkout.

`environment.json` records provenance and flags. `summary.json` retains every sample and all twelve median/pair results. `SHA256SUMS` covers every other package member; `verification.json` independently checks member hashes and verifies timing-row counts, zero allocation fields and PASS markers. No new benchmarks or builds were run while packaging.
