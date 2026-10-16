# Go 1.27 native SIMD migration

The migration starts from `main` at `691cd051` and removes all 96 tracked
hand-written assembly files, their assembly declarations, and the custom
CPUID package. The module requires Go 1.27 and selects toolchain Go 1.27.2.
`make` and `scripts/go.sh` enable `GOEXPERIMENT=simd` by default. Applications
using govpx should also build with that experiment to select native kernels:

```sh
GOEXPERIMENT=simd go build ./...
```

ARM64 uses `simd/archsimd`; amd64 dispatch checks `archsimd.X86.AVX2()`.
`-tags purego`, `GOEXPERIMENT=`, unsupported architectures, and amd64 CPUs
without AVX2 retain scalar Go implementations. No module dependencies were
added. Compiler-generated machine code and standard-library implementations
are outside the repository's assembly exclusion check.

## Kernel coverage

| Family | Go implementation |
| --- | --- |
| Pixel search/statistics | SAD, row-limited SAD, four-reference SAD, variance/SSE, fused bilinear subpixel variance, projections and vector variance |
| Prediction/filtering | VP8 six-tap and bilinear, VP9 eight-tap and averaging, DC/TM/V/H and directional intra, VP8/VP9 loop filters |
| Residuals | Subtraction, nonzero detection, VP8 block-scan gather, coefficient error/energy/absolute sums |
| VP8 transforms | Forward/inverse DCT and Walsh, DC-only and paired reconstruction, dequantization and quantization |
| VP9 transforms | Forward DCT/ADST 4/8/16, inverse DCT 4/8/16/32 and ADST hybrids, sparse/DC variants, lossless Walsh, Hadamard and quantization |
| Temporal processing | ARNR accumulation and VP8 denoiser first pass |
| Bitstream packing | Serial Go bit-packing replaces the former assembly helper; its bit-order contract remains independently tested |

The new 32-point inverse DCT and directional intra paths extend SIMD
coverage beyond the former assembly dispatch. Forward 32-point transforms,
scaled-coordinate convolution, and other existing scalar operations retain
Go implementations.

The self-contained `internal/pixel` package holds copied goav1 SAD,
pixel-statistics, convolution, and widening/reduction kernels plus local
VP8/VP9 glue. [UPSTREAM.md](internal/pixel/UPSTREAM.md) pins the source commit
and adaptations; `third_party/goav1` retains its license and notices. The
available gopus audio kernels were inspected but are not required here.
Codec-specific rounding, clipping after each filter pass, scan orders, and
wrapping int16 narrowing remain govpx's responsibility.

`scripts/gensimd.go` generates four-lane VP9 transform networks from the
independent scalar functions and emits specialized DCT pipelines with
fixed-size scratch space. Regenerate with:

```sh
scripts/go.sh run scripts/gensimd.go
```

ARM64 vertical eight-tap filtering reuses seven rows in registers as the
window advances. Bounds-checked entry points cover exact final-row lengths;
unsafe loads read only four, eight, or sixteen bytes within the proven
source allocation. Pixel tests cover misalignment, minimal buffers, signed
filter lobes, coefficient extremes, tails, permitted aliases, and zero
allocations. Codec tests retain independent scalar/libvpx references.

## Validation

The validation environment is official Go 1.27.2, Apple M4 Max,
darwin/arm64, with `GOEXPERIMENT=simd`. AVX2 parity was also executed under
Rosetta with AVX advertisement enabled; that is correctness evidence, not
native x86 performance evidence. Runtime AVX2-disabled execution selects
the scalar backend separately from experiment-disabled builds.

The local checks include:

- Full native and purego package suites; experiment-disabled short suite.
- Native SIMD selection, six architecture/build-tag selection checks, and
  repository assembly exclusion.
- Race and checkptr=1 tests on pixel and both codecs' DSP/encoder packages.
- `go vet`, Linux amd64/arm64 builds, and a scalar Linux riscv64 build.
- Oracle trace compilation, decoder conformance, quality gates, libvpx DSP
  and codec oracles, pinned encoder byte parity, and parity regression report.
