# Rejected hoisted DC caller guard

Baseline is exact `10caaa67ed53f8c3c1248f20070e5a4117245852`. Only luma in decoder/reconstruct.go changes; chroma and all backend kernels remain identical. This is a separate experiment from the initial per-pair guard, whose frozen package remains unchanged.

The guard is safe: stride>=8, len>=16 and stride<=(len-16)/15 prove every off+3*stride+8 is within len without signed overflow. The original qflat construction precedes the pure guard. Zero DC produces byte identity, including read/write aliases and overlapping full luma rows at stride8..15; all skipped original slices are already proven valid. Other paths and panic order remain original.

All six candidate correctness modes pass 41,920 frozen-parent comparisons each (251,520 total): native ordinary, native race/checkptr2, purego, nosimd, AMD selected/off checkptr2. The tests compare full backing, qcoeff/EOB/dequant and exact panic messages, including full-transform mutation, seeded/floor cases, input aliases, nil/invalid indices, lengths/caps, negative and extreme strides. AMD execution is translated correctness only. Native LLVM luma grows176→194 static instructions with96-byte frame unchanged. Chroma source and normalized112 instructions are identical to parent.

The Apple M4 Max screen used official Go1.27.2 SIMD, CGO1, PGOoff, GOMAXPROCS1, six alternating AB/BA rounds75ms, 18 shared caller workloads:216 zero-allocation samples. Every timed invocation includes frozen-parent pre/post full-state checks and live output. Root confirms no new hoisted one-shot invocations were run: zero excluded preflight rows are claimed for this window. The initial parent18 one-shot rows are retained only as historical validation. Initial candidate one-shots are not attributed to hoisted; no post-screen rows were created.

Luma seeded-zero medians improve1.311–1.369x and mixed1.112–1.123x, all36 pairs faster. Luma nonzero ratios0.958–0.981x,17/18 pairs slower. All54 unchanged chroma control pairs are slower, medians0.892–0.917x; identical chroma source/codegen does not erase these measured losses or establish their cause. No case reaches1.5x. Root rejects this tradeoff. No whole-codec, x64 throughput or assembly inference is made.

All sample arrays, medians, paired ratios and outliers are in compact.json; actual commands/timestamps, source maps, relevant source/patch context, overlays, logs and LLVM are archived. Executables/full snapshots are omitted with original paths/SHA256/size recorded. No ROOT/Git/README/dependency/assembly changes, no test threshold weakening, and no new timing or source edits occurred during this audit.


Publication: `raw-evidence.tar.gz` preserves all 5 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
