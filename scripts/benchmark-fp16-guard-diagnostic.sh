#!/bin/sh
# Reversible guard trial: current-only diagnostics, outside codec timing pairs.
set -eu
if [ "$(uname -s)" != Linux ] || [ "$(uname -m)" != x86_64 ]; then
 echo 'FP16 guard timing requires native Linux x64; Rosetta is correctness-only.' >&2
 exit 2
fi
binary=${1:?usage: benchmark-fp16-guard-diagnostic.sh CURRENT_PIXEL_TEST OUTPUT_LOG}
output=${2:?usage: benchmark-fp16-guard-diagnostic.sh CURRENT_PIXEL_TEST OUTPUT_LOG}
case "$binary" in /*) ;; *) binary="$(pwd)/$binary" ;; esac
case "$output" in /*) ;; *) output="$(pwd)/$output" ;; esac
[ -x "$binary" ]
export GOMAXPROCS=1 GOVPX_REQUIRE_NATIVE_GUARD=1 GOVPX_REQUIRE_SIMD=1
# The actual trial guard must remain the exact frozen unrolled production implementation.
actual_sha=$(sha256sum internal/pixel/quantize_fp_amd64.go | awk '{print $1}')
expected_sha=6402500232ee21792b751cad72aa24417e695f4bc0b49749308b839d93b1efdf
if [ "$actual_sha" != "$expected_sha" ]; then
 echo 'Trial FP16 source changed; this frozen diagnostic needs explicit revalidation.' >&2
 exit 2
fi
{
 printf '# component diagnostic only; no codec speed claim\n'
 printf '# head '; git rev-parse HEAD
 printf '# GOMAXPROCS=1 GOEXPERIMENT=simd SDK=go1.27.2 pgo=off\n'
 printf '# uname '; uname -a
 lscpu
 scripts/go.sh version
 scripts/go.sh env GOARCH GOOS GOAMD64 GOEXPERIMENT GOCACHE
 sha256sum "$binary" internal/pixel/quantize_fp_amd64.go internal/pixel/fp16_guard_component_native_test.go scripts/benchmark-fp16-guard-diagnostic.sh
 printf '# independent interval/full-input oracle and selected backend\n'
 "$binary" -test.run '^Test(FP16GuardComponentOracle|BackendSelected)$' -test.count=1 -test.v
 printf '# excluded one-shot oracle inventory: 66 cases\n'
 "$binary" -test.run '^$' -test.bench '^BenchmarkFP16GuardComponent$' -test.benchtime=1x -test.count=1 -test.cpu=1 -test.benchmem
} > "$output" 2>&1
for round in 1 2 3 4 5 6; do
 case "$round" in
  1) order='parent single unrolled' ;;
  2) order='single unrolled parent' ;;
  3) order='unrolled parent single' ;;
  4) order='parent unrolled single' ;;
  5) order='unrolled single parent' ;;
  6) order='single parent unrolled' ;;
 esac
 printf '# round %s order %s\n' "$round" "$order" >> "$output"
 for variant in $order; do
  {
   printf '# round %s guard %s\n' "$round" "$variant"
   "$binary" -test.run '^$' -test.bench "^BenchmarkFP16GuardComponent$/^$variant$" -test.benchtime=80ms -test.count=1 -test.cpu=1 -test.benchmem
  } >> "$output" 2>&1
 done
done
# This is an inventory/allocation check, never a speed threshold.
python3 - "$output" >> "$output" 2>&1 <<'PY'
import sys,re,collections
lines=open(sys.argv[1]).read().splitlines();phase='preflight';counts=collections.Counter();names=set();alloc=[];pername=collections.Counter()
for line in lines:
 if line.startswith('# round '):phase='timed'
 if not line.startswith('BenchmarkFP16GuardComponent/'):continue
 m=re.match(r'^(Benchmark\S+)\s+\d+\s+(?:[\d.eE+\-]+ ns/op\s+)?(\d+) B/op\s+(\d+) allocs/op',line)
 if not m:raise SystemExit('Unparseable component result: '+line)
 name,b,a=m.groups();counts[phase]+=1;names.add(name);pername[(phase,name)]+=1
 if (int(b),int(a))!=(0,0):alloc.append(line)
assert counts=={'preflight':66,'timed':396},counts
assert len(names)==66,len(names)
assert all(pername[('preflight',n)]==1 and pername[('timed',n)]==6 for n in names),pername
assert not alloc,alloc
print('component inventory PASS:66 names,396 timed samples,66 excluded one-shots,all0B/0alloc')
PY
