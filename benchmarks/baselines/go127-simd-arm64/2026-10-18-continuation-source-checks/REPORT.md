# Continuation source and local check binding

Official Go1.27.2 archsimd on Apple M4 Max, CGO enabled. The exact 1952 Go-source hashes before generator replay match the final staged sources byte for byte. Full ordinary Go tests, generator replay, portable token purego/nosimd tests, PGO refresh, make pre-commit, no-assembly and SIMD build-selection gates all exited zero; full logs and source/profile/harness hashes are preserved. README files are unchanged. Local module stat-cache permission warnings are retained; they did not fail these commands.

Changed-path race/checkptr, translated x64 selected/off correctness, independent whole-buffer oracles, native kernel timing and whole-codec evidence are in the separate component archives. The full unscoped race/checkptr trial had eight allocation-assertion failures reproduced with the same values on the exact strong SAT parent; it is not reported as passing. Translated x64 checks are correctness only. Native Linux x64 throughput remains a CI gate. The provisional horizontal-loop unroll is archived but excluded from production sources.


Publication: `raw-evidence.tar.gz` preserves all 10 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
