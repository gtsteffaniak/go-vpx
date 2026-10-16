# Exact final fused ARM validation

All local commands use Go 1.27.2 through the project cache wrapper, with
SIMD enabled, `GOFLAGS=-p=2` and an agent-specific temporary directory.
Six exact logs are retained as `validation-*.log`.

- Full native suite passes; root package 45.486 s.
- Pixel and both codecs' DSP/encoder packages pass race and checkptr=1.
- Strict encoder byte parity passes; root package 34.636 s.
- Regression report passes 18/18.
- No-assembly, amd64/arm64 native/purego/experiment-off selection, formatter
  and refreshed-PGO checks pass.

The fused ARM change selects no new production or test source in purego:
all changed Go files are tagged SIMD and !purego or build-ignore. The prior
full purego 117.286-s suite therefore remains valid for identical selected
sources. Hosted CI independently reruns that full suite on the checkpoint.
Five-package native ARM and Rosetta AVX2 correctness also pass; Rosetta
is not used for performance measurements.
