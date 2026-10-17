import subprocess,os
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-next-kernels')
wrapper='/Users/thesyncim/.codex/bin/project-env'
variants=['projection-parent','projection','average','projection-average']
orders=[[(i+s)%4 for i in [0,1,3,2]] for s in range(4)];orders += [list(reversed(o)) for o in orders]
with open('/private/tmp/govpx-projection-average-kernels.txt','w') as f:
 for r,order in enumerate(orders):
  for i in order:
   name=variants[i];f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,f'/private/tmp/govpx-{name}.test','-test.run=^$','-test.bench=^Benchmark(Avg8x8Quad|ProjectionColumnsNative)$','-test.benchtime=100ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('projection kernel round',r+1,'done',flush=True)
with open('/private/tmp/govpx-six-roll-kernels.txt','w') as f:
 for r in range(6):
  for name,binary in ([('parent','govpx-internal-final.test'),('rolling','govpx-six-roll.test')] if r%2==0 else [('rolling','govpx-six-roll.test'),('parent','govpx-internal-final.test')]):
   f.write(f'ROUND {r} VARIANT {name}\n');f.flush()
   subprocess.run([wrapper,f'/private/tmp/{binary}','-test.run=^$','-test.bench=^BenchmarkSixTapPhases$','-test.benchtime=100ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('rolling kernel round',r+1,'done',flush=True)
