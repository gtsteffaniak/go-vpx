# VP9 non-RD coefficient-cost table setup — preliminary codec screen

Production: `vp9_encoder_state.go`; exact private source in this directory and
`costgate.patch`. Baseline `7e0e37b7bbbcda65523f5581e7879c7d2d93827c`.
No root production or Git edits were made by this lane.

Frame-context preparation and adaptation previously rebuilt all coefficient
cost tables even when non-RD inter mode selection used modeled rates. Pinned
libvpx `vp9/encoder/vp9_rd.c:425–426` fills token costs only for full-RD mode
selection or a keyframe. The candidate skips that setup for non-RD,
non-compound inter frames, invalidates stale cost tables, and leaves the
existing current-probability scalar fallback available to unusual consumers.
Keyframes, intra-only, full-RD and compound frames keep table preparation.
There is no lazy mutation from parallel block consumers.

Native Apple M4 Max / Darwin arm64, Go1.27.2 SIMD, CGO=1, GOMAXPROCS=1,
PGO off. Two balanced short Go-only canonical VP9 CPU8 CBR encode-only
samples; one encoder thread, no C/quality/profile/full-suite/matrix run.

| Workload | Baseline ns/frame | Candidate ns/frame | Speedup |
|---|---:|---:|---:|
| 720p, 60f, 2500kbps, pair1 | 11,060,568 | 10,662,951 | 1.03729× |
| 720p, 60f, 2500kbps, pair2 | 10,785,124 | 10,612,889 | 1.01623× |
| 320×180, 120f, 500kbps, pair1 | 1,474,995 | 1,464,106 | 1.00744× |
| 320×180, 120f, 500kbps, pair2 | 2,136,403 | 1,551,451 | 1.37704× |

720p median paired speedup is **1.02676× (+2.68% throughput)**, preliminary.
The second 320p baseline is an obvious outlier; **do not claim or aggregate
the apparent 37.7% gain**. Its raw data remains present. No component timing
claim is made from this local run. The shared fixture includes the actual
frame-context setup operation and key/full-RD/compound controls for requested
native x64 feedback.

All non-timing CLI fields matched within pairs, including output bytes,
encoded/dropped frames, quantizers and allocation counts. Exact stream hashes
matched unchanged baseline for all96 VP8 packets and72 VP9 ML packets (nine
configurations). Focused tests cover scalar coefficient costs across mutated
key/intra/non-RD/full-RD/compound contexts, both costing modes, all transform
sizes/planes/ref types, plus entropy-reset/error-resilient controls. The final
fixture additionally asserts cache validity per transition so stale-table
fallback coverage cannot silently disappear. Baseline compatibility is
selected by the presence of the candidate's internal frame-gate method.

`binding.json` binds production source, CLI binaries, fixture and stream
hashes. `native/commands.json` records every command. `summary.json` retains
all measurements. `manifest.json` and `overlays.tar.gz` are independent native
x64 variant202; `whole-commands.json` contains the canonical codec commands.
