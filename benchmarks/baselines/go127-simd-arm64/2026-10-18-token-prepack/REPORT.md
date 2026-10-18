# ARM VP8 coefficient-token record-base CLZ kernel

Decision: retain the correctness-verified standalone ARM kernel. Leave default
`probability_tokens.go` byte-for-byte unchanged. Both attempted public integration
paths regress dense accumulation, and neither is enabled in the kernel-only overlay.
No whole-codec performance claim is made.

## Scope and provenance

Parent: `c638b7c76eac857affa0780820e86c211cdd84d3`. Native machine: Apple M4 Max,
16 cores, macOS 26.3.1(a), Darwin 25.3.0 arm64. Compiler: official Go 1.27.2,
`GOEXPERIMENT=simd`, `-pgo=off`, `GOMAXPROCS=1`. Builds use the root `scripts/go.sh`
project-env wrapper and shared project cache with agent-specific temporary paths.
No root repository files, README files, Git state, PRs, or CI heads were modified.

The helper computes all 16 raster-order static record bases:
`token | magnitude<<15 | sign<<27`. It does not count tokens, mutate contexts,
append records, or compute probability offsets/skip-EOB bits. Tokens follow the
existing VP8 thresholds. For magnitude >=5, `min(19-CLZ16(magnitude-3),10)` gives
the category; lower magnitudes retain their token ID. Two exact eight-lane loads,
unsigned absolute-magnitude rejection, and four exact uint32-vector stores are
used. Invalid coefficients anywhere, nil inputs, and overlapping input/output
return false before stores. INT16_MIN is rejected by the unsigned reduction.

The new arithmetic is derived from VP8's thresholds and the Go1.27.2 SDK APIs.
goav1's unsigned interpretation of 16-bit ABS informed the INT16_MIN check; its
clip-to127 coefficient-level behavior is deliberately not copied. No equivalent
VP8 token classifier was found in the inspected goav1/gopus implementations.

The exact measured ARM helper and globals are byte-identical in the kernel-only
overlay, verified in `evidence/kernel-source-identity.json`. Scalar references
are test-only there. Architecture selection is ARM64+simd+non-purego; all other
builds use an unchanged-output false stub. Mandatory ARM64 NEON availability
matches the existing pixel.Enabled architecture contract.

## Fresh-six standalone component proof

The acceptance denominator is **ScalarLUTStream**, the original one-pass
abs/range/LUT classification plus the same 16 record-base stores. It does not add
a whole-input validation pass. Both references and SIMD produce exactly the
same 16 descriptors for every valid timed fixture; fixtures are disjoint. The
scalar stream implementation is inlined into its timed caller, while SIMD
retains its leaf call, so SIMD call/store costs are included. The atomic-output
two-pass scalar reference is diagnostic only and never the acceptance denominator.

Fresh six rounds use 200ms per sample, balanced AB/BA method order and rotated
fixture positions. Every one of 36 samples reports0 B/op and0 allocs/op. Full
independent record-base checks run before and after each benchmark. Timed work
is live (the optimized caller contains the native helper call and scalar LUT
load/store loops); no no-op shortcut or unobserved output is used.

| Fixture | Scalar stream ns | Native CLZ ns | Speedup | All six paired ratios |
|---|---:|---:|---:|---|
| DenseMixed | 9.5275 | 3.6860 | 2.5848× | 2.600218, 2.570384, 2.596844, 2.565194, 2.562634, 2.590823 |
| DenseSmall | 9.5535 | 3.6990 | 2.5827× | 2.536605, 2.587503, 2.575520, 2.582837, 2.582908, 2.572535 |
| Gappy | 9.5840 | 3.6940 | 2.5945× | 2.543461, 2.599564, 2.573417, 2.618674, 2.603462, 2.595079 |

Raw 36 fresh logs and commands are under `raw/fresh6`; exact sample values and
all paired ratios are in `evidence/fresh6-summary.json`. The original three-round
screen and the slower atomic-reference measurements remain under `raw/screen`.
Mixed covers signs and category boundaries; Small covers -4..4; Gappy has one
nonzero at the final coefficient. This proves the standalone prepack component
on these native fixtures; it does not prove an encoder speedup or an x64 result.

## Public integration trials (all slower controls retained)

