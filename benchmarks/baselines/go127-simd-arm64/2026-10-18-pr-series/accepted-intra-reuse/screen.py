from pathlib import Path
import os,subprocess,json,re,statistics,hashlib
p=Path('/private/tmp/govpx-frame-recon-next');b=Path('/private/tmp/govpx-frame-recon-base');out=p/'screen';out.mkdir(exist_ok=True);env=os.environ|{'GOMAXPROCS':'1'};refs={'base':b,'candidate':p};rows=[]
for rnd in [1,2]:
 for kind,cpu in [('fixed',-8),('adaptive',8)]:
  for ref in (['base','candidate'] if rnd==1 else ['candidate','base']):
   cmd=[str(refs[ref]/'whole'),'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used='+str(cpu),'-threads=1','-encode-only','-auto-libvpx=false','-format=json'];r=subprocess.run(cmd,env=env,capture_output=True,text=True);(out/f'{kind}-r{rnd}-{ref}.json').write_text(r.stdout);(out/f'{kind}-r{rnd}-{ref}.stderr').write_text(r.stderr);assert r.returncode==0,(r.returncode,r.stderr);v=json.loads(r.stdout);rows.append({'kind':kind,'round':rnd,'ref':ref,'value':v,'cmd':cmd});print(kind,rnd,ref,v['ns_per_frame'],flush=True)
summary={};perf={'ns_per_frame','encode_fps','macroblocks_per_second','allocs_per_frame','latency_ns'}
for kind in ['fixed','adaptive']:
 ratios=[];allocs=[]
 for rnd in [1,2]:
  a=next(x['value']for x in rows if x['kind']==kind and x['round']==rnd and x['ref']=='base');z=next(x['value']for x in rows if x['kind']==kind and x['round']==rnd and x['ref']=='candidate');assert {k:v for k,v in a.items()if k not in perf}=={k:v for k,v in z.items()if k not in perf};ratios.append(a['ns_per_frame']/z['ns_per_frame']);allocs.append([a['allocs_per_frame'],z['allocs_per_frame']])
 summary[kind]={'paired_speedups':ratios,'median_speedup':statistics.median(ratios),'semantic_equal':True,'allocs_pairs':allocs}
components={'base':[],'candidate':[]}
for rnd in [1,2]:
 for ref in (['base','candidate'] if rnd==1 else ['candidate','base']):
  cmd=[str(refs[ref]/'root.test'),'-test.run=^$','-test.bench=^BenchmarkAcceptedWholeIntraReuse$','-test.benchtime=40ms','-test.cpu=1','-test.count=1','-test.benchmem'];r=subprocess.run(cmd,env=env,capture_output=True,text=True);(out/f'component-r{rnd}-{ref}.txt').write_text(r.stdout+r.stderr);assert r.returncode==0;v={}
  for line in r.stdout.splitlines():
   m=re.match(r'(Benchmark\S+)\s+\d+\s+([\d.]+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op',line)
   if m:v[m[1]]={'ns':float(m[2]),'bytes':int(m[3]),'allocs':int(m[4])}
  assert v;components[ref].append(v)
summary['components']={}
for name in components['base'][0]:
 ratio=[components['base'][r][name]['ns']/components['candidate'][r][name]['ns']for r in [0,1]];summary['components'][name]={'paired_speedups':ratio,'median_speedup':statistics.median(ratio)}
summary['preliminary']=True;summary['binding']={ref:{'whole_sha256':hashlib.sha256((v/'whole').read_bytes()).hexdigest(),'test_sha256':hashlib.sha256((v/'root.test').read_bytes()).hexdigest()}for ref,v in refs.items()};(out/'summary.json').write_text(json.dumps(summary,indent=2)+'\n');(out/'rows.json').write_text(json.dumps(rows,indent=2)+'\n');print(json.dumps(summary),flush=True)
