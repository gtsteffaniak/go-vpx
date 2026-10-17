import os,subprocess
from pathlib import Path
env=dict(os.environ,GOFLAGS='-p=2',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-gather-design',GOTOOLCHAIN='local',GO='/Users/thesyncim/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.2.darwin-arm64/bin/go')
for name in ('base','pairs','quads'):
 root='/private/tmp/govpx-gather-'+name+'-reference'
 with open('/private/tmp/govpx-gather-'+name+'-tests.txt','w') as log:
  subprocess.run(['scripts/go.sh','-C',root,'test','-pgo=off','./internal/pixel','./internal/vp8/dsp'],env=env,stdout=log,stderr=subprocess.STDOUT,check=True)
  subprocess.run(['scripts/go.sh','-C',root,'test','-c','-pgo=off','-o','/private/tmp/govpx-gather-'+name+'.test','./internal/pixel'],env=env,stdout=log,stderr=subprocess.STDOUT,check=True)
  if name!='base':subprocess.run(['scripts/go.sh','-C',root,'build','-pgo=off','-o','/private/tmp/govpx-gather-'+name+'-bench','./cmd/govpx-bench'],env=env,stdout=log,stderr=subprocess.STDOUT,check=True)
 print('prepared',name,flush=True)