Actual `AccumulateInterMacroblockTokenCountsAndRecords` runs a 4×4 macroblock
frame including counts reset, record reset/row markers, Y2/Y1/UV traversal, and
sequential context updates. Full counts, complete records, row starts, above and
left contexts match an independent scalar oracle before and after timing.
The original sparse fixture remains a control. Dense8 includes zero uncoded tail;
Gappy16 exercises a long zero run. All67 original screen samples and all 40
flattened-screen samples report0 B/op and0 allocs/op.

Original wrapper trial: three calls per dense block (try wrapper, alias guard,
prepack). Its countBlock caller spills 7 arguments, then wrapper spills all 7 again;
the wrapper has a 112 B frame, calls the 92-instruction alias checker, zeroes 64 B
scratch, and calls the 137-instruction leaf classifier. There are no per-token
append helper calls: replay stores packed records inline and EOB append remains
inline. Scalar recurrence, scatter updates, scratch and guard costs dominate.

Flattened trial: alias checks and replay move into the existing countBlock
function. Only the prepack call remains on the dense fast path; frame size is
160 B. This removes wrapper calls but cannot amortize the full alias check/scratch
across blocks. It improves the trial relative to its own paired parent, but still
regresses dense cases. The two trials are separate timing windows; do not use
their absolute differences as another controlled comparison.

| Public 4×4 MB fixture | Original wrapper parent ns | Wrapper ns | Wrapper time change | Flatten parent ns | Flatten ns | Flatten time change |
|---|---:|---:|---:|---:|---:|---:|
| Dense16 | 11795.5 | 12944.5 | +9.74% | 11787.0 | 12093.5 | +2.60% |
| Dense8 | 6596.0 | 8431.5 | +27.83% | 6548.5 | 7777.0 | +18.76% |
| Dense16Four | 11784.5 | 12788.5 | +8.52% | 11787.5 | 11931.5 | +1.22% |
| Gappy16 | 11608.0 | 12836.5 | +10.58% | 11606.0 | 11815.0 | +1.80% |
| Sparse | 1551.5 | 1528.0 | -1.51% | 1533.5 | 1512.0 | -1.40% |

Positive time change is a slowdown. Sparse differences are small/noise and do not
justify enabling slower dense glue. Both trial sources, full optimized callers,
raw ABBA logs, sample counts, allocation fields and summaries are preserved.
No further integration variants or whole-codec timing were run.

## Correctness

The independent oracle uses literal VP8 scan/band/category tables and explicit
record packing. Tests cover all signed16-bit values, every valid magnitude/sign,
threshold extremes and INT16_MIN; invalid descriptors; EOB0/1/2/7/8/9/15/16;
skipDC0/1; all block types/contexts; nil-record count-only behavior; invalid active,
uncoded-tail and skipped-DC coefficients; insufficient record capacity with exact
partial counts/records and panic text; Y2/Y1/UV macroblock context updates; and
qcoeff aliases to counts or record backing. Independent root numeric proof checks
all 65536 signed values (4097 valid), rejects all invalid unsigned magnitudes, and
verifies token/magnitude/sign fields with probability/skip fields zero.

Trial native encoder short suite, independent ordinary tests, targeted
race+checkptr2 tests, purego tests and parent oracles passed. The all-signed SIMD
test was additionally strengthened so other lanes stay valid and cannot mask the
single potentially-invalid lane; its final race/checkptr2 run passed. Earlier test
source is retained for reproducibility. Optimized LLVM confirms two CLZ.8H,
no helper calls/spills inside the leaf, and four exact output stores. Reports/logs
and numeric proof are under `evidence`.

## Delivery

Apply only `kernel-only-overlay/` for retention; it contains no default accumulator
edit. Failed wrappers and replay sources are isolated in `trials/` and must not
be enabled by default. Original build/screen/summarization drivers are preserved
under `drivers/`; paths target the original private experiment directory.
`SHA256SUMS` hashes every delivered file except itself. No binaries or complete
source snapshots are included; original binaries remain in the private trial
folder and their hashes are recorded in the environment metadata.

Final minimal overlay verification against exact c638 passed native encoder short,
full encoder purego short, kernel race+checkptr2, and nosimd kernel tests. Retained
codegen matches the measured 137-instruction helper exactly after stripping only
addresses and global relocations (normalized SHA256
`633ce5468921c33f0dd80c7861f53e23647cebf0d719edb34b457c878579e018`).
The default accumulator source identity is independently checked against git.
`trials/` contains exact measured sources; the final stronger single-invalid-lane
oracle is also saved separately and used in the retained overlay tests. See
`retention-validation.json`, `default-source-identity.json`, and
`codegen-identity.json` under `evidence/`.
