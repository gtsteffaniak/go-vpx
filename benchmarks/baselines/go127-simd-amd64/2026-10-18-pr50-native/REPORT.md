# PR50 exact10c native full audit

Run39122180995/job117659832036 SUCCESS at2026-10-18T11:55:00Z. Hardware AMD EPYC7763, Go1.27.2/SIMD/PGOoff/GOMAX1. Original reference691cd051 and previouscb9bbdcd are verified from raw environment. All127source hashes come from explicit10caaa67 Git blobs, never a continuation worktree.

Strict independent audit passes96raw/127hashes/34commonfixtures/549cases×18=9882timings/491excluded/guard396+66/optional132excluded+216Stats timings. Token15one-shots and ordinary/race independent state oracle pass on all3variants; the in-place reference adapter is used only on old691, preserving cb9originalAPI. ActualSDK512=false here; no selected512 throughput acceptance is claimed. The feature-first private experiment is not applied.

The18new6-slot in-place medians improve1.708–1.870× versuscb9; all108pairs≥1.698×. Generic18pairedDC medians1.467–1.669×,9medians≥1.5 and minimum pair1.447×. All15token workloads are retained: allEOB p0/p230 improve2.313×/2.050×; denseSmall controls regress to0.9837–0.9927×. No zero/alias shortcut is used as the sole acceptance denominator.

Whole medians vs cb9: VP8720−3.435% (all6faster), VP9320−0.550% (4/6faster), VP9720+0.753% (4/6slower). Versus original691:−24.068%,−40.901%,−51.267% respectively, all6faster. All paired deltas remain. No crossCPU/window pooling or source-identity excuse.

All22semantic fields match. Current/base allocation samples match across every workload/profile. One old691 VP9720round3 sample is76/60allocs/frame versus75/60 usual; this exception is preserved, not normalized away. All549kernel allocations arezero.246median regressions and23all-six-slower cases remain in fullproof.

No executables or full repository snapshots are included. Exact127tracked sources, raw logs/profiles, complete auditors/manifests and all result tables are preserved. This package does not authorize merge.


Publication: `raw-evidence.tar.gz` preserves all 242 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
