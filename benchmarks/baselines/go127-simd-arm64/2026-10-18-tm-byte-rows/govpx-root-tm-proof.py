from pathlib import Path
import hashlib,json,re,statistics,subprocess,datetime
P=Path('/private/tmp/govpx-next-tm-byte-arm64');HEAD='ac259fce071f3ac5fc45603c6bd7e7926bdfcf73';OLD='691cd0512c4829d9ee40d483205dd223f9eb028d';V=('parent','candidate','old691')
def git_blob(data):return hashlib.sha1(b'blob '+str(len(data)).encode()+b'\0'+data).hexdigest()
def tree(commit):
 out={}
 for l in subprocess.check_output(['git','ls-tree','-r',commit]).splitlines():
  m,p=l.split(b'\t',1);out[p.decode()]=m.split()[2].decode()
 return out
t=tree(HEAD);cnt=0
for f,h in t.items():
 if (f.endswith('.go') and not f.endswith('_test.go'))or f in ('go.mod','go.sum'):
  assert git_blob((P/'parent'/f).read_bytes())==h,f;cnt+=1
assert cnt==784,cnt
prod=lambda d:{str(x.relative_to(d)):hashlib.sha256(x.read_bytes()).hexdigest() for x in d.rglob('*.go') if not x.name.endswith('_test.go')}
a,b=prod(P/'parent'),prod(P/'candidate');delta=sorted(k for k in a.keys()|b.keys() if a.get(k)!=b.get(k));assert delta==['internal/pixel/kernels_simd.go','internal/pixel/tmpredict_bytes_arm64.go','internal/pixel/tmpredict_bytes_other.go'],delta
# Existing module and all README bytes remain exact; no dependencies added.
for f in t:
 if f in ('go.mod','go.sum') or (Path(f).name.lower().startswith('readme') and ((P/'parent'/f).exists() or (P/'candidate'/f).exists())):
  assert (P/'parent'/f).read_bytes()==(P/'candidate'/f).read_bytes(),f
pat=re.compile(r'^BenchmarkTMPredictRowsNative/(\S+)\s+(\d+)\s+(\S+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M)
summary=json.loads((P/'fresh-dsp/summary.json').read_text());raw={v:{}for v in V}
for v in V:
 for r in range(1,7):
  s=(P/'fresh-dsp'/('r%d-%s.txt'%(r,v))).read_text();assert re.search(r'^PASS$',s,re.M)
  vals=pat.findall(s);assert len(vals)==27
  for k,n,ns,by,al in vals:assert int(n)>0 and int(by)==int(al)==0;raw[v].setdefault(k,[]).append(float(ns))
assert all(raw[v].keys()==raw['candidate'].keys()for v in V)
for k in raw['candidate']:
 x=summary['cases'][k];med={v:statistics.median(raw[v][k])for v in V}
 assert med==x['median_ns'] and all(raw[v][k]==x['ns'][v]for v in V)
 assert med['parent']/med['candidate']==x['parent_over_candidate'];assert med['old691']/med['candidate']==x['old_over_candidate']
 assert [p/c for p,c in zip(raw['parent'][k],raw['candidate'][k])]==x['paired_parent_over_candidate']
 assert [p/c for p,c in zip(raw['old691'][k],raw['candidate'][k])]==x['paired_old_over_candidate']
 if k.startswith(('8x8/','16x16/')):assert min(x['paired_parent_over_candidate'])>=1.5,k
for phase in ('fresh-dsp','fresh-whole'):
 c=json.loads((P/phase/'commands.json').read_text());orders=c['orders'];assert len(orders)==len(set(map(tuple,orders)))==6
 for col in range(3):assert sorted(o[col]for o in orders)==sorted(V*2)
 assert len(c['measurements'])==18
 for i,r in enumerate(c['measurements']):
  assert r['exit_code']==0 and r['round']==i//3+1 and r['variant']==orders[i//3][i%3]
  if i:assert c['measurements'][i-1]['end_utc']<=r['start_utc']
ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
w={v:[json.loads((P/'fresh-whole'/('r%d-%s.json'%(r,v))).read_text())for r in range(1,7)]for v in V}
rr=[x for a in w.values()for x in a];fields=set(rr[0])-ignore;assert len(fields)==22
assert all({k:r[k]for k in fields}=={k:rr[0][k]for k in fields}and r['allocs_per_frame']==0 for r in rr)
med={v:statistics.median(x['ns_per_frame']for x in w[v])for v in V};reported=json.loads((P/'fresh-whole/summary.json').read_text());assert med==reported['median_ns_per_frame']
pairs={v:[100*(a['ns_per_frame']/b['ns_per_frame']-1)for a,b in zip(w['candidate'],w[v])]for v in ('parent','old691')}
profiles={v:json.loads((P/'profiles600'/(v+'.json')).read_text())for v in V};pf=set(profiles['candidate'])-ignore;assert len(pf)==22
assert all({k:r[k]for k in pf}=={k:profiles['candidate'][k]for k in pf}and r['allocs_per_frame']==0 for r in profiles.values())
def group(prefix):
 a=[x for k,x in summary['cases'].items()if k.startswith(prefix)];return {'median_speedup_range':[min(x['parent_over_candidate']for x in a),max(x['parent_over_candidate']for x in a)],'all_paired_speedup_range':[min(min(x['paired_parent_over_candidate'])for x in a),max(max(x['paired_parent_over_candidate'])for x in a)],'old_over_current_range':[min(x['old_over_candidate']for x in a),max(x['old_over_candidate']for x in a)]}
out={'base':HEAD,'parent_Git_production_modules_verified':cnt,'only3_production_changes':delta,'fresh486_raw_samples_and27_all6_reported_medians_pairs_verified':True,'all18_target_cases_all6_pairs_at_least1_5':True,'TM8':group('8x8/'),'TM16':group('16x16/'),'TM4_controls':group('4x4/'),'whole18_all22_semantics_and0alloc_exact':True,'whole_median_ns_per_input':med,'whole_parent_time_difference_percent':100*(med['candidate']/med['parent']-1),'whole_old_time_difference_percent':100*(med['candidate']/med['old691']-1),'whole_all6_paired_time_difference_percent':pairs,'profiles3_all22_semantics_and0alloc_exact':True,'scope':'Kernel retention qualifies independently; whole flat/mixed, old still faster. Initial mixed-CGO243screen separate and not acceptance; no cross-window multiplication or gap-attribution. Profiles diagnostic only.'}
Path('/private/tmp/govpx-root-tm-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
