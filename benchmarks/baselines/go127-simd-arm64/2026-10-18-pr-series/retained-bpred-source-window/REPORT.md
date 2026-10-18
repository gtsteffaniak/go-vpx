# Retained B_PRED source-window preparation probe

Prepare/validate a complete raw16x16 source window once for all64 mode-statistics and16 accepted residual calls; partial blocks retain existing clamping. Root frozen algorithm comparisons, clamped edges and96 VP8/72 ML streams pass. No new math or unsafe entrypoints.

Native M4 Max/Go1.27.2 SIMD, PGO off, CGO1, GOMAXPROCS1. Component B_PRED p0 throughput1.16629/1.16251x in two40ms pairs (qualifies user's15% component rule); p1 is1.13677/1.08484x. All component allocations0. Codec720p60 fixed−8 is0.98547/0.99457x; adaptive+8 is0.99612/0.96277x. All codec non-performance fields match and allocations stay0.

Code and raw native component win are retained here for scoped integration/shared-library use. Global dispatch is not enabled because this short screen does not improve codec throughput. Keep this qualifying source; do not lose it as an unsuccessful experiment. No full CI/C/profile/audit ran; README unchanged.
