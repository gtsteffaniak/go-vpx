# VP9 batched reductions and VP8 rolling prediction

Parent 649237be, original assembly 691cd051. Apple M4 Max, official Go 1.27.2,
GOEXPERIMENT=simd, GOMAXPROCS=1, PGO off. Timings are native ARM; Rosetta is
used only for correctness. Local builds, tests, profiling and other timing
experiments were serialized. Drivers, raw samples, independent tests,
profiles, disassembly and a production patch are retained here.

Column projection copies goav1's four-independent-row dataflow from commit
ff46346e73f063ff1ea3a8d7f06ee26cd9baaad4, with safe slice reads, VP9 widths
16/32/64, variable normalization and remainder rows. No dependency is added.
Each accepted row totals at most 64*255=16320, preserving signed int16
conversion before the right shift. ARM uses uint16 lane accumulation; the
local x64 adaptation uses independent PSADBW accumulators. Other shapes and
normalizations retain the existing scalar/general fallback.

The four rounded 8x8 averages of a 16x16 region now share full row loads.
ARM keeps four uint16 accumulators (each lane at most 8*255=2040); x64 uses
two PSADBW packets, whose low/high lanes represent left/right blocks.
Divide-before-multiply guards prove complete source/output bounds without
integer overflow. Exact/minimal and unaligned tests include zero/overlapping
read strides, extreme pixels, remainder rows and unsupported fallbacks.

VP8 vertical six/four-tap prediction retains five/three widened rows in
registers and loads one new row per output. Disjoint input/output spans and
non-overlapping output rows are required for this traversal. Aliased buffers
retain the previous eight-pixel sequencing, checked against an independent
packet oracle. Four-pixel widths and single-row packets retain the existing
path. An initial unconditional helper call regressed 4x4 vertical kernels
about 8%; the final guard removes that cost. Fresh final 4x4 cases vary only
-1.04% to +1.19%, while horizontal-only 8/16 cases vary +0.85%/+0.55%.

Eight balanced kernel quartets compare parent, projection alone, average
alone and both. Six alternating fresh parent/final pairs confirm:

| Public DSP operation | Parent ns | Final ns | Time change |
| --- | ---: | ---: | ---: |
| Four 8x8 averages | 51.430 | 8.748 | -82.99% |
| Projection 16 columns × 128 rows | 390.850 | 90.415 | -76.87% |
| Projection 32 columns × 128 rows | 518.950 | 145.100 | -72.04% |
| Projection 64 columns × 128 rows | 781.400 | 245.100 | -68.63% |
| Six-tap 8x8 vertical | 42.720 | 21.845 | -48.86% |
| Six-tap 16x16 vertical | 144.450 | 53.665 | -62.85% |
| Six-tap 16x16 both axes | 223.150 | 133.200 | -40.31% |
| Odd-phase 16x16 both axes | 184.800 | 118.900 | -35.66% |

All samples report zero bytes and allocations. The single-row projection
cases also improve 8.21–17.96%; short seven-row packets improve 55.54–59.42%.
B.Loop fixtures call production DSP dispatch and expose results.

Ten balanced public quintets isolate projection, average and their
combination against parent and original assembly. VP9 720p medians:
assembly 8.476611 ms, parent 12.417141 ms, projection 12.366217 ms,
average 12.351885 ms, combined 12.104621 ms (-2.517% from parent).
The combination is faster in nine of ten pairs (-5.60..+1.39%). Its 42.80%
extra time versus original ARM assembly remains unresolved. Six separate
VP8 triples select rolling rows: assembly 7.086106 ms, parent 9.440474 ms,
prototype 9.265941 ms (-1.849%, every pair faster). That prototype still
trails assembly by 30.76%. Separate-session absolute timings cannot isolate
changes; these assembly comparisons belong to their respective same-run
prototype experiments.

After integration, six fresh balanced parent/prototype/final triples per
codec compare the exact final combination, realtime CPU-used 8, 1280x720,
60 inputs, 2500 kbps, one thread, encode-only:

| Codec | Parent ms/input | Prototype | Final | Final vs parent |
| --- | ---: | ---: | ---: | ---: |
| VP8 | 9.398132 | 9.235976 | 9.288409 | -1.167% |
| VP9 | 12.550278 | 12.338014 | 12.187932 | -2.887% |

VP9 is faster in all six pairs (2.36–7.20%). VP8 has mixed paired results
(-2.06..+2.40%, four of six faster), so a stable final whole-encoder VP8 gain
is not established by this fresh set. The large predictor kernel gains are
scoped DSP results. All 104 public/design/fresh samples preserve every
recorded non-timing field within each codec/workload, including quantizer,
byte/frame counts and allocation fields. VP8 remains 685957 bytes, 58/2
encoded/dropped and zero allocations/input; VP9 642654 bytes, 48/12 and
1.25 allocations/input. These counters are distinct from the strict bitstream
parity gate, which remains active.

Five pixel/DSP/encoder packages pass on ARM and Rosetta AVX2; three pixel/DSP
packages also pass with the SIMD experiment off. Direct backend tests prove
native selection for valid shapes and untouched rejection of invalid spans.
The 600-frame live profiles confirm ProjectionCols16 (1.16% cumulative),
AverageQuad16 (0.46%) and rolling six-tap prediction execute in the encoders.
ARM LLVM output shows widening/add/reduction and rolling NEON arithmetic;
x64 LLVM output shows VPSADBW on XMM registers, with no YMM body. No native
x64 timing is claimed yet; the expanded PR benchmark measures both new VP9
operations against parent and original assembly.

Native x64 649237be evidence arrived before this checkpoint was pushed.
Horizontal internal filtering improves 13.09%, but its ARM-tuned vertical
stores regress 2.85% in every pair. The final sources therefore retain the
measured previous generic vertical SIMD body on x64, while keeping the
specialized horizontal body and ARM's measured vertical specialization.
Five-package ARM/Rosetta tests and a fresh exact-source PGO/pre-commit pass
again. The architecture constant removes this branch completely on ARM:
all 3,435,188 instruction-region bytes are identical before/after the scope
change (SHA-256 bc524538640d180e76abc2daa0bde4b4c334f8052100846af1a836c6f3761c38).
Only the 112-byte Go build-ID header preceding the first function differs.
The parser, raw hashes and matching entry-point disassembly are preserved;
ARM timing above remains evidence for the final executable instructions.
Native timing of the x64 scope change and reductions is still pending.