- Fresh VP8 encode/quality and VP9 realtime PGO profiles and source fingerprint.

Checkptr=2's numerical tests also pass. Its forced heap escapes change the
allocation assertions in pointer-heavy VP8 entry points; release and
checkptr=1 zero-allocation checks pass.

Use `scripts/go.sh` for Go commands: it invokes Codex's `project-env` wrapper
when installed, otherwise derives a project-wide cache from `git-common-dir`.
`CODEX_AGENT_ID` keeps concurrent agents' temporary build directories separate.

```sh
make verify-production
GOEXPERIMENT= scripts/go.sh test -short ./...
scripts/go.sh test -race ./internal/pixel ./internal/vp8/dsp ./internal/vp9/dsp ./internal/vp8/encoder ./internal/vp9/encoder
GOOS=linux GOARCH=amd64 scripts/go.sh build ./...
GOOS=linux GOARCH=arm64 scripts/go.sh build ./...
```

CI fixes preserve the existing quality thresholds: BD-rate now rejects
duplicates on its independent PSNR axis instead of rejecting valid equal
encoded rates; the ARNR fixture holds reference scheduling fixed while
toggling filter strength; the VP8 oracle trace helper disables libvpx's
interactive warning prompt. Both original quality failures were reproduced
on unchanged main with exactly the same curves.

## Performance and remaining limits

The migration removes architecture-specific assembly maintenance; it does
not establish performance parity with the removed assembly. Common inverse
DCT setup and the rolling vertical filter are substantially faster than the
initial generic SIMD ports, and native SIMD beats scalar Go on the measured
hot kernels. Some kernels and the representative encoder workload remain
slower than the assembly baseline. The PR remains draft while that gap is
visible, and the historical README numbers should not be interpreted as
measurements of this implementation.

Measurements below compare `691cd051` with this implementation using the
same Go 1.27.2 compiler, experiment, CPU and `-pgo=off`. Kernel medians use
five runs with 150ms per benchmark. The public encoder comparison alternates
six baseline and six SIMD runs: VP9 realtime, 320x180, 120 frames, 30 fps,
500 kbps, CPU-used 8, one thread, encode-only, external comparison disabled.

| Operation | Assembly-era dispatch | Go SIMD dispatch |
| --- | ---: | ---: |
| VP8 SAD 16x16 | 7.96 ns | 9.69 ns |
| VP9 inverse DCT 4x4 | 31.26 ns | 15.05 ns |
| VP9 inverse DCT 8x8 | 16.71 ns | 104.60 ns |
| VP9 forward DCT 16x16 | 83.03 ns | 217.30 ns |
| VP9 horizontal filter 16x16 | 20.71 ns | 164.40 ns |
| VP9 vertical filter 16x16 | 22.01 ns | 112.60 ns |
| VP9 subpixel variance 16x16 | 26.36 ns | 86.32 ns |
| VP9 public encoder | 1.474 ms/frame | 2.013 ms/frame |

The public encoder is **36.6% slower** on this workload. All twelve runs
produce the same 246,893 output bytes, 117 encoded/3 dropped frames,
quantizer statistics, and 0.617 allocations per input frame. The measured
range is 1.466–1.482 ms/frame for assembly and 2.008–2.043 ms/frame for SIMD.
This is one CPU/workload, not a general throughput or statistical confidence
claim. The forward 16x16 scalar reference measures about 545 ns; the SIMD
port remains faster than scalar while falling behind the removed assembly.

Raw results and commands are retained in
[benchmarks/baselines/go127-simd-arm64](benchmarks/baselines/go127-simd-arm64).
The public workload's CPU profile confirms live calls through the pixel
kernels. Generated ARM64 code also contains widening vector multiplies,
vector additions, and byte narrowing in the selected filter body. Its
multiply/add instruction sequence and the four-lane transform layout leave
further tuning opportunities; no x86 performance claim is made from Rosetta.
