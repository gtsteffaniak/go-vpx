from pathlib import Path
import json,re,hashlib,statistics
P=Path('/private/tmp/govpx-next-dc-pair-native');strong=Path('/private/tmp/govpx-tm4-47-screen/candidate');V=('parent','candidate','old691')
def hashes(p):return{str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest()for f in p.rglob('*.go')}
source=json.loads((P/'evidence/screen-source-proof.json').read_text())['variants']
for v in V:assert hashes(P/v)==source[v]['all_go_sha256'],v
a,b=hashes(strong),hashes(P/'parent');delta=sorted(k for k in a.keys()|b.keys()if a.get(k)!=b.get(k));assert delta==['internal/vp8/dsp/dc_pair_native_bench_test.go'],delta
assert b[delta[0]]==source['candidate']['all_go_sha256'][delta[0]]==source['old691']['all_go_sha256'][delta[0]]
c=source['candidate']['all_go_sha256'];changes=sorted(k for k in b.keys()|c.keys()if b.get(k)!=c.get(k));prod=[k for k in changes if not k.endswith('_test.go')]
assert prod==['internal/pixel/dc_pair_native_arm64.go','internal/pixel/dc_pair_native_fallback.go','internal/pixel/dc_pair_native_other.go','internal/pixel/dc_pair_simd.go','internal/vp8/dsp/idct.go'],prod
pat=re.compile(r'^BenchmarkDCPairRowsNative/(\S+)\s+(\d+)\s+(\S+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M)
phases={}
for phase,rounds in [('screen',3),('fresh',6)]:
 d=P/phase;data={v:{}for v in V};summary=json.loads((d/'summary.json').read_text());j=json.loads((d/'commands.json').read_text());orders=j['orders'];assert len(orders)==len(set(map(tuple,orders)))==rounds
 for col in range(3):assert sorted(o[col]for o in orders)==sorted(V*(rounds//3))
 assert len(j['measurements'])==3*rounds
 for i,row in enumerate(j['measurements']):
  assert row['exit_code']==0 and row['round']==i//3+1 and row['variant']==orders[i//3][i%3]
  if i:assert j['measurements'][i-1]['end_utc']<=row['start_utc']
 for v in V:
  for r in range(1,rounds+1):
   text=(d/f'r{r}-{v}.txt').read_text();assert re.search(r'^PASS$',text,re.M);assert (d/f'r{r}-{v}.stderr').read_text()==''
   rows=pat.findall(text);assert len(rows)==18
   for k,n,ns,byt,alloc in rows:
    assert int(n)>1 and float(ns)>0 and int(byt)==int(alloc)==0;data[v].setdefault(k,[]).append(float(ns))
 assert len(summary['cases'])==18
 for k,row in summary['cases'].items():
  assert row['ns']=={v:data[v][k]for v in V};med={v:statistics.median(data[v][k])for v in V};assert med==row['median_ns']
  assert med['parent']/med['candidate']==row['parent_over_candidate'];assert med['old691']/med['candidate']==row['old_over_candidate']
  for v,key in [('parent','paired_parent_over_candidate'),('old691','paired_old_over_candidate')]:assert [x/y for x,y in zip(data[v][k],data['candidate'][k])]==row[key]
 vals=list(summary['cases'].values());phases[phase]={'cases':18,'samples':18*3*rounds,'all_zero_alloc_and_scalar_full_backing_oracles_pass':True,'all_median_parent_speedup_range':[min(r['parent_over_candidate']for r in vals),max(r['parent_over_candidate']for r in vals)],'all_paired_parent_speedup_range':[min(min(r['paired_parent_over_candidate'])for r in vals),max(max(r['paired_parent_over_candidate'])for r in vals)],'all_parent_pairs_faster':all(min(r['paired_parent_over_candidate'])>1 for r in vals),'all_median_old_speedup_range':[min(r['old_over_candidate']for r in vals),max(r['old_over_candidate']for r in vals)],'all_old_pairs_slower':all(max(r['paired_old_over_candidate'])<1 for r in vals),'all18_cases':summary['cases']}
assert phases['fresh']['all_parent_pairs_faster'] and phases['fresh']['all_old_pairs_slower']
d=P/'whole300';j=json.loads((d/'commands.json').read_text());orders=j['orders'];assert len(orders)==len(set(map(tuple,orders)))==6 and len(j['measurements'])==18
for col in range(3):assert sorted(o[col]for o in orders)==sorted(V*2)
for i,row in enumerate(j['measurements']):
 assert row['exit_code']==0 and row['round']==i//3+1 and row['variant']==orders[i//3][i%3]
 if i:assert j['measurements'][i-1]['end_utc']<=row['start_utc']
ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'};data={v:[json.loads((d/f'r{r}-{v}.json').read_text())for r in range(1,7)]for v in V};fields=set(data['candidate'][0])-ignore;assert len(fields)==22
reference={k:data['candidate'][0][k]for k in fields};assert all({k:r[k]for k in fields}==reference and r['allocs_per_frame']==0 for a in data.values()for r in a)
assert all((d/f'r{r}-{v}.stderr').read_text()==''for v in V for r in range(1,7))
med={v:statistics.median(r['ns_per_frame']for r in data[v])for v in V};assert med==json.loads((d/'summary.json').read_text())['median_ns_per_frame']
out={'strong_parent_production_exact_frozen_TM4_verified':True,'parent_only_added_shared18case_fixture':delta,'all_actual_source_hashes_verified':True,'production_changes':prod,'kernel_phases':phases,'whole18_all22_fields_equal_zero_alloc_verified':True,'whole_ns':{v:[r['ns_per_frame']for r in data[v]]for v in V},'whole_median_ns':med,'whole_median_time_differences_percent':{v:100*(med['candidate']/med[v]-1)for v in ['parent','old691']},'whole_all6_paired_time_differences_percent':{v:[100*(a['ns_per_frame']/b['ns_per_frame']-1)for a,b in zip(data['candidate'],data[v])]for v in ['parent','old691']},'scope':'M4Max SDK1.27.2 CGO1 SIMD PGOoff GOMAX1. Smaller robust 1.207-1.311x fresh public kernel wins retained by engineering judgment, below mandatory1.5 threshold. Whole mixed/flat, old assembly faster in every kernel and whole pair; no gap-closure claim. In-place timings include identical32Breset. No profiles or bitstream identity claim.'}
Path('/private/tmp/govpx-root-dc-pair-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps({k:v for k,v in out.items()if k!='kernel_phases'},indent=2))
