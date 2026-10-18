# VP9 ML classifier reuse — short native screen

Exact parent: bc5e178048d18bdb2ea7801b7830a1e1f24ad0f6. Production candidate is the three files in cache.patch; source hashes are in cache-source.json. No ROOT changes. Reference: pinned libvpx 1.16.0 1024874c5919305883187e2953de8fcb4c3d7fa6, vp9/encoder/vp9_encodeframe.c:4662–4673. C computes the ML partition decision and passes pred_pixel_ready onward. Go re-ran the same five variance scans and NN during vp9NonrdReuseInterPredReady after the partition classifier. The per-SB source, estimated predictor, qindex and speed remain immutable; live pickPred is separate. Cache the 21 aligned square nodes; invalidate with the existing per-frame context reset. Unsupported and non-tree-aligned windows retain the existing classifier.

Correctness: both variants pass the affected ML/non-RD/predictor tests. New tests cover every cache key across five changed frame inputs, repeated queries, live pickPred mutations, unsupported/nonaligned windows, unready state and frame invalidation. Existing scalar-derived classifier oracle is unchanged. 72 public encoded packets across 64×64, 320×180 and 352×288 have identical length-delimited SHA256 in cache-{parent,candidate}-targeted.log. No new coding heuristic.

Apple M4 Max / Darwin arm64, Go1.27.2 archsimd, CGO1, PGOoff, GOMAXPROCS1. Three alternating-order 50ms component pairs: actual partition decision + subsequent predictor-reuse query medians 1.954× (64), 1.914× (32), 1.818× (16). All nine pairs faster, zero B/op and allocs/op; 16's minimum pair is 1.484×. This is a short screen of the repeated-query component, not the entire encoder.

Whole short screen: 320×180/120 inputs/500kbps CPU8 median +3.788% time, 2/3 slower; 1280×720/60 inputs/2500kbps negative control median +2.222%, 2/3 slower. The 720 lane does not use ML; no cause for the noisy common losses is established. All 22 reported semantic fields and allocation counts match. No whole throughput gain, C-throughput gain, or target closure is claimed. Preserve the component result under the user’s >=1.5× policy and retain every whole loss.

Raw measurements, command order, binary/source bindings and summary: cache-screen/. Unchanged fixtures: candidate/vp9_ml_cache_native_test.go and candidate/vp9_ml_cache_stream_native_test.go (not included in production patch). No broad retest or elaborate packaging; next effort targets BlockYrd scoring work.


Publication: `raw-evidence.tar.gz` preserves all 47 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
