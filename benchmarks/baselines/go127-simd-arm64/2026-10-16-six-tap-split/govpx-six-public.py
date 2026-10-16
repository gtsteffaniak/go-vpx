import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-six-public')
wrapper='/Users/thesyncim/.codex/bin/project-env'
variants=[('assembly','govpx-asm-bench'),('parent','govpx-bench-perf33'),('six','govpx-sixhybrid-bench'),('split','govpx-split-subpel-bench'),('combined','govpx-six-combined-bench')]
orders=[[(i+s)%5 for i in [0,1,4,2,3]] for s in range(5)]
orders += [list(reversed(o)) for o in orders]
rows=[]
for r,order in enumerate(orders):
 for idx in order:
  name,binary=variants[idx]
  cmd=[wrapper,'/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-six-public-ab.json').write_text(json.dumps(dict(orders=[[variants[i][0] for i in order] for order in orders],samples=rows),indent=2)+'\n')
 print('public round',r+1,'done',flush=True)
