# Retain ARM four-row FillBlock unroll: component gain over the locked winner

Native Apple M4 Max, official Go 1.27.2, GOEXPERIMENT=simd, PGO off,
GOMAXPROCS=1. Strong parent is the exact previously verified fixed-store
FillBlock winner, based on 768f2fc2. PR43 merge bba90959 has the identical
Git tree; this equality is independently checked in the package. Comparing
against bare 768 would inflate this new improvement and is not done here.
Old reference is immutable 691cd051, with actual ARM NEON DSP kernels.

Retain under the user's explicit component rule: fresh 8x8 FillBlock improves
1.5908–1.6163× across packed,37,1312 pitches over the already optimized parent;
packed public IntraDCPredict16x16 with no neighbors improves 1.6274×. Whole
VP8 is flat/mixed (−0.091934% median,4/6 faster). Candidate remains8.670979% slower than original assembly in this run. Component gains and whole results are
separate; no additive gains, universal assembly closure, or x64 claim.

## Exact single variant and semantic contract

Only the proven ARM fixed-store path changes, and only when h%4==0. One row
base is computed per four stores; ds,2*ds,3*ds select successive rows. Full
existing Enabled/width/height/bounded-stride/complete-LEN guard stays intact.
The original current loop handles other heights. Stores remain exact4/8/16B,
row0,row1,row2,row3 followed by the next group, preserving overlap semantics
and writes of the same value. Maximum row remains h−1; all offsets lie within
the already-proved window. Decline occurs before stores; invalid/short-LEN/
spare-CAP/late-capacity panic and partial writes remain original behavior.
No allocation, new unsafe contract, dependency, assembly or extra helper call.
AMD selection remains constant false; its source helper/other runtime functions
are unchanged. Source inventory proves one production file and only its
FillBlock body differ. Exact parent and candidate bytes are preserved.

Old NEON fact check: intraFill8x8NEON and intraFill16x16NEON are NOT unrolled.
They use one VST1, pointer ADD, counter SUB, and branch per row. The earlier
unroll hypothesis about old code is withdrawn; source is included. Actual new
LLVM shows one MADD per four stores, no extra calls/vector spills. Full function
entries grow166→218 because both four-row and original fallback loops exist;
counts are diagnostic and not the reason for acceptance.

## Correctness before timing

Both variants passed live ARM backend selection, ordinary/race+checkptr2 pixel
and DSP oracles, purego and no-SIMD state tests. Pixel oracle compares full
backing and panic/partial state against independently coded prior behavior,
including clipped capacities, short LEN with spare CAP, exact windows, negative,
zero and overlapping strides, h3/7/63 plus64/65 boundaries, stride2^30±boundary,
MaxInt/MinInt, misalignment, and every0..255 payload. Added nonmultiple heights
prove the original current loop remains correct. DSP scalar oracle covers
8/16 none/both neighbors, exact windows/pitches, misalignments and contexts
aliasing destination. All one-shot public fixtures pass before timing, including
pinned old NEON. Complete validation commands and logs are preserved.

Call-scope correction to the earlier frozen171 archive: n4 production also
calls intraDCGo via intra4x4_other.go:10. Actual measured DSP callers are n8/n16;
4x4 FillBlock component is measured, while actual4x4 DSP caller performance
is unmeasured. See INTEGRATION-NOTES.md. The earlier archive remains byte-for-byte
unchanged, including its source and raw evidence; this outer correction fixes
the overly narrow caller statement without editing frozen history.

## Exploratory screen and fresh confirmation

The preliminary three75ms pairs/cyclic permutations are preserved, including
all slower one-row controls. They are not acceptance data. Fresh confirmation
uses six balanced AB/BA pairs for pixel and all six three-variant permutations
for DSP,100ms per case. All648 screen+fresh component samples have0B/op and
0allocs/op. Each B.Loop invokes the real public operation; Fill changes payload
each iteration, and DSP checks scalar result, complete output/canaries and
immutable contexts. All samples, paired ranges and order commands are retained.

### Fresh pixel component

