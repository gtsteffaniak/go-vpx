# ARM wider subpixel variance and non-flat VP9 filtering

Parent ef35b4ba, original assembly691cd051; Go1.27.2, SIMD experiment,
PGO off, Apple M4 Max, GOMAXPROCS=1. No concurrent local Go validation or
benchmark jobs. `govpx-wide-public.py` preserves the balanced schedule:
ten quintets per codec, each variant occupies each position twice, with
reversed orders balancing predecessors. All100 public samples are retained.

Variants are original assembly, parent, non-flat-only, safe subpixel16-only,
and the retained combination. Production references remain pinned; only
identical flat-filter test/benchmark fixtures are added to the parent archive.

| 720p realtime CPU-used8 | Assembly | Parent | Non-flat only | Subpixel16 only | Combined |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP8, ms/input frame | 6.981 | 9.956 | 10.034 | 10.034 | 10.006 |
| VP9, ms/input frame | 8.685 | 13.825 | 13.052 | 13.395 | 12.754 |

VP9 improves7.745% versus parent; all ten pairs improve5.87–9.65%.
VP8's +0.503% median change is not an established regression or gain:
paired changes range from -7.34% to +11.74%. Every non-timing JSON field
matches across all variants, including output bytes, quantizer histograms,
encoded/dropped counts and allocations. VP8 uses0 allocations/input frame,
VP9 uses1.25. These are encoder benchmarks, not substitutes for strict
byte parity. ARM still trails assembly: VP8 +43.33%, VP9 +46.86% time.

The non-flat byte filter performs the standard VP9 narrow update only
when no active lane requires a flat filter. Rejected packets are untouched
and use the existing7/15-tap path. The copied arithmetic is generated from
the canonical goav1-derived byte filter; constants specialize eight lanes.
Focused tests cover sparse-flat packets, exact/minimal allocations,
misalignment, overlapping layouts/order and the255 threshold fallback.
Independent DSP scalar/reference comparisons also pass.

Six initial ARM microbenchmark pairs improve existing non-flat cases43–53%.
A second six-pair comparison specializes the packet width and improves
another2.9–6.7%. Real flat packets incur an additional test before fallback;
`govpx-wide-final-flat-ab.txt` records the final-code tradeoff separately: H8 +14.08%, V8 +17.04%,
H16 +12.02%, V16 +4.77%, with zero allocations.

The new subpixel kernel processes sixteen pixels per iteration: paired
NEON int16 vectors on ARM, one AVX2 int16 vector on amd64, retaining the
horizontal rows in registers. Safe slice reads suffice. Signed differences
and squared error remain exact; final modulo32-bit accumulation matches
the previous kernel. Four/eight-pixel cases retain their existing path.
Six ARM pairs improve16x16 49.195→39.040ns (-20.64%),32x32
189.2→150.75ns (-20.32%),64x64 750.2→597.35ns (-20.37%), all0alloc.
An unsafe ARM load variant changes medians only0.64–1.15%; it was rejected
because this evidence does not establish enough benefit for extra unsafe
access. The retained subpixel implementation adds no unsafe accesses.

Five pixel/DSP/encoder packages pass on ARM and Rosetta AVX2. Tests include
all bilinear phases, byte extrema, zero/overlapping strides, minimal final
rows and misalignment. Rosetta establishes x64 correctness only. Attached
LLVM disassembly confirms live AVX2 VPMADDWD and VEX vector operations,
with no legacy XMM moves in the YMM loop and ClearAVXUpperBits on return.
The fresh360-frame ARM profile confirms both new functions on the public
encode path (subpixel16 6.50% cumulative; byte non-flat3.15%). PGO is
refreshed and pre-commit passes. Native x64 performance of this checkpoint
is pending; its paired-kernel parent is separately fully validated.
