from pathlib import Path
import hashlib, json, statistics, re
p = Path('/private/tmp/govpx-libvpx-amd64-46a08f8e-artifact')
manifest = json.loads(Path('/private/tmp/govpx-libvpx-amd64-46a08f8e-files.json').read_text())
assert len(manifest) == 74
for name, digest in manifest.items():
    assert hashlib.sha256((p / name).read_bytes()).hexdigest() == digest, name
env = (p / 'environment.txt').read_text()
assert env.startswith('govpx commit: 46a08f8e6eb03f6ee4549d975384bd6c952a81fa\n')
cpu = re.search(r'^Model name:\s+(.+)$', env, re.M).group(1)
assert 'libvpx commit: 1024874c5919305883187e2953de8fcb4c3d7fa6' in env
summary = json.loads((p / 'summary.json').read_text())
fields = ['output_bytes', 'output_bitrate_kbps', 'bitrate_error_pct', 'encoded_frames', 'dropped_frames', 'keyframe_bytes', 'avg_interframe_bytes']
results = {}
for label in ['vp8-720p-fixed', 'vp8-720p-adaptive', 'vp9-320p', 'vp9-720p']:
    rows = [json.loads((p / (label + '-round-' + str(i) + '.json')).read_text()) for i in range(1, 7)]
    quality = json.loads((p / (label + '-quality.json')).read_text())
    for side in ['govpx', 'libvpx']:
        timed = [r if side == 'govpx' else r['reference'] for r in rows]
        qr = quality if side == 'govpx' else quality['reference']
        assert all({k: r[k] for k in fields} == {k: qr[k] for k in fields} for r in timed)
    assert all(r['reference']['timing_source'] == 'vpxenc-stats' for r in rows)
    go = statistics.median(r['ns_per_frame'] for r in rows)
    c = statistics.median(r['reference']['ns_per_frame'] for r in rows)
    assert go == summary['cases'][label]['govpx_median_ns_per_input']
    assert c == summary['cases'][label]['libvpx_median_ns_per_input']
    results[label] = {'go_over_libvpx_time': go/c, 'paired_time_ratios': [r['ns_per_frame']/r['reference']['ns_per_frame'] for r in rows]}
    if label == 'vp8-720p-fixed':
        assert all(r['allocs_per_frame'] == 0 for r in rows)
        for r in rows + [quality]:
            assert all(r[k] == r['reference'][k] for k in fields)
            assert '--cpu-used=-8' in r['reference']['parity_flags']
            assert '--min-q=2' in r['reference']['parity_flags']
            assert '--max-q=56' in r['reference']['parity_flags']
        assert not quality.get('quality_skipped', False)
        assert not quality['reference'].get('quality_skipped', False)
        assert quality['psnr'] == quality['reference']['psnr']
        assert quality['ssim'] == quality['reference']['ssim']
        assert quality['quality_frames'] == quality['reference']['quality_frames'] == 58
config = (p / 'vpx_config.h').read_text()
for setting in ['#define CONFIG_RUNTIME_CPU_DETECT 1', '#define HAVE_SSE2 1', '#define HAVE_AVX2 1']:
    assert setting in config, setting
proof = {'head': '46a08f8e6eb03f6ee4549d975384bd6c952a81fa', 'raw_files_verified': 74, 'cpu': cpu, 'scope': 'Separate native runner optimized libvpx v1.16 fixed CPU-used -8 VP8 case has matched measured output and quality; adaptive VP8 and VP9 differ in quality/output work. No cross-job change attribution or bitstream identity claim.', 'all24_encode_timers_are_vpxenc_stats': True, 'quality_pass_output_fields_match_each_timing_run': True, 'results': results}
Path('/private/tmp/govpx-root-46a08f8e-C-proof.json').write_text(json.dumps(proof, indent=2)+'\n')
print(json.dumps(proof, indent=2))
