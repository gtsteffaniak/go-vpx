import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-next-public')
wrapper='/Users/thesyncim/.codex/bin/project-env'
variants=[('assembly','govpx-asm-bench'),('parent','govpx-projection-parent-bench'),('projection','govpx-projection-bench'),('average','govpx-average-bench'),('combined','govpx-projection-average-bench')]
orders=[[(i+s)%5 for i in [0,1,4,2,3]] for s in range(5)];orders += [list(reversed(o)) for o in orders]
rows=[]
for r,order in enumerate(orders):
 for i in order:
  name,binary=variants[i]
  cmd=[wrapper,'/private/tmp/'+binary,'-codec=vp9','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-projection-average-public.json').write_text(json.dumps(dict(orders=[[variants[i][0] for i in order] for order in orders],samples=rows),indent=2)+'\n')
 print('projection public round',r+1,'done',flush=True)
variants=[('assembly','govpx-asm-bench'),('parent','govpx-internal-final-bench'),('rolling','govpx-six-roll-bench')]
orders=[[0,1,2],[1,2,0],[2,0,1],[2,1,0],[0,2,1],[1,0,2]];rows=[]
for r,order in enumerate(orders):
 for i in order:
  name,binary=variants[i]
  cmd=[wrapper,'/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-six-roll-public.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
 print('rolling public round',r+1,'done',flush=True)
