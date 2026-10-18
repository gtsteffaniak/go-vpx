from pathlib import Path
import hashlib, json, re, statistics, subprocess

P=Path('/private/tmp/govpx-amd64-46a08f8e-artifact')
HEAD='46a08f8e6eb03f6ee4549d975384bd6c952a81fa'
BASE='78ce74d39806d74ecf91bc3d082899a74d7bc85f'
OLD='691cd0512c4829d9ee40d483205dd223f9eb028d'
V=('assembly','previous','simd')
MODULE='github.com/thesyncim/govpx'
manifest=json.loads(Path('/private/tmp/govpx-amd64-46a08f8e-files.json').read_text())
assert len(manifest)==91
for name,digest in manifest.items():
    assert hashlib.sha256((P/name).read_bytes()).hexdigest()==digest,name
env=(P/'environment.txt').read_text()
assert env.splitlines()[:3]==['assembly commit: '+OLD,'previous SIMD commit: '+BASE,'SIMD commit: '+HEAD]
cpu=re.search(r'^Model name:\s+(.+)$',env,re.M).group(1)
assert 'GOMAXPROCS=1 PGO=off rounds=6 kernel-benchtime=150ms' in env
assert 'go version go1.27.2 linux/amd64' in env
proposal=Path('/private/tmp/govpx-pr46-native-harness-498-proposal')
paths=json.loads((proposal/'source-sha-paths.json').read_text())
if isinstance(paths,dict):
    paths=paths.get('paths',paths.get('source_SHAs'))
if paths and isinstance(paths[0],dict):
    paths=[row['path'] for row in paths]
assert len(paths)==len(set(paths))==77
source={}
for path in paths:
    found=[m[1] for line in env.splitlines() for m in [re.match(r'([0-9a-f]{64})  .*/govpx/govpx/(.+)$',line)] if m and m[2]==path]
    assert len(found)==1,(path,found)
    want=hashlib.sha256(subprocess.check_output(['git','show',HEAD+':'+path])).hexdigest()
    assert found[0]==want,path
    source[path]=want
BENCH=re.compile(r'^(Benchmark\S+)\s+(\d+)\s+(\S+) ns/op.*?\s+(\S+) B/op\s+(\S+) allocs/op$')
def samples(s):
    pkg='';out={}
    for line in s.splitlines():
        if line.startswith('pkg: '): pkg=line[5:]
        match=BENCH.match(line)
        if match:
            assert int(match[2])>0 and float(match[4])==float(match[5])==0,line
            out.setdefault(pkg+'/'+match[1],[]).append({'n':int(match[2]),'ns':float(match[3])})
    return out
kernels={v:samples((P/(v+'-kernels.txt')).read_text()) for v in V}
assert all(len(k)==498 and all(len(a)==6 for a in k.values()) for k in kernels.values())
assert all(kernels[v].keys()==kernels['simd'].keys() for v in V)
inventory=json.loads((proposal/'expected-timed-inventory.json').read_text())['all_names']
expected={MODULE+'/'+name.removeprefix('./') for name in inventory}
assert set(kernels['simd'])==expected
orders=[line.split(': ',1)[1].split() for line in (P/'orders.txt').read_text().splitlines()]
assert len(orders)==len(set(map(tuple,orders)))==6
for col in range(3): assert sorted(row[col] for row in orders)==sorted(V*2)
results={}
for key in kernels['simd']:
    med={v:statistics.median(x['ns'] for x in kernels[v][key]) for v in V}
    pairs={v:[a['ns']/b['ns'] for a,b in zip(kernels[v][key],kernels['simd'][key])] for v in ('previous','assembly')}
    results[key]={'median_ns':med,'parent_speedup':med['previous']/med['simd'],'old_speedup':med['assembly']/med['simd'],'all6_paired_speedups':pairs}
