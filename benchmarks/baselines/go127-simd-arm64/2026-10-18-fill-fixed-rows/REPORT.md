# ARM exact-width FillBlock: retain scoped component improvement

Pinned parent768f2fc2aba7bbd82981494473cbadc026a95570 (PR43); original real
NEON reference691cd0512c4829d9ee40d483205dd223f9eb028d. Native Apple M4 Max,
official Go1.27.2, GOEXPERIMENT=simd, PGOoff, GOMAXPROCS1 for all timings.
All binaries freshly built with the same SDK/root scripts/go.sh/project-env
cache scope and unique agent IDs. Build concurrency maxGOMAX2, no builds/Go/
heavy hashing during native timing windows. START/END coordinated with root,
ci_validation and merge_checkpoint. Every result/outlier/slower control retained.

Retain under the user's component rule: fresh six-pair public FillBlock square
medians improve1.544–2.090×; actual public IntraDCPredict8x8 with no neighbors
improves1.623–1.814×. Whole VP8 remains small/mixed(−0.363%median,5/6faster),
and original NEON remains faster on every measured DSP case and whole sample.
No additive, whole-codec speedup, cross-hardware or x64 performance claim.

## Frozen production and exact contract

Only FillBlock gains an ARM constant-selected fixed-width path. With SIMD
Enabled,2<=h<=64,nonnegative stride<2^30,w4/8/16 and full (h−1)*stride+w<=LEN,
each row uses exactly one4B/8B integer or16B SIMD store. No overstore, no new
allocation/alias restriction, positive endpoints fit int64, complete window
proven before stores. ds0/overlapping rows retain sequential writes of the same
value. Invalid/short-LEN/negative/extreme strides or other geometry decline
before stores and execute original row slicing/StorePart fallback unchanged.
Capacity remains usable by original fallback; malformed late-row capacity
retains original partial writes and panic. h1 stays in the original body but
added dispatch changes codegen/latency, with slower controls disclosed below.
AMD constantfalse removes specialization; purego/noexperiment remain scalar.
Production byte lock verified before and after all measurements; only test-only
fixtures/drivers changed during this confirmation. Root/PR/README untouched.

Private initial separate-helper and later inline27-case screens are exploratory,
not fresh acceptance. All initial/inline logs, build/correctness outputs, LLVM,
commands, summaries and h1 regressions remain under initial-and-inline.
The exact separately timed initial helper source was not independently archived
before root replaced it; its raw LLVM/commands remain. Final fresh-measured
inline source is fully preserved and hashed. No source identity claim equates
that initial helper with the final inline implementation.

## Independent correctness before fresh timing

The original-state oracle implements scalar writes directly and, for native
SIMD, reproduces prior per-row slice plus StorePart's min(width,16) contract.
Full backing buffers and panic booleans are compared, including clipped3-index
capacity, one-byte-short last rows, shortLEN with spareCAP, zeroLEN/height,
negative/zero/overlapping strides, h64/65, stride2^30−1/2^30/MaxInt/MinInt,
misalignments and every0..255 payload on exact windows. This closes the first
fixture's gap where all positive shortLEN cases had ample spare capacity.

Parent/candidate ordinary and race/checkptr2 pixel tests PASS; purego and
GOEXPERIMENT-empty state oracles PASS independently. Live backend logged native
ARM SIMD=true. DSP8/16 oracle computes DC via independent scalar sums, covers
none/both neighbor availability, exact capacities, packed/37/1312pitches,
misalignments and above/left aliasing destination. Parent/candidate/old691
ordinary DSP oracles PASS; currentparent/candidate DSP race/checkptr2 PASS.
All fresh benchmark fixtures pass1x output/canary/context oracles first.
Old invalid unsafe/short-capacity calls are not benchmarked as supported inputs.

Final LLVM is captured from measured binaries. Current FillBlock exact-store
bodies and old actual intraFill8x8NEON/intraFill16x16NEON calls are preserved,
not inferred merely from build tags. Component API scope: current production
FillBlock caller is VP8 intraDCGo(n8/n16);4x4 is standalone public component
coverage, not a claimed current codec caller.

## Fresh public pixel confirmation

