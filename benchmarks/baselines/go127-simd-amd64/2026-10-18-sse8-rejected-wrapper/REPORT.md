# Rejected AVX2 SSE8 wrapper integration

PR42 head c638b7c76eac857affa0780820e86c211cdd84d3 is not accepted for merge. All three CI jobs pass correctness, but native EPYC7763 performance rejects the implementation. Every live SSE8 median regresses versus parent: slice8x4 +7.06–7.24%, slice8x8 +36.85–38.44%, slice8x16 +41.32–42.56%, pointer8x8 +6.40–7.09%;35of36 are slower in every paired sample. No live case qualifies for a1.5× gain. The legacy B.N diagnostic is also50.21% slower than parent. These are separate from future repaired-head measurements.

All338×6×3 kernel samples allocate zero.32 exact source hashes, all22 public/profile semantic fields and every allocation sample match. WholeVP8 is -1.245% by median with mixed signs; VP9 controls regress0.953% and2.416%. That wholeVP8 result does not establish a kernel gain. All188 median regressions and62 all-six-slower cases remain in the full audit.

raw-evidence.tar.gz preserves the entire frozen seven-file bundle: all160 native/C raw files, complete audit, codegen scope correction and original/member manifests. bundle-sha256.json verifies every bundle member. The independent root proof is beside this report. The separate C reference job uses an independently optimized libvpx build and its own CPU; it is not evidence of AVX512 execution.
