# VP9 single RD model — preliminary codec screen

Native Apple M4 Max/Darwin arm64, Go1.27.2 SIMD, CGO_ENABLED=1, PGO off, GOMAXPROCS=1; root project wrapper. Source baseline7e0e37b7; current PR main base preserves identical VP9 picker source.

The non-RD picker used to run ModelRdForSbY before the large-block model and overwrite every small-model result whenever the large model was valid. Keep prediction/variance handling and the existing large-model code unchanged; run the small model only when the large model does not provide a valid result. Apply to filter sweep and ordinary candidate path. Both models are pure value computations.

720p60 canonical CPU8 realtime2500kbps one-thread codec paired throughput ratios1.037342/1.007396 (median1.022369x, preliminary). All non-performance fields match, codec allocations unchanged at1.25/frame. A320p second-parent outlier is preserved and excluded from gain claims. No component speed claim for this variant.

All96 VP8 and72 VP9 ML packet hashes match, and existing prepared-predictor differential oracle passes. Exact source and binary bindings, logs, patch and every raw measurement are retained here. This is the smaller alternative to the retained larger prediction/variance-scan elision experiment; that component-qualified code and its losses remain on the performance evidence branch. No full suite/CI, profile, C comparison or README changes.
