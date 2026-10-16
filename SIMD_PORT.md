# Go 1.27 SIMD migration

This branch replaces govpx's hand-written assembly with Go implementations.
The module defaults to Go 1.27.2. `make` and `scripts/go.sh` enable
`GOEXPERIMENT=simd`; consumers enable that experiment when building their
application. `purego`, disabled SIMD, and unsupported architectures retain
scalar Go implementations. amd64 SIMD dispatch requires AVX2.

The migration is in progress. Pixel arithmetic, VP8 forward transforms and
quantization, and VP9 forward/inverse transform networks and quantization
have Go SIMD paths. Remaining scalar dispatches and performance regressions
are being addressed before the PR is ready for review.

## Validation

Run commands through `scripts/go.sh`, which uses the installed Codex
`project-env` wrapper, or derives a project cache from `git-common-dir` on
machines without Codex. Each agent can set `CODEX_AGENT_ID` for its temporary
build directory.

```
scripts/go.sh test ./...
scripts/go.sh test -tags purego ./...
GOEXPERIMENT= scripts/go.sh test ./...
make test-noasm
GOOS=linux GOARCH=amd64 scripts/go.sh build ./...
GOOS=linux GOARCH=arm64 scripts/go.sh build ./...
```

The scalar libvpx v1.16.0 ports remain independent differential oracles.
`scripts/gensimd.go` translates the one-dimensional VP9 transform networks;
run `scripts/go.sh run scripts/gensimd.go` after changing those networks.
Integer narrowing must wrap at exactly the same stages as the scalar code.
Pixel filters round and clip after each pass, using int32 products to avoid
premature saturation from signed filter lobes.

Baseline: `origin/main` at `691cd051`, official Go 1.27.2, darwin/arm64,
Apple M4 Max. Pre-migration internal VP8/VP9 tests passed. Assembly kernel
benchmarks and a VP9 realtime encoder CPU profile were captured before
removal. Performance measurements will be recorded after specialization;
the initial generic SIMD transform dispatch has measurable overhead and is
not an accepted performance result.
