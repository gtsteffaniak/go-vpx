# Rejected ARM horizontal partial-mask zero-correction pre-gate

Decision: **do not promote**. The one candidate passed scalar/vector, public-output and race/checkptr2 correctness. Fresh public picker flat/ramp improves 1.054–1.074×, below the user's 1.5× component rule, with mixed-pattern regressions. Six balanced 300-input whole-codec rounds are flat/mixed. No source was transplanted into the root project, and no README, branch or dependency was changed.

## Pinned source and numeric proof

Parent `f8dd8c7472e91c1f147511b04ceec716aca95aa7` has the exact tested tree subsequently merged as PR42. Old whole-codec reference `691cd0512c4829d9ee40d483205dd223f9eb028d` selects handwritten NEON where present. Apple M4 Max, official Go1.27.2 SIMD darwin/arm64, PGOoff, GOMAXPROCS=1, one encoder thread. Existing parent/old CLI hashes were checked against the frozen Copy8 experiment; only the candidate CLI was built for this follow-up.

Only ARM fused horizontal macro+three internal filters change. The generator performs the transformation after the existing rounded-zero and vertical specializations, matching exactly two horizontal blocks. Each computes threshold-mask need, HEV and signed f before the remaining four neighbor absolute differences. If every partial f has magnitude below3 macro/below4 internal, it skips those differences and full-neighbor need. Otherwise it computes the original full need, masks the same f again and retains the original final rounded-zero guard/correction arithmetic. It removes the original two-GPR needAny check. This is a separate arithmetic/dataflow pre-gate, not the rejected shift-constant, ReduceMax replacement, tile or partial-C/D cache variants.

Full need is a lane-wise subset of partial threshold need, so final f can only retain partial f or become zero. For macro |f|≤2, `(weight*f+63)>>7` is zero for weights9/18/27, and saturating `(f+3)>>3` and `(f+4)>>3` are zero. For internal |f|≤3 those shifts and the derived outer correction are zero. The unsigned reinterpretation after ABS treats -128 as128. Therefore a skipped block would not change any stored pixel or carried neighbor under either HEV branch. Row loads, macro/internal ordering and stores on changing paths are preserved.

`TestHorizontalPreGateSignedProof` enumerates every signed-byte f, both kinds, all16 lane positions, keep/clear full masks and both HEV choices. Existing independent 7200-case fusion scalar oracle, 640 zero/cancel/late-lane cases and complete filter/UV/bounds/alias tests pass in pixel, VP8 DSP and decoder, both ordinary and race+checkptr2. Both decoder benchmarks also pass one-shot complete-backing-output comparison. Full command/raw output logs are retained.

Regeneration after formatting is byte-identical to the measured generated source. Entire AMD generated source and the ARM vertical function remain byte-identical to parent. The snapshot has only the generator, ARM generated horizontal body and one test added/changed. Optimized ordinary H body: parent364/candidate368 instructions, both zero calls and zero vector stack loads/stores. Counts are diagnostics rather than acceptance evidence.

## Public picker screen and fresh confirmation

The unchanged public BenchmarkVP8PickerFilterNative fixture restores the exact union of horizontal/vertical read windows per B.Loop iteration and checks complete output against four separate edge calls. Independent scalar pixel oracles run before timing. Three exploratory 75ms alternating pairs cover all10 cases; six fresh100ms pairs start with reversed order. Every sample has0B/0alloc. All60 exploratory and120 fresh samples, slower controls and all pair ratios are retained.

| Fresh case | Parent/candidate median ratio | Faster pairs |
|---|---:|---:|
| pitch32/flat | 1.063725× | 6/6 |
| pitch32/ramp | 1.054258× | 6/6 |
| pitch32/smooth | 1.001448× | 2/6 |
| pitch32/mixed | 0.994491× | 2/6 |
| pitch32/noop | 1.004254× | 4/6 |
| pitch1312/flat | 1.066263× | 6/6 |
| pitch1312/ramp | 1.073521× | 6/6 |
| pitch1312/smooth | 0.989139× | 2/6 |
| pitch1312/mixed | 0.987764× | 1/6 |
| pitch1312/noop | 1.006708× | 5/6 |

Ratios are ratios of medians, not medians of paired ratios. Fresh smooth pitch32 median is near-flat despite4/6 individual regressions; the earlier screen regressed2.50%. Mixed and padded smooth controls remain explicitly retained. No1.5× component or assembly-component gain is claimed.

## Whole-codec confirmation

Six permutations of parent/candidate/old691 were run serially with prebuilt CLIs: VP8,1280×720,300 inputs,30fps,2500kbps,realtime CPU-used8,one thread,encode-only. Root and other agents held Go/build/CPU-heavy work during the exclusive measurement window. Positive CPU8 is adaptive, so recorded workload parity was checked rather than assumed.

| Variant | Median measured encode ns/input |
|---|---:|
| f8 parent | 7,330,001 |
| candidate | 7,333,907 |
| old691 | 6,850,321.5 |

Candidate time is **+0.053288% vs parent**,4/6 pairs faster; individual changes range **−1.004332% to+0.417984%**. This is flat/mixed and fails the agreed requirement for a clear consistent whole gain. Against old691 in this same run, candidate is **+7.059311%** by ratio of medians and slower in all6 pairs (gap3.682864–11.468166%). Old-reference variance is retained; these figures are this experiment's scope and cannot be subtracted from other runs/CPUs or treated as universal assembly gaps.

All18 JSONs match **all22 semantic fields and separate zero allocations**. Raw JSON, stderr, command order, every per-input timing, output/frame/drop/quantizer fields and complete comparison records are in final-summary.json. Encode-only quality is skipped; matching recorded statistics is not bitstream-identity proof. CLI timer measures encoding rather than subprocess startup/allocation-sampling wall time.

## Reproduction and archive

Use the pinned parent, apply candidate.patch, then run scripts/gensimd.go through the root scripts/go.sh/project-env wrapper with an agent-specific cache ID. Sources include exact parent/candidate generator/generated files, new numeric test, unchanged public fixture and relevant shared helpers/oracles. No complete snapshots or binaries are copied; binary/source hashes identify measured inputs. run_final.py preserves exact benchmark arguments/order and summarize.py preserves audit calculations. The generator and ordinary/race commands are recorded in reproduction-commands.txt.

SHA256SUMS covers every delivered file except itself. raw-evidence.tar.gz contains all other delivered members except the outer manifest/archive itself; its inner hash list verifies those members. Keep this rejected variant only as research evidence, and retain the independently qualifying Copy8 candidate separately.