Six balanced AB/BA pairs,100ms,18cases:4x4/8x8/16x16 +each width's one-row
control, packed/37/1312pitches. B.Loop varies payload each iteration; complete
backing/canary output checked after loop. All216samples0B/op,0allocs/op.
Ratios below are ratio of six-sample medians; paired ranges retain all rounds.

| Case | Parent ns | Candidate ns | Median ratio | Six paired ratio range | Faster pairs |
|---|---:|---:|---:|---:|---:|
| w16/h1/pitch1312 | 1.8575 | 1.8465 | 1.00596× | 0.96844–1.03471× | 3/6 |
| w16/h1/pitch16 | 1.8595 | 1.8610 | 0.99919× | 0.97885–1.01297× | 3/6 |
| w16/h1/pitch37 | 1.8610 | 1.8535 | 1.00405× | 0.98878–1.02642× | 3/6 |
| w16/h16/pitch1312 | 9.8410 | 6.3720 | 1.54441× | 1.47483–1.59297× | 6/6 |
| w16/h16/pitch16 | 9.9000 | 6.3170 | 1.56720× | 1.50581–1.59828× | 6/6 |
| w16/h16/pitch37 | 9.8800 | 6.3065 | 1.56664× | 1.53453–1.60413× | 6/6 |
| w4/h1/pitch1312 | 2.2185 | 2.3035 | 0.96310× | 0.85374–0.97982× | 0/6 |
| w4/h1/pitch37 | 2.1695 | 2.2100 | 0.98167× | 0.87741–1.09848× | 3/6 |
| w4/h1/pitch4 | 2.2115 | 2.4080 | 0.91840× | 0.87418–1.08787× | 2/6 |
| w4/h4/pitch1312 | 4.3210 | 2.4630 | 1.75436× | 1.70348–2.26778× | 6/6 |
| w4/h4/pitch37 | 4.8165 | 2.4525 | 1.96391× | 1.72213–2.27313× | 6/6 |
| w4/h4/pitch4 | 5.0035 | 2.4670 | 2.02817× | 1.72561–2.27469× | 6/6 |
| w8/h1/pitch1312 | 2.0525 | 2.2035 | 0.93147× | 0.86481–0.99727× | 0/6 |
| w8/h1/pitch37 | 1.9960 | 2.1350 | 0.93489× | 0.86072–1.07980× | 1/6 |
| w8/h1/pitch8 | 2.2415 | 2.4345 | 0.92072× | 0.81176–1.09754× | 1/6 |
| w8/h8/pitch1312 | 8.2335 | 3.9950 | 2.06095× | 1.84051–2.45101× | 6/6 |
| w8/h8/pitch37 | 8.2170 | 3.9530 | 2.07867× | 1.93729–2.43959× | 6/6 |
| w8/h8/pitch8 | 8.4575 | 4.0470 | 2.08982× | 1.99151–2.26106× | 6/6 |

All square cases are faster in6/6 pairs; not every individual pair reaches1.5×.
In particular16x16/pitch1312 round4 is1.47483×, though median1.54441× qualifies.
One-row4B/8B regressions remain: median ratios0.918–0.982× and0.921–0.935×;
16B one-row is near-flat0.999–1.006×. No controls are excluded from reports.

## Fresh actual DSP prediction confirmation versus old NEON

Six balanced permutations of all three variants on the same host,100ms,
12cases. Timed functions directly call public IntraDCPredict8x8/16x16, none or
both neighbors; oracle checks entire output, canaries, immutable contexts.
All216samples0B/op,0allocs/op. Every case improves against parent in6/6pairs,
but old NEON remains faster in every case. old/candidate<1 means candidate
is slower than old; this explicitly records the remaining gap.

