import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-internal-fresh')
wrapper='/Users/thesyncim/.codex/bin/project-env'
variants=[('parent','govpx-bench-perf34'),('prototype','govpx-internal-store-proof-bench'),('final','govpx-internal-final-bench')]
orders=[[0,1,2],[1,2,0],[2,0,1],[2,1,0],[0,2,1],[1,0,2]]
rows=[]
for r,order in enumerate(orders):
 for i in order:
  name,binary=variants[i]
  cmd=[wrapper,'/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-internal-fresh-public.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
 print('fresh public round',r+1,'done',flush=True)
rows=[]
for r in range(6):
 for i in ([0,2] if r%2==0 else [2,0]):
  name,binary=variants[i]
  cmd=[wrapper,'/private/tmp/'+binary,'-codec=vp9','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-internal-vp9-control.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
 print('vp9 control round',r+1,'done',flush=True)

with open('/private/tmp/govpx-internal-final-kernels-fresh.txt','w') as f:
 for r in range(6):
  for name,binary in ([('parent','govpx-internal-before.test'),('final','govpx-internal-final.test')] if r%2==0 else [('final','govpx-internal-final.test'),('parent','govpx-internal-before.test')]):
   f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,'/private/tmp/'+binary,'-test.run=^$','-test.bench=BenchmarkLoopFilter(Horizontal|Vertical)EdgesY$','-test.benchtime=150ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('fresh kernel round',r+1,'done',flush=True)
