# AMD FP class-index candidate: validated transplant

Root source is unchanged at 63c2f188. Candidate source is in
`/private/tmp/govpx-vp9-fpclass-reference`.

The three production files and two test files in
`/private/tmp/govpx-fpclass-transplant.patch` are the complete transplant.
The patch is based on root 63c2f188; do not copy the candidate's other files
or its profiling/experimental commands. SHA256 manifest:
`/private/tmp/govpx-fpclass-transplant-sha256.json`.

## Why it is exact

QuantizeFPBlock guards nonnegative round and multiplier, positive dequant.
The quantizer clips unsigned abs(coeff)+round at 32767, then uses a signed
high multiply with a nonnegative multiplier. The resulting magnitude is
nonnegative; masking preserves this. The existing min(11) therefore proves
the packet passed to fpClassIndices is in [0,11]. AMD can extract each low
byte with the existing PSHUFB mask without the general [0,255] pixel clamp.
ARM uses the identical direct saturated-narrowing expression as packBytes.
No public bounds, aliasing, tail, optional output or feature guard changed.

## Compiled evidence

PGOoff Linux/AMD64 v1, Go1.27.2 simd: the actual encoder-reachable quantizer
removes five instructions from the class branch in every eight-coefficient
iteration: VPMAXSW zero, MOV $255, VMOVQ, VPBROADCASTW and VPMINSW255.
The required min(11), low-byte PSHUFB and class-table PSHUFB remain. There
is no helper call in the loop and no YMM/ZMM or added ISA requirement.
The only non-runtime call remains the existing tail handler.

The final ARM quantizer is byte-identical: 800 bytes, SHA256
2f560bffa7d93d0d1e19a8c1701399313d75cea17909593b68c57b71ef9b751e.
See `/private/tmp/govpx-fpclass-codegen.json` and the parent/candidate
`.asm` files; exact AMD diff is `/private/tmp/govpx-fpclass-amd.diff`.

## Validation

All pixel, VP9 encoder and VP9 DSP package tests pass on ARM and on Rosetta
with SIMD required. This includes the independent 600-case scalar quantizer
oracle (q/dq/EOB/classes, random/full signed extremes, dirty buffers) and
new 144 mixed-lane index packets. Experiment-off and purego pixel+encoder
suites pass. The final ARM direct-helper expression was followed by a fresh
three-package native test run.

Logs: `/private/tmp/govpx-fpclass-{arm-final,rosetta,off,purego}.log`.
All build/test calls used scripts/go.sh, project scope cache and unique IDs.

## Native measurement fixture

The exact same public B.Loop fixture cross-compiles for original assembly
691cd051, parent63c2f188 and candidate. Nine size/distribution cases cover
4x4/8x8/16x16 and zero/sparse/dense coefficients, plus one no-classes control.
Observable q/dq/EOB/classes and ReportAllocs keep the work live. Assembly
reports produced=false, so the fixture supplies its scalar token class walk:
the measured contract is full quantization plus every raster output class.
This is a defined equivalent-output composite, not a claim that every
production caller needs a full class walk after zero/EOB clipping. Public
encoder timing remains the acceptance criterion.

To run in scripts/benchmark-amd64.sh, copy the new encoder fixture to both
old-reference directories beside dct32_bench_test.go, and add the anchored
`BenchmarkVP9(FusedFPClassesNative|FPNoClassesNative)$` alternatives to the
kernel regex. This adds ten cases (roughly 30-50 seconds across18 variants,
plus overhead) to the existing native matrix.

Required next evidence: repeated native AMD parent/candidate320p+720p encoder
comparisons, equivalent composite micro cases, semantic/allocation parity,
and original ASM composite baseline. No Rosetta timings were taken and no
speedup is claimed yet. The +0.479% control regression is still open until
native measurements establish the selected change.
