from pathlib import Path
import hashlib,json,subprocess,re,statistics
p=Path('/private/tmp/govpx-arm-fused-dc4-publication-ready')
manifest=json.loads((p/'MANIFEST.json').read_text());assert len(manifest)==2554
for item in manifest:
 f=p/item['path'];assert f.stat().st_size==item['size'];assert hashlib.sha256(f.read_bytes()).hexdigest()==item['sha256'],item['path']
assert {str(f.relative_to(p))for f in p.rglob('*')if f.is_file()}=={i['path']for i in manifest}|{'MANIFEST.json'}
d8=Path('/private/tmp/govpx-arm-fused-dc-publication-ready')
def production(root):
 return {str(f.relative_to(root)):f.read_bytes() for f in root.rglob('*') if f.is_file() and ((f.suffix=='.go' and not f.name.endswith('_test.go'))or f.suffix in ['.s','.S'])}
parent=production(p/'source/parent');old=production(p/'source/old691');candidate=production(p/'source/candidate')
assert parent==production(d8/'source/candidate')
assert old==production(d8/'source/old691')
changed=sorted(n for n in parent if parent[n]!=candidate[n]);assert changed==['internal/vp8/dsp/intra4x4_other.go'],changed
added_production=sorted(set(candidate)-set(parent));assert added_production==['internal/pixel/intra4x4_dc_fill_arm64.go','internal/pixel/intra4x4_dc_fill_other.go']
added_tests=[]
source_counts={'strong_D8_parent_production_byte_identical':len(parent),'old691_production_byte_identical_to_prior_Git_verified_D8_archive':len(old)}
fresh=p
d=fresh/'fresh-dsp';summary=json.loads((d/'summary.json').read_text());orders=summary['orders'];assert len(orders)==len(set(map(tuple,orders)))==6
for col in range(3):assert sorted(o[col]for o in orders)==['candidate','candidate','old691','old691','parent','parent']
expected=set()
for n in [8,16]:
 for pitch in [n,37,1312]:
  for mode in ['none','both']:expected.add('BenchmarkIntraDCFillExact/n%d/%s/pitch%d'%(n,mode,pitch))
  for mode in ['up-only','left-only']:expected.add('BenchmarkIntraDCFillSupplement/n%d/%s/pitch%d'%(n,mode,pitch))
for pitch in [4,37,1312]:expected.add('BenchmarkIntraDCFillSupplement/n4/both/pitch%d'%pitch)
pat=re.compile(r'^(Benchmark\S+)\s+(\d+)\s+([0-9.]+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M)
rows={};reported={r['case']:r for r in summary['cases']};assert set(reported)==expected
for v in ['parent','candidate','old691']:
 rows[v]=[]
 for i in range(1,7):
  text=(d/('r%d-%s.log'%(i,v))).read_text();assert re.search(r'^PASS$',text,re.M);assert 'cpu: Apple M4 Max' in text
  a={}
  for name,it,ns,b,alloc in pat.findall(text):assert int(it)>0 and int(b)==int(alloc)==0;a[name]=float(ns)
  assert set(a)==expected;rows[v].append(a)
results={}
for name in sorted(expected):
 med={v:statistics.median(r[name]for r in rows[v])for v in rows};pairs={v:[a[name]/b[name]for a,b in zip(rows[v],rows['candidate'])]for v in ['parent','old691']}
 r=reported[name];assert med==r['medians_ns'];assert med['parent']/med['candidate']==r['parent_over_candidate'];assert med['old691']/med['candidate']==r['old_over_candidate'];assert pairs['parent']==r['all6paired_parent_ratios'];assert pairs['old691']==r['all6paired_old_ratios']
 results[name]={'median_ns':med,'parent_speedup':med['parent']/med['candidate'],'old_speedup':med['old691']/med['candidate'],'all6_paired_speedups':pairs}
both=[r for n,r in results.items()if n.startswith('BenchmarkIntraDCFillSupplement/n4/')];assert len(both)==3 and all(min(r['all6_paired_speedups']['parent'])>=1.5 for r in both)
commands=json.loads((d/'commands.json').read_text());assert len(commands)==18
assert [[c['variant']for c in commands if c['round']==i]for i in range(1,7)]==orders;assert all(c['exit']==0 for c in commands)
ignored={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'};whole={};times={};w=fresh/'fresh-whole';ws=json.loads((w/'summary.json').read_text())
for v in rows:
 whole[v]=[json.loads((w/('r%d-%s.json'%(i,v))).read_text())for i in range(1,7)];times[v]=[r['ns_per_frame']for r in whole[v]]
allrows=[r for vv in whole.values()for r in vv];fields=set(allrows[0])-ignored;assert len(fields)==22;assert all({k:r[k]for k in fields}=={k:allrows[0][k]for k in fields}and r['allocs_per_frame']==0 for r in allrows)
med={v:statistics.median(t)for v,t in times.items()};assert med==ws['medians_ns_per_frame'];comparison={}
for v in ['parent','old691']:
 pairs=[100*(a/b-1)for a,b in zip(times['candidate'],times[v])];ratio=100*(med['candidate']/med[v]-1);assert ratio==ws['comparisons'][v]['candidate_time_change_percent_ratio_of_medians'];assert pairs==ws['comparisons'][v]['paired_candidate_time_changes_percent'];comparison[v]={'time_change_percent_ratio_of_medians':ratio,'all6_paired_time_changes_percent':pairs}
wc=json.loads((w/'commands.json').read_text());assert len(wc)==18 and all(c['exit']==0 for c in wc);assert [[c['variant']for c in wc if c['round']==i]for i in range(1,7)]==orders
pr=[json.loads((fresh/'profiles600'/(v+'.json')).read_text())for v in rows];pf=set(pr[0])-ignored;assert len(pf)==22 and all({k:r[k]for k in pf}=={k:pr[0][k]for k in pf}and r['allocs_per_frame']==0 for r in pr)
for cmd in wc:
 for flag in ['-frames=300','-codec=vp8','-width=1280','-height=720','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false']:assert flag in cmd['command']
out={'publication_manifest_nonself_files':len(manifest),'all_publication_hashes_verified':True,'source_Git_blob_checks':source_counts,'sole_changed_existing_production':changed,'added_production':added_production,'added_test_sources_in_source_snapshot':added_tests,'fresh486_kernel_samples_allzeroalloc':True,'all6_orders_balanced':True,'all27_results':results,'all3_n4_cases_all6pairs_ge1_5':True,'whole_median_ns':med,'whole_comparisons':comparison,'whole18_all22semantic_fields_equal_andzeroalloc':True,'profiles3_all22semantic_fields_equal_andzeroalloc':True,'scope':'M4Max strongest frozenD8/fusedD4/oldNEON, separate local300-input encode-only window; profile600 diagnostic. No multiplied gains, bitstream identity, or assembly gap closure claim.'}
Path('/private/tmp/govpx-root-d4-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps({k:v for k,v in out.items()if k!='all27_results'},indent=2))
