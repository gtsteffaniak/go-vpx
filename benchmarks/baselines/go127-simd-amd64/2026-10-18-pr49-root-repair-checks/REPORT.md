# PR49 final x64 stats and CI repair checks

Restore the original AMD scalar extraction using a compile-time architecture constant and width alias. ARM keeps its 32-bit vector reduction; the odd-row ADD operands are commuted with identical dataflow. Separate native/translated correctness and LLVM evidence is in the architecture repair archive. No x64 performance is inferred from Rosetta.

Both full ARM test runs and PGO refresh/precommit checks pass. The first int64-only repair is retained as an intermediate trial; final aliases avoid its three extra ARM tail instructions. Source-binding records every tracked Go file and proves staged bytes match tested working files. Script syntax and the complete failed native log parser replay pass. Nonfatal module stat-cache write warnings remain in successful build logs.

The full failed41596 native artifact remains a historical checkpoint: its timing completed, its final parser failed, and AVX512 was unavailable. The final repaired head needs fresh CI and native performance audit before merge. README files are unchanged.


Publication: `raw-evidence.tar.gz` preserves all 12 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
