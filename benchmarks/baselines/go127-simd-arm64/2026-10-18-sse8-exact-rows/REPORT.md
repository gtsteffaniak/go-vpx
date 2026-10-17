# Exact-row SIMD squared error

Apple M4 Max, Go 1.27.2 darwin/arm64 archsimd, PGO off, GOMAXPROCS=1.
Immediate parent 1e8809b63498b77fb6a347b0c085ebfe3c9aab49 (PR37), old assembly
691cd0512c4829d9ee40d483205dd223f9eb028d. Native speed only; Rosetta results
are correctness checks, never x64 performance evidence.

SSE8's bounded path proves h=1..64 and nonnegative strides<2^30, then checks
both final eight-byte row endpoints. Exact uint64 row loads need no padding,
including unaligned, zero-stride, overlapping or mutually aliased read-only
rows. Typed unsafe pointers stay live and are never retained or written.
Offsets and unused final increments fit int64. Larger heights and unusual
strides keep the original checked-slice fallback byte for byte, including
h1 with an unused negative or maximal stride and empty geometry behavior.

Pack two rows by inserting both uint64 lanes directly. Compute unsigned byte
Max-Min before widening: this is |a-b| in 0..255 for every ordered byte pair,
so squaring is exact. Separate low/high int32 sums cover each row once; an odd
tail accumulates only low8 bytes. At 64 rows total SSE<=33,292,800 fits int32.
No new AVX256 state, CPU requirements, assembly, dependencies or production cgo.

Fresh six balanced parent/assembly/current orders,100ms per 36 identical public
B.Loop fixtures. Four APIs, three strides 8/8,32/37,1312/1320, and random,
maximum-error, identical-pixel patterns. Exact visible spans, origins3/5,
independent per-pixel scalar SSE before/after timing and full backing canaries.
All samples allocate zero. Ranges below cover all three pixel patterns.

| Public API / strides | Parent/current time | Assembly/current time |
|---|---:|---:|
| slice8x4 / 8-8 | 1.2961–1.3112x | 2.2662–3.2790x |
| slice8x4 / 32-37 | 1.3021–1.3104x | 2.4183–3.1861x |
| slice8x4 / 1312-1320 | 1.3148–1.3170x | 2.5106–2.9633x |
| slice8x8 / 8-8 | 1.4807–1.4975x | 0.9947–1.0612x |
| slice8x8 / 32-37 | 1.4769–1.4926x | 0.9975–1.0682x |
| slice8x8 / 1312-1320 | 1.4875–1.4986x | 1.0009–1.0360x |
| slice8x16 / 8-8 | 1.6253–1.6354x | 0.8515–0.8718x |
| slice8x16 / 32-37 | 1.6294–1.6389x | 0.8551–0.8696x |
| slice8x16 / 1312-1320 | 1.6262–1.6423x | 0.8497–0.8646x |
| pointer8x8 / 8-8 | 1.4389–1.4465x | 0.6988–0.7076x |
| pointer8x8 / 32-37 | 1.4336–1.4452x | 0.6952–0.7108x |
| pointer8x8 / 1312-1320 | 1.4372–1.4535x | 0.6990–0.7016x |

All nine 8x16 cases qualify: median 1.6253–1.6423x parent;
all 54 paired samples exceed 1.5802x. Keep this correctness-verified
component gain independently of whole-codec changes. Other sizes remain
below 1.5x parent; eight/sixteen-row and pointer paths still include assembly
deficits. This does not establish blanket assembly parity.

Whole synthetic 1280x720 VP8, realtime CPU 8, 2500 kbps, one thread, six balanced
orders, separate from kernel timing (positive means more time/slower):

* 300 frames: -0.1075% versus parent (3/6 faster), +5.9503% versus old assembly (0/6 faster).

* 60 frames: +1.4459% versus parent (0/6 faster), +12.5059% versus old assembly (0/6 faster).

Fresh separate 60-frame follow-up: -0.6109% versus parent (5/6 faster), +9.4058% versus old assembly (0/6 faster). Retain both runs; do not replace the first regression with this follow-up.

All 23 non-timing fields and every allocation sample match, with zero allocations.
These output statistics do not assert bitstream identity. Shared-desktop timing
and varied assembly samples do not support causal/additive claims across runs.
The whole ARM assembly gap remains open. No libvpx C result is inferred from
these assembly comparisons; C has its own native CI job and quality scope.

Correctness: full native suite passes, ARM race/checkptr2 passes, x64 checkptr2
under Rosetta passes, and actual AVX-off fallback passes (SIMD-only tests
explicitly skip there). New tests exercise 720 disjoint +720 mutual-alias calls,
14 short-window panics, unused extreme strides, odd/boundary/fallback heights,
and every one of 65,536 ordered byte pairs in 512 eight-by-sixteen calls.
Independent static review verified arithmetic, endpoints, pointer lifetime,
odd tails, fallback equivalence and live benchmark scope. The first private
full-suite run lacked two upstream C control-table fixtures; both its failure
and the successful rerun with exact pinned fixture hashes are retained.
Refreshed PGO and repository gates are logged. Native x64 performance awaits
this successor PR's exact-head CI; its harness retains 246 controls and adds 36
public fixtures, both new source hashes and focused selected/race checks.

Every diagnostic, rejected variant, raw sample and source snapshot is retained.
The initial row-pointer-only change topped out below 1.5x. Four-row unrolling
alone did not qualify; byte absolute differences helped 16-row cases, and
removing the overwritten broadcast produced the selected winner. Fresh runs
confirm it independently of those exploratory measurements. Generated code
is diagnostic evidence, never a substitute for timing. No competing local
Go builds/tests/profiles ran during native measurements. README files are
untouched; experiment records belong here.
