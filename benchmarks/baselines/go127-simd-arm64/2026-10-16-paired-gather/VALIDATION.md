# Exact final gather validation

Go 1.27.2, SIMD enabled, Apple M4 Max. All Go commands run through
scripts/go.sh / the project cache wrapper, serialized with other local Go
jobs and benchmark windows. Source remained frozen after PGO refresh.

- Full native suite passes; root package 45.872 s.
- Focused purego pixel suite passes all new repeated/overlapping-row cases.
- Production purego sources are unchanged from the prior full 112.845-s
  suite; hosted CI independently reruns that full suite.
- Pixel and both codecs' DSP/encoder packages pass race and checkptr=1.
- Strict encoder byte parity passes; root package 34.267 s.
- Regression report passes 18/18.
- No-assembly, amd64/arm64 native/purego/experiment-off selection, formatter
  and exact refreshed-PGO checks pass.

Raw logs: /private/tmp/govpx-ci-gather-*.log
Final PGO source fingerprint: 9ef8ad9eb8de1186550791d0c6dc0c74230b9750