copy={k:x for k,x in results.items() if '/BenchmarkCopy16RowsNative/' in k}
tm={k:x for k,x in results.items() if '/BenchmarkTMPredictRowsNative/' in k}
assert len(copy)==16 and len(tm)==27
assert sum(k.startswith(MODULE+'/BenchmarkCopy16RowsNative/') for k in copy)==8
ordinary={k:x for k,x in copy.items() if any('/'+p+'/aliasfalse' in k for p in ('pitch16_16','pitch1312_1320','pitch17_31'))}
assert len(ordinary)==6
controls={k:x for k,x in results.items() if k not in copy}
quant=(P/'quant-batch-validation.txt').read_text()
assert '--- FAIL:' not in quant and not re.search(r'^FAIL(?:\s|$)',quant,re.M)
qs=samples(quant)
assert sum(len(x) for x in qs.values())==194
assert all(row['n']==1 for rows in qs.values() for row in rows)
counts={name:sum(len(x) for k,x in qs.items() if '/Benchmark'+name in k) for name in ('Copy16RowsNative/','TMPredictRowsNative/','IntraDCFill')}
assert counts=={'Copy16RowsNative/':32,'TMPredictRowsNative/':81,'IntraDCFill':81},counts
assert quant.count('# DC preflight: ')==quant.count('# TM preflight: ')==3
assert quant.count('--- PASS: TestFusedDCSupplementState ')==2
assert quant.count('--- PASS: TestTMPredictRowsNativeState ')==7
for name in ('TestCopy16AMDProbeRowsAndAliases','TestCopy16AMDProbeEveryAlignment','TestCopy16AMDProbeBoundsAndDisabled'):
    assert quant.count('--- PASS: '+name+' ')==4,name
assert quant.count('--- PASS: TestCopy16FrozenPartialPanic ')==8
for name in ('TestTMPredictBytesSelected','TestTMPredictBytesFullBackingAliases','TestTMPredictBytesIndependentAlignments','TestTMPredictBytesFallbackPanicState','TestTMPredictBytesUnsupportedWidths','TestTMPredictBytesDeclinesBeforeWrites','TestTMPredictPublicFullState'):
    assert quant.count('--- PASS: '+name+' ')==4,name
assert quant.count('--- SKIP: TestTMPredictBytesArithmetic ')==4
fp_counts={}
for name in [v+'-fp16-fixture-validation.txt' for v in V]+['simd-fp16-cpuoff-fixture-validation.txt']:
    rows=samples((P/name).read_text())
    assert len(rows)==36 and all(len(x)==1 and x[0]['n']==1 for x in rows.values())
    fp_counts[name]=36
flog=(P/'fp16-validation.txt').read_text()
assert '--- FAIL:' not in flog and flog.count('--- PASS: TestFP16Selected ')==4
assert flog.count('Enabled=true AVX2=true')==2 and flog.count('Enabled=false AVX2=false')==2
for name in ('TestFP16DisjointScalar','TestFP16AllSignedCoefficients','TestFP16UnalignedMinimum','TestFP16AliasedStoreOrder','TestFP16OriginalFallback','TestFP16Disabled'):
    assert flog.count('--- PASS: '+name+' ')==4,name
for name in ('TestCopy16RowsPublicState','TestCopy8RowsPublicState'):
    assert quant.count('--- PASS: '+name+' ')==4,name
for name in ('TestSSE8PtrExactWindows','TestSSE8PtrFallbackPanics','TestSSE8RowsWideDispatch','TestSSE8EntryDisabledScalar','TestDenseTokenPrivateScalarComponent','TestDenseTokenPrivateStreamReference','TestDenseTokenPrepackAvailability'):
    assert quant.count('--- PASS: '+name+' ')==2,name
