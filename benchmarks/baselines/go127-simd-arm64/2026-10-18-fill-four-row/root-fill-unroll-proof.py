from pathlib import Path
import hashlib,json,re,statistics,subprocess,collections
p=Path('/private/tmp/govpx-arm-fill-unroll4-ready')
manifest=(p/'SHA256SUMS').read_text().splitlines()
assert len(manifest)==166
for line in manifest:
    digest,name=line.split('  ',1)
    assert hashlib.sha256((p/name).read_bytes()).hexdigest()==digest,name
head='768f2fc2aba7bbd82981494473cbadc026a95570'
env=json.loads((p/'environment.json').read_text());assert env['parent_base_commit']==head
assert subprocess.check_output(['git','rev-parse',head+'^{tree}'],text=True).strip()==subprocess.check_output(['git','rev-parse','bba90959^{tree}'],text=True).strip()
locked=Path('/private/tmp/govpx-arm-fill-exact-ready')
for name in ['internal/pixel/kernels_simd.go','internal/pixel/fill_fixed_arm64.go','internal/pixel/fill_fixed_other.go']:
    assert (p/'parent-source'/name).read_bytes()==(locked/'overlay'/name).read_bytes(),name
    if name!='internal/pixel/kernels_simd.go':assert (p/'overlay'/name).read_bytes()==(locked/'overlay'/name).read_bytes()
for name,digest in env['source_sha256'].items():
    assert hashlib.sha256((p/name).read_bytes()).hexdigest()==digest,name
assert (p/'overlay/internal/pixel/kernels_simd.go').read_bytes()==Path('internal/pixel/kernels_simd.go').read_bytes()
lock=env['source_sha256']
pattern=re.compile(r'^(Benchmark\S+)\s+(\d+)\s+([0-9.eE+-]+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$',re.M)
out={}
for scope,variants,count in [('pixel',['parent','candidate'],18),('dsp',['parent','candidate','old691'],12)]:
    a={}
    for variant in variants:
        rounds=[]
        for r in range(1,7):
            text=(p/('fresh-'+scope)/('r'+str(r)+'-'+variant+'.log')).read_text()
            assert re.search(r'^PASS$',text,re.M)
            values={}
            for name,its,ns,bs,alloc in pattern.findall(text):
                assert name not in values
                assert int(its)>0 and int(bs)==int(alloc)==0
                values[name]=float(ns)
            assert len(values)==count
            rounds.append(values)
        assert all(set(v)==set(rounds[0]) for v in rounds)
        a[variant]=rounds
    assert all(set(rows[0])==set(a['parent'][0]) for rows in a.values())
    reported={r['case']:r for r in json.loads((p/('fresh-'+scope)/'summary.json').read_text())['cases']}
    results={}
    for name in a['parent'][0]:
        med={v:statistics.median(row[name] for row in rows) for v,rows in a.items()}
        pairs=[b[name]/c[name] for b,c in zip(a['parent'],a['candidate'])]
        assert med==reported[name]['medians_ns'],name
        assert pairs==reported[name]['pair_ratios_parent_over_candidate'],name
        results[name]={'medians_ns':med,'parent_speedup':med['parent']/med['candidate'],'all6_paired_speedups':pairs}
    commands=json.loads((p/('fresh-'+scope)/'commands.json').read_text())
    orders=[tuple(c['variant'] for c in commands if c['round']==r) for r in range(1,7)]
    assert all(c['exit']==0 for c in commands)
    if len(variants)==3:
        assert len(set(orders))==6
        for pos in range(3):assert collections.Counter(order[pos] for order in orders)=={'parent':2,'candidate':2,'old691':2}
    else:
        assert collections.Counter(orders)=={('parent','candidate'):3,('candidate','parent'):3}
    out[scope]=results
ignored={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
rows={v:[json.loads((p/'fresh-whole'/('r'+str(r)+'-'+v+'.json')).read_text()) for r in range(1,7)] for v in ['parent','candidate','old691']}
first=rows['parent'][0];fields=set(first)-ignored;assert len(fields)==22
assert all({k:r[k] for k in fields}=={k:first[k] for k in fields} for a in rows.values() for r in a)
assert all(r['allocs_per_frame']==0 for a in rows.values() for r in a)
med={v:statistics.median(r['ns_per_frame'] for r in a) for v,a in rows.items()}
out['whole']={'medians_ns':med,'vs_parent_time_percent':100*(med['candidate']/med['parent']-1),'vs_old_assembly_time_percent':100*(med['candidate']/med['old691']-1),'all22_fields_and_allocation_samples_exact':True}
for v in ['parent','old691']:
    out['whole'][v+'_paired_time_percent']=[100*(c['ns_per_frame']/b['ns_per_frame']-1) for b,c in zip(rows[v],rows['candidate'])]
screen_count=0
for scope,count,variants in [('pixel',18,['parent','candidate']),('dsp',12,['parent','candidate','old691'])]:
    for variant in variants:
        for r in range(1,4):
            text=(p/('screen-'+scope)/('r'+str(r)+'-'+variant+'.log')).read_text()
            values=pattern.findall(text);assert len(values)==count
            assert re.search(r'^PASS$',text,re.M)
            assert all(int(bs)==int(alloc)==0 for name,it,ns,bs,alloc in values)
            screen_count+=len(values)
assert screen_count==216
proof={'parent':head,'manifest_files_verified':len(manifest),'production_source_lock_entries_verified':len(lock),'432_fresh_native_component_samples_all_zero_alloc':True,'18_whole_jsons_verified':True,'results':out}
Path('/private/tmp/govpx-root-fill-unroll-proof.json').write_text(json.dumps(proof,indent=2)+'\n')
for scope in ['pixel','dsp']:
    qualifying={k:v['parent_speedup'] for k,v in out[scope].items() if v['parent_speedup']>=1.5}
    print(scope,'qualifying medians',json.dumps(qualifying,indent=2))
print(json.dumps(out['whole'],indent=2))
