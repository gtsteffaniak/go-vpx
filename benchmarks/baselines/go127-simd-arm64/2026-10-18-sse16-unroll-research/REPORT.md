# Rejected SSE16 byte-square/two-row interaction probe

Parent: `34c4e4a17b231a24a3d3774c9621b7e5b80bac03`. Apple M4 Max, Go1.27.2 darwin/arm64, SIMD enabled, PGO off, GOMAXPROCS=1. All variants rebuild the same parent snapshot, changing only the SSE16 ARM file and applying identical private test fixtures.

Six 80ms samples per public fixture/variant cover all six variant permutations. Every returned squared error is checked against an independent scalar fixture and escapes to an observable package sink. All byte/allocation samples are zero.

| Public fixture | Parent ns | Byte square ns | Two-row ns | Two-row / parent speedup | Two-row / byte-square speedup |
| --- | ---: | ---: | ---: | ---: | ---: |
| h1/pitch32/mode0 | 3.219 | 2.978 | 3.462 | 0.930x | 0.860x |
| h1/pitch32/mode1 | 3.221 | 2.974 | 3.465 | 0.930x | 0.859x |
| h1/pitch1312/mode0 | 3.223 | 2.983 | 3.466 | 0.930x | 0.861x |
| h1/pitch1312/mode1 | 3.260 | 3.006 | 3.478 | 0.937x | 0.864x |
| h8/pitch32/mode0 | 9.677 | 6.774 | 6.982 | 1.386x | 0.970x |
| h8/pitch32/mode1 | 9.666 | 6.697 | 6.998 | 1.381x | 0.957x |
| h8/pitch1312/mode0 | 9.702 | 6.758 | 6.998 | 1.386x | 0.966x |
| h8/pitch1312/mode1 | 9.702 | 6.710 | 6.985 | 1.389x | 0.961x |
| h16/pitch32/mode0 | 17.120 | 13.210 | 13.280 | 1.289x | 0.995x |
| h16/pitch32/mode1 | 17.120 | 13.155 | 13.235 | 1.294x | 0.994x |
| h16/pitch1312/mode0 | 17.100 | 13.080 | 13.250 | 1.291x | 0.987x |
| h16/pitch1312/mode1 | 17.060 | 13.100 | 13.355 | 1.277x | 0.981x |
| h720/pitch32/mode0 | 667.150 | 575.400 | 565.850 | 1.179x | 1.017x |
| h720/pitch32/mode1 | 668.850 | 576.900 | 571.400 | 1.171x | 1.010x |
| h720/pitch1312/mode0 | 670.900 | 580.900 | 579.500 | 1.158x | 1.002x |
| h720/pitch1312/mode1 | 669.600 | 580.250 | 573.200 | 1.168x | 1.012x |

The candidate retains the exact parent safe loop for h<8 and unsupported geometry. For 8<=h<=4096 and nonnegative strides below2^30, final-window checks bound exact16-byte pointer loads before two-row traversal. Four independent int32 vector accumulators preserve modulo2^32 arithmetic, including odd tail rows.

All variants pass `TestSquaredError|TestSSE` in pixel and the DSP package (DSP has no matching named tests; benchmark1x verifies its public fixtures). Private pixel boundary cases add heights7/8/9 and4096/4097, disjoint/shifted alias spans, zero/overlapping/asymmetric strides, exact minimal buffers and one-byte-short rejection, and full255/0 differences.

Existing byte-square codegen has inline UMULL/UMULL2 and no live-loop calls/spills. Two-row codegen retains four independent accumulators and no live-loop calls/spills; the two-row branch covers38instructions. Instruction count alone is not acceptance evidence.

Rejected: no representative case reaches1.5x production speedup. Two-row is slower than byte-square alone at h8/16, only1.002–1.017x it at h720, and h1 loses about7% to guard overhead despite retaining parent arithmetic. No end-to-end codec speedup was measured or claimed. All raw logs and candidate sources remain private; no root checkout, README or Git changes were made.

## Reproduction and archive layout

This package intentionally contains no binaries or complete source snapshots.
Reconstruct three private archives from the recorded parent commit. Copy the
matching `sources/<variant>/sse16_simd_arm64.go` into each archive's
`internal/pixel/` directory and the common benchmark/test fixtures to their
original paths. The exact original build and measurement drivers are retained;
they refer to `/private/tmp/govpx-sse16-u8-unroll2-research/{parent,u8square,unroll2}`.
Run builds from the root checkout via its `scripts/go.sh -C` wrapper, as in the
driver, to retain the shared project cache scope. Common production dependencies
are included for inspection; the parent commit remains authoritative.

The three aggregate benchmark logs are unchanged. Eighteen separate round logs
are exact text slices of those aggregates, covering6rounds×3variants. One-shot
fixture oracle output, all three ordinary correctness/build logs, original
byte-square and fresh unroll disassembly, environment, sources and summaries
are retained. Empty build logs indicate successful quiet compilation.
`SHA256SUMS` covers every packaged file except itself.
