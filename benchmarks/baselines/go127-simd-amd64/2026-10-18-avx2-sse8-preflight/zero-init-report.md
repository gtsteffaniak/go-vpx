# Private common SSE8 VEX-zero assessment

No performance claim is established. No rootedits, PR or CIhead changes were
made. This source-only candidate removes a measured-codegen hazard; it has
not been timed natively or under emulation.

Parent commit: ab3abd6b335242d5732fa1956a63b7591eaf781b. Official Go1.27.2,
SIMD experiment, GOAMD64=v1, -pgo=off, GOFLAGS=-p=2. All builds run from the
root scripts/go.sh -C private snapshot, project-env, distinct agent scope.
Optimized Linux/amd64 and ARM64 codegen is decoded by Apple LLVM21.
Rosetta correctness uses ROSETTA_ADVERTISE_AVX=1 and GOVPX_REQUIRE_SIMD=1;
actual cpu.avx/avx2-off fallback is repeated separately. No native timing
or emulated speed inference was performed.

## Confirmed baseline context

Common SSE8 has four STATIC legacy MOVUPS xmm15-to-zero-accumulator copies:
two execute on the bounded fast path, the other two on generic fallback.
Its vector arithmetic is VEX128 and the complete function has no YMM write
or ClearAVXUpperBits. These facts alone do not show that an AVX-to-SSE
transition caused the reported native legacy benchmark latency.

The actual legacy BenchmarkSSE8x8 B.N caller initializes two stack buffers
using legacy MOVUPS stores once, then runs scalar integer loop/control and
calls pixel.SSE8 each iteration. It has no YMM-producing instruction. The
live B.Loop public fixture also reaches common SSE8 through API forwarding
or an indirect function call; its own saved caller body does not establish a
dirty-YMM source. Wider harness/prior-callee state might be inherited, and
legacy-to-VEX alternation could preserve that state between calls, but the
151ns versus11.9ns cause remains unproven without a controlled native test.
The zero-input legacy benchmark differs in API/data/setup from live fixtures;
its timing alone must not be assigned to one instruction sequence.

## Private change and generated code

AMD defines TWO DISTINCT private zero arrays of four int32 each and loads
them with LoadInt32x4Array. Distinct arrays prevent the compiler from reusing
one loaded value with a legacy register copy. A compile-time true AMD flag
replaces each accumulator initialization; ARM's false flag removes the new
block entirely. This follows existing helpers_amd64.go packZero16 precedent.

All four static legacy copies become VMOVDQU loads from the two zero arrays.
No legacy MOVUPS remain anywhere in the optimized common SSE8 body. No YMM
write, upper-bit fence, hot function call or vector stack spill is added.
Enabled's CPU branch precedes every AVX load, so unavailable AVX2 retains the
original scalar fallback. No memory geometry, aliasing, row arithmetic,
external loads/stores, numerical semantics or panic contract changes.

All existing ARM-specific source files are byte-identical. ARM SSE8 emitted
machine code is also byte-identical, including addresses:228instructions,
SHA2560bb310fef008c3578c646c3ca9a89b438b04e3f92f23c46ce17fe9d4af449b21.
arm-audit.json and the two complete ARM listings retain the proof.

## Correctness and applicability

Existing TestSquaredError and all SSE tests pass under live Rosetta AVX2
with checkptr1, including byte pairs, exact/short buffers, zero/overlapping
strides, immutable aliases, extreme differences and native fallback heights.
Actual AVX/AVX2-off test logs show false feature flags and passing scalar
results. ARM native correctness, purego, and experiment-off suites pass.
The AMD Enabled guard is preserved before AVX intrinsics.

This candidate is useful to assess alongside the wide h>=8 slice and trusted
pointer prototypes: their direct paths may bypass common SSE8, but h4 controls
and legitimate generic fallback still retain these legacy initializations.
An all-VEX common baseline removes that mixed-encoding variable without
assuming a performance win. Native comparisons remain required, ideally with
identical data/callers and controlled clean/dirty upper-register context.

Only sources, exact driver, correctness logs, source patch, disassembly and
audits are packaged. Binaries/full snapshots are excluded. Recreate private
archives from parent-commit.txt and apply source-change.patch, then run the
driver from the root checkout/cache scope. SHA256SUMS covers every file except
itself. README files were not changed.
