# Private AMD SSE8 direct return-fence assessment

Base: `c638b7c76eac857affa0780820e86c211cdd84d3` (PR42 remains frozen). Go 1.27.2, GOEXPERIMENT=simd, GOAMD64=v1, PGO off. This package changes no root files. It contains no native or emulated performance result or acceptance claim. One-shot fixture timing columns are correctness-run output only.

Exactly two of 1,876 checked Go/module files differ from the base. Each change removes a deferred ClearAVXUpperBits and extracts the reduction into a scalar result, executes ClearAVXUpperBits directly, then returns. Arithmetic, row order, geometry, dispatch and fallbacks are unchanged.

## Contract and cleanup review

The slice entry checks Enabled and bounds `h` to 1..64 and nonnegative strides below 2^30. Both original last-byte checks execute before the exact-window helper. Monotone row offsets are below 2^36, and each load reads exactly eight bytes. Read-only overlapping/zero-stride rows and unaligned addresses remain supported. The pointer entry retains the existing trusted complete-8x8-window contract and declines unsupported CPU, nil pointers or unusual strides through the exact original unsafe.Slice/SSE8 expression before any YMM arithmetic. Integer terminal offsets are not dereferenced as pointers.

Eligible valid fast paths contain no Go bounds checks or other language-level panic sites after wide-vector work starts. Direct SDK lowering is a memory-effect VZEROUPPER. Both optimized targets extract the result into EAX before the fence; VZEROUPPER does not alter EAX or other general-purpose registers, and no vector value is live afterward. There is one fence on each normal SIMD return. Fallback exceptions occur before any YMM instruction in the pointer path.

Removing defer also removes its exceptional-unwind cleanup. A nonnil but undersized, dangling or otherwise invalid pointer violates the preexisting trusted pointer contract. A memory fault or checkptr failure after entering the unsafe SIMD path is not promised a fence during unwinding by this variant. The valid-window contract is not widened. Nil/invalid-stride fallback diagnostics remain tested and unchanged. The internal slice helper similarly requires its guarded caller; direct calls with invalid geometry are outside its contract.

## Correctness and code generation

The root-supplied selected Rosetta pixel log passes exact windows, exhaustive byte pairs, aliases, short bounds/fallback, invalid-geometry subprocess diagnostics and squared-error tests under checkptr=1 (confirmed in binary build metadata). Independent CPU AVX/AVX2-off execution passes pointer/fallback and public squared-error tests with AVX2=false and Enabled=false; SIMD-only tests skip as intended.

A fresh public DSP binary uses checkptr=2. All 36 identical scalar/canary fixtures pass both with AVX2 selected and with actual CPU AVX/AVX2 disabled, one operation per fixture, with zero B/op and allocs/op readings. This is correctness and dispatch evidence, not a throughput measurement.

Optimized Linux original/candidate total decoded instruction counts (including cold paths/alignment padding) are SSE8Ptr 180 -> 149 and sse8RowsWide 125 -> 90. Fresh optimized Darwin candidate matches 149/90. Each candidate body has one inline VZEROUPPER, no ClearAVXUpperBits call, no deferreturn call and no legacy vector instructions. The rows helper has no calls; pointer calls are confined to scalar fallback, panic and stack-growth paths before SIMD. Parent code has actual ClearAVXUpperBits and deferreturn calls plus XMM15 stack-zero bookkeeping. Counts diagnose code generation only; native speed remains unmeasured.

## Contents and reproduction

`parent/` and `candidate/` pin the two changed sources; `context/` pins the public call chain, test fixtures and wrapper. `candidate.patch` applies to the exact base. `audit.json`, SDK snippet, original and candidate LLVM logs, build metadata and complete raw correctness logs preserve the proof. Binary hashes are recorded; large binaries are not copied. `reproduce.sh` uses the root project wrapper and creates a private source snapshot from the pinned base. Run only in a coordinated correctness window. `SHA256SUMS` covers every package member except itself.
