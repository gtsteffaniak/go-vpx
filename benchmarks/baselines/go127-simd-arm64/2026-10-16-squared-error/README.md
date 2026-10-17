# Squared-error-only scoring and full-range Walsh

Apple M4 Max, darwin/arm64, official Go 1.27.2, GOEXPERIMENT=simd,
GOMAXPROCS=1 / test.cpu=1, PGO=off. Six alternating kernel comparisons
measure VP8 squared-error calls before/after removing the unused signed
sum: 16x720 approximately 1,160→750 ns, 16x16 approximately 30→19 ns, and
8x8 approximately 11.9→8.2 ns, all zero allocations.

The separate six-triple encoder run compares perf18 (3c94d2cc production
source), perf19 (dedicated 8/16-wide SSE), and perf20 (SSE plus full-signed
native Walsh). Medians are 10.783,10.697, and 10.640 ms/frame. These small and
noisy differences do not establish a public encoder speedup. Output bytes,
encoded/dropped frames, quantizers and zero allocations match in every run.
No competing local builds/tests ran during either timing batch.

SSE tests cover exact final-row buffers, misalignment, odd/zero heights,
read-only aliasing, independent random data, and modulo 2^32 overflow.
Full-signed Walsh tests force live native dispatch across 6000 cases on
ARM and Rosetta amd64, including int16 extremes and multiple strides.
All five pixel/DSP/encoder suites pass on both execution targets. Native
x64 performance remains to be established by the hosted benchmark.

The separate filter-kind specialization was rejected and reverted: its
vertical-edge change was small/noisy and horizontal edges regressed.
