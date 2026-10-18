# Provisional horizontal ARM loop unroll

Only the three fixed horizontal internal edges are unrolled in a private generated-file prototype. ROOT does not include this change; a generator update and fresh kernel/codec confirmation are still pending. This does not satisfy the 1.5x retention threshold.

Apple M4 Max, official Go1.27.2 archsimd, CGO1, PGOoff, GOMAXPROCS1: six balanced parent/current rounds preserve all360timed samples and60excluded oracle checks, all0B/0alloc. Horizontal10cases improve1.012306–1.034412x medians (55/60pairs faster,5slower). Actual decoder-picker10cases improve1.000118–1.047390x (42/60faster,18slower). Unchanged vertical controls range0.989628–1.006513x, including slower samples; no causal claim for their movement. Full three-package and scoped race/checkptr tests pass.

Actual CLI LLVM horizontal static code grows364→557instructions, including47→73NOPs (317→484non-NOP). Both bodies are leaf functions with no vector stack accesses. These are static code-size diagnostics, not dynamic executed instruction counts or a codec speed claim. Prior source snapshots remain immutable. Binaries are omitted; exact hashes remain in binding/codegen receipts. Every raw slower sample is retained.


Publication: `raw-evidence.tar.gz` preserves all 69 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
