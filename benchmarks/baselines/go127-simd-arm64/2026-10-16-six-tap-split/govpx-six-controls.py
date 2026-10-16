import subprocess,os,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-six-controls')
wrapper='/Users/thesyncim/.codex/bin/project-env'
with open('/private/tmp/govpx-split-observed-ab.txt','w') as f:
 for r in range(6):
  for name in (['before','after'] if r%2==0 else ['after','before']):
   f.write(f'{r+1} {name}\n');f.flush()
   subprocess.run([wrapper,f'/private/tmp/govpx-split-observed-{name}.test','-test.run=^$','-test.bench=^BenchmarkSplitSubpelPhases$','-test.benchtime=100ms','-test.count=1','-test.cpu=1','-test.benchmem'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
 print('observable split microbenchmarks done',flush=True)
rows=[]
for r in range(6):
 for codec,width,height,frames,bitrate,speed in [('vp9',1280,720,60,2500,8),('vp8',320,180,60,500,4)]:
  for name,binary in ([('before','govpx-bench-perf33'),('after','govpx-six-combined-bench')] if r%2==0 else [('after','govpx-six-combined-bench'),('before','govpx-bench-perf33')]):
   cmd=[wrapper,'/private/tmp/'+binary,'-codec='+codec,f'-width={width}',f'-height={height}',f'-frames={frames}','-fps=30',f'-bitrate={bitrate}','-mode=realtime',f'-cpu-used={speed}','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
   rows.append(dict(round=r,codec=codec,speed=speed,variant=name,result=json.loads(subprocess.check_output(cmd,env=env,text=True))))
 Path('/private/tmp/govpx-six-controls-ab.json').write_text(json.dumps(rows,indent=2)+'\n')
 print('control round',r+1,'done',flush=True)
