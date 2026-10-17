# Empty projection regression check

Final review found that the batched projection glue evaluated `ref[off:]`
before checking whether there were any rows. Previously, zero or negative
row counts returned without touching either buffer. The public `IntProCols`
entry now preserves that no-op before SIMD dispatch.

`before.txt` reproduces the panic with a nil source and unused offset.
The regression test covers zero/negative rows and unused positive, negative,
and maximum integer offsets, preserving sentinel output values.

Go 1.27.2, GOEXPERIMENT=simd, Apple M4 Max, native arm64; targeted ARM suites,
AVX2-enabled Rosetta VP9 DSP correctness, experiment-off and purego VP9 DSP
checks pass. Rosetta is correctness evidence only. The PGO profile was
refreshed after the final source edit, and pre-commit passes with fingerprint
`bd5031e5889589f04c4d01da296645535ee6d44b`.

This correctness fix makes no additional performance claim. Broader final
validation logs are retained alongside these targeted checks.

Final exact-source gates pass: full native suite (45.181s root package),
full purego (112.845s), five-package race/checkptr1, strict encoder byte
parity (34.302s), 18/18 parity report, no-assembly audit, all six architecture
selection modes, formatting and refreshed PGO. Raw logs are retained here.
