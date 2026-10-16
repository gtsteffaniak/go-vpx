import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-internal-fresh')
wrapper='/Users/thesyncim/.codex/bin/project-env'
with open('/private/tmp/govpx-internal-final-kernels-fresh.txt','w') as f:
 for r in range(6):
  for name,binary in ([('parent','govpx-internal-before.test'),('final','govpx-internal-final.test')] if r%2==0 else [('final','govpx-internal-final.test'),('parent','govpx-internal-before.test')]):
   f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,'/private/tmp/'+binary,'-test.run=^$','-test.bench=BenchmarkLoopFilter(Horizontal|Vertical)EdgesY$','-test.benchtime=150ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('fresh kernel round',r+1,'done',flush=True)
