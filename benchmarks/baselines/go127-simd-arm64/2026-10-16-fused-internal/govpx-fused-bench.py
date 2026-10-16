import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-fused-bench')
wrapper='/Users/thesyncim/.codex/bin/project-env'
variants=[('parent','govpx-reductions-final-vp8.test'),('calls','govpx-fused-calls.test'),('inline','govpx-fused-inline.test')]
orders=[[0,1,2],[1,2,0],[2,0,1],[2,1,0],[0,2,1],[1,0,2]]
with open('/private/tmp/govpx-fused-kernels.txt','w') as f:
 for r,order in enumerate(orders):
  for i in order:
   name,binary=variants[i];f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,'/private/tmp/'+binary,'-test.run=^$','-test.bench=^BenchmarkLoopFilter(Horizontal|Vertical)EdgesY$','-test.benchtime=150ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('fused kernel round',r+1,'done',flush=True)
variants=[('assembly','govpx-asm-bench'),('parent','govpx-reductions-scoped-bench'),('calls','govpx-fused-calls-bench'),('inline','govpx-fused-inline-bench')]
orders=[[(i+s)%4 for i in [0,1,3,2]] for s in range(4)];orders += [list(reversed(o)) for o in orders]
rows=[]
for r,order in enumerate(orders):
 for i in order:
  name,binary=variants[i]
  cmd=[wrapper,'/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-fused-public.json').write_text(json.dumps(dict(orders=[[variants[i][0] for i in order] for order in orders],samples=rows),indent=2)+'\n')
 print('fused public round',r+1,'done',flush=True)
