# CI49 Stats architecture repair verification

All requested correctness checks pass. No timing or performance claim is made.

ARM native ordinary and race/checkptr=2 pixel, VP8 DSP and VP9 DSP Stats/extreme/variance tests pass. Darwin AMD64 v1 checkptr=2 tests pass through Rosetta with required SIMD selected; TestStats16NativeSelected passes. Actual AVX and AVX2 disabled runs pass all public tests, with the native-selection test correctly skipping. Full commands, environment and raw outputs are preserved.

The seven selected AMD functions exactly restore the prior 31123 extraction instruction sequences after removing only PCs and RIP relocation displacements. ARM Rows16, Quads4, dispatcher, narrow16, public Stats and tall Stats are exactly unchanged from 41596. The ARM Pairs8 vector loop before reduction is unchanged, including both ADDV reductions; its full body is **not identical**: 169 to 172 instructions, with sign extension at reduction and odd-row merge caused by the restored int64 helper shape. This is a codegen limitation to disclose, not evidence of a measured regression. Sign extension and 64-bit odd-row addition preserve the low 32-bit public result, including modular sums, as checked by tall/extreme tests.

The comparison uses exact baseline source overlays in the same ROOT link context, not different full binaries with unrelated link layouts. All three ROOT Stats source hashes remain equal to the tested receipt. No ROOT files, PGO assets, or old ARM50 snapshots were changed. Binaries are omitted from publication but their private paths, sizes and SHA-256 hashes are recorded.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
