import os,subprocess,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-gather-measure')
wrapper='/Users/thesyncim/.codex/bin/project-env'
names=['base','pairs','quads'];orders=[[0,1,2],[1,2,0],[2,0,1],[2,1,0],[0,2,1],[1,0,2]]
with open('/private/tmp/govpx-gather-both-design-kernels.txt','w') as f:
 for r,order in enumerate(orders):
  for i in order:
   name=names[i];f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,'/private/tmp/govpx-gather-'+name+'.test','-test.run=^$','-test.bench=^BenchmarkGatherBlockLayout$','-test.benchtime=200ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('gather kernel round',r+1,'done',flush=True)
rows=[]
for r,order in enumerate(orders):
 for i in order:
  name=names[i];binary='govpx-fused-final-bench' if name=='base' else 'govpx-gather-'+name+'-bench'
  args=[wrapper,'/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(args,env=env,text=True))))
 Path('/private/tmp/govpx-gather-both-design-public.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
 print('gather public round',r+1,'done',flush=True)
