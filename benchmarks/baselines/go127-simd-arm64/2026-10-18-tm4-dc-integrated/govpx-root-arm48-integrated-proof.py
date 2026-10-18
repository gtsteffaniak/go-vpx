from pathlib import Path
import hashlib,json,re,statistics,subprocess,io
P=Path('/private/tmp/govpx-arm48-integrated');ROOT=Path('/Users/thesyncim/GolandProjects/govpx');V=('parent','candidate','old691');head='d2610cd0747edd72f7656068dfd4604d29e45307';old='691cd0512c4829d9ee40d483205dd223f9eb028d'
def hashes(p):return{str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest()for f in p.rglob('*.go')}
proof=json.loads((P/'evidence/source-proof.json').read_text());variants=proof['variants']
for v in V:assert hashes(P/v)==variants[v]['all_go_sha256'],v
root_go=hashes(ROOT);assert hashes(P/'candidate')==root_go,'Frozen candidate all-Go differs from integrated ROOT'
sourcecounts={}
for v,c in [('parent',head),('old691',old)]:
 listing=subprocess.check_output(['git','ls-tree','-r',c]).splitlines();entries=[]
 for line in listing:
  meta,path=line.split(b'\t',1);mode,typ,oid=meta.split();assert typ==b'blob';entries.append((path.decode(),oid))
 cmd=subprocess.Popen(['git','cat-file','--batch'],stdin=subprocess.PIPE,stdout=subprocess.PIPE);out,_=cmd.communicate(b''.join(oid+b'\n'for path,oid in entries));bio=io.BytesIO(out);diff=[]
 for name,oid in entries:
  header=bio.readline().split();assert header[:2]==[oid,b'blob'];n=int(header[2]);b=bio.read(n);assert bio.read(1)==b'\n'
  if (P/v/name).read_bytes()!=b:diff.append(name)
 assert not bio.read()
 # All reference production stays pinned; only shared standalone public
 # benchmark fixtures may replace tests on the older era.
 assert all(k.endswith('_test.go') for k in diff),diff
 sourcecounts[v]={'tracked_files':len(entries),'only_changed_tracked_files':diff}
a=variants['parent']['all_go_sha256'];b=variants['candidate']['all_go_sha256'];changes=sorted(k for k in a.keys()|b.keys()if a.get(k)!=b.get(k));prod=[k for k in changes if not k.endswith('_test.go')];assert len(prod)==13
hashes_bin=json.loads((P/'evidence/binary-hashes.json').read_text())
assert set(hashes_bin)=={v+'-'+kind+('.test' if kind!='whole' else '-bench') for v in V for kind in ('dsp','decoder','whole')}
for name,row in hashes_bin.items():
 f=P/name;assert f.stat().st_size==row['size'] and hashlib.sha256(f.read_bytes()).hexdigest()==row['sha256'],name
