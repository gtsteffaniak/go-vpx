# Fixed VP8 internal-edge packets

Parent 44090994 has the same runtime sources as 3792c617. Apple M4 Max,
official Go 1.27.2, GOEXPERIMENT=simd, PGO off, GOMAXPROCS=1. Timings are
native ARM. Local Go builds, tests, profiling and other experiments were
serialized with these measurements. Rosetta is used only for correctness.

The live VP8 path filters three internal edges in order. The new generated
packet fixes sixteen positions and four output columns, retaining canonical
byte arithmetic from the copied goav1 network. General kind/width branches
are eliminated. Horizontal output uses four sixteen-byte stores; vertical
output transposes two eight-position halves into sixteen exact four-byte
stores. These are the same memory windows and ordering as the prior helper.
The checked entry point proves 15*(step+outer) < len(p) by dividing before
multiplication; oversized strides reject without overflow or writes.

Four designs were compared: A fixes kind/lanes, B also inlines exact-width
stores, C splits orientation, D combines A's horizontal helper with C's
vertical specialization. Initial six-triple kernel rounds narrowed the
choices. Six final A/B/D triples favor D's vertical microbenchmark, but eight
balanced public quartets favor B for the complete operation:

| VP8 720p variant | Median ms/input | Change from parent |
| --- | ---: | ---: |
| Parent | 9.769700 | — |
| A fixed kind/lanes | 9.591729 | -1.82% |
| B exact-width stores | 9.517802 | -2.58% |
| D separate orientation | 9.593962 | -1.80% |

Every B/parent pair is faster (0.06–3.95%). C was not retained after its
horizontal kernel regressed relative to A and its extra vertical difference
from B was small. D's microbenchmark win did not establish a stronger public
result. Patches preserve all rejected designs against the final source.

After integration and generation, six fresh balanced parent/prototype/final
triples select the final B build: 9.703888 / 9.499927 / 9.492774 ms/input.
The final build improves 2.176% from parent, all six pairs faster (1.12–2.87%).
Prototype/final differs by -0.075%, within noise. These are VP8 realtime,
CPU-used 8, 1280x720, 60 input frames, 2500 kbps, one thread, encode-only.
Every run retains 685957 encoded bytes, 58 encoded / 2 dropped frames and
zero allocations per input. All recorded non-timing fields match across the
32 design and 18 fresh samples; summary.json asserts the full field set.

Six fresh alternating parent/final kernel pairs measure horizontal
22.075→20.750 ns (-6.00%) and vertical 73.600→63.465 ns (-13.77%). Every
sample reports zero bytes and allocations. These public DSP benchmarks call
the actual dispatch and retain the resulting mutable output.

Six VP9 720p control pairs measure 12.908325→12.693277 ms (-1.666%); the
changed VP8 entry point does not establish the cause of this control
variation, so no VP9 algorithmic speedup is claimed. All twelve controls
retain every recorded output/allocation field. Separate-session absolute
timings cannot isolate a change. This checkpoint does not remeasure the
original ARM assembly gap, which remains an unresolved migration limit.

Five pixel/DSP/encoder packages pass on ARM and Rosetta AVX2. Expanded exact
window tests compare all backing bytes for unaligned and overlapping logical
rows (strides 8/9/15), and reject nil, short and overflow-sized spans without
writing. Existing randomized scalar and codec oracle tests remain active.
The 600-frame live profile confirms loopFilterInternalEdge is executed in
the complete encoder (14.90% cumulative). LLVM disassembly shows NEON byte
saturation and only the vertical load/transpose call plus stack growth; the
general storeLFBytes call is gone from this specialized body. These are
code-generation diagnostics, not a separate performance claim. Native x64
timing of this checkpoint is pending the PR workflow.
