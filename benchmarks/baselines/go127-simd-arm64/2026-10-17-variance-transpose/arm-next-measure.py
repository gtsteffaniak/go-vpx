import os,json,subprocess,statistics
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-arm-next-measure');w='/Users/thesyncim/.codex/bin/project-env'
orders=[[0,1,2,3,4],[1,3,0,4,2],[4,2,1,0,3],[2,0,4,3,1],[3,4,2,1,0],[4,3,1,2,0]]
names=['base','trn','dup','fpdup','denoise'];rows=[]
for r,order in enumerate(orders):
 for i in order:
  n=names[i];b='govpx-gather-final-bench' if n=='base' else 'govpx-denoise-bytes-bench' if n=='denoise' else 'govpx-columns-'+n+'-bench'
  args=[w,'/private/tmp/'+b,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=n,result=json.loads(subprocess.check_output(args,env=env,text=True))))
 Path('/private/tmp/govpx-arm-next-public.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n');print('arm next public round',r+1,'done',flush=True)
with open('/private/tmp/govpx-denoise-stats-kernels.txt','w') as f:
 for r in range(6):
  names=['base','bytes'] if r%2==0 else ['bytes','base']
  for n in names:
   f.write(f'ROUND {r} DENOISE {n}\n');f.flush();subprocess.run([w,'/private/tmp/govpx-denoise-'+n+'.test','-test.run=^$','-test.bench=^BenchmarkDenoiseBytes16$','-test.benchtime=150ms','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  names=['base','safe','ptr'] if r%2==0 else ['ptr','safe','base']
  for n in names:
   f.write(f'ROUND {r} STATS {n}\n');f.flush();subprocess.run([w,'/private/tmp/govpx-stats-'+n+'.test','-test.run=^$','-test.bench=^BenchmarkVariance16x16Native$','-test.benchtime=150ms','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('arm next kernel round',r+1,'done',flush=True)