| Case | Parent ns | Candidate ns | Old NEON ns | Parent/candidate median | Six paired ratio range | Old/candidate median |
|---|---:|---:|---:|---:|---:|---:|
| n16/both/pitch1312 | 16.9200 | 13.5300 | 7.7430 | 1.25055× | 1.22384–1.31907× | 0.57228× |
| n16/both/pitch16 | 16.7900 | 13.5400 | 7.7880 | 1.24003× | 1.14696–1.28146× | 0.57518× |
| n16/both/pitch37 | 16.9800 | 13.3500 | 7.7845 | 1.27191× | 1.22964–1.29680× | 0.58311× |
| n16/none/pitch1312 | 12.0500 | 8.4485 | 6.8395 | 1.42629× | 1.36448–1.54588× | 0.80955× |
| n16/none/pitch16 | 12.0150 | 8.4175 | 6.8465 | 1.42738× | 1.40260–1.44016× | 0.81337× |
| n16/none/pitch37 | 11.9500 | 8.2555 | 6.8645 | 1.44752× | 1.38242–1.48181× | 0.83151× |
| n8/both/pitch1312 | 14.2050 | 9.8695 | 6.0745 | 1.43928× | 1.33280–1.55243× | 0.61548× |
| n8/both/pitch37 | 14.8500 | 9.9160 | 6.0610 | 1.49758× | 1.33239–1.58172× | 0.61123× |
| n8/both/pitch8 | 14.0650 | 10.0610 | 6.0540 | 1.39797× | 1.25543–1.48816× | 0.60173× |
| n8/none/pitch1312 | 10.1650 | 5.8130 | 4.1330 | 1.74867× | 1.65782–1.98451× | 0.71099× |
| n8/none/pitch37 | 9.6850 | 5.9670 | 4.2885 | 1.62309× | 1.49775–1.87490× | 0.71870× |
| n8/none/pitch8 | 11.0400 | 6.0850 | 4.1250 | 1.81430× | 1.64238–2.01071× | 0.67790× |

8x8/no-neighbor medians1.623–1.814×; its pitch37 round4 pair is below1.5×
(the exact minimum1.49775× is retained), so no all-pairs1.5 claim. Other DSP
availability/size cases improve1.240–1.498×, reported separately from qualifying
pixel/8x8-none components.

## Fresh whole VP8 quantification

Six balanced three-way permutations,300inputs each,1280x720,30fps,2500kbps,
realtime CPU-used8,one thread,encode-only,auto-libvpx=false. Positive CPU8 is
adaptive, so all recorded workload fields are checked rather than assumed.
All22 semantic fields and separate zero allocations match across all18JSONs:
298encoded/2dropped,3120743outputbytes,133700keyframebytes and identical
quantizer metadata/histogram. Quality is skipped; recorded equality does not
establish bitstream identity. CPU timer excludes fixture setup/alloc sampling.

- parent: median7,327,148.0ns/input; six samples[7364637, 7281993, 7260401, 7332436, 7327884, 7326412].

- candidate: median7,300,580.5ns/input; six samples[7181989, 7246922, 7296509, 7310648, 7327618, 7304652].

- old691: median6,739,792.0ns/input; six samples[6971599, 6949542, 6787461, 6664437, 6686159, 6692123].

Candidate vsparent: -0.362590% by ratio of medians;5/6faster; paired range-2.480068..+0.497328%.

Candidate vsold691: +8.320561% by ratio of medians;0/6faster; paired range+3.017816..+9.696408%.

Whole result is small/mixed, not a strong codec improvement claim. Candidate
remains8.32056% slower than old691 in this run,0/6faster, with3.0178–9.6964%
paired gap; do not subtract this from another experiment or runner's figures.
Keep the verified component win under the user's explicit retention rule.

## Package and integration scope

Exact final source/overlay, independent fixtures, full raw initial/inline/fresh
logs, failed/initial build evidence if present, all paired samples/ratios,
18wholeJSONs/stderr, commands, LLVM, production/source/binary hashes retained.
No binaries/full snapshots/dependencies are bundled. SHA256SUMS covers all other
members; verification.json checks inventories, oracle/PASS markers and hashes.
The overlay is private/reviewable; root owns any integration, PGO refresh and full CI gates. Root confirmed
the per-file citation gate does not cover these pixel const files; preserve
the exact measured source bytes rather than adding unnecessary comments.
No README edit.

All six candidate percent changes, retained separately from medians:

- vsparent: [-2.4800679245969626, -0.48161265741397674, 0.49732790241199254, -0.2971454507069704, -0.0036299701250741556, -0.29700759389452713].

- vsold691: [3.017815568566129, 4.27913091251193, 7.4998294649501585, 9.696407963643438, 9.593834068259511, 9.152984785246776].
