# Private AMD FP16 overlap guard variants

Base: 46a08f8e6eb03f6ee4549d975384bd6c952a81fa. No repository writes, PR changes or native AMD timing. Both candidates remain UNMEASURED performance proposals. Rosetta runs are correctness-only; their incidental ns/op values are not speed evidence.

Only fp16WindowsDisjoint changes. Public validation/eligibility order, all ten pair policies including read-only overlaps, optional nil omission, 16-lane arithmetic/store order, ClearAVXUpperBits and narrow fallback source are exact. Existing ARM sources are byte-identical. The live ARM quantizer normalized 200-instruction body is identical (encoder binary; initial pixel binary dead-stripped this function).

Independent proof uses occupied-byte sets for 242,010 valid small-ring cases (966/14,620/226,424) plus 5,604 wide-boundary interval cases. For individually nonwrapping positive windows with length sum<=uintptr modulus, overlap iff uintptr(a-b+LA-1)<LA+LB-1. LA belongs to a. Both adjacencies are excluded; same pointer included. Zero lengths are outside the identity; existing n64..1024 eligibility remains. No fabricated pointer is dereferenced.

## Generated code

Darwin and Linux match these counts, excluding trailing INT3 padding:

| Guard | Instructions | Conditional jumps | Local frame | Calls | Inline cost |
| --- | ---: | ---: | ---: | --- | ---: |
| Parent | 83 | 9 | 80 bytes | panicBounds cold path | 146 |
| Single compare, original array/loops | 75 | 6 | 80 bytes | panicBounds cold path | 135 |
| Explicit ten pairs | 88 | 13 | 0 bytes | none | 185 |

All remain noninline. Unrolled guard still has three argument-home stores into the caller frame and stack argument loads; it is not stack-free. Unrolling removes loop/array setup, not every branch. Static counts cannot select the winner: single compare has fewer static instructions, unrolled executes fixed checks without the nested-loop work. Native six-balanced public36 measurements must compare both against this strong base. No 1.5x claim is made.

## Correctness

All three variants pass existing and new FP tests with selected AVX2 and actual AVX/AVX2 CPU-off, both ordinary and checkptr2. New tests cover all ten pairs at all61 eligible sizes (26,840 direct guard comparisons), 1,320 aliased public full-backing comparisons, 37 nil/empty guard cases and160 public independently short/nil/spare-cap/clipped-cap cases. Frozen public oracle compares EOB, exact panic strings and every backing byte. Existing full-signed, alignment/minimal-window, invalid-table and disabled tests remain.

216 ordinary public36 one-shot samples are all0B/0alloc.216 checkptr2 one-shot samples preserve outputs; checked classes cases allocate16B/1alloc on parent and both candidates, and all allocation vectors match parent exactly. Initial parser incorrectly required checked zeroalloc; runtime assertions passed, parser was corrected and failure output retained in session history. Purego and noSIMD encoder tests pass for parent/unrolled. Reconstructed complete-command-inventory.json enumerates68 commands from executed driver logs; two additional live ARM encoder builds have their own exact command file. This reconstruction avoids concurrent-driver journal overwrite; raw per-command logs are retained.

## Preservation and limits

Variants, all tests, exact drivers, proof, LLVM bodies, compile diagnostics and raw correctness/one-shot logs are included. Binaries are retained at immutable private source paths listed with size/SHA256 in binary-inventory.json and omitted from this lean package only. Full snapshot trees remain at /private/tmp/govpx-fp16-offset-guard. No SDK/library dependencies added. No native throughput or whole-codec claim. Root must select any later CI experiment; this package does not change the active PR.


Publication: `raw-evidence.tar.gz` preserves all 126 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
