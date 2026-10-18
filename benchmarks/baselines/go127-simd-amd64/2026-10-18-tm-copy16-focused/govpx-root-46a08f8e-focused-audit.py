from pathlib import Path
import hashlib,json,re,statistics,subprocess,sys
p=Path(sys.argv[1]);head=sys.argv[2];base='78ce74d39806d74ecf91bc3d082899a74d7bc85f';old='691cd0512c4829d9ee40d483205dd223f9eb028d'
for line in (p/'SHA256SUMS').read_text().splitlines():
 h,name=line.split('  ',1);assert hashlib.sha256((p/name).read_bytes()).hexdigest()==h,name
commits=json.loads((p/'commits.json').read_text());assert commits=={'head':head,'base':base,'assembly':old},commits
source_counts={}
for v,c in commits.items():
 lines=(p/(v+'-source-sha256.txt')).read_text().splitlines();source_counts[v]=len(lines)
 for line in lines:
  h,path=line.split('  ',1);blob=subprocess.check_output(['git','show',c+':'+path]);assert hashlib.sha256(blob).hexdigest()==h,(v,path)
expected={f'BenchmarkSSE8RowsNative/{api}/pitch{a}-{b}/pattern{q}'for api in ('slice8x4','slice8x8','slice8x16','pointer8x8') for a,b in ((8,8),(32,37),(1312,1320))for q in range(3)}
legacy={'BenchmarkSSE8x4','BenchmarkSSE8x8','BenchmarkSSE8x16','BenchmarkSSE8x8PtrFast'}
pat=re.compile(r'^(Benchmark\S+)\s+(\d+)\s+(?:([0-9.eE+-]+) ns/op\s+)?(?:[0-9.eE+-]+ MB/s\s+)?(\d+) B/op\s+(\d+) allocs/op$',re.M)
def parse(f,timed=True):
 text=f.read_text();assert re.search(r'^PASS$',text,re.M);a={}
 for n,i,ns,b,alloc in pat.findall(text):
  assert n not in a;assert int(i)>0
  if timed or n in expected:assert int(b)==int(alloc)==0
  if timed:assert ns and float(ns)>0
  a[n]=float(ns)if ns else None
 assert set(a)==expected|legacy,(f,set(a));return a
for v in commits:parse(p/(v+'-one-shot.txt'),False)
a={v:[parse(p/(v+f'-round-{r}.txt')) for r in range(1,7)] for v in commits};out={}
summary=json.loads((p/'summary.json').read_text());reported={r['case']:r for r in summary['results']}
for n in sorted(expected|legacy):
 med={v:statistics.median(x[n]for x in a[v])for v in a};pairs=[b[n]/h[n]for b,h in zip(a['base'],a['head'])]
 for v,val in med.items():assert val==reported[n]['median_ns'][v]
 assert reported[n]['base_over_head']==med['base']/med['head'];assert reported[n]['assembly_over_head']==med['assembly']/med['head'];assert reported[n]['base_head_paired_ratios']==pairs
 out[n]={'median_ns':med,'base_speedup':med['base']/med['head'],'old_assembly_speedup':med['assembly']/med['head'],'all6_paired_speedups':pairs,'all6_faster':min(pairs)>1,'median_ge1_5':med['base']/med['head']>=1.5}
orders=[line.split(': ')[1].split()for line in (p/'orders.txt').read_text().splitlines()];assert len(orders)==6 and len(set(map(tuple,orders)))==6
for k in range(3):assert sorted(x[k]for x in orders)==['assembly','assembly','base','base','head','head']
assert json.loads((p/'status.json').read_text())=={'exit_status':0,'phase':'complete'}
assert '--- PASS: TestBackendSelected ' in (p/'head-dispatch.txt').read_text()
for f in ['head-correctness.txt','head-race-checkptr2.txt']:
 text=(p/f).read_text();assert '--- PASS: TestSSE8EntryDisabledScalar ' in text;assert not re.search(r'^(FAIL|--- FAIL:)',text,re.M)
public=[out[n]for n in expected]
proof={'commits':commits,'manifest_files':len((p/'SHA256SUMS').read_text().splitlines()),'source_counts_verified':source_counts,'native_scope':'See retained environment.txt CPU/flags/toolchain; no whole-codec inference','720_timed_samples_and120_excluded_oracles_verified':True,'all_measured_allocations_zero':True,'all36_median_faster_base':sum(x['base_speedup']>1 for x in public),'all36_median_ge1_5_base':sum(x['median_ge1_5']for x in public),'all36_median_faster_old_ASM':sum(x['old_assembly_speedup']>1 for x in public),'all40_results':out}
Path(sys.argv[3]).write_text(json.dumps(proof,indent=2)+'\n');print(json.dumps({k:v for k,v in proof.items()if k!='all40_results'},indent=2))
for api in ['slice8x4','slice8x8','slice8x16','pointer8x8']:
 vals=[out[n]for n in expected if '/'+api+'/'in n];print(api,'base range',min(x['base_speedup']for x in vals),max(x['base_speedup']for x in vals),'ASM range',min(x['old_assembly_speedup']for x in vals),max(x['old_assembly_speedup']for x in vals))
