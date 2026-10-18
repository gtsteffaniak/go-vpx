# ARM zero-packet tuning screen

Apple M4 Max, Darwin ARM64, Go 1.27.2, SIMD enabled, CGO=1, PGO off, GOMAXPROCS=1. Two balanced 40ms rounds across the same 48 public optional-output fixtures, zero allocations. This is an iteration screen, not a final end-to-end acceptance measurement.

The original packet skip passed numeric/alias correctness but regressed every dense median. Removing repeated slice checks inside the ARM helper (the public entry validates every window first) improved all 48 medians versus the PR52 parent: zero 1.371–2.464x, sparse 1.086–1.549x, dense 1.043–1.410x. All raw screen rows and binary hashes are included. Two rounds are insufficient for strong stability claims. Store order and exact packet bounds were independently reviewed; targeted pixel/VP8/VP9 tests passed.

Original six-pair packet results remain in root-arm-summary.json and independent-arm-audit.json; those describe the earlier slice-based implementation, not this pointer refinement. Original packet full raw evidence is preserved in the [initial slice-gate report](../2026-10-18-zero-packet-initial/REPORT.md). The initial broad combined test run had one source-citation failure, fixed by adding the libvpx baseline comment; only the failed citation and targeted affected tests were rerun. No codec gap closure is claimed here.
