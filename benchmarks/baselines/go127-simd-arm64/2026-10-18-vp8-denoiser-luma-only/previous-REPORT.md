# VP8 useless-work pass, merged results

PR107 and PR108 were merged after exact-head short correctness and native component feedback passed. Main now499ef24ee5f30dbb44890e3fb19b88bbe2e0407b. README unchanged; no full validation, profile audit or fresh libvpx comparison.

## Removed work

- PR107: accepted fast quantization produced dequantized coefficients that no caller consumed; reconstruction computed them again. QCoeff/EOB-only SIMD paths omit those multiplies/stores on ARM and x64, and write EOBs directly rather than staging them through empty context walks. RD scoring and fallback semantics remain intact.
- PR108: an exact motion predictor already attains the globally smallest SAD-plus-motion cost, but the old code still searched surrounding neighborhoods. Stop at that floor, retaining each method’s return-cost convention and fractional refinement.

## Scope of measured gains

All figures are preliminary short native screens versus715645fd, PGOoff/Go1.27.2SIMD/zero allocations.

| Component | Apple M4 Max ARM64 | Intel Xeon8370C x64 VM |
|---|---|---|
| Accepted coefficient pipeline | Zero-residual1.153×; DC1.127–1.163×; dense1.088–1.092× | Zero-residual1.210×; DC1.135×; dense1.139× |
| Eight-block quantization with discarded dequantization removed |1.21–1.32×| Covered through public accepted-MB component above |
| Exact/flat/shifted full-pel motion search |1.77–2.65×|1.71–2.43×|
| Complete HEX motion search on exact/flat inputs |1.20–1.24×|1.24–1.27×|

Dense/improved-start motion controls and all observed losses are retained in the per-run logs, including the approximately1% x64 DIAMOND dense loss. There is no claim of a2× whole-codec speedup. Individual ARM canonical whole-codec screens remain near flat for fixed-speed VP8 and around1–2% for adaptive VP8; small deltas are preliminary and are not added together. No x64 whole-codec speed claim is made. Intel component measurements are not combined with earlier EPYC libvpx timings to calculate a new gap.

Affected scalar/overlap/pure-Go tests, frozen motion dispatch differentials,96 VP8 packet/flag hashes and72 root ML packet hashes pass. The x64 two-variant run verifies actual AVX2/AVX512 native availability, both source overlays and both oracle sets. The1.5× immutablec4 whole-codec gate remains unmet.

The extra table-initialization experiment is retained privately in `inconclusive-token-table-probe.tar.gz`: it skips an unused RD cost-table fill in the fast picker, but the short whole-codec screen is too noisy to establish a qualifying gain. The patch, losses and limitations are preserved; it is not merged or presented as a win.

See `x64-REPORT.md`, raw archive, compact binding and merge receipts for the native run39155920162. Main source/report links remain in the two merged PRs. The denoiser luma-only follow-up is separate ongoing work.
