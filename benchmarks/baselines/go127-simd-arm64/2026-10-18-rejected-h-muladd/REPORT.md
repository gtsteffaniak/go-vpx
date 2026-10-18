# ARM horizontal-only MulAdd macro correction: rejected after screen

Parent `768f2fc2aba7bbd82981494473cbadc026a95570` (PR43). Stop/archive:
near-noise/mixed public changes, no meaningful representative gain. No root,
README or PR changes; no fresh-six or whole-codec measurements.

The sole candidate changes only ARM fused-horizontal macro correction calls
to a helper computing exactly `(f*c+63)>>7` with Int16x8.MulAdd, coefficients
9/18/27. Full source inventory and source-scope-proof.json establish unchanged
ARM vertical source, AMD source, generic helper, dispatch/gates/window checks,
loads/stores and all alias semantics. The common helper remains unchanged.
A test-only numeric fixture is byte-for-byte equivalent except helper name.

Pinned SDK MulAddInt16x8 lowers to destructive VMLA8H with accumulator in
arg0. Before timing both variants passed all768 signed-byte/coefficient
combinations, live ARM backend detection, five-package short suites and
four-package loop-filter race/checkptr2 oracles. Existing7200 fused-filter
cases, zero-correction/later edges, aliases, minimal windows and sentinels pass.
All commands/build/test outputs are retained. No assertions/thresholds weakened.

Optimized LLVM horizontal body changes6MUL.8H+6ADD.8H to6MLA.8H+5additional
MOV.16B bias copies. It has no calls/vector spills/stack-vector loads in either
variant. Disassembly entries364→362 include one removed trailing UDF alignment
word: the EXECUTABLE body saves only ONE instruction(363→362). No claim of two
hot instructions removed. Earlier conversational speculation about another
eliminated copy is withdrawn; opcode-count-diff.json records the exact delta.
The arithmetic dependency becomes bias-copy→MLA→shift, instead of MUL→ADD→shift;
this alone does not establish M4 latency/throughput benefit. Vertical698 and
generic416 entries/opcode mixes remain unchanged. Complete LLVM is retained.

Native Apple M4 Max, official Go1.27.2 SIMD, PGOoff,GOMAXPROCS1. After exclusive
host release from root, three75ms alternating AB/BA/AB rounds run the existing
actual public10-picker B.Loop fixtures; full output buffer must equal four
separate public edge operations. All60samples0B/op/0allocs; all30paired ratios
and slower controls retained. No artificial correction-only benchmark used.

| Public case | Parent ns | Candidate ns | Parent/candidate | Candidate time change |
|---|---:|---:|---:|---:|
| pitch1312/flat | 45.50 | 45.91 | 0.99107× | +0.901% |
| pitch1312/mixed | 86.98 | 85.98 | 1.01163× | -1.150% |
| pitch1312/noop | 40.42 | 40.90 | 0.98826× | +1.188% |
| pitch1312/ramp | 45.53 | 45.92 | 0.99151× | +0.857% |
| pitch1312/smooth | 87.00 | 86.33 | 1.00776× | -0.770% |
| pitch32/flat | 45.87 | 45.97 | 0.99782× | +0.218% |
| pitch32/mixed | 86.59 | 85.66 | 1.01086× | -1.074% |
| pitch32/noop | 40.48 | 40.51 | 0.99926× | +0.074% |
| pitch32/ramp | 45.86 | 45.39 | 1.01035× | -1.025% |
| pitch32/smooth | 86.13 | 85.50 | 1.00737× | -0.731% |

Median ratios0.9883–1.0116×; candidate−1.15% to+1.19%time. Mixed cases improve
about1.1% but flat/ramp/noop controls regress or are mixed; this is insufficient
to establish useful gain. No kernel1.5×/whole-codec/x64 performance claim.
All Go/build/test/timing jobs stopped; host explicitly released after screen.

Exact sources, scoped patch/proof, public/numeric fixtures, full optimized LLVM,
SDK rule, ordinary/race/build logs, every raw timing, commands, paired summary,
environment/source/binary hashes retained. Full snapshots/binaries excluded.
Drivers retain original private paths. SHA256SUMS covers every other member.
The prior rejected UQSUB/OR archive remains unchanged.