| Case | Locked parent ns | Unroll ns | Ratio of medians | All six paired ratios range | Faster pairs |
|---|---:|---:|---:|---:|---:|
| w16/h1/pitch1312 | 1.8480 | 1.8500 | 0.99892× | 0.98244–1.00433× | 3/6 |
| w16/h1/pitch16 | 1.8635 | 1.8510 | 1.00675× | 0.99676–1.02071× | 4/6 |
| w16/h1/pitch37 | 1.8605 | 1.8605 | 1.00000× | 0.98932–1.03911× | 3/6 |
| w16/h16/pitch1312 | 6.4550 | 5.0750 | 1.27192× | 1.25068–1.29380× | 6/6 |
| w16/h16/pitch16 | 6.5355 | 4.7185 | 1.38508× | 1.31757–1.40722× | 6/6 |
| w16/h16/pitch37 | 6.4595 | 5.0315 | 1.28381× | 1.24797–1.31429× | 6/6 |
| w4/h1/pitch1312 | 2.3280 | 2.2920 | 1.01571× | 0.90369–1.24636× | 2/6 |
| w4/h1/pitch37 | 2.4330 | 2.3845 | 1.02034× | 0.97685–1.15720× | 3/6 |
| w4/h1/pitch4 | 2.4555 | 2.4560 | 0.99980× | 0.97273–1.18282× | 2/6 |
| w4/h4/pitch1312 | 2.5015 | 1.9150 | 1.30627× | 1.27606–1.66684× | 6/6 |
| w4/h4/pitch37 | 2.4900 | 1.9095 | 1.30401× | 1.29243–1.39739× | 6/6 |
| w4/h4/pitch4 | 2.4855 | 1.9065 | 1.30370× | 1.29056–1.34349× | 6/6 |
| w8/h1/pitch1312 | 2.2475 | 2.3535 | 0.95496× | 0.81291–1.08933× | 2/6 |
| w8/h1/pitch37 | 2.3580 | 2.3415 | 1.00705× | 0.85439–1.09511× | 4/6 |
| w8/h1/pitch8 | 2.3545 | 2.4635 | 0.95575× | 0.90032–1.21104× | 2/6 |
| w8/h8/pitch1312 | 4.0165 | 2.4945 | 1.61014× | 1.55301–1.73598× | 6/6 |
| w8/h8/pitch37 | 4.0505 | 2.5060 | 1.61632× | 1.56341–1.67836× | 6/6 |
| w8/h8/pitch8 | 3.9890 | 2.5075 | 1.59083× | 1.50616–1.66575× | 6/6 |

All square cases improve in6/6 pairs. Only8x8 medians cross1.5× against the
strong parent;4x4 is1.304–1.306× and16x16 is1.272–1.385×. One-row8-byte
packed/large-pitch controls regress to~0.955×; they remain explicit. Other
one-row controls are near-flat/mixed. No bad sample or control is dropped.

### Fresh actual DSP caller, with original NEON comparator

| Case | Locked parent ns | Unroll ns | Old NEON ns | Parent/unroll median | Six-pair range | Old/unroll median |
|---|---:|---:|---:|---:|---:|---:|
| n16/both/pitch1312 | 14.0750 | 11.0300 | 7.9175 | 1.27607× | 1.25227–1.35351× | 0.71782× |
| n16/both/pitch16 | 13.6250 | 10.7900 | 7.8190 | 1.26274× | 1.24593–1.32464× | 0.72465× |
| n16/both/pitch37 | 13.8400 | 10.9050 | 7.9985 | 1.26914× | 1.21715–1.32811× | 0.73347× |
| n16/none/pitch1312 | 8.4805 | 6.2780 | 7.0075 | 1.35083× | 1.31442–1.36955× | 1.11620× |
| n16/none/pitch16 | 8.5090 | 5.2285 | 7.0495 | 1.62743× | 1.56785–1.65886× | 1.34828× |
| n16/none/pitch37 | 8.5250 | 6.0905 | 7.0180 | 1.39972× | 1.37731–1.41634× | 1.15229× |
| n8/both/pitch1312 | 9.9155 | 8.1865 | 6.1710 | 1.21120× | 1.12173–1.29695× | 0.75380× |
| n8/both/pitch37 | 9.8420 | 8.2915 | 6.1755 | 1.18700× | 1.10491–1.27922× | 0.74480× |
| n8/both/pitch8 | 9.6515 | 8.2505 | 6.1280 | 1.16981× | 1.15499–1.27701× | 0.74274× |
| n8/none/pitch1312 | 6.0625 | 4.2785 | 4.1745 | 1.41697× | 1.29037–1.58875× | 0.97569× |
| n8/none/pitch37 | 5.9865 | 4.3490 | 4.1915 | 1.37652× | 1.27885–1.53125× | 0.96378× |
| n8/none/pitch8 | 6.0480 | 4.4340 | 4.1930 | 1.36401× | 1.32323–1.45175× | 0.94565× |

