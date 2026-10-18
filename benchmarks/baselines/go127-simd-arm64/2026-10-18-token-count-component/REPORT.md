# Independent token-count51 component audit

PASS: exact10caaa67 baseline1952tracked Go files plus one shared fixture; initial candidate/countonly each change only probability_tokens.go. All1953-source maps and6bound encoder/CLI hashes verify. Six buildinfos confirm Go1.27.2/CGO1/DarwinARM64; explicitPGOoff is verified in saved commands/bindings because Go buildinfo omits the off setting. No new native tests, builds or timing were run by this audit.

Saved original144native samples and revised216native samples are independently parsed with all12cases, six balancedAB/BA orders then every three-way permutation, monotonic non-overlapping commands, successful process exits and0alloc. Native hardware Apple M4Max/GOMAX1. No cross-window pooling; the original candidate and all initial count-only losses remain separate.

The frozen macro oracle is byte-exact original10c except functionname. Both block workers are byte-exact unchanged. The final dedicated count-only walk is exactly the original macro walk with record forwarding removed; records-nil is routed once. The three prepared zero-EOB paths preserve counter increment before record append and context updates afterward. Shared fixture executes live B.Loop with fresh contexts/recordlength each iteration, accumulating counters checked against independent reference×b.N; full records/context checks run before/after and a sink retains output.

All12288full-state/error/panic/backing cases pass ordinary on3variants and race/checkptr2 on both candidates. Full encoder-suite success logs preserved. Coverage includes nil entries, invalid contexts/EOBs/coefficients, above-left alias, backingcapacity and partial error state; arbitrary unsafe cross-type storage alias is not claimed.

In the216sample window final prepared-zero medians improve1.866× (non4x4) and1.798× (4x4). Initial count-only controls at0.915–0.964× are replaced by1.230–1.326× for zero/sparse. The 4x4 sparse+records control still regresses to0.982×; every slower sample remains. This is a portable token-count component gain, not SIMD arithmetic, overall-codec or x64 performance. Root will measure combined whole-codec behavior separately.

Original failed fixture syntax, setup failure log and failed-attempt source are preserved; the corrected fixture has sharedSHA5d42b86f48f9257b4d0bb34ed9f18257db0386cf1770e21937ae47f83489c9c9. No failed measurement is passed off as acceptance evidence. Final native codegen is captured from already-bound binaries read-only.

Only binary blobs and full source snapshots are omitted; exact6binary SHA/size/private paths/buildinfos and all-source hashes remain. All raw360timed samples,36excluded one-shots, logs, drivers, bindings, original/revised summaries, exactsource diffs and selected overlay are preserved. ROOT/Git/PRs unchanged.


Publication: `raw-evidence.tar.gz` preserves all 129 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
