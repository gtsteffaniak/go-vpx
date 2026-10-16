# Rejected ARM eight-byte broadcast loads

Parent: ef35b4ba, Go1.27.2, GOEXPERIMENT=simd, PGO off, Apple M4 Max,
GOMAXPROCS=1. Six alternating native A/B pairs, 150ms per microbenchmark,
no concurrent Go validation or other benchmark jobs.

Candidate replaced the ARM `loadNPtr(p,8)` zero-vector insertion with
`BroadcastUint64x2(*(*uint64)(p))`, leaving amd64 and four-byte loads unchanged.
Every consumer observes only the low eight bytes. Five codec/kernel packages
passed, including independent scalar parity and minimal/misaligned buffers.

Rejected: VP8/VP9 subpixel variance16x16 regressed from about49.1 to54.9ns,
and VP9 sizes8/32/64 regressed9.5–12.4%. Allocations remained zero.
No public encode or x64 performance claim; the candidate was not retained.

The proposed elimination of vector-zero copies did not occur in actual Go
code generation. The attached ARM disassembly retains the zero-vector copy
and inserts an additional VDUP per load. A source-level broadcast therefore
adds instructions in this pinned compiler. Raw benchmark and parity outputs
are preserved beside this record.
