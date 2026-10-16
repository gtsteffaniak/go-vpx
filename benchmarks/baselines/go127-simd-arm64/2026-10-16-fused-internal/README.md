# Three-edge ARM internal-filter packets

Parent dff92311, original assembly 691cd051. Apple M4 Max, official Go 1.27.2,
GOEXPERIMENT=simd, GOMAXPROCS=1, PGO off. Timings are native ARM; Rosetta is
used only for correctness. Local tests/builds/profiling and other timing
experiments were serialized with the benchmark windows.

The three internal VP8 edges now share four neighboring rows/columns between
iterations. After each edge's sequential stores, four new rows/columns are
loaded for the next edge. Both orientations read 256 source bytes per tile
instead of 384, while preserving the same four-column output windows and
edge order. Loading new columns after the stores preserves overlapping
logical rows at pitches 8/9/15. Updated q0/q1 remain in the carried packet;
untouched q2/q3 cannot alias those stores at accepted pitches >=8.

Arithmetic is generated from the canonical copied goav1 byte network via the
existing exact-width internal specialization. The four-column transpose is
copied from the adjacent column loader and expanded inline. An alternative
calls that loader between edges. Six balanced triples measure:

| Kernel | Parent ns | Loader calls | Inline transpose |
| --- | ---: | ---: | ---: |
| Horizontal internal edges | 21.225 | 19.555 | 19.590 |
| Vertical internal edges | 62.085 | 48.435 | 39.855 |

Every sample has zero bytes/allocations. LLVM shows a next-column loader call
and fourteen static vector stack-load/store instructions in the call variant;
the inline variant removes that call and those vector spills. It retains only
the initial eight-column load call plus stack growth. Horizontal fusion has
no load/store helper calls. These explain the observed design difference on
this compiler/CPU; instruction counts are diagnostic evidence.

Eight balanced public quartets compare assembly, parent and both designs:

| VP8 720p | Assembly | Parent | Loader calls | Inline transpose |
| --- | ---: | ---: | ---: | ---: |
| ms/input | 6.807422 | 9.214026 | 8.788180 | 8.666572 |

Inline fusion improves 5.942% versus parent, every pair faster (4.20–6.73%).
The call variant improves 4.622%, every pair faster, establishing a stronger
compatible baseline than the parent alone. The inline design improves 1.384%
from that alternative. Original ARM assembly still leads by 27.31% in this
same-run comparison; the remaining migration gap is unresolved.

After integration/generation, six fresh balanced parent/prototype/final
triples confirm 9.192055 / 8.673118 / 8.629128 ms/input. The final build
improves 6.124% versus parent, all six pairs faster (4.22–6.69%). Its prototype
change is -0.507%; the design comparison above carries the main selection
evidence. Six alternating final kernel pairs confirm horizontal
21.140→19.620 ns (-7.19%) and vertical 62.035→39.815 ns (-35.82%), all zero
allocations. These realtime CPU-used 8, 1280x720, 60-input, 2500-kbps,
one-thread, encode-only runs preserve 685957 encoded bytes, 58 encoded / 2
dropped frames and zero allocations per input.

Six VP9 720p control pairs measure 12.040086→12.004477 ms (-0.296%), with
mixed paired changes (-0.62..+4.87%). No stable VP9 gain is established by
this VP8-specific change. All 62 design/fresh/control public samples preserve
every recorded non-timing field within their codec/workload, including
quantizer, frame-size and allocation fields. summary.json asserts the complete
field set rather than selected counters.

Five pixel/DSP/encoder packages pass on ARM and Rosetta AVX2. An additional
3456 exact-window cases compare the selected dispatch with the independent
int16 network, executing original edges in order across unaligned and
overlapping logical rows, varied masks/HEV thresholds, near-flat and random
bytes. Existing scalar, minimal/overflow-span and codec oracle checks remain.
The live 600-frame profile confirms fused vertical/horizontal bodies at
6.93%/2.38% cumulative time. Exact-source PGO refresh and pre-commit pass.

Fusion is enabled only on ARM. The x64 compiled VP8 DSP test instruction
region is byte-identical to dff92311 after excluding its embedded 128-byte
Go build-ID header: 2,375,377 bytes, SHA-256
ba6699cdde98da9a27b8f1ab42728ce914ff6f33170d24a0b40b6f521b5504b4.
The architecture constant keeps the measured x64 horizontal specialization
and restored general vertical path. Native x64 timings of the preceding
reduction/scope checkpoint are collected separately. Patches retain both
prototype designs against the final generated code.

Final validation also passes the full native suite (45.486s root package),
five-package race/checkptr1, strict encoder byte parity (34.636s), the 18/18
parity report and all source/architecture/fmt/PGO gates. The selected purego
sources are unchanged from the preceding full 117.286s passing suite; the
hosted workflow runs it again. Raw gate logs and VALIDATION.md are retained.
