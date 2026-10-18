# PR49 nonAVX512 quantizer dispatch investigation

Pinned baselinec8bd67df and currentcb9bbdcd were reconstructed for Linux/amd64 and Darwin/amd64 using officialGo1.27.2,GOAMD64v1,PGOoff,CGO_ENABLED0. Original nativeCI usedCGO1; generated-code comparisons are scoped reconstructed evidence. No native or emulated performance timing was run. The measured VP9 regression cause remains UNPROVEN. Exact full49 pair counts:320p6/6slower,720p4/6slower; VP85/6slower.

The frame-growth hypothesis is rejected by these binaries. Public pixel.QuantizeFPBlock is136LLVM instructions with0x98 local frame and normalized EXACT across baseline/current/candidate. quantizeFPBlockSIMD local frame remains0xa8 and early prologue/argument homes are identical. Its static instruction count grows475→497 includingNOPs; the extra wide-call argument marshaling is after the taken-on-512 branch and is not executed when512is unavailable.

Concrete eligible n256+ non512 overhead: geometry compare/branch,mod32test/branch,thenSDK512compare/branch—six nonNOP instructions. The one-line private candidate evaluates SDK512 first, reducing that path to one compare/branch—two instructions. The exact Linux/Darwin normalized current→candidate diff is only that ordering and paddingNOP placement. Public quantizer, all10pair disjoint guard and32lane helper are normalized EXACT; vector arithmetic, call sites, spills and frames otherwise unchanged. All conjunction predicates remain present inside the original eligibility/alias guard. No public check, store order or fallback is moved.

Controls remain essential: non512 n64has two instructions in both versions, but the candidate uses a feature-byte load rather than an integer-n compare. On512 hardware,n64adds a feature compare/branch before its size rejection. These changes may win or regress; no throughput benefit or measured-cause claim is made.

Six translated correctness runs—ordinary/checkptr2 × selected/512off/allAVXoff—each PASS15,SKIP1. ActualSDK512=false locally, so selected32 native oracle remains skipped; no fakecaps. Existing n256exact-window/alias/partialpanic tests pass on fallback. Actual512 reachability is structurally identical, and selected512 execution/performance requires future nativeCI. Existing native49 predecessor selected512 proof does not substitute for running this changed dispatch.

Exact sources, command lines, codegen/normalized diffs, all test logs and7binarySHA/size/buildinfo receipts are preserved. Binary blobs and full source snapshots are omitted only from this publication package; private originals remain in the work directory. One candidate only; ROOT50 and all Git/PRs are unchanged.


Publication: `raw-evidence.tar.gz` preserves all 90 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
