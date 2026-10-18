from pathlib import Path
import json,hashlib,re,statistics,tarfile
p=Path('/private/tmp/govpx-copy8-final-ready');root=Path.cwd()
manifest={}
for line in (p/'SHA256SUMS').read_text().splitlines():
 h,name=line.split('  ',1);assert hashlib.sha256((p/name).read_bytes()).hexdigest()==h,name;manifest[name]=h
inner=json.loads((p/'raw-member-sha256.json').read_text())
with tarfile.open(p/'raw-evidence.tar.gz') as t:
 members=t.getmembers();assert len(members)==len(inner)+1
 assert {m.name for m in members}==set(inner)|{'raw-member-sha256.json'}
 for m in members:
  assert m.isfile()
  h=inner[m.name] if m.name in inner else manifest[m.name]
  assert hashlib.sha256(t.extractfile(m).read()).hexdigest()==h,m.name
sources=json.loads((p/'overlay-sha256.json').read_text())
for n,h in sources.items():assert hashlib.sha256((root/n).read_bytes()).hexdigest()==h,n
pattern=re.compile(r'^(Benchmark\S+)\s+\d+\s+(\S+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M)
allsets={};counts={}
for phase,rounds in [('screen',3),('final',6)]:
 data={}
 for api in ['root','dsp']:
  for v in ['parent','candidate','oldbuild']:
   for r in range(1,rounds+1):
    f=p/phase/f'round{r}-{v}-{api}.log';text=f.read_text();assert re.search(r'^PASS$',text,re.M)
    a=pattern.findall(text);assert len(a)==8
    for n,ns,b,alloc in a:
     assert int(b)==int(alloc)==0;data.setdefault(api+'/'+n,{}).setdefault(v,[]).append(float(ns))
 assert len(data)==16 and all(len(z)==rounds for x in data.values()for z in x.values());allsets[phase]=data;counts[phase]=16*rounds*3
reported=json.loads((p/'final-kernel-summary.json').read_text());results={}
for n,values in allsets['final'].items():
 med={v:statistics.median(x)for v,x in values.items()};pairs=[b/h for b,h in zip(values['parent'],values['candidate'])];assert med==reported[n]['medians_ns'];assert pairs==reported[n]['paired_parent_over_candidate']
 results[n]={'median_ns':med,'parent_speedup':med['parent']/med['candidate'],'old_compiler_go_copy_speedup':med['oldbuild']/med['candidate'],'all6_paired_speedups':pairs};assert min(pairs)>1,n
main=results['root/BenchmarkCopy8RowsNative/pitch8_8/aliasfalse'];assert main['parent_speedup']>=1.5;assert sum(x<1.5 for x in main['all6_paired_speedups'])==2
public={};ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
for frames in ['60','300']:
 rows={v:[json.loads((p/'public'/f'{frames}-round{r}-{v}.json').read_text())for r in range(1,7)]for v in ['parent','candidate','oldbuild']};rr=[x for a in rows.values()for x in a];fields=set(rr[0])-ignore;assert len(fields)==22
 assert all({k:x[k]for k in fields}=={k:rr[0][k]for k in fields} and x['allocs_per_frame']==0 for x in rr)
 med={v:statistics.median(x['ns_per_frame']for x in a)for v,a in rows.items()};public[frames]={'median_ns_per_frame':med,'vs_parent_time_percent':100*(med['candidate']/med['parent']-1),'vs_old_assembly_codec_time_percent':100*(med['candidate']/med['oldbuild']-1),'all6_parent_time_change_percent':[100*(h['ns_per_frame']/b['ns_per_frame']-1)for h,b in zip(rows['candidate'],rows['parent'])],'all22_semantic_fields_and_zero_allocation_samples_match':True}
proof={'parent':'f8dd8c7472e91c1f147511b04ceec716aca95aa7','outer_verified_files':len(manifest),'tar_verified_members':len(members),'root_overlay_exact_measured7_source_sha256':sources,'component_samples_verified':counts,'all_final16_cases_faster6of6_zero_allocations':True,'ordinary_packed_root_qualifies_median_not_all_pairs':True,'component_results':results,'whole_encoder_results':public,'scope':'Apple M4 Max Go1.27.2 simd PGOoff; old copy is compiler Go, complete old codec has original project assembly; no whole gain claim'}
Path('/private/tmp/govpx-root-copy8-proof.json').write_text(json.dumps(proof,indent=2)+'\n');print(json.dumps({k:v for k,v in proof.items()if k not in ['component_results','root_overlay_exact_measured7_source_sha256']},indent=2))
