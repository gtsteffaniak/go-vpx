# Combined ARM whole-codec tuning screen

Apple M4 Max, Go1.27.2, CGO=1, SIMD, PGO off, GOMAXPROCS=1. Three cyclic balanced rounds per workload across PR52 parent, combined PR53 production sources and old assembly691cd051. CPU8, one thread; VP8 720p/300 inputs/2500kbps, VP9 320x180/120 inputs/500kbps, VP9 720p/60 inputs/2500kbps. All22 semantic fields and allocation counts match. Full JSON, binary hashes and commands are retained. The candidate executable was built before the subsequent comment-only baseline-citation correction and PGO refresh; timing has PGO disabled.

Preliminary median time changes:
- VP8 720p: -5.693% versus parent, -1.092% versus old assembly.
- VP9 320x180: -1.288% versus parent, +8.211% versus old assembly.
- VP9 720p: -1.053% versus parent, +33.201% versus old assembly.

This short screen supports continued tuning; three rounds do not establish stable gap closure. These are combined results; do not attribute their entire change to any single kernel. Absolute libvpx comparisons are separately scoped in the published PR52 libvpx report. Targeted race/checkptr packet oracles also pass with the pointer refinement.
