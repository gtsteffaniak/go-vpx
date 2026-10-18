# ARM DC pair staged evidence

Scope: Apple M4 Max, Go 1.27.2 SIMD, PGO off. Independent component windows; do not multiply cross-window ratios. Runtime samples are 216 pointer + 576 signed + 1,152 in-place + 648 rejected-guard = 2,592. Every retained sample passed its public oracle with 0 B/op and 0 allocs/op. Full private originals remain unchanged. No new whole-codec timing has run for this staged candidate.

| Stage | Same-window median speedups | Decision |
|---|---|---|
| Exact pointer windows | 1.180–1.298× vs exact PR48 | Retain intermediate evidence |
| Signed DUP/ZIP + SQXTUN | 1.230–1.384× vs PR48; 1.033–1.073× vs pointer | Retain in chosen candidate |
| Six-argument in-place owner158 | New nine actual in-place cases 1.4751299183–1.4881319851× vs PR48, all72 pairs faster; 1.184–1.200× vs signed | Keep158; robust smaller engineering choice, no1.5× claim |
| Guard156 | New targets1.00265–1.02399× vs158,46/54faster8slower; generic/disjoint controls mixed/regressing | Reject; keep158 |

The new in-place fixture is shared identically by all references; reference-only adapters forward to the old generic API. Their inlining is proven. The historical original old fixture used two indirect BLR calls; the new fixture uses direct BL calls to the same resolved original API entry. This is a stronger old denominator, not adapter overhead. In the1,152-sample window the selected owner remained6.91–8.47% slower than that old direct-call implementation across target medians. The goal of beating old assembly remains open.

The fresh guard window directly measured PR48→158 target medians1.46329–1.48336×, PR48→1561.48445–1.50721×; it does not establish every target≥1.5× and the rejected guard paired minimum was1.46112×. All mixed signs and slower outliers are preserved, including the severe generic-control ratio.70261 and earlier.69272. No sample was dropped.

Owner158 keeps one exact in-place window guard, signed deltas, ordered two-row load/store packets, original invalid-layout block fallback, disabled scalar semantics, and unchanged generic public API. Runtime validation includes55,296 new in-place full-backing/panic cases and156,000 general full-state cases, exhaustive bytes/deltas, race/checkptr2, purego/nosimd, artificial ARM disabled selection, AMD selected/actualCPU-off correctness, and normalized other-architecture callers. Only two reconstructed callers survive the ordinary CLI link; all five source callers were separately retained and proven in decoder test binaries. Generic owner225 remains instruction-identical between158 and rejected156; actual CLI callers are also normalized identical.

The rejected guard changed only ARM owner eligibility to8≤stride≤2^30+7. Independent integer boundary proof covers186 cases, maximum endpoint3,221,225,501. Larger complete windows take unchanged fallback. All validations passed, but the small/mixed performance change did not justify adoption.

Diagnostic failures are preserved: dead-stripped generic CLI symbol scope assertion, old indirect-call matching assertion, and signed SQXTUN counting scope correction. These were proof/parser errors, not codec output failures. Source checks/build logs/commands/oracles/LLVM and all raw timings are included. Executables alone are omitted, with private paths, sizes and SHA256 receipts. No entire source snapshots are copied. No cross-CPU comparisons, bitstream-identity claims, or whole-codec gains are inferred.


Publication: `raw-evidence.tar.gz` preserves all 7 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
