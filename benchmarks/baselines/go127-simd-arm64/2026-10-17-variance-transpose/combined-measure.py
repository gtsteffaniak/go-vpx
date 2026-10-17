import os,json,subprocess
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-macro-measure')
wrapper='/Users/thesyncim/.codex/bin/project-env'
names=['base','stats','trn','dup']
orders=[[0,1,2,3],[1,3,0,2],[3,2,1,0],[2,0,3,1],[3,1,0,2],[2,3,1,0],[1,0,2,3],[0,2,3,1]]
with open('/private/tmp/govpx-combined-design-kernels.txt','w') as f:
 for r,order in enumerate(orders[:4]):
  for i in order:
   name=names[i];f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,('/private/tmp/govpx-stats-base.test' if name=='base' else '/private/tmp/govpx-arm-combined-'+name+'.test'),'-test.run=^$','-test.bench=^BenchmarkVariance16x16Native$','-test.benchtime=150ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('macro kernel round',r+1,'done',flush=True)
rows=[]
for r,order in enumerate(orders[:6]):
 for i in order:
  name=names[i];binary='govpx-gather-final-bench' if name=='base' else 'govpx-arm-combined-'+name+'-bench'
  args=[wrapper,'/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(args,env=env,text=True))))
 Path('/private/tmp/govpx-combined-design-public.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
 print('macro public round',r+1,'done',flush=True)
