# One-axis SIMD subpixel statistics

Native Apple M4 Max, Go1.27.2 archsimd, PGO off, GOMAXPROCS1. Parent eb820478,
old assembly691cd051. This is an immediate-parent comparison, not a comparison
with the earlier PR35 source. Skip an identity bilinear stage, retaining both
original extra-row/column endpoint checks and read-only alias semantics.
The bounded dispatcher proves w/h<=64, w%16==0 and nonnegative strides<2^30.
Every pointer offset is below2^36+64 and each exact load reads16bytes.

ARM uses wrapping int16 multiply-add, separate per-tile narrow difference sums,
independent int32 squared-error sums and paired rows with an odd-row tail.
Even arbitrary int16 coefficients produce a rounded stage in[-256,255]; the
removed identity stage maps that entire range exactly to itself. Differences
are[-511,255], so each separate lane sum over64rows lies[-32704,16320]. Widen
each half before merging. Total SSE<=1,069,551,616 and |sum|<=2,093,056 fit
signedint32. AMD uses one AVX2 packet and retains deferred upper-state cleanup.
Neither input is written or retained. CPU-off and other geometry/filters keep
their original fallbacks. Phase checks precede axis setup to protect controls.

Fresh six balanced parent/assembly/current permutations,100ms per40 public
fixtures, all zero B/op/allocations. Identical B.Loop fixtures, independent
per-pixel two-rounding scalar variance/SSE checks, unaligned3/5 origins, exact
visible spans, minimal/unequal padded strides and full backing canaries.

| One-axis fixture | Parent/current time | Assembly/current time | Minimum paired parent/current |
|---|---:|---:|---:|
| 16x8/phase3-0/paddedfalse | 1.4805x | 2.1805x | 1.4733x |
| 16x8/phase3-0/paddedtrue | 1.4778x | 2.1776x | 1.4741x |
| 16x8/phase0-5/paddedfalse | 1.5344x | 2.2659x | 1.5131x |
| 16x8/phase0-5/paddedtrue | 1.5359x | 2.2986x | 1.5244x |
| 16x16/phase3-0/paddedfalse | 1.4286x | 0.9197x | 1.4164x |
| 16x16/phase3-0/paddedtrue | 1.4400x | 0.9273x | 1.4324x |
| 16x16/phase0-5/paddedfalse | 1.4874x | 0.9555x | 1.4821x |
| 16x16/phase0-5/paddedtrue | 1.4934x | 0.9583x | 1.4796x |
| 32x16/phase3-0/paddedfalse | 1.4373x | 1.3551x | 1.3987x |
| 32x16/phase3-0/paddedtrue | 1.4410x | 1.3511x | 1.4291x |
| 32x16/phase0-5/paddedfalse | 1.4917x | 1.3997x | 1.4417x |
| 32x16/phase0-5/paddedtrue | 1.4824x | 1.3794x | 1.4523x |
| 16x32/phase3-0/paddedfalse | 1.4388x | 1.0387x | 1.4257x |
| 16x32/phase3-0/paddedtrue | 1.4431x | 1.0375x | 1.4288x |
| 16x32/phase0-5/paddedfalse | 1.4996x | 1.0822x | 1.4855x |
| 16x32/phase0-5/paddedtrue | 1.4991x | 1.0826x | 1.4842x |

The two vertical16x8 fixtures qualify independently of whole-codec results:
all12 paired samples exceed 1.5131x parent.
Other one-axis results and every slower control remain in the summaries/raw
logs; do not round a sub1.5 result into a qualifying gain. The24 identity/two-axis
controls range 0.9947–1.0169x parent.
Some sizes still trail assembly; this is not a blanket assembly-parity claim.

Whole synthetic720p VP8,realtime CPU8,2500kbps,one thread, six balanced rounds:

* 300frames: -0.0246% time versus parent (4/6 faster), +4.3681% versus assembly (0/6 faster).

* 60frames: +0.0005% time versus parent (3/6 faster), +11.5610% versus assembly (0/6 faster).

All23 non-timing fields match and every allocation sample is zero. Six fixed
speed packet/drop digests match the validated parent exactly. Preserve the
verified >=1.5x component gain even if whole-codec performance is unchanged;
the overall ARM assembly gap remains open. Shared-desktop timings and varied
assembly samples do not justify causal/additive claims across checkpoints.

Validation includes the full final native suite, dispatcher-focused tests,
ARM race/checkptr2, AMD correctness/checkptr under
Rosetta, actual AVX-off fallback, independent arithmetic/bounds/alias/dispatch
review, refreshed PGO and repository gates. New coverage includes864 wrapped
coefficient/alias/stride/shape cases,48 short-window cases and40 directed narrow
sum boundary cases including h65 fallback and w48 three-packet geometry.
All64 legal public phases and existing numeric extremes pass scalar oracles.
Selected root/prototype generated-code evidence and final source hashes remain
with the raw logs. Native x64 speed awaits the successor's exact-head CI.

The archive retains every diagnostic, all six rounds, all outliers, the earlier
qualifying unrolled checkpoint, unselected weighted-row cache/shift/multiply
variants and source/driver snapshots. The first packet-digest command matched
no tests; its warning is retained, then the explicit fixture was copied and
six real digests verified. No concurrent local builds/profiles during timing.
README files, dependencies and handwritten-assembly production policy remain
unchanged. No production cgo added.
