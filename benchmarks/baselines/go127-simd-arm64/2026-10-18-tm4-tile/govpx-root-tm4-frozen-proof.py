from pathlib import Path
import hashlib,json,re,statistics,subprocess,tarfile,io
P=Path('/private/tmp/govpx-tm4-47-screen'); F=Path('/private/tmp/govpx-arm-tm4-tile-ready')
head='b6209532fbc30cf9567e0cdc78a21d61a60fe900';old='691cd0512c4829d9ee40d483205dd223f9eb028d';V=('parent','candidate','old691')
def sha(b):return hashlib.sha256(b).hexdigest()
outer=json.loads((F/'sha256.json').read_text())
for name,row in outer['files'].items():
 b=(F/name).read_bytes();assert len(b)==row['size'] and sha(b)==row['sha256'],name
raw=json.loads((F/'raw-sha256.json').read_text())
with tarfile.open(F/'raw.tar.gz') as t:
 members={m.name:m for m in t if m.isfile()}
 assert len(members)==raw['member_count']==292
 assert members.keys()==raw['members'].keys()
 for name,row in raw['members'].items():
  b=t.extractfile(members[name]).read();assert len(b)==row['size'] and sha(b)==row['sha256'],name
source=json.loads((P/'evidence/source-proof.json').read_text());overlay=source['overlay']
source_verified={}
for v,c in [('parent',head),('old691',old)]:
 listing=subprocess.check_output(['git','ls-tree','-r','--full-tree',c]).splitlines()
 entries=[]
 for line in listing:
  meta,path=line.split(b'\t',1);mode,typ,oid=meta.split();assert typ==b'blob';entries.append((path.decode(),oid))
 cmd=subprocess.Popen(['git','cat-file','--batch'],stdin=subprocess.PIPE,stdout=subprocess.PIPE)
 out,_=cmd.communicate(b''.join(oid+b'\n' for path,oid in entries));bio=io.BytesIO(out)
 for name,oid in entries:
  header=bio.readline().split();assert header[:2]==[oid,b'blob'];n=int(header[2]);blob=bio.read(n);assert bio.read(1)==b'\n'
  assert (P/v/name).read_bytes()==blob,(v,name)
 assert not bio.read();source_verified[v]=len(entries)
 assert len(entries)==source[v]['tracked_files']
