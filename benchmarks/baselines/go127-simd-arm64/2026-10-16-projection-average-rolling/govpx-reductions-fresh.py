import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-reductions-fresh')
wrapper='/Users/thesyncim/.codex/bin/project-env'
variants={'vp8':[('parent','govpx-internal-final-bench'),('prototype','govpx-six-roll-bench'),('final','govpx-reductions-final-bench')],'vp9':[('parent','govpx-projection-parent-bench'),('prototype','govpx-projection-average-bench'),('final','govpx-reductions-final-bench')]}
orders=[[0,1,2],[1,2,0],[2,0,1],[2,1,0],[0,2,1],[1,0,2]]
for codec,items in variants.items():
 rows=[]
 for r,order in enumerate(orders):
  for i in order:
   name,binary=items[i]
   cmd=[wrapper,'/private/tmp/'+binary,f'-codec={codec}','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
   rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
  Path(f'/private/tmp/govpx-reductions-fresh-{codec}.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
  print(codec,'fresh round',r+1,'done',flush=True)
with open('/private/tmp/govpx-reductions-final-kernels.txt','w') as f:
 for r in range(6):
  for name in (['parent','final'] if r%2==0 else ['final','parent']):
   for codec in ['vp8','vp9']:
    binary=('govpx-internal-final.test' if codec=='vp8' else 'govpx-projection-parent.test') if name=='parent' else f'govpx-reductions-final-{codec}.test'
    regex='^BenchmarkSixTapPhases$' if codec=='vp8' else '^Benchmark(Avg8x8Quad|ProjectionColumnsNative)$'
    f.write(f'ROUND {r} VARIANT {name} CODEC {codec}\n');f.flush()
    subprocess.run([wrapper,'/private/tmp/'+binary,'-test.run=^$','-test.bench='+regex,'-test.benchtime=100ms','-test.count=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
  print('final kernels round',r+1,'done',flush=True)
