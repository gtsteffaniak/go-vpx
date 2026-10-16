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

Hosted CI also exposed adaptive-timing noise in VP8 oracle fixtures.
The threaded byte matrix now pins its realtime floor with `CpuUsed=-4`.
The inter-mode distribution report pins speed 8 with `CpuUsed=-8`, using
fresh expected curves generated from unchanged `691cd051` with the same
fixed-speed fixture. Its percentage-point, L1, and EOB tolerances remain
unchanged; the independent 720p baseline now has L1=0 and EOB ratio=1.
Adaptive-control coverage remains in the other oracle tests. CI retains
the complete parity log as an artifact for future failed assertions.

## Performance and remaining limits

The migration removes architecture-specific assembly maintenance; it does
not establish performance parity with the removed assembly. Common inverse
DCT setup and the rolling vertical filter are substantially faster than the
initial generic SIMD ports, and native SIMD beats scalar Go on the measured
hot kernels. Some kernels and the representative encoder workload remain
slower than the assembly baseline. The PR remains draft while that gap is
visible, and the historical README numbers should not be interpreted as
measurements of this implementation.

The initial migration measurements below compare `691cd051` with `d5374fe0` using the
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

### Subsequent optimization checkpoint

The next pass adds actual eight-lane AVX2 DCT8/16 pipelines, an ARM64
int16 inverse DCT8 pipeline, sixteen-pixel AVX2 convolution, exact narrowed
int16 convolution on both architectures, and fused FP quantization/EOB/token
classes. Four-pixel convolution now reaches the native backend. The wider
transforms include sparse and full signed-int16 parity tests; FP token-class
tests require the fused native path when SIMD is enabled.

On ARM64, inverse DCT8 now measures about 32 ns, forward DCT16 about 180 ns,
horizontal filtering about 98 ns and vertical filtering about 70 ns, with
zero allocations. Six new alternating public VP9 encoder pairs measure
1.889 ms/frame for SIMD versus 1.518 ms/frame for the assembly baseline
(24.4% more time); every run still emits the same 246,893 bytes. The initial
36.6% regression above is historical. ARM performance work remains open.

Native Linux/amd64 measurements on an AMD EPYC 7763, before these latest
wider kernels (commit `818e4630`), measured 5.408 ms/frame versus 6.995
ms/frame for the same assembly-era public VP9 workload: 22.7% less time.
All twelve runs agreed on bytes, encoded/dropped frame counts and allocations.
This result is specific to that CPU/workload; the old x64 dispatch had less
assembly coverage than ARM64. Horizontal convolution was still slower in
that checkpoint.

### Native amd64 measurements at `bb1a6ef3`

The next hosted run used an AMD EPYC 9V74, which is a different CPU from
the previous EPYC 7763 runner. Each comparison below uses baseline
`691cd051` and SIMD `bb1a6ef3` on the same runner, official Go 1.27.2,
`GOEXPERIMENT=simd`, `GOAMD64=v1`, `GOMAXPROCS=1`, and `-pgo=off`. Native
SIMD selection was asserted with `GOVPX_REQUIRE_SIMD=1`. Six pairs alternate
assembly/SIMD and SIMD/assembly order; kernel samples run for 150 ms.

| Operation | Assembly-era dispatch | Go SIMD dispatch | Time change |
| --- | ---: | ---: | ---: |
| VP8 SAD 16x16 | 8.204 ns | 12.505 ns | +52.4% |
| VP9 inverse DCT 4x4 | 55.84 ns | 17.76 ns | -68.2% |
| VP9 inverse DCT 8x8 | 289.25 ns | 46.46 ns | -83.9% |
| VP9 forward DCT 16x16 | 1,074.5 ns | 187.90 ns | -82.5% |
| VP9 horizontal filter 16x16 | 233.10 ns | 87.38 ns | -62.5% |
| VP9 vertical filter 16x16 | 618.70 ns | 87.665 ns | -85.8% |
| VP9 subpixel variance 16x16 | 146.65 ns | 85.92 ns | -41.4% |
| VP9 realtime 320x180 | 5.498 ms/frame | 3.607 ms/frame | -34.4% |
| VP8 realtime 720p | 14.040 ms/frame | 45.187 ms/frame | +221.9% |
| VP9 realtime 720p | 47.558 ms/frame | 28.475 ms/frame | -40.1% |

Both 720p workloads encode 60 input frames at 2,500 kbps, CPU-used 8,
one thread, with encode-only mode and external comparison disabled. The
320x180 case retains the 120-frame, 500 kbps configuration described above.
All kernels report zero allocations. The forward DCT scalar control is
unchanged within these samples: 1,074.0 ns versus 1,073.5 ns.

The VP8 720p regression is substantial and unresolved at this checkpoint:
every paired SIMD run takes 220.8–223.1% more time. The VP9 paired time
changes are -34.1 to -34.7% at 320x180 and -39.8 to -40.6% at 720p.
These results support workload-specific VP9 gains; they do not establish
overall performance parity with the removed assembly.

Every run agrees on output byte counts, encoded/dropped frame counts,
quantizer statistics/histograms, and allocations within each workload.
These JSON reports retain byte counts, not bitstream hashes; exact codec
bytes are checked separately by the oracle parity gates. VP8 720p
produces 685,957 bytes, 58 encoded/2 dropped frames, and zero allocations
per input frame. VP9 720p produces 642,654 bytes, 48 encoded/12 dropped
frames, and 1.25 allocations per input frame. The 320x180 output remains
246,893 bytes, 117 encoded/3 dropped frames, and 0.617 allocations per
input frame.

[The native amd64 workflow run](https://github.com/thesyncim/govpx/actions/runs/38947952120)
retains raw samples, encoder JSON, CPU/compiler metadata, and the dispatch
assertion. Do not compare absolute timings between its EPYC 9V74 and the
earlier EPYC 7763 run. The workflow now captures separate VP8/VP9 720p CPU
profiles after all timed pairs, with raw profiles and flat/cumulative top
tables, to guide the remaining work without changing the timing samples.

### Loop-filter and metric follow-up

The next checkpoint replaces the generic narrow loop-filter arrays with
named registers, adds a sixteen-position AVX2 filter, and copies goav1's
ARM64 byte-domain filter with full VP8/VP9 masks and macroblock corrections.
Vertical loads/stores use register transposes and touch exactly the original
modified columns. Entry bounds and minimal/overlapping-layout tests cover
unsafe row I/O; the block-limit-255 case retains exact int16 arithmetic.

Local ARM64 comparisons reduced VP8 720p from 18.74 to 14.40 ms/frame with
the narrow-register filter, then to 13.71 ms/frame with the copied byte
pipeline before the final bounded row-I/O tuning. These are intermediate
comparisons; the assembly baseline remains faster. Native x64 results for
this checkpoint are pending. Variance64x64 improved from 508 to 309 ns,
and a fused SAD/sum/SSE pass reduced bottom-edge source metrics from 915
to 203 ns. All measured kernels retain zero allocations.
