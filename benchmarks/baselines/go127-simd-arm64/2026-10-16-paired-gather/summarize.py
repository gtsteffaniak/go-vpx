import json,re,statistics
from pathlib import Path
root=Path(__file__).parent
skip={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns'}
summary={'cpu':'Apple M4 Max','go':'go1.27.2','experiment':'simd','gomaxprocs':1,'pgo':'off','vp8_baseline':'c49e9d60 (a6c1f016 changes only VP9 empty-row behavior)','architecture':'native arm64; Rosetta only correctness'}
def public(name, baseline="base"):
 rows=json.loads((root/name).read_text())['samples'];values={}
 reference={k:v for k,v in rows[0]['result'].items() if k not in skip}
 for row in rows:
  assert {k:v for k,v in row['result'].items() if k not in skip}==reference
  values.setdefault(row['variant'],{})[row['round']]=row['result']['ns_per_frame']
 return dict(samples=len(rows),non_timing_fields_match=True,median_ns={v:statistics.median(a.values()) for v,a in values.items()},paired_pct_vs_base={v:[100*(a[i]/values[baseline][i]-1) for i in a] for v,a in values.items()})
for name in ['design-public','both-design-public','fresh-public','bounds-fresh-public','vp9-control']:
 summary[name]=public('govpx-gather-'+name+'.json')
def kernels(name):
 values={};variant=''
 for line in (root/name).read_text().splitlines():
  if line.startswith('ROUND'):variant=line.split()[-1]
  if line.startswith('Benchmark'):
   assert re.search(r'0 B/op\s+0 allocs/op',line)
   values.setdefault(variant,{}).setdefault(line.split()[0],[]).append(float(re.search(r'([\d.]+) ns/op',line).group(1)))
 return {v:{k:statistics.median(a) for k,a in cases.items()} for v,cases in values.items()}
summary['assembly-control']=public('govpx-gather-assembly-control.json','assembly')
summary['final_kernels_ns']=kernels('govpx-gather-bounds-fresh-kernels.txt')
(root/'summary.json').write_text(json.dumps(summary,indent=2)+'\n')
