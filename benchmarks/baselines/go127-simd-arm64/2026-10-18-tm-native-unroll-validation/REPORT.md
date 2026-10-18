# Integrated ARM true-motion and x64 SIMD trial validation

Base merged PR46:2df3eaf5fb547db07bd970cd65bf504ecb287590. Seven changed production Go paths match frozen corrected ARM, FP16 guard and pointer-SSE candidates byte for byte. Copy16, README and module dependencies remain unchanged. The source ledger binds all changed Go files, final diagnostic scripts/workflow and refreshed PGO.

Native Apple M4 Max, official Go1.27.2: full41-package suite, ARM race/checkptr2, six native/purego/scalar build selections, assembly absence, refreshed PGO and pre-commit gates PASS. Commands use scripts/go.sh/project-env, git-common-dir cache identity, explicit project root /private/tmp/codex-projects, distinct agent IDs, GOFLAGS=-p=2/GOMAXPROCS2. These are correctness/build checks, separate from the single-thread PGO-off acceptance windows in frozen ARM reports.

Integrated darwin/amd64 checkptr2 binaries pass actual AVX/AVX2-selected and actual-disabled tests under Rosetta. Frozen FP full-signed/alias/fallback/component-oracle tests, SSE pointer/state tests, TMPredict scalar/unsupported/full-backing/panic tests pass. All71 actual DSP public one-shots per mode (27TM+36SSE8+8Copy16) pass with0B/0alloc. No translated throughput claim; displayed1x timing values are correctness artifacts. Historical instrumentation allocation caveats remain in frozen FP research; no assertion was weakened here.

The final focused-code capture executes only after timing, uses GNU binutils on native Linux and records the actual measured executable SHA before/after disassembly. Local sh-n/Python AST checks pass; GNU runtime capture remains pending native CI. Full coverage retains498threeway cases/8964samples/338excluded/31commonfixtures,86sourcehashes, plus separate current-only66guard names/396samples/66excluded,92full artifacts. Focused coverage retains720samples/120excluded,55artifacts. Mandatory kernel retention is independent of whole-codec results. New AMD production trials remain unmeasured pending CI.

The scalar-contract ordering supplement keeps the measured364-instruction ARM native body and actual public caller code unchanged. The original frozen report's below-threshold attribution is corrected separately: usermandated>=1.5 retention; smaller assembly-beating wins are retained by engineering judgment. Every raw ARM target/control/slower pair and separate whole/profile semantic/alloc proof is preserved. No source formatting change, test failure, or performance measurement was hidden during this integration.


Publication: `raw-evidence.tar.gz` preserves all 21 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.

Final pre-commit after staging also passes. The retained log includes a read-only module-stat-cache warning; Go completed the profile-backed build and the Make target removed its check binary normally. Final harness edits are shell/Python/workflow diagnostics; production Go bytes remain those tested above.
