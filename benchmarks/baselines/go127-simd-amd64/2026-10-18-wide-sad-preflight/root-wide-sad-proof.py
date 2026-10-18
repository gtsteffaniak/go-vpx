from pathlib import Path
import hashlib,json,re,subprocess,tarfile,shutil
p=Path('/private/tmp/govpx-wide-sad-avx2-ready')
lines=(p/'SHA256SUMS').read_text().splitlines();assert len(lines)==62
for line in lines:
    digest,name=line.split('  ',1)
    assert hashlib.sha256((p/name).read_bytes()).hexdigest()==digest,name
sources=json.loads((p/'source-sha256.json').read_text())
head='768f2fc2aba7bbd82981494473cbadc026a95570'
for name,ident in sources.items():
    assert hashlib.sha256((p/'overlay'/name).read_bytes()).hexdigest()==ident['candidate'],name
    if ident['parent']:
        assert hashlib.sha256(subprocess.check_output(['git','show',head+':'+name])).hexdigest()==ident['parent'],name
    else:
        assert subprocess.run(['git','cat-file','-e',head+':'+name],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL).returncode!=0
wide={f'BenchmarkVP9WideSADNative/{w}x{h}/{layout}/{pattern}' for w,h in [(32,16),(32,32),(32,64),(64,32),(64,64)] for layout in ['packed','padded','largepitch'] for pattern in ['patterned','zero','max']}
controls={f'BenchmarkVP9WideSADNative/{w}x{h}/{layout}/patterned' for w,h in [(8,8),(8,16),(16,16)] for layout in ['packed','padded','largepitch']}
expected=wide|controls;assert len(expected)==54
one_shots=[]
for mode in ['ordinary','checkptr2','actual-avx-off','parent768','old691','purego','nosimd']:
    text=(p/'evidence'/('rosetta-'+mode+'-public54.log')).read_text()
    assert re.search(r'^PASS$',text,re.M)
    cases=set()
    for name,iterations,bs,alloc in re.findall(r'^(Benchmark\S+)\s+(\d+)\s+(?:[0-9.eE+-]+ ns/op\s+)?(?:[0-9.eE+-]+ MB/s\s+)?(\d+) B/op\s+(\d+) allocs/op$',text,re.M):
        name=re.sub(r'-\d+$','',name)
        assert name not in cases and int(iterations)==1 and int(bs)==int(alloc)==0
        cases.add(name)
    assert cases==expected,mode
    one_shots.append(mode)
off=(p/'evidence/rosetta-actual-avx-off.log').read_text()
assert 'Enabled=false actual AVX2=false' in off
assert '--- SKIP: TestSADWideAVX2Selected ' in off
for test in ['Rows','FallbackAndBounds','NilInputs','Disabled']:
    assert '--- PASS: TestSADWideAVX2'+test+' ' in off,test
codegen={}
for target in ['linux','darwin']:
    text=(p/'evidence'/(target+'-SAD-llvm.txt')).read_text()
    ops=[]
    for line in text.splitlines():
        match=re.match(r'^\s*[0-9a-f]+:\s+(?:(?:[0-9a-f]{2})\s+)*\s*([a-z][a-z0-9.]*)\s*(.*)$',line)
        if match:ops.append(match.groups())
    assert len(ops)==194,(target,len(ops))
    assert sum(op=='vzeroUpper'.lower() for op,_ in ops)==1
    assert sum(op=='vpsadbw' and '%ymm' in args for op,args in ops)==2
    assert not any(('%xmm' in args or '%ymm' in args or '%zmm' in args) and '(%rsp)' in args for op,args in ops)
    assert not any(('%xmm' in args or '%ymm' in args or '%zmm' in args) and not op.startswith('v') for op,args in ops)
    fence=next(i for i,(op,_) in enumerate(ops) if op=='vzeroupper')
    assert [op for op,_ in ops[fence+1:fence+4]]==['addq','popq','retq']
    codegen[target]={'instructions':len(ops),'two_ymm_psadbw_sites':True,'vector_spills_and_legacy_vector_ops':0,'fence_then_frame_teardown_return':True}
proof={'parent':head,'manifest_members_verified':len(lines),'overlay_git_and_source_hashes_verified':len(sources),'378_one_shot_samples_are_correctness_only_not_perf':True,'all54_cases_and_zero_allocations_verified_for_modes':one_shots,'wide_cases':45,'narrow_controls':9,'codegen':codegen,'native_performance':'unmeasured; CI required before acceptance'}
Path('/private/tmp/govpx-root-wide-sad-proof.json').write_text(json.dumps(proof,indent=2)+'\n')
print(json.dumps(proof,indent=2))
