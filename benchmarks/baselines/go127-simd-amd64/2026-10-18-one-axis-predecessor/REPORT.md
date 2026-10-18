# One-axis native x64 predecessor and repair

Unmerged PR #37 head 1e8809b63498b77fb6a347b0c085ebfe3c9aab49, compared with
main 983024613abf3b427bd3c1a99fa4c3fa0ce21ea1 and old assembly 691cd051.
Native Intel Xeon Platinum 8370C, Go 1.27.2, PGO off, GOMAXPROCS=1. All three
checks passed, all nine recorded source hashes match exact Git objects,
246 public kernels have six samples per variant with zero B/op/allocations,
109 targeted executions produce 254 PASS lines, and all 22 non-timing output
fields and every public/profile allocation sample match. All raw 160 members
(86 native+74 optimized-C) and their hashes are preserved without filtering.

The changed vertical path has four real regressions, not control noise:

| Fixture | More time than main | Slower pairs |
|---|---:|---:|
| 16x16/phase0-5/paddedfalse | +20.8217% | 6/6 |
| 16x16/phase0-5/paddedtrue | +21.4604% | 6/6 |
| 16x8/phase0-5/paddedfalse | +45.8849% | 6/6 |
| 16x8/phase0-5/paddedtrue | +46.5283% | 6/6 |

The eight horizontal fixtures are 1.197–1.224x main, all six pairs faster each.
All 16 one-axis candidates, 24 identity/two-axis controls, 116 median regressions
and 40 all-six-slower cases remain in the audit. Do not discard the four changed
vertical regressions among otherwise unchanged controls.

Whole-codec results on this SAME native runner (negative means less time):

* VP9 320x180: -0.1747% main (4/6 faster), -40.2994% old assembly (6/6 faster).

* VP8 720p: -0.6287% main (5/6 faster), +9.5227% old assembly (0/6 faster).

* VP9 720p: -0.2920% main (6/6 faster), -48.9087% old assembly (6/6 faster).

Separate optimized libvpx v1.16.0 job on EPYC 7763: fixed-speed VP8 Go takes
3.589323x libvpx TIME (+258.9323%), with matching measured PSNR/SSIM, payload
and frame/drop counts. Matching statistics do not assert bitstream identity.
Adaptive/VP9 rows have unequal output/quality work and are not equivalent-speed
claims. Keep C and Go-assembly baselines separate, and do not attribute differences
between separate runner jobs or CPU models to a source change.

Merge is held despite green checks. A separate saved optimized Darwin/amd64
build reproduces a legacy MOVUPS xmm7→xmm6 while YMM state is live in the byte-row
carrier. This is a code-generation clue, not native causal proof. The repair
carries already-zero-extended Int16x16 rows instead: same raw bytes, same
wrapping arithmetic for every int16 coefficient, same bounds and read-only
aliases, same deferred upper-state cleanup. It removes a repeated widening and
the saved corrected hot body uses VMOVDQU ymm7→ymm6 with no legacy XMM operation
in that region. A legacy prologue stack initialization remains before AVX use;
this is not a whole-function all-VEX claim.

ARM qualifying code, shared dispatch and AMD horizontal code remain unchanged.
AMD scalar/extreme/checkptr correctness and actual AVX-off fallback pass;
independent arithmetic/bounds/codegen review passes. Corrected native Linux
performance must be remeasured at its own exact head before merging. The
validated two-stage wrapper remains a fallback research source if this repair
fails native timing. Its unmeasured performance is not claimed. All predecessor
measurements remain scoped to 1e8809b6, not the corrected head. README files are
untouched; the >=1.5x ARM kernel win is retained independently of overall time.
