import os,json,subprocess
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-arm-final-measure');w='/Users/thesyncim/.codex/bin/project-env'
names=['assembly','parent','current'];bins={'assembly':'govpx-asm-bench','parent':'govpx-arm-parent611f-bench','current':'govpx-arm-final-bench'}
orders=[[0,1,2],[1,2,0],[2,0,1],[0,2,1],[2,1,0],[1,0,2]];rows=[]
for codec in ['vp8','vp9']:
 for r,order in enumerate(orders):
  for i in order:
   n=names[i]
   if codec=='vp9' and n=='assembly':continue
   args=[w,'/private/tmp/'+bins[n],'-codec='+codec,'-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
   rows.append(dict(codec=codec,round=r,variant=n,result=json.loads(subprocess.check_output(args,env=env,text=True))))
  Path('/private/tmp/govpx-arm-final-public.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n');print(codec,'final public round',r+1,'done',flush=True)
with open('/private/tmp/govpx-arm-final-stats-kernels.txt','w') as f:
 for r in range(6):
  for n in (['parent','current'] if r%2==0 else ['current','parent']):
   binary='govpx-arm-parent611f-pixel.test' if n=='parent' else 'govpx-arm-final-pixel.test'
   f.write(f'ROUND {r} VARIANT {n}\n');f.flush();subprocess.run([w,'/private/tmp/'+binary,'-test.run=^$','-test.bench=^BenchmarkStats16NarrowNative$','-test.benchtime=150ms','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('stats final kernel round',r+1,'done',flush=True)
