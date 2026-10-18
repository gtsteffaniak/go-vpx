import pathlib,hashlib,json,re,subprocess,statistics
p=pathlib.Path('/private/tmp/govpx-amd64-f8dd8c74-artifact');head='f8dd8c7472e91c1f147511b04ceec716aca95aa7'
env=(p/'environment.txt').read_text(); assert env.splitlines()[:3]==['assembly commit: 691cd0512c4829d9ee40d483205dd223f9eb028d','previous SIMD commit: cf3117194f0590fa0af35951a32e8fadb6da9ab3','SIMD commit: '+head]
shas={}
for line in env.splitlines():
 m=re.match(r'([0-9a-f]{64})  .*/govpx/govpx/(.+)',line)
 if m:
  want=hashlib.sha256(subprocess.check_output(['git','show',head+':'+m[2]])).hexdigest();assert m[1]==want,m[2];shas[m[2]]=want
assert len(shas)==36
ks={}
for v in ['assembly','previous','simd']:
 d={};pkg=''
 for line in (p/(v+'-kernels.txt')).read_text().splitlines():
  if line.startswith('pkg: '):pkg=line[5:]
  m=re.match(r'(Benchmark\S+)\s+\d+\s+(\S+)\s+ns/op.*?\s+(\S+)\s+B/op\s+(\S+)\s+allocs/op',line)
  if m:
   assert float(m[3])==float(m[4])==0;d.setdefault(pkg+'/'+m[1],[]).append(float(m[2]))
 assert len(d)==338;assert all(len(a)==6 for a in d.values());ks[v]=d
assert ks['assembly'].keys()==ks['previous'].keys()==ks['simd'].keys()
ratios={}
for key in ks['simd']:
 if 'BenchmarkSSE8RowsNative/' in key:
  pairs=[a/b for a,b in zip(ks['previous'][key],ks['simd'][key])]
  ratios[key]={'median_speedup':statistics.median(ks['previous'][key])/statistics.median(ks['simd'][key]),'paired_speedups':pairs}
assert len(ratios)==36
ignored={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'};public={}
for label,pat in [('VP9320','{v}-encoder-{i}.json'),('VP8720','{v}-vp8-720p-encoder-{i}.json'),('VP9720','{v}-vp9-720p-encoder-{i}.json')]:
 rows={v:[json.loads((p/pat.format(v=v,i=i)).read_text()) for i in range(1,7)] for v in ks}
 assert all(r['allocs_per_frame']==rows['simd'][0]['allocs_per_frame'] for a in rows.values() for r in a)
 rr=[r for a in rows.values() for r in a];fields=set(rr[0])-ignored;assert len(fields)==22;assert all({k:r[k] for k in fields}=={k:rr[0][k] for k in fields} for r in rr)
 med={v:statistics.median(r['ns_per_frame'] for r in a) for v,a in rows.items()}
 public[label]={'median_ns':med,'vs_parent_percent':100*(med['simd']/med['previous']-1),'vs_old_binary_percent':100*(med['simd']/med['assembly']-1),'allocations':{v:[r['allocs_per_frame'] for r in a] for v,a in rows.items()}}
for codec in ['vp8','vp9']:
 rows=[json.loads((p/(v+'-'+codec+'-720p-profile.json')).read_text()) for v in ks];fields=set(rows[0])-ignored;assert len(fields)==22;assert all({k:r[k] for k in fields}=={k:rows[0][k] for k in fields} for r in rows);assert all(r['allocs_per_frame']==rows[0]['allocs_per_frame'] for r in rows)
manifest=json.loads(pathlib.Path('/private/tmp/govpx-amd64-f8dd8c74-files.json').read_text());assert len(manifest)==86
for f,h in manifest.items():assert hashlib.sha256((p/f).read_bytes()).hexdigest()==h,f
c=pathlib.Path('/private/tmp/govpx-libvpx-amd64-f8dd8c74-artifact');cm=json.loads(pathlib.Path('/private/tmp/govpx-libvpx-amd64-f8dd8c74-files.json').read_text());assert len(cm)==74
for f,h in cm.items():assert hashlib.sha256((c/f).read_bytes()).hexdigest()==h,f
log=(p/'quant-batch-validation.txt').read_text()
assert log.count('--- PASS: TestCopy16RowsPublicState ')==4
assert log.count('--- SKIP: TestDenseTokenPrivatePrepack ')==2
assert '--- FAIL' not in log
for name in ['TestSSE8PtrExactWindows','TestSSE8PtrFallbackPanics','TestSSE8RowsWideDispatch','TestSSE8EntryDisabledScalar','TestDenseTokenPrivateScalarComponent','TestDenseTokenPrivateStreamReference','TestDenseTokenPrepackAvailability']:
 assert log.count('--- PASS: '+name+' ')==2,name
# PR42 inherits AVX512 unchanged, and its runtime
# state must be reported honestly rather than inferred from the preceding host.
active=log.count('AVX2=true AVX512=true;')==2
if active:assert log.count('--- PASS: TestSSEPlane16AVX512 ')==2
else:assert log.count('--- SKIP: TestSSEPlane16AVX512 ')==2
assert subprocess.check_output(['git','show',head+':internal/pixel/sse_plane_amd64.go'])==subprocess.check_output(['git','show','ab3abd6b335242d5732fa1956a63b7591eaf781b:internal/pixel/sse_plane_amd64.go'])
out={'head':head,'native_files_verified':len(manifest),'c_files_verified':len(cm),'git_source_shas_verified':shas,'kernel_cases':338,'samples_per_variant':6,'all_kernel_allocations_zero':True,'sse8_all36':ratios,'sse8_medians_at_least_1_5':sum(r['median_speedup']>=1.5 for r in ratios.values()),'sse8_all6_pairs_at_least_1_5':sum(min(r['paired_speedups'])>=1.5 for r in ratios.values()),'avx512_active_this_host':active,'inherited_avx512_source_unchanged':True,'all22_public_and_profile_semantics_equal':True,'all_public_and_profile_allocation_samples_equal':True,'public':public}
pathlib.Path('/private/tmp/govpx-root-f8dd-full-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps({k:v for k,v in out.items() if k not in ['git_source_shas_verified','sse8_all36']},indent=2))
