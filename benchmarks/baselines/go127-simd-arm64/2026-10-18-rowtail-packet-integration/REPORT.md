Private ARM row-tail experiment against PR52 5c5cf397, Apple M4 Max, Go1.27.2/SIMD/CGO1/PGOoff/GOMAX1. Not integrated into production at freeze.

The caller already proves disjoint rows and a complete20×16 tile. Replace only the generated value helper's16checked scalar tail writes with unsafe byte writes at identical offsets, preserving row order and all pixel writes. The pointer-array helper is unchanged. Exactly2production Go files differ; the new shared scalar/full-backing fixture adds960cases and12native store workloads. Ordinary pixel/DSP/decoder suites and scopedrace/checkptr PASS. Independent safety review checks every production caller and overflow proof.

LLVM actual timed helper:217→177instructions throughfirstreturn,280→180total includingcold/padding.16panicBounds sites become0; the helper becomes frameless. Cold morestack vector spills were NOT hot-path spills. The first parser incorrectly searched panicIndex rather than ARM panicBounds; preserved and corrected before claims. Actual raw instruction captures and executable/build-info hashes are included.

Native six balanced pairs:552component samples+92excluded one-shots, all0alloc. The12store medians improve1.1304–1.1560×; ALL72pairs faster. This is a discretionary engineering keeper with a robust component gain, NOT a1.5× threshold claim. Vertical12medians improve1.0012–1.0291× with22slowerpairs. Picker10medians .99044–1.01016×,3median losses/35slowerpairs. Horizontal controls improve1.0026–1.1504× despite byte-for-byte SAME364-instruction Hkernel at SAMEaddress in DSP and whole executables; these are diagnosticcaller/layout controls, not Hkernel gains or a proven mechanism. Preserve all losses and raw pairs.

Whole VP8 720p/300inputs/CPU8/2500kbps/one-thread, six all-permutation triples:parent7.0826535ms,candidate7.0571785ms,old7.260224ms medians. Candidate−0.3596816%parent,5/6faster, one+.7969348%slower; candidate−2.7966837%old,4/6faster,two slower. Old samples vary6.867771–7.455178ms: NO stable whole gain, assembly-gap closure or cross-window claim. All22reportednon-timing semantic fields match,all18whole samples0alloc; no bitstream-identity claim.

Frozen raw files preserve both complete Go source maps,9binary bindings (8new+oldbaseline), commands, source captures, benchmark logs, whole JSON, corrections and independent reproducer. Binaries/full trees excluded. No README changes.


Publication: `raw-evidence.tar.gz` preserves all 176 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
