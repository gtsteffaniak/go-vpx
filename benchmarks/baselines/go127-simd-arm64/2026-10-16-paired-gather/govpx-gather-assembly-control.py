import os,subprocess,json
from pathlib import Path
env=dict(os.environ,GOMAXPROCS='1',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-gather-assembly')
rows=[]
for r in range(6):
 for name,binary in ([('assembly','govpx-asm-bench'),('final','govpx-gather-final-bench')] if r%2==0 else [('final','govpx-gather-final-bench'),('assembly','govpx-asm-bench')]):
  args=['/Users/thesyncim/.codex/bin/project-env','/private/tmp/'+binary,'-codec=vp8','-width=1280','-height=720','-frames=60','-fps=30','-bitrate=2500','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
  rows.append(dict(round=r,variant=name,result=json.loads(subprocess.check_output(args,env=env,text=True))))
 Path('/private/tmp/govpx-gather-assembly-control.json').write_text(json.dumps(dict(samples=rows),indent=2)+'\n')
 print('assembly control round',r+1,'done',flush=True)
