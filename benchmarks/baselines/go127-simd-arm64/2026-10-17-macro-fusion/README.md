# ARM VP8 macroblock/internal edge fusion

The ARM SIMD path now filters a 16-lane macroblock boundary and its three
internal edges in their original order. The generated arithmetic comes from
canonical goav1-derived byte filtering; it retains updated neighbors, loads
only the incoming four rows/columns, and stores exactly six or four filtered
pixels. Vertical incoming loads pack four exact uint32 rows into a vector and
recover columns with eight TRNs. Other architectures compile out the whole
fused call, including slice construction.

Apple M4 Max, Go 1.27.2, SIMD on, PGO off, GOMAXPROCS=1 and encoder threads=1.
Native timings run without competing local Go builds or tests. Production
parent is 9a442f66; the old Go assembly control is 691cd051. Inputs are the
public 1280x720 realtime VP8 synthetic benchmark, 2500 kbps, cpu-used=8.

| Final integration run | Parent ms/input | Current ms/input | Current time delta | Pairs faster |
|---|---:|---:|---:|---:|
| Initial 60-frame, three variants | 8.441168 | 8.322421 | -1.41% | 4/6 |
| Fresh 60-frame confirmation | 8.351395 | 8.110483 | -2.88% | 6/6 |
| Fresh 300-frame confirmation | 8.368106 | 8.105590 | -3.14% | 6/6 |

The initial six-permutation run includes same-run old assembly at 7.283229 ms;
current remains 14.27% slower, with current taking more time in all six pairs. The older
26.6% and 17.75% gaps were measured in different runs; do not subtract them to
attribute a speedup. Fresh confirmations compare current/parent only, and
cannot establish a new absolute gap against assembly. All non-timing JSON
fields, quantizer distributions, payload/frame/drop counts and allocations
match in every final comparison. VP8 stays at zero allocations/input.

The unchanged VP9 control changes +0.456% by medians with paired changes
-0.355..+1.108%; no VP9 speedup is claimed. It retains 1.25 allocations/input.
The selected x64 decoder function is byte-identical to the parent: 850 bytes,
217 instructions, same four DSP calls and bounds checks. See the codegen JSON.

Design exploration is retained, including vertical-only/horizontal-only,
inlined stores, four-column transposes, packed 64-bit vector inserts and
helper-expression expansion. Prototype confirmation favored the selected
packed input layout by 3.05% over the old production path, all eight pairs.
Further helper expansion reduced hardware NOPs and instruction count but did
not establish a reliable extra public-call improvement, so it is not included.
Prototype and integrated-source timings are labeled separately in raw files.
The integrated-source confirmation, not instruction count, supports keeping
the final change.

The independent int16 SIMD network compares the sequential macroblock plus
three internal filters across both orientations, five pitches, eight byte
alignments, 18 flat/random trials and five threshold sets (7200 cases).
Minimal accepted windows retain prefix/suffix sentinels. Invalid thresholds,
short windows, negative strides and maximum-int strides decline without stores.

An additional 9216 cases compare the public fused kernel directly against the
independent scalar libvpx v1.16.0 port, including active filters near byte
extremes, exact minimum windows, eight alignments and full-allocation sentinels.
Both oracle suites pass with race detection and checkptr enabled.

Full native suite, experiment-off/purego codec suites, strict byte parity,
18/18 parity report, refreshed PGO, formatting, assembly absence and all six
architecture/build selections pass. PGO fingerprint: `8e62a32941a321454f75cf470b9787f85a9a2257`.