assert quant.count('--- SKIP: TestDenseTokenPrivatePrepack ')==2
active=quant.count('AVX2=true AVX512=true;')==2
assert quant.count(('--- PASS: ' if active else '--- SKIP: ')+'TestSSEPlane16AVX512 ')==2
for name in ('TestSADWideAVX2Rows','TestSADWideAVX2FallbackAndBounds','TestSADWideAVX2NilInputs','TestSADWideAVX2Disabled'):
    assert quant.count('--- PASS: '+name+' ')==4,name
assert quant.count('--- PASS: TestSADWideAVX2Selected ')==2 and quant.count('--- SKIP: TestSADWideAVX2Selected ')==2
assert quant.count('Enabled=false actual AVX2=false')==2
for name in ('dispatch.txt','previous-dispatch.txt'):
    assert '--- PASS: TestBackendSelected ' in (P/name).read_text()
ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
whole={};profiles={}
for label,pattern in [('VP9320','{v}-encoder-{i}.json'),('VP8720','{v}-vp8-720p-encoder-{i}.json'),('VP9720','{v}-vp9-720p-encoder-{i}.json')]:
    rows={v:[json.loads((P/pattern.format(v=v,i=i)).read_text()) for i in range(1,7)] for v in V}
    allrows=[row for group in rows.values() for row in group];fields=set(allrows[0])-ignore
    assert len(fields)==22 and all({k:r[k] for k in fields}=={k:allrows[0][k] for k in fields} for r in allrows)
    med={v:statistics.median(r['ns_per_frame'] for r in group) for v,group in rows.items()}
    whole[label]={'median_ns':med,'vs_parent_percent':100*(med['simd']/med['previous']-1),'vs_old_percent':100*(med['simd']/med['assembly']-1),'all6_paired_time_differences_percent':{v:[100*(a['ns_per_frame']/b['ns_per_frame']-1) for a,b in zip(rows['simd'],rows[v])] for v in ('previous','assembly')},'all_allocations':{v:[r['allocs_per_frame'] for r in group] for v,group in rows.items()},'all22_semantics_equal':True}
for codec in ('vp8','vp9'):
    rows={v:json.loads((P/(v+'-'+codec+'-720p-profile.json')).read_text()) for v in V}
    fields=set(rows['simd'])-ignore
    assert len(fields)==22 and all({k:r[k] for k in fields}=={k:rows['simd'][k] for k in fields} for r in rows.values())
    profiles[codec]={'all22_semantics_equal':True,'all_allocations':{v:r['allocs_per_frame'] for v,r in rows.items()}}
for path in ('internal/pixel/quantize_fp_amd64.go','internal/pixel/sad_amd64.go','internal/pixel/sse8_ptr_amd64.go','internal/pixel/sse_plane_amd64.go'):
    assert subprocess.check_output(['git','show',HEAD+':'+path])==subprocess.check_output(['git','show',BASE+':'+path])
out={'head':HEAD,'base':BASE,'cpu':cpu,'native_files_verified':91,'source_SHA_verified':source,'kernel_cases':498,'kernel_samples':8964,'all_timed_kernel_allocations_zero':True,'all498_results':results,'Copy16_all16_actual_public_results':copy,'Copy16_six_disjoint_ordinary_results':ordinary,'TM27_AMD_controls':tm,'all482_nonCopy16_controls':controls,'controls_median_slower':sum(x['parent_speedup']<1 for x in controls.values()),'controls_all6_slower':sum(max(x['all6_paired_speedups']['previous'])<1 for x in controls.values()),'excluded338_oneshots':{'quant_log':counts,'FP':fp_counts},'public':whole,'profiles':profiles,'scope':'Old Copy16 is compiler-generated Go, not handwritten assembly. ARM TM helper is inactive on AMD. Keep every control regression and all observed allocation differences; do not infer cause from source identity.'}
Path('/private/tmp/govpx-root-46a08f8e-full-proof.json').write_text(json.dumps(out,indent=2)+'\n')
print(json.dumps({k:v for k,v in out.items() if k not in ('source_SHA_verified','all498_results','all482_nonCopy16_controls','TM27_AMD_controls')},indent=2))
