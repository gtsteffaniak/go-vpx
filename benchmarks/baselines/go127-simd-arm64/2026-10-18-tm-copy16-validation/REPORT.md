# Integrated true-motion and Copy16 validation

Base: merged PR45, 78ce74d39806d74ecf91bc3d082899a74d7bc85f. The five changed production Go files match their frozen candidate bytes. No module dependencies or README files change. The source manifest binds the integrated Go tests, harness and refreshed PGO profile.

Native Apple M4 Max/arm64: full Go suite, focused race/checkptr=2, assembly absence, native/purego/scalar build selection, refreshed PGO and final pre-commit gates pass. Go1.27.2 runs through scripts/go.sh/project-env with project-common-dir cache scope, distinct agent IDs and explicit /private/tmp/codex-projects cache root. Compilation/test GOMAXPROCS=2 and GOFLAGS=-p=2 are separate from the single-thread native acceptance protocols in the ARM evidence packages.

Translated darwin/amd64 GOAMD64=v1/checkptr=2: actual SIMD selected and actual AVX/AVX2 disabled pixel guards/frozen semantics pass. Both actual Copy16 callers and all27 public TM one-shots pass with zero bytes/allocations:35 DSP and8 root cases per mode. These are correctness-only runs; their one-shot ns/op fields are not performance measurements. Frozen candidate packages preserve independent Linux/Darwin generated code and reference proofs. Native amd64 throughput remains pending CI.

The standalone common TM state test and27 one-shot fixtures also pass on all three ARM reference snapshots using Go overlays without modifying frozen source trees. Each state test exercises324 independent valid, disjoint calls and full input/output backing canaries. Current-only malformed-window/alias/partial-panic tests are never imposed on historical assembly.

Audit adjustments are preserved: the first final-format gate found one unformatted three-index slice in the new common test; gofmt changed only spaces, then the final gate passed. Three-reference compatibility logs retain the original unformatted fixture hash; the integrated-source manifest records the formatted hash. An initial root packet auditor assumed go.sum existed; all three snapshots consistently omit it. The auditor now verifies file presence before comparing bytes. No production source or measurement changed for either correction. Harness metadata separately corrects FP one-shots to108 reference+36 CPU-off (144 total); all338 excluded preflights remain distinct from498 timed cases.

ARM component and whole-codec results are in separate frozen baseline folders. Both stages retain qualifying kernel wins despite mixed/flat whole-codec outcomes; no product of gains across separate windows is reported. README files contain no experiment logs. The independent review preserves the originally identified missing root-caller timing and the corrected repeated native coverage.


Publication: `raw-evidence.tar.gz` preserves all 44 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
