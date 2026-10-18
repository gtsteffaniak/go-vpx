from pathlib import Path
import json,hashlib,statistics
sha=lambda b:hashlib.sha256(b).hexdigest()
p=Path('/private/tmp/govpx-pr48-arm-libvpx-verified/verification-proof.json');proof=json.loads(p.read_text());ROOT=Path('/Users/thesyncim/GolandProjects/govpx');snap=Path('/private/tmp/govpx-arm48-integrated/candidate')
assert len(proof['all_Go_sha256'])==1937
for n,h in proof['all_Go_sha256'].items():assert sha((ROOT/n).read_bytes())==sha((snap/n).read_bytes())==h,n
assert sha(Path(proof['candidate_binary']).read_bytes())==proof['candidate_binary_sha256']=='c5c4395d55ece52bfeb2c1d9dd4d19cef90337477117629d0085bf50f853b7ae'
assert len(proof['production_changes'])==13 and len(proof['test_changes'])==6
assert proof['Go_build_environment']['CODEX_PROJECT_CACHE_ROOT']=='/private/tmp/codex-projects'
up=Path('/private/tmp/govpx-libvpx-arm-upstream');prior=json.loads(Path('/private/tmp/govpx-libvpx-arm-evidence/source-validation.json').read_text());assert prior['peeled_commit']==proof['C_pin']=='1024874c5919305883187e2953de8fcb4c3d7fa6'
assert len(proof['C1250current_source_hashes'])==1250
for n,h in proof['C1250current_source_hashes'].items():assert h==prior['files'][n]['commit_download']==sha((up/n).read_bytes()),n
assert sha((up/'vpxenc').read_bytes())==proof['vpxenc_binary_sha256']=='522cae8963b7c86caf76ed6c4611c7f8d9d550c612d7c50b96dd50e1787da568'
assert sha((up/'libvpx.a').read_bytes())==proof['libvpx_static_archive_sha256']=='84e97dc48cd382ada5c17f77820dc8aacbe65525aa55688f939da4013f821a00'
assert proof['fresh_caps_sandbox']==proof['fresh_caps_native']=='caps=7 NEON=1 DOTPROD=1 I8MM=1'
r=Path('/private/tmp/govpx-pr48-arm-libvpx-raw');summary=json.loads((r/'summary.json').read_text());cmds=json.loads((r/'commands.json').read_text());assert len(cmds)==32
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
out={'compiled_candidate_SHA256':proof['candidate_binary_sha256'],'source_base_Git':proof['parent_Git'],'compiled_currentROOT1937_Go_hashes_verified':True,'production_overlay_files':13,'test_overlay_files':6,'final_Git_binding':'pending publication; exact compiled Go bytes verified against current ROOT','C1250_exact_pinned_commit_current_source_verified':True,'C_optimized_binary_static_archive_caps7_verified':True,'all32helpers_C_codec_timer_vpxenc_stats_verified':True,'results':results,'scope':'Apple M4 Max/Mac16,6 SDKGo1.27.2 CGO1 SIMD PGOoff GOMAX1; exact integrated candidate. BothVP8 workloads match measured output and separately decoded58-frame quality. VP9unequal output/quality. Quality only separate passes, timed outputstats checked versus those passes. No bitstream identity or causal inference from separate windows. C wall/legacy latency not used for performance ratios.'}
Path('/private/tmp/govpx-root-arm48-C-proof.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
