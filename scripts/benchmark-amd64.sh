#!/bin/sh
# Run both implementations serially on one native amd64 runner. Each pair
# alternates A/B and B/A to expose ordering drift in the raw measurements.
set -eu
case "$(uname -s)/$(uname -m)" in
	Linux/x86_64) ;;
	*) printf '%s\n' 'This benchmark requires native Linux/amd64.' >&2; exit 1 ;;
esac

root=$(git rev-parse --show-toplevel)
cd "$root"
baseline_ref=${1:-691cd0512c4829d9ee40d483205dd223f9eb028d}
output_dir=${2:?provide an output directory}
mkdir -p "$output_dir"
output_dir=$(CDPATH= cd -- "$output_dir" && pwd)
reference_dir=$(mktemp -d "${TMPDIR:-/tmp}/govpx-assembly-reference.XXXXXX")
git worktree add --detach "$reference_dir" "$baseline_ref"
trap 'git worktree remove --force "$reference_dir" >/dev/null 2>&1 || true' EXIT

# Apply the same residual-domain benchmark and independent scalar test to
# both implementations. The assembly reference's production files stay at
# the pinned baseline commit.
cp "$root/internal/vp8/encoder/dct_batch_native_test.go" "$reference_dir/internal/vp8/encoder/dct_batch_native_test.go"

export GOTOOLCHAIN=local GOEXPERIMENT=simd GOAMD64=v1 GOMAXPROCS=1
{
	printf 'assembly commit: %s\n' "$(git rev-parse "$baseline_ref^{commit}")"
	printf 'SIMD commit: %s\n' "$(git rev-parse HEAD)"
	printf '%s\n' 'GOMAXPROCS=1 PGO=off rounds=6 kernel-benchtime=150ms'
	scripts/go.sh version
	scripts/go.sh env GOOS GOARCH GOAMD64 GOEXPERIMENT GOTOOLCHAIN GOCACHE
	lscpu
} > "$output_dir/environment.txt"

# This test fails if runtime feature detection selected the scalar backend.
GOVPX_REQUIRE_SIMD=1 scripts/go.sh test -pgo=off ./internal/pixel -run '^TestBackendSelected$' -count=1 -v > "$output_dir/dispatch.txt" 2>&1

packages='internal/vp9/dsp internal/vp8/dsp internal/vp9/encoder internal/vp8/encoder'
benchmarks='BenchmarkVP9(Idct4x4_16Add|Idct8x8_64Add|Convolve8Horiz16x16|Convolve8Vert16x16|SubPixelVariance16x16)$|BenchmarkSAD16x16$|BenchmarkForwardDCT16x16$|Benchmark(FastQuantizeBlockSIMD|ForwardDCT4x4SIMD|ForwardDCTBatchResidual25|LoopFilterHorizontalEdgesY|LoopFilterVerticalEdgesY|SAD16x16LimitEarly|SAD16x16x4LimitPtrFastFull)$'
for variant in assembly simd; do
	source_dir=$root
	if [ "$variant" = assembly ]; then source_dir=$reference_dir; fi
	for package in $packages; do
		name=$(printf '%s' "$package" | tr / -)
		scripts/go.sh -C "$source_dir" test -pgo=off -c -o "$output_dir/$variant-$name.test" "./$package"
	done
	scripts/go.sh -C "$source_dir" build -pgo=off -o "$output_dir/$variant-govpx-bench" ./cmd/govpx-bench
done

for round in 1 2 3 4 5 6; do
	order='assembly simd'
	case "$round" in 2|4|6) order='simd assembly' ;; esac
	for variant in $order; do
		printf 'round %s: %s\n' "$round" "$variant"
		for package in $packages; do
			name=$(printf '%s' "$package" | tr / -)
			{
				printf '# round %s, %s\n' "$round" "$package"
				GOVPX_REQUIRE_SIMD=1 "$output_dir/$variant-$name.test" -test.run '^$' -test.bench "$benchmarks" -test.benchtime=150ms -test.count=1 -test.cpu=1 -test.benchmem
			} >> "$output_dir/$variant-kernels.txt" 2>&1
		done
		"$output_dir/$variant-govpx-bench" -codec=vp9 -width=320 -height=180 -frames=120 -fps=30 -bitrate=500 -mode=realtime -cpu-used=8 -threads=1 -encode-only -auto-libvpx=false -format=json > "$output_dir/$variant-encoder-$round.json"
		for codec in vp8 vp9; do
			"$output_dir/$variant-govpx-bench" -codec="$codec" -width=1280 -height=720 -frames=60 -fps=30 -bitrate=2500 -mode=realtime -cpu-used=8 -threads=1 -encode-only -auto-libvpx=false -format=json > "$output_dir/$variant-$codec-720p-encoder-$round.json"
		done
	done
done

# Profile separately after every timing pair has finished. Retain both flat
# and cumulative views so a native encoder regression can be traced through
# its public callers without contaminating the benchmark samples.
for variant in assembly simd; do
	for codec in vp8 vp9; do
		profile="$output_dir/$variant-$codec-720p-profile.pgo"
		"$output_dir/$variant-govpx-bench" -codec="$codec" -width=1280 -height=720 -frames=120 -fps=30 -bitrate=2500 -mode=realtime -cpu-used=8 -threads=1 -encode-only -auto-libvpx=false -format=json -cpuprofile="$profile" > "$output_dir/$variant-$codec-720p-profile.json"
		scripts/go.sh tool pprof -top -nodecount=40 "$output_dir/$variant-govpx-bench" "$profile" > "$output_dir/$variant-$codec-720p-profile-top.txt"
		scripts/go.sh tool pprof -top -cum -nodecount=40 "$output_dir/$variant-govpx-bench" "$profile" > "$output_dir/$variant-$codec-720p-profile-cumulative.txt"
	done
done

# Keep artifacts compact; raw text/JSON and environment details are sufficient
# to reproduce and analyze these measurements.
rm -f "$output_dir/"*.test "$output_dir/"*-govpx-bench
