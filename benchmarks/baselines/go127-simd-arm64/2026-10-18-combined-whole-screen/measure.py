from pathlib import Path
import os,subprocess,json,statistics,hashlib
p=Path('/private/tmp/govpx-pr53-short-whole');p.mkdir(exist_ok=False)
bins={'parent':'/private/tmp/govpx-astra-quant-packet54/root-arm-parent-whole','current':'/private/tmp/govpx-pr53-final-whole','old691':'/private/tmp/govpx-arm48-integrated/old691-whole-bench'}
meta={'hardware':'Apple M4 Max/Darwin ARM64','compiler':'Go1.27.2','PGO':'off','CGO_ENABLED':1,'GOMAXPROCS':1,'scope':'three-round preliminary combined whole-codec screen','binary_sha256':{k:hashlib.sha256(Path(v).read_bytes()).hexdigest()for k,v in bins.items()}}
env=os.environ.copy();env['GOMAXPROCS']='1';ignore={'ns_per_frame','encode_fps','macroblocks_per_second','latency_ns','allocs_per_frame'}
for case,codec,w,h,frames,bitrate in [('vp8-720','vp8',1280,720,300,2500),('vp9-320','vp9',320,180,120,500),('vp9-720','vp9',1280,720,60,2500)]:
 rows={k:[]for k in bins};ref=None
 for rnd,order in enumerate([['parent','current','old691'],['current','old691','parent'],['old691','parent','current']],1):
  for k in order:
   cmd=[bins[k],f'-codec={codec}',f'-width={w}',f'-height={h}',f'-frames={frames}','-fps=30',f'-bitrate={bitrate}','-mode=realtime','-cpu-used=8','-threads=1','-encode-only','-auto-libvpx=false','-format=json']
   t=subprocess.run(cmd,env=env,capture_output=True,text=True);(p/f'{case}-r{rnd}-{k}.json').write_text(t.stdout);(p/f'{case}-r{rnd}-{k}.stderr').write_text(t.stderr);assert t.returncode==0
   r=json.loads(t.stdout);sem={a:b for a,b in r.items()if a not in ignore};assert len(sem)==22
   if ref is None:ref=sem
   assert sem==ref,(case,rnd,k);rows[k].append(r)
 med={k:statistics.median(x['ns_per_frame']for x in rs)for k,rs in rows.items()}
 meta[case]={'median_ns_per_frame':med,'change_vs_parent_pct':100*(med['current']/med['parent']-1),'change_vs_old_asm_pct':100*(med['current']/med['old691']-1),'allocs_per_frame':{k:[r['allocs_per_frame']for r in rs]for k,rs in rows.items()},'semantic_fields_equal':True}
 (p/'summary.json').write_text(json.dumps(meta,indent=2)+'\n');print(case,meta[case],flush=True)
print('END all timers stopped',flush=True)