parent_files={str(f.relative_to(P/'parent')):sha(f.read_bytes()) for f in (P/'parent').rglob('*.go')}
candidate_files={str(f.relative_to(P/'candidate')):sha(f.read_bytes()) for f in (P/'candidate').rglob('*.go')}
changed=sorted(k for k in parent_files.keys()|candidate_files.keys() if parent_files.get(k)!=candidate_files.get(k))
assert changed==sorted(overlay),changed
for k,h in overlay.items():assert candidate_files[k]==h and sha((F/'overlay'/k).read_bytes())==h,k
assert len([x for x in changed if not x.endswith('_test.go')])==5
pat=re.compile(r'^BenchmarkTMPredictRowsNative/(\S+)\s+(\d+)\s+(\S+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M)
summaries={}
for phase,nrounds in [('screen',3),('fresh',6)]:
 summary=json.loads((P/phase/'summary.json').read_text());data={v:{} for v in V}
 for v in V:
  for r in range(1,nrounds+1):
   text=(P/phase/f'r{r}-{v}.txt').read_text();assert re.search(r'^PASS$',text,re.M);rows=pat.findall(text);assert len(rows)==27
   assert (P/phase/f'r{r}-{v}.stderr').read_text()==''
   for k,n,ns,b,a in rows:
    assert int(n)>0 and float(ns)>0 and int(b)==int(a)==0;data[v].setdefault(k,[]).append(float(ns))
 for k,row in summary['cases'].items():
  assert {v:data[v][k] for v in V}==row['ns']
  med={v:statistics.median(data[v][k])for v in V};assert med==row['median_ns']
  assert med['parent']/med['candidate']==row['parent_over_candidate']
  assert med['old691']/med['candidate']==row['old_over_candidate']
  for v,key in [('parent','paired_parent_over_candidate'),('old691','paired_old_over_candidate')]:
   assert [a/b for a,b in zip(data[v][k],data['candidate'][k])]==row[key]
 assert len(summary['cases'])==27
 summaries[phase]=summary
for phase,rounds in [('screen',3),('fresh',6),('whole300',6)]:
 commands=json.loads((P/phase/'commands.json').read_text());orders=commands['orders'];assert len(orders)==rounds and len(set(map(tuple,orders)))==rounds
 for col in range(3):assert sorted(o[col]for o in orders)==sorted(V*(rounds//3))
 assert len(commands['measurements'])==3*rounds
 for i,row in enumerate(commands['measurements']):
  assert row['exit_code']==0 and row['round']==i//3+1 and row['variant']==orders[i//3][i%3]
  if i:assert commands['measurements'][i-1]['end_utc']<=row['start_utc']
ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
w={v:[json.loads((P/'whole300'/f'r{r}-{v}.json').read_text())for r in range(1,7)]for v in V}
fields=set(w['candidate'][0])-ignore;assert len(fields)==22
reference={k:w['candidate'][0][k]for k in fields}
for v in V:
 for row in w[v]:assert {k:row[k]for k in fields}==reference and row['allocs_per_frame']==0
profiles={v:json.loads((P/'profiles600'/(v+'.json')).read_text())for v in V}
pfields=set(profiles['candidate'])-ignore;assert len(pfields)==22
assert all({k:r[k]for k in pfields}=={k:profiles['candidate'][k]for k in pfields} and r['allocs_per_frame']==0 for r in profiles.values())
def group(prefix):
 cases=[r for k,r in summaries['fresh']['cases'].items()if k.startswith(prefix)]
 return {'median_parent_speedup_range':[min(r['parent_over_candidate']for r in cases),max(r['parent_over_candidate']for r in cases)],'all54_parent_pair_range':[min(min(r['paired_parent_over_candidate'])for r in cases),max(max(r['paired_parent_over_candidate'])for r in cases)],'median_old_speedup_range':[min(r['old_over_candidate']for r in cases),max(r['old_over_candidate']for r in cases)],'parent_faster_pairs':sum(x>1 for r in cases for x in r['paired_parent_over_candidate']),'old_faster_pairs':sum(x>1 for r in cases for x in r['paired_old_over_candidate'])}
TM4=group('4x4/');assert TM4['all54_parent_pair_range'][0]>=1.5
med={v:statistics.median(r['ns_per_frame']for r in w[v])for v in V}
out={'head':head,'old':old,'hardware':'Apple M4 Max','outer_files_verified':len(outer['files']),'raw_members_verified':len(members),'exact_Git_tracked_source_files_verified':source_verified,'candidate_only6_Go_changes_verified':changed,'screen243_and_fresh486_all_raw_medians_pairs_zero_alloc_verified':True,'balanced_orders_all_commands_success_verified':True,'TM4':TM4,'TM8_controls':group('8x8/'),'TM16_controls':group('16x16/'),'whole18_all22_fields_zero_alloc_verified':True,'whole_median_ns_per_frame':med,'whole_time_difference_percent':{v:100*(med['candidate']/med[v]-1)for v in ('parent','old691')},'whole_all6_pair_time_differences_percent':{v:[100*(a['ns_per_frame']/b['ns_per_frame']-1)for a,b in zip(w['candidate'],w[v])]for v in ('parent','old691')},'profiles3_all22_fields_zero_alloc_verified':True,'acceptance':'Mandatory retained >=1.5x native kernel gain against exact47. Whole codec flat; old reference slowed across windows, no causal gap closure claim. All control regressions and mixed old comparisons retained. No bitstream identity claim.'}
Path('/private/tmp/govpx-root-tm4-frozen-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