pat=re.compile(r'^(Benchmark\S+)\s+(\d+)\s+(\S+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M);data={v:{}for v in V};D=P/'fresh';summary=json.loads((D/'summary.json').read_text());j=json.loads((D/'commands.json').read_text());orders=j['orders'];assert len(orders)==len(set(map(tuple,orders)))==6 and len(j['measurements'])==36
for col in range(3):assert sorted(o[col]for o in orders)==sorted(V*2)
for i,row in enumerate(j['measurements']):
 assert row['exit_code']==0 and row['round']==i//6+1 and row['variant']==orders[i//6][(i%6)//2] and row['kind']==('dsp'if i%2==0 else'decoder')
 if i:assert j['measurements'][i-1]['end_utc']<=row['start_utc']
for v in V:
 for rnd in range(1,7):
  for kind,count in [('dsp',45),('decoder',10)]:
   text=(D/f'r{rnd}-{v}-{kind}.txt').read_text();assert re.search(r'^PASS$',text,re.M);assert (D/f'r{rnd}-{v}-{kind}.stderr').read_text()=='';rows=pat.findall(text);assert len(rows)==count
   for k,n,ns,byt,alloc in rows:
    assert int(n)>1 and float(ns)>0 and int(byt)==int(alloc)==0;data[v].setdefault(k,[]).append(float(ns))
assert len(summary['cases'])==55
for k,row in summary['cases'].items():
 assert row['ns']=={v:data[v][k]for v in V};med={v:statistics.median(data[v][k])for v in V};assert med==row['median_ns'];assert med['parent']/med['candidate']==row['parent_over_candidate'] and med['old691']/med['candidate']==row['old_over_candidate']
 for v,key in [('parent','paired_parent_over_candidate'),('old691','paired_old_over_candidate')]:assert [x/y for x,y in zip(data[v][k],data['candidate'][k])]==row[key]
TM4=[r for k,r in summary['cases'].items()if k.startswith('BenchmarkTMPredictRowsNative/4x4/')];assert len(TM4)==9 and all(min(r['paired_parent_over_candidate'])>=1.5 and min(r['paired_old_over_candidate'])>1 for r in TM4)
D=P/'whole300';j=json.loads((D/'commands.json').read_text());assert j['orders']==orders and len(j['measurements'])==18
for i,row in enumerate(j['measurements']):
 assert row['exit_code']==0 and row['round']==i//3+1 and row['variant']==orders[i//3][i%3]
 if i:assert j['measurements'][i-1]['end_utc']<=row['start_utc']
ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'};w={v:[json.loads((D/f'r{r}-{v}.json').read_text())for r in range(1,7)]for v in V};fields=set(w['candidate'][0])-ignore;assert len(fields)==22;reference={k:w['candidate'][0][k]for k in fields}
assert all({k:r[k]for k in fields}==reference and r['allocs_per_frame']==0 for a in w.values()for r in a)
assert all((D/f'r{r}-{v}.stderr').read_text()==''for v in V for r in range(1,7))
med={v:statistics.median(r['ns_per_frame']for r in w[v])for v in V}
def group(prefix):
 a=[r for k,r in summary['cases'].items()if k.startswith(prefix)];return {'case_count':len(a),'median_parent_speedup_range':[min(r['parent_over_candidate']for r in a),max(r['parent_over_candidate']for r in a)],'all_parent_pair_speedup_range':[min(min(r['paired_parent_over_candidate'])for r in a),max(max(r['paired_parent_over_candidate'])for r in a)],'candidate_faster_vs_parent_pairs':sum(x>1 for r in a for x in r['paired_parent_over_candidate']),'median_old_speedup_range':[min(r['old_over_candidate']for r in a),max(r['old_over_candidate']for r in a)],'all_old_pair_speedup_range':[min(min(r['paired_old_over_candidate'])for r in a),max(max(r['paired_old_over_candidate'])for r in a)],'candidate_faster_vs_old_pairs':sum(x>1 for r in a for x in r['paired_old_over_candidate'])}
out={'base':head,'old':old,'hardware':'Apple M4 Max','officialSDK':'1.27.2','CGO':1,'PGO':'off','GOMAXPROCS':1,'exact_reference_tracked_source_files':sourcecounts,'all13changed_production_Go':prod,'all_candidate_Go_snapshot_bytes_equal_ROOT':True,'all9uniform_executable_hashes_verified':True,'all990raw_samples_all55_medians_pairs_oracles_zero_alloc_verified':True,'TM4':group('BenchmarkTMPredictRowsNative/4x4/'),'TM8_controls':group('BenchmarkTMPredictRowsNative/8x8/'),'TM16_controls':group('BenchmarkTMPredictRowsNative/16x16/'),'DC':group('BenchmarkDCPairRowsNative/'),'picker':group('BenchmarkVP8PickerFilterNative/'),'whole18_all22_fields_equal_zero_alloc_verified':True,'whole_all_ns':{v:[r['ns_per_frame']for r in w[v]]for v in V},'whole_median_ns':med,'whole_time_differences_percent':{v:100*(med['candidate']/med[v]-1)for v in ['parent','old691']},'whole_all6_paired_time_differences_percent':{v:[100*(a['ns_per_frame']/b['ns_per_frame']-1)for a,b in zip(w['candidate'],w[v])]for v in ['parent','old691']},'scope':'Exact integrated candidate, separate from earlier experiments; allcontrolregressions and sloweroldDC retained. NativeTM4>=1.5 mandatorykeep irrespectivewholecodec. No bitstream identity claim. No pooling acrosswindows/CPUs or inference from multiplied stage ratios.'}
Path('/private/tmp/govpx-root-arm48-integrated-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