All12 DSP cases improve in6/6 pairs against the locked parent. Packed16x16
with no neighbors qualifies at1.6274×; other caller gains are below1.5×.
All16x16 no-neighbor median cases beat old NEON (old/unroll1.116–1.348×).
Fresh8x8 no-neighbor remains2.5–5.8% slower than old NEON; preliminary padded
screen wins over old do not survive fresh confirmation. Both-neighbor cases
remain slower than old. Every screen and fresh result is preserved separately.

## Fresh whole-codec quantification

Six balanced three-way permutations; VP8,1280x720,300 inputs,30fps,2500kbps,
realtime CPU-used8,one thread,encode-only,auto-libvpx=false. Strong locked
parent and candidate CLIs were freshly built with the same SDK/PGO-off flags;
old691 was reused only after verifying its previously built SDK/PGO-off hash.
Builds use the root project wrapper, unique IDs, GOMAX2; all builds and heavy
parsing/hashing stop during announced native windows. Root and other agents
confirm exclusive START/END; no competing local benchmark/build/profile.

All22 semantic fields and separate zero allocations match across18 JSONs:
298 encoded,2 dropped,3120743 output bytes,133700 keyframe bytes, identical
quantizer metadata/histogram. Positive CPU8 is adaptive, so work parity is
checked rather than assumed. Quality is skipped; these fields do not prove
bitstream identity. Timer is encode-only and excludes fixture setup.

- parent: median 7,305,768.0 ns/input; six values [7209094, 7304173, 7307363, 7301504, 7355257, 7526557].

- candidate: median 7,299,051.5 ns/input; six values [7253179, 7291333, 7301061, 7297042, 7370514, 7339571].

- old691: median 6,716,652.0 ns/input; six values [6600274, 6814188, 6700665, 6667184, 6732639, 7097822].

Candidate vs parent: -0.091934% by ratio of medians, 4/6 faster; all paired percent changes [0.6115192838378913, -0.17578992173378216, -0.08624178106383384, -0.061110697193345054, 0.20742986954773723, -2.4843497498258538].

Candidate vs old691: +8.670979% by ratio of medians, 0/6 faster; all paired percent changes [9.892089328412723, 7.002228291910928, 8.960244990609144, 9.4471369021764, 9.474368074688089, 3.4059603072604583].

The whole change is flat/mixed, with no claimed codec gain. Candidate
remains8.670979% slower than original assembly by this run's ratio of medians (slower in all
six rounds; paired gap3.40596–9.89209%). This gap is not comparable by subtraction
with another run or runner. Retain the verified component wins despite flat
whole behavior, as explicitly required by the user's acceptance policy.

## Delivery

Overlay contains exact final source and two strengthened fixtures; root owns
integration, whitespace-only test formatting, fresh PGO and combined full gates.
No root, README or PR edits were made by this agent. Exact source, scoped patch,
LLVM, all ordinary/race/fallback logs, preliminary and fresh raw measurements,
18 whole JSONs/stderr, commands, all paired ratios, hardware/compiler/build/
source/binary hashes and independent scope verification are bundled. No full
snapshots, binaries, dependency copies or source changes to the frozen171 archive.
SHA256SUMS covers every other member. No more timing or experiments are planned.
