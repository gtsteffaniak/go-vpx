from pathlib import Path
import json,re,hashlib,statistics
p=Path('/private/tmp/govpx-lf-axis-private');strong=Path('/private/tmp/govpx-tm4-47-screen/candidate');V=('parent','candidate','old691')
def hashes(d):return {str(f.relative_to(d)):hashlib.sha256(f.read_bytes()).hexdigest()for f in d.rglob('*.go')}
source=json.loads((p/'evidence/source-proof.json').read_text())
for v in V:assert hashes(p/v)==source[v]['all_go_sha256'],v
parent_all=hashes(p/'parent');strong_all=hashes(strong)
parent_extra=sorted(k for k in parent_all.keys()|strong_all.keys()if parent_all.get(k)!=strong_all.get(k))
assert parent_extra==['internal/vp8/decoder/loopfilter_axis_frozen_test.go'],parent_extra
assert parent_all[parent_extra[0]]==hashes(p/'candidate')[parent_extra[0]]
assert {k:v for k,v in parent_all.items()if not k.endswith('_test.go')}=={k:v for k,v in strong_all.items()if not k.endswith('_test.go')}
a=source['parent']['all_go_sha256'];b=source['candidate']['all_go_sha256'];changed=sorted(k for k in a.keys()|b.keys()if a.get(k)!=b.get(k))
prod=[k for k in changed if not k.endswith('_test.go')]
assert prod==['internal/pixel/loopfilter_macro_axis_arm64.go','internal/pixel/loopfilter_macro_axis_other.go','internal/vp8/decoder/loopfilter.go'],prod
pat=re.compile(r'^(BenchmarkVP8PickerFilterNative/\S+?)(?:-\d+)?\s+(\d+)\s+([\d.]+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op',re.M)
phases={}
for phase in ['screen','fresh-screen']:
 d=p/phase;j=json.loads((d/'run-journal.json').read_text());orders=j['orders'];assert len(orders)==len(set(map(tuple,orders)))==6;assert len(j['records'])==18
 for c in range(3):assert sorted(o[c]for o in orders)==sorted(V*2)
 data={v:{}for v in V}
 for i,row in enumerate(j['records']):
  assert row['returncode']==0 and row['round']==i//3+1 and row['position']==i%3+1 and row['variant']==orders[i//3][i%3]
  if i:assert j['records'][i-1]['ended_unix']<=row['started_unix']
  s=(d/row['log']).read_text();assert '\nPASS\n' in s;rows=pat.findall(s);assert len(rows)==10
  for k,n,ns,byt,alloc in rows:
   assert int(n)>1 and float(ns)>0 and int(byt)==int(alloc)==0;data[row['variant']].setdefault(k,[]).append(float(ns))
 cases={}
 for k in data['candidate']:
  vals={v:data[v][k]for v in V};assert all(len(a)==6 for a in vals.values());med={v:statistics.median(x)for v,x in vals.items()}
  cases[k]={'ns':vals,'median_ns':med,'parent_speedup':med['parent']/med['candidate'],'old_speedup':med['old691']/med['candidate'],'paired_parent_speedups':[a/b for a,b in zip(vals['parent'],vals['candidate'])],'paired_old_speedups':[a/b for a,b in zip(vals['old691'],vals['candidate'])]}
 phases[phase]={'benchtime':j['benchtime'],'all180timings_zero_alloc_full_output_oracles_pass':True,'all10results':cases}
ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
whole={}
for phase in ['whole300','fresh-whole300']:
 d=p/phase;j=json.loads((d/'commands.json').read_text());orders=j['orders'];assert len(orders)==len(set(map(tuple,orders)))==6
 for c in range(3):assert sorted(o[c]for o in orders)==sorted(V*2)
 assert len(j['measurements'])==18
 for i,row in enumerate(j['measurements']):
  assert row['exit_code']==0 and row['round']==i//3+1 and row['variant']==orders[i//3][i%3]
  if i:assert j['measurements'][i-1]['end_utc']<=row['start_utc']
 data={v:[json.loads((d/f'r{r}-{v}.json').read_text())for r in range(1,7)]for v in V};fields=set(data['candidate'][0])-ignore;assert len(fields)==22
 reference={k:data['candidate'][0][k]for k in fields}
 assert all({k:r[k]for k in fields}==reference and r['allocs_per_frame']==0 for a in data.values()for r in a)
 assert all((d/f'r{r}-{v}.stderr').read_text()=='' for v in V for r in range(1,7))
 med={v:statistics.median(r['ns_per_frame']for r in data[v])for v in V}
 whole[phase]={'all18timed_runs22_fields_equal_zero_alloc':True,'ns_per_frame':{v:[r['ns_per_frame']for r in data[v]]for v in V},'median_ns_per_frame':med,'median_time_difference_percent':{v:100*(med['candidate']/med[v]-1)for v in ['parent','old691']},'all6pair_time_differences_percent':{v:[100*(a['ns_per_frame']/b['ns_per_frame']-1)for a,b in zip(data['candidate'],data[v])]for v in ['parent','old691']}}
out={'strong_parent_all_production_Go_sources_exact_frozen_TM4_verified':True,'parent_only_added_common_frozen_oracle':parent_extra,'production_changes':prod,'all_candidate_source_hashes_verified':True,'kernel_phases':phases,'whole_phases':whole,'scope':'M4 Max, official Go1.27.2 CGO1 PGOoff SIMD GOMAXPROCS1, separate timing windows. No whole/old gap inference from kernel ratios or separate windows. No bitstream identity claim; all regressions retained.'}
Path('/private/tmp/govpx-root-lf-axis-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps({'production_changes':prod,'whole':whole},indent=2))
