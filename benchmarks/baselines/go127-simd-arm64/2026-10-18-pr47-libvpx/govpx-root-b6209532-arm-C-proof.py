from pathlib import Path
import json,hashlib,tarfile,subprocess,statistics
p=Path('/private/tmp/govpx-pr47-arm-libvpx-ready');sha=lambda b:hashlib.sha256(b).hexdigest()
outer=json.loads((p/'sha256.json').read_text())
for n,h in outer.items():assert sha((p/n).read_bytes())==h,n
raw=json.loads((p/'raw-sha256.json').read_text())
with tarfile.open(p/'raw.tar.gz') as t:
 fs={m.name:m for m in t.getmembers() if m.isfile()};assert fs.keys()==raw.keys()
 for n,m in fs.items():assert sha(t.extractfile(m).read())==raw[n],n
proof=json.loads((p/'preparation-proof.json').read_text());head=proof['Go_head'];assert head=='b6209532fbc30cf9567e0cdc78a21d61a60fe900'
assert proof['Go_tree']==subprocess.check_output(['git','rev-parse',head+'^{tree}'],text=True).strip()
names=list(proof['Go_source_sha256']);assert len(names)==proof['Go_source_file_count']==1925
blobs=subprocess.check_output(['git','cat-file','--batch'],input=''.join(head+':'+n+'\n' for n in names).encode());pos=0
for n in names:
 end=blobs.index(b'\n',pos);header=blobs[pos:end].split();assert header[1]==b'blob';size=int(header[2]);pos=end+1;b=blobs[pos:pos+size];pos+=size+1;assert sha(b)==proof['Go_source_sha256'][n],n
assert pos==len(blobs)
up=Path('/private/tmp/govpx-libvpx-arm-upstream');prior=json.loads(Path('/private/tmp/govpx-libvpx-arm-evidence/source-validation.json').read_text());assert prior['peeled_commit']==proof['C_pin']=='1024874c5919305883187e2953de8fcb4c3d7fa6'
assert len(proof['C_source_sha256'])==1250
for n,h in proof['C_source_sha256'].items():assert h==prior['files'][n]['commit_download']==sha((up/n).read_bytes()),n
assert sha((up/'vpxenc').read_bytes())==proof['vpxenc_sha256']=='522cae8963b7c86caf76ed6c4611c7f8d9d550c612d7c50b96dd50e1787da568'
assert sha((up/'libvpx.a').read_bytes())==proof['static_archive_sha256']=='84e97dc48cd382ada5c17f77820dc8aacbe65525aa55688f939da4013f821a00'
assert sha(Path(proof['Go_binary_path']).read_bytes())==proof['Go_binary_sha256']
for mode in ('sandbox','native'):assert proof['fresh_runtime_caps'][mode]=='caps=7 NEON=1 DOTPROD=1 I8MM=1'
r=Path('/private/tmp/govpx-pr47-arm-libvpx-raw');summary=json.loads((r/'summary.json').read_text());cmds=json.loads((r/'commands.json').read_text());assert len(cmds)==32
for call in cmds:
 report=json.loads((r/call['file']).read_text());assert report['reference']['timing_source']=='vpxenc-stats'
fields=['output_bytes','output_bitrate_kbps','bitrate_error_pct','encoded_frames','dropped_frames','keyframe_bytes','avg_interframe_bytes'];results={}
for label in ('vp8-720p-fixed','vp8-720p-adaptive','vp9-320p','vp9-720p'):
 rows=[json.loads((r/(label+'-round-'+str(i)+'.json')).read_text()) for i in range(1,7)];quality=json.loads((r/(label+'-quality.json')).read_text())
 calls=[x for x in cmds if x['file'].startswith(label+'-round-')];assert len(calls)==6 and sum(x['reference_first'] for x in calls)==3
 assert all('-encode-only' in x['arguments'] for x in calls)
 assert all(x['quality_skipped'] for x in rows)
 for side in ('govpx','libvpx'):
  timed=rows if side=='govpx' else [x['reference'] for x in rows];q=quality if side=='govpx' else quality['reference']
  assert all({k:x[k] for k in fields}=={k:q[k] for k in fields} for x in timed)
 assert all(x['reference']['timing_source']=='vpxenc-stats' for x in rows+[quality])
 med=statistics.median(x['ns_per_frame'] for x in rows)/statistics.median(x['reference']['ns_per_frame'] for x in rows)
 assert med==summary['cases'][label]['govpx_median_ns_per_input']/summary['cases'][label]['libvpx_median_ns_per_input']
 assert abs((med-1)*100-summary['cases'][label]['govpx_time_delta_percent'])<1e-10
 result={'go_over_libvpx_time':med,'paired_time_ratios':[x['ns_per_frame']/x['reference']['ns_per_frame'] for x in rows],'all_allocations':[x['allocs_per_frame'] for x in rows],'timed_output_matches_quality_per_side':True,'PSNR_SSIM_measured_only_in_separate_quality_pass':True}
 if label.startswith('vp8'):
  assert all({k:x[k] for k in fields}=={k:x['reference'][k] for k in fields} for x in rows+[quality])
  assert all(quality[k]==quality['reference'][k] for k in ('psnr','ssim','quality_frames')) and quality['quality_frames']==58
  assert all(x['allocs_per_frame']==0 for x in rows)
  result['measured_quality_and_output_match']=True
 else:result['measured_quality_and_output_match']=False
 results[label]=result
out={'head':head,'outer6_raw83_SHA_verified':True,'Go1925_exact_Git_blobs':True,'C1250_exact_prior_pinned_commit_and_current_source':True,'actual_bins_and_static_archive_verified':True,'physical_CPU':'Apple M4 Max / Mac16,6','results':results,'quality_scope':'PSNR/SSIM measured only in separate quality passes; timed output stats match those passes. No bitstream identity or cross-stage causal claim. VP9 unequal measured output/quality. Compilation/cache location is excluded from throughput denominators. Actual Cref Go build used wrapper default cache prefix with common project identity; supplemental note preserves this.'}
Path('/private/tmp/govpx-root-b6209532-arm-C-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
