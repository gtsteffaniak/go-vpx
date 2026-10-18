# ARM zero-packet initial slice-gate experiment

Initial six-pair Apple M4 Max / Go1.27.2 SIMD / CGO1 / PGOoff / one-thread evidence, before the separate pointer refinement. The slice-based packet gate preserves numeric and alias correctness but regresses all 28 dense-case medians: parent/current ratios 0.876–0.971×. The optional qnil n256 zero fixture reaches approximately 1.533×; this narrow optional-output component result is not an overall codec win. All losses and original individual measurements remain visible in root-arm-summary.json and the independent audit.

Preserved raw directories contain component rounds, excluded preflights and three whole-codec workloads. Binding/summary/driver/build/test/race logs and original optimized quantizer LLVM evidence are included unchanged. Original failed checks and audit attempts are retained. Executables and full source snapshots are deliberately omitted; existing source and binary hash receipts remain. This publication adds no new measurements, re-audit or acceptance claim. The later pointer-based two-round screen is a separate window and must not be pooled with this initial implementation. No overall speedup or assembly-gap closure is claimed.


Publication: `raw-evidence.tar.gz` preserves all 194 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
