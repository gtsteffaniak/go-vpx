#!/bin/sh
# Additive two-way pixel API diagnostics; old691 has no pixel package.
set -eu
case "$(uname -s)/$(uname -m)" in Linux/x86_64) ;; *) echo 'Native Linux/x64 required; Rosetta is correctness-only.' >&2; exit 2 ;; esac
parent=${1:?usage: benchmark-pixel-optional.sh PARENT_PIXEL_TEST CURRENT_PIXEL_TEST OUTPUT_LOG}
current=${2:?usage: benchmark-pixel-optional.sh PARENT_PIXEL_TEST CURRENT_PIXEL_TEST OUTPUT_LOG}
output=${3:?usage: benchmark-pixel-optional.sh PARENT_PIXEL_TEST CURRENT_PIXEL_TEST OUTPUT_LOG}
[ -x "$parent" ] && [ -x "$current" ]
source_dir=${4:-$(pwd)}
head_commit=${5:-$(git rev-parse HEAD)}
base_commit=${6:-}
export GOMAXPROCS=1 GOVPX_REQUIRE_SIMD=1
{
 printf '# additive component API only; qnil has no actual encoder caller claim\n'
 printf '# head %s\n' "$head_commit"
 if [ -n "$base_commit" ]; then printf '# base %s\n' "$base_commit"; fi
 printf '# driver checkout '; git rev-parse HEAD
 printf '# SDK=go1.27.2 GOEXPERIMENT=simd PGO=off GOMAXPROCS=1 rounds=6 benchtime=80ms\n'
 uname -a; lscpu
 scripts/go.sh version
 scripts/go.sh env GOOS GOARCH GOAMD64 GOEXPERIMENT GOCACHE
 sha256sum "$parent" "$current" "$source_dir/internal/pixel/fp32_optional_native_test.go" "$source_dir/internal/pixel/stats_reduction_native_test.go" "$source_dir/internal/pixel/quantize_fp32_avx512_amd64.go" "$source_dir/internal/pixel/goav1_pixelstats_simd.go" scripts/benchmark-pixel-optional.sh
 printf '# actual current SDK feature; unsupported AVX512 is explicitly unavailable\n'
 GOVPX_REQUIRE_FP32_NATIVE=0 "$current" -test.run '^TestFP32PublicHost$' -test.v -test.count=1
 for variant in parent current; do
  case "$variant" in parent) binary=$parent ;; current) binary=$current ;; esac
  printf '# excluded preflight %s: 48 FP32 API and 18 Stats cases\n' "$variant"
  GOVPX_REQUIRE_FP32_NATIVE=0 "$binary" -test.run '^TestFP32PublicHost$' -test.bench '^Benchmark(FP32PublicOptionalNative|StatsReductionNative)$' -test.benchtime=1x -test.count=1 -test.cpu=1 -test.benchmem -test.v
 done
} > "$output" 2>&1
available=0
host=$(grep 'FP32PublicHost AVX512=' "$output" | head -n 1)
if printf '%s\n' "$host" | grep -q 'FP32PublicHost AVX512=true AVX2=true Enabled=true'; then
 available=1
 {
  printf '# SELECTED_AVX512_AVAILABLE=1; actual native selected oracle required\n'
  GOVPX_REQUIRE_FP32_NATIVE=1 "$current" -test.run '^TestFP32(ActualFeature|NativeSelectedOracles|EligibleExactWindows|EligibleAliasedDeclines|EligiblePanicWindows)$' -test.count=1 -test.v
 } >> "$output" 2>&1
else
 printf '# SELECTED_AVX512_UNAVAILABLE; FP32 fallback timings omitted; cannot accept AVX512 correctness/performance\n' >> "$output"
fi
for round in 1 2 3 4 5 6; do
 order='parent current'; case "$round" in 2|4|6) order='current parent' ;; esac
 printf '# round %s order %s\n' "$round" "$order" >> "$output"
 for variant in $order; do
  case "$variant" in parent) binary=$parent ;; current) binary=$current ;; esac
  {
   printf '# round %s pixel %s\n' "$round" "$variant"
   pattern='^BenchmarkStatsReductionNative$'
   if [ "$available" = 1 ]; then pattern='^Benchmark(FP32PublicOptionalNative|StatsReductionNative)$'; fi
   "$binary" -test.run '^$' -test.bench "$pattern" -test.benchtime=80ms -test.count=1 -test.cpu=1 -test.benchmem
  } >> "$output" 2>&1
 done
done
python3 - "$output" "$available" >> "$output" 2>&1 <<'PY'
import sys,re,collections
phase='preflight';variant=None;counts=collections.Counter();names=collections.defaultdict(set);per=collections.Counter()
for line in open(sys.argv[1]):
 if line.startswith('# excluded preflight '):variant=line.split()[3].rstrip(':')
 if line.startswith('# round '):
  phase='timed'
  if ' pixel ' in line:variant=line.split()[4]
 if not line.startswith(('BenchmarkFP32PublicOptionalNative/','BenchmarkStatsReductionNative/')):continue
 # Verbose one-shot preflights print a leaf name before its result row.
 if phase=='preflight' and re.fullmatch(r'(?:BenchmarkFP32PublicOptionalNative|BenchmarkStatsReductionNative)/\S+\s*',line):continue
 m=re.match(r'^(Benchmark\S+)\s+\d+\s+(?:([\d.eE+\-]+) ns/op\s+)?(\d+) B/op\s+(\d+) allocs/op',line)
 assert m,('unparseable',line)
 name,ns,b,a=m.groups();assert (int(b),int(a))==(0,0),line
 if phase=='timed':assert ns is not None and float(ns)>0,line
 counts[(phase,variant)]+=1;names[(phase,variant)].add(name);per[(phase,variant,name)]+=1
expected=66 if sys.argv[2]=='1' else 18
for v in ('parent','current'):
 assert counts[('preflight',v)]==66,counts
 assert len(names[('preflight',v)])==66,names
 assert counts[('timed',v)]==6*expected,counts
 assert len(names[('timed',v)])==expected,names
 for n in names[('timed',v)]:assert per[('timed',v,n)]==6,per
assert names[('preflight','parent')]==names[('preflight','current')]
assert names[('timed','parent')]==names[('timed','current')]
print('optional pixel inventory PASS:132 excluded one-shots; timed samples=%d; actual selected512=%s; all0B/0alloc'%(12*expected,sys.argv[2]))
PY
