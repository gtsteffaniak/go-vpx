# Go 1.27 SIMD migration measurements

Measured 2026-10-16 on Apple M4 Max, darwin/arm64, official Go 1.27.2,
GOEXPERIMENT=simd, -pgo=off. Assembly baseline: main at 691cd051.
SIMD implementation: PR #23. See ../../../SIMD_PORT.md for interpretation.
The final recorded kernel runs and alternating encoder runs were serialized
with no competing benchmark experiments.

`assembly.txt` and `simd.txt` contain five 150ms runs of existing public
kernel benchmarks, with identical inputs and setup. Run against each source
tree (Go commands through the project cache wrapper):

```sh
scripts/go.sh test -pgo=off ./internal/vp9/dsp ./internal/vp8/dsp ./internal/vp9/encoder \
  -run '^$' \
  -bench 'BenchmarkVP9(Idct4x4_16Add|Idct8x8_64Add|Convolve8Horiz16x16|Convolve8Vert16x16|SubPixelVariance16x16)$|BenchmarkSAD16x16$|BenchmarkForwardDCT16x16$' \
  -benchtime=150ms -count=5
```

`encoder.json` retains all six runs per implementation, alternating A/B then
B/A. Build each command with `-pgo=off`, then execute:

```sh
govpx-bench -codec=vp9 -width=320 -height=180 -frames=120 -fps=30 \
  -bitrate=500 -mode=realtime -cpu-used=8 -threads=1 -encode-only \
  -auto-libvpx=false -format=json
```
