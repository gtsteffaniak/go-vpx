import os,subprocess
from pathlib import Path
env=dict(os.environ,GOFLAGS='-p=2',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',CODEX_AGENT_ID='root-gather-final',GOTOOLCHAIN='local',GO='/Users/thesyncim/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.2.darwin-arm64/bin/go')
with open('/private/tmp/govpx-gather-final-build.txt','w') as f:
 for args in (['test','-c','-pgo=off','-o','/private/tmp/govpx-gather-final.test','./internal/pixel'],['build','-pgo=off','-o','/private/tmp/govpx-gather-final-bench','./cmd/govpx-bench']):subprocess.run(['scripts/go.sh']+args,env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
 env.update(GOOS='darwin',GOARCH='amd64')
 subprocess.run(['scripts/go.sh','test','-c','-pgo=off','-o','/private/tmp/govpx-gather-final-amd64.test','./internal/pixel'],env=env,stdout=f,stderr=subprocess.STDOUT,check=True)
