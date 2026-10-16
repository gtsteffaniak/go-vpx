# Final local validation

All commands run through `scripts/go.sh` (or Make targets using that wrapper)
with Go 1.27.2, `GOEXPERIMENT=simd`, `GOFLAGS=-p=2`, an agent-specific temporary
scope and the common Git project build cache. The seven raw logs are retained
as `validation-*.log`.

- Full native suite passes; root package 47.483 s.
- Full purego suite passes; root package 117.286 s.
- Pixel and both codecs' DSP/encoder packages pass race and checkptr=1.
- Strict encoder byte parity passes; root package 36.086 s.
- Regression report passes 18/18.
- No-assembly, amd64/arm64 native/purego/experiment-off selection, formatter
  and refreshed-PGO checks pass.

The full native/purego suites precede the final amd64 vertical-filter scope
change. That change restores the previously measured generic vertical SIMD
body only on amd64. The architecture constant makes all 3,435,188 ARM
instruction-region bytes identical before/after; the neighboring stored
parser/hash/disassembly evidence proves this equivalence. Purego excludes
the changed SIMD files. Five-package native ARM and Rosetta AVX2 parity pass
after the scope change. Race/checkptr, strict bytes, report and source/PGO
checks run on the exact committed dff92311 sources.
