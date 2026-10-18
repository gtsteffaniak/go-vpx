from pathlib import Path
import re,subprocess,os,json,hashlib
root=Path('/Users/thesyncim/GolandProjects/govpx');out=Path('/private/tmp/govpx-denoise-fixed-fixture');base=Path('/private/tmp/govpx-frame-inter-base');head=Path('/private/tmp/govpx-frame-inter-next');name='vp8_denoiser_average_stream_test.go'
s=(out/'original-adaptive-stream_test.go.source').read_text();s=s.replace('''				opts.NoiseSensitivity = 4''','''                // Packet pins require a fixed speed: small adaptive frames use
                // real wall-clock feedback and can change picker dispatch.
                opts.CpuUsed = -8
				opts.NoiseSensitivity = 4''',1)
s=s.replace('''					if err != nil {
						t.Fatal(err)
					}
					binary.LittleEndian''','''					if err != nil {
						t.Fatal(err)
					}
                    if e.libvpxCPUUsed() != 8 { t.Fatalf("fixed speed changed to %d",e.libvpxCPUUsed()) }
					binary.LittleEndian''',1)
a=s.index('\n\t\t\t\twant := map[string]string');b=s.index('\n\t\t\t\tt.Logf(',a);probe=s[:a]+s[b:];(base/name).write_text(probe);(out/'candidate-template.source').write_text(s)
env=os.environ|dict(CODEX_AGENT_ID='denoise-fixed',CODEX_PROJECT_CACHE_ROOT='/private/tmp/codex-projects',GOTOOLCHAIN='local',GO='/Users/thesyncim/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.2.darwin-arm64/bin/go',GOEXPERIMENT='simd',GOMAXPROCS='2',GOFLAGS='-p=2',CGO_ENABLED='1')
def build(p,name):
 with(out/(name+'-build.log')).open('w')as f:r=subprocess.run([str(root/'scripts/go.sh'),'test','-C',str(p),'-c','-pgo=off','-o',str(out/(name+'.test')),'.'],cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT)
 print(name,'build',r.returncode,flush=True)
 if r.returncode:print((out/(name+'-build.log')).read_text());raise SystemExit(1)
def run(name,log):
 with(out/log).open('w')as f:r=subprocess.run([str(out/(name+'.test')),'-test.run=^TestDenoiserAverageStream$','-test.count=3','-test.v'],env=dict(env,GOMAXPROCS='1'),stdout=f,stderr=subprocess.STDOUT)
 print(name,'repeat3',r.returncode,flush=True)
 if r.returncode:print((out/log).read_text());raise SystemExit(1)
build(base,'base');run('base','baseline-repeat3.log');expected={};current=None
for line in (out/'baseline-repeat3.log').read_text().splitlines():
 if line.startswith('=== RUN   TestDenoiserAverageStream/'):current=line.split('TestDenoiserAverageStream/',1)[1]
 if 'denoiser_stream_sha256='in line:
  value=line.split('denoiser_stream_sha256=',1)[1];assert current not in expected or expected[current]==value,(current,expected.get(current),value);expected[current]=value
assert len(expected)==4;assert (out/'baseline-repeat3.log').read_text().count('denoiser_stream_sha256=')==12
(out/'baseline-pins.json').write_text(json.dumps(expected,indent=2)+'\n');literal='map[string]string{'+','.join(json.dumps(k)+':'+json.dumps(v)for k,v in expected.items())+'}'
s=re.sub(r'map\[string\]string\{[^\n]+\}',lambda _:literal,s,count=1);(head/name).write_text(s)
subprocess.run([env['GO'].replace('/bin/go','/bin/gofmt'),'-w',str(head/name)],check=True)
build(head,'candidate');run('candidate','candidate-repeat3.log')
(out/'fixture.go.source').write_bytes((head/name).read_bytes());(out/'binding.json').write_text(json.dumps({'baseline_runtime':'88fe29f77914a38e96a23adb28bd0275b827ce87 (private2bbc6431 runtime-equivalent PR110 tree)','base_tree':str(base),'candidate_tree':str(head),'fixture':name,'CpuUsed':-8,'GOMAXPROCS':1,'baseline_before_candidate':True,'repeats_per_case':3,'pins':expected,'fixture_sha256':hashlib.sha256((head/name).read_bytes()).hexdigest(),'baseline_binary_sha256':hashlib.sha256((out/'base.test').read_bytes()).hexdigest(),'candidate_binary_sha256':hashlib.sha256((out/'candidate.test').read_bytes()).hexdigest(),'production_changes':'none; fixture-only fixed-speed correction'},indent=2)+'\n')
print('FIXED PINS READY',json.dumps(expected),flush=True)
