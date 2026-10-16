# Paired VP8 residual gather

Apple M4 Max, official Go 1.27.2, GOEXPERIMENT=simd, GOMAXPROCS=1, PGO off.
All timing is native ARM; Rosetta AVX2 is correctness evidence only. Local
Go builds/tests/profiles were serialized with timing windows. The pre-gather
VP8 benchmark binary is c49e9d60; parent a6c1f016 changes only VP9's empty-row
projection behavior. Integrated candidates include that correctness fix.

The previous gather extracted each four-coefficient packet to scalar uint64
and stored it separately. The new layout interleaves two rows, directly
storing eight int16 values at a time. Two pairs form each four-row group,
retaining VP8's exact block scan order and halving output stores. The 8x8
path packs two exact eight-byte input rows using the copied goav1 helper.
No dependencies or new unsafe operations are added. Read-only rows may
repeat or overlap; output/input overlap remains outside the existing contract.

Plausible designs compare one-row stores, a two-row loop and a four-row
loop with two explicit packets. First optimizing only 16x16 halves its kernel
time but gives mixed public pairs (-0.30% median); this does not establish
an encoder win. Including 8x8 gives six balanced public triples:

| VP8 720p ms/input | Parent | Two-row loop | Four-row loop |
| --- | ---: | ---: | ---: |
| Both widths |8.566673|8.439965|8.412200|

Both designs improve every parent pair. Four-row grouping improves 1.803%;
six fresh integration triples confirm 1.462%, again every pair faster.
A final safe-indexing design checks output size once and uses a bounded
block index. Go's BCE diagnostics retain the initial checks and input
checks, but remove the repeated vector output slice/store checks. No new
unsafe proof is needed. Six alternating kernel pairs compare the stronger
four-row design with this bounded variant before final integration.

Six fresh balanced parent/four-row/final triples confirm:

| Operation | Parent | Four-row design | Final bounded design |
| --- | ---: | ---: | ---: |
| Gather 16x16 ns |40.230|21.645|16.010|
| Gather 8x8 ns |18.760|9.2045|8.084|
| VP8 720p ms/input |8.574611|8.465418|8.406684|

Final kernels use 60.20%/56.91% less time than parent, all zero bytes and
allocations. Public frame time improves 1.958%, every pair faster
(0.716–4.326%). Against the stronger four-row design, the kernel improvement
is 26.03%/12.17%; public time improves 0.694% with mixed pairs (five of six
faster), so that small public difference is not treated as a stable gain.

These realtime CPU-used 8, 1280x720, 60-input, 2500-kbps, one-thread, encode-only
runs preserve 685957 encoded bytes, 58 encoded/2 dropped frames, zero public
allocations/input and every other recorded non-timing field. All 72 VP8 and
12 VP9 design/fresh/control samples preserve their full non-timing field set,
including allocation and quantizer fields. Six VP9 control pairs are mixed
(-0.302% median), establishing no stable VP9 change. `summarize.py` checks
all fields except the four explicit timing fields and retains raw medians.

Existing independent scalar/minimal/unaligned/extreme tests and 80 new cases
cover repeated and overlapping read rows. Five pixel/DSP/encoder packages
pass on ARM and Rosetta AVX2. LLVM evidence shows NEON ZIP and VEX128
interleaves/stores; no new YMM work is introduced. A 600-input live profile
records Gather at 2.44% cumulative time. Both prototype loops and the original
source are retained as patches against the final implementation.

Production sources were frozen before refreshing PGO; fingerprint
`9ef8ad9eb8de1186550791d0c6dc0c74230b9750`. Broader exact-source gates and
native x64 results are tracked in the validation record and PR. Native x64
timing is pending; the earlier 9V45 VP8 assembly gap remains open until measured.
The native runner now copies an identical B.Loop DSP fixture into all three
references, avoiding benchmark-loop differences across assembly/previous/current.

Final exact-source gates pass: full native suite (45.872s root package),
focused purego pixel tests, five-package race/checkptr1, strict encoder
byte parity (34.267s), 18/18 parity report, no-assembly audit, all six build
selection modes, formatting and refreshed PGO. Production purego sources
are unchanged from the preceding full 112.845s passing suite; hosted CI
also reruns the full fallback suite. Raw logs and VALIDATION.md retain coverage.

A final six-pair native ARM assembly control measures original assembly
6.651471 ms versus current 8.420126 ms/input: 26.590% more time remains
(all six pairs slower, 19.265–27.910%). All 12 control samples preserve every
recorded non-timing/allocation field. This is the current measured ARM VP8
limit; the kernel/parent gains do not establish assembly parity.
