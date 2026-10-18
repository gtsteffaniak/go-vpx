#!/bin/sh
# Additive live public paired-DC timing, after the original focused SSE8 phase.
set -eu
case "$(uname -s)/$(uname -m)" in Linux/x86_64) ;; *) echo 'Native Linux/amd64 required.' >&2; exit 2 ;; esac
assembly=${1:?provide assembly DSP test binary}
base=${2:?provide base DSP test binary}
head=${3:?provide head DSP test binary}
out=${4:?provide output directory}
mkdir -p "$out"
export GOMAXPROCS=1 GOVPX_REQUIRE_SIMD=1
pattern='^BenchmarkDCPair(RowsNative|InPlaceRowsNative)$'
for variant in assembly base head; do
 case "$variant" in assembly) binary=$assembly;; base) binary=$base;; head) binary=$head;; esac
 "$binary" -test.run '^$' -test.bench "$pattern" -test.benchtime=1x -test.count=1 -test.cpu=1 -test.benchmem > "$out/$variant-preflight.txt" 2>&1
 sha256sum "$binary" >> "$out/binary-sha256.txt"
done
for round in 1 2 3 4 5 6; do
 case "$round" in 1) order='assembly base head';; 2) order='base head assembly';; 3) order='head assembly base';; 4) order='assembly head base';; 5) order='head base assembly';; 6) order='base assembly head';; esac
 printf 'round %s: %s\n' "$round" "$order" >> "$out/orders.txt"
 for variant in $order; do
  case "$variant" in assembly) binary=$assembly;; base) binary=$base;; head) binary=$head;; esac
  "$binary" -test.run '^$' -test.bench "$pattern" -test.benchtime=80ms -test.count=1 -test.cpu=1 -test.benchmem > "$out/$variant-round-$round.txt" 2>&1
 done
done
# Native GNU code capture follows every timing process. Missing old/inlined
# symbols are expected; B.Loop callers retain the actual public operation.
for variant in assembly base head; do
 case "$variant" in assembly) binary=$assembly;; base) binary=$base;; head) binary=$head;; esac
 {
  objdump --version | head -n 1
  for symbol in github.com/thesyncim/govpx/internal/pixel.DCPairInt32Native github.com/thesyncim/govpx/internal/pixel.DCPairInPlaceInt32Native github.com/thesyncim/govpx/internal/vp8/dsp.DCOnlyIDCT4x4AddPairInt32 github.com/thesyncim/govpx/internal/vp8/dsp.DCOnlyIDCT4x4AddPairInPlaceInt32 github.com/thesyncim/govpx/internal/vp8/dsp.BenchmarkDCPairRowsNative.func1 github.com/thesyncim/govpx/internal/vp8/dsp.BenchmarkDCPairInPlaceRowsNative.func1; do
   printf '\n# exact symbol %s\n' "$symbol"
   objdump -d --no-show-raw-insn --disassemble="$symbol" "$binary"
  done
 } > "$out/$variant-codegen.txt" 2>&1
done
python3 - "$out" <<'PY'
from pathlib import Path
import collections,json,math,re,statistics,sys
p=Path(sys.argv[1]);expected={f'BenchmarkDCPair{api}/pitch{pitch}/inplace{alias}/pair{pair}' for api in ('RowsNative','InPlaceRowsNative') for pitch in (8,32,1312) for alias in ('false','true') for pair in range(3)}
orders=[x.split(': ',1)[1].split() for x in(p/'orders.txt').read_text().splitlines()]
assert len(orders)==6 and len({tuple(x)for x in orders})==6
assert all(len(order)==3 and set(order)=={'assembly','base','head'}for order in orders)
for pos in range(3):assert collections.Counter(x[pos]for x in orders)=={'assembly':2,'base':2,'head':2}
pattern=re.compile(r'^(Benchmark\S+)\s+(\d+)\s+(?:([\d.eE+-]+) ns/op\s+)?(\d+) B/op\s+(\d+) allocs/op$',re.M)
def parse(path,timed):
 text=path.read_text();out={}
 for name,n,ns,b,a in pattern.findall(text):
  assert name not in out and int(n)>0 and int(b)==int(a)==0,(path,name)
  value=float(ns)if ns else None
  if timed:assert value is not None and math.isfinite(value)and value>0,(path,name)
  out[name]=value
 assert set(out)==expected,(path,expected-set(out),set(out)-expected)
 assert re.search(r'^PASS$',text,re.M),path
 return out
for v in ('assembly','base','head'):parse(p/f'{v}-preflight.txt',False)
rounds={v:[parse(p/f'{v}-round-{r}.txt',True)for r in range(1,7)]for v in ('assembly','base','head')};results=[]
for name in sorted(expected):
 samples={v:[r[name]for r in rs]for v,rs in rounds.items()};med={v:statistics.median(x)for v,x in samples.items()};paired={v:[a/b for a,b in zip(samples[v],samples['head'])]for v in ('base','assembly')}
 results.append({'case':name,'samples_ns':samples,'median_ns':med,'base_over_head':med['base']/med['head'],'assembly_over_head':med['assembly']/med['head'],'paired_reference_over_head':paired,'head_slower_pairs':{v:sum(x<1 for x in ratios)for v,ratios in paired.items()}})
summary={'scope':'Native Linux/x64 public paired-DC component only; no codec-wide claim','public_cases':36,'timed_samples':648,'excluded_one_shots':108,'round_orders':orders,'all_zeroalloc':True,'results':results}
(p/'summary.json').write_text(json.dumps(summary,indent=2)+'\n')
lines=['# Native x64 paired DC results','','36 public cases, six balanced rounds; every sample retained. No whole-codec claim.','','| Case | Head ns | Base/head | Old/head | Base pairs slower |','|---|---:|---:|---:|---:|']
for r in results:lines.append(f"| {r['case']} | {r['median_ns']['head']:.4f} | {r['base_over_head']:.4f} | {r['assembly_over_head']:.4f} | {r['head_slower_pairs']['base']}/6 |")
(p/'summary.md').write_text('\n'.join(lines)+'\n')
print('paired DC inventory PASS:36 cases/648 timings/108 excluded, all0B/0alloc')
PY
