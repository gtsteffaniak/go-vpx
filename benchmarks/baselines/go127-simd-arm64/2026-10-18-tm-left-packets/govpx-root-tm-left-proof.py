from pathlib import Path
import hashlib, json, re, statistics

P = Path('/private/tmp/govpx-next-tm-left-packet')
STRONG = Path('/private/tmp/govpx-next-tm-byte-arm64/candidate')
V = ('parent', 'candidate', 'old691')
def prod(d):
    return {str(f.relative_to(d)): hashlib.sha256(f.read_bytes()).hexdigest() for f in d.rglob('*.go') if not f.name.endswith('_test.go')}
a, b, c = prod(STRONG), prod(P/'parent'), prod(P/'candidate')
assert a == b and len(a) == 785
delta = sorted(k for k in b.keys() | c.keys() if b.get(k) != c.get(k))
assert delta == ['internal/pixel/tmpredict_bytes_arm64.go'], delta
assert c[delta[0]] == '2c701ece6a5c3e3335623d1e2a07674dc71b958af06aac9dae365f3a54eae954'
for f in ('go.mod', 'go.sum', 'README.md'):
    paths = [STRONG/f, P/'parent'/f, P/'candidate'/f]
    assert len({p.exists() for p in paths}) == 1, f
    if paths[0].exists():
        assert len({p.read_bytes() for p in paths}) == 1, f
pat = re.compile(r'^BenchmarkTMPredictRowsNative/(\S+)\s+(\d+)\s+(\S+) ns/op\s+(\d+) B/op\s+(\d+) allocs/op$', re.M)
summary = json.loads((P/'fresh-dsp/summary.json').read_text())
raw = {v: {} for v in V}
for v in V:
    for r in range(1, 7):
        s = (P/'fresh-dsp'/('r%d-%s.txt' % (r, v))).read_text()
        assert re.search(r'^PASS$', s, re.M)
        vals = pat.findall(s)
        assert len(vals) == 27
        for k, n, ns, by, al in vals:
            assert int(n) > 0 and int(by) == int(al) == 0
            raw[v].setdefault(k, []).append(float(ns))
assert all(raw[v].keys() == raw['candidate'].keys() for v in V)
for k in raw['candidate']:
    x = summary['cases'][k]
    med = {v: statistics.median(raw[v][k]) for v in V}
    assert med == x['median_ns']
    assert all(raw[v][k] == x['ns'][v] for v in V)
    assert med['parent']/med['candidate'] == x['parent_over_candidate']
    assert med['old691']/med['candidate'] == x['old_over_candidate']
    assert [p/c for p, c in zip(raw['parent'][k], raw['candidate'][k])] == x['paired_parent_over_candidate']
    assert [p/c for p, c in zip(raw['old691'][k], raw['candidate'][k])] == x['paired_old_over_candidate']
    if k.startswith('8x8/'):
        assert min(x['paired_parent_over_candidate']) >= 1.5, k
for phase in ('fresh-dsp', 'fresh-whole'):
    commands = json.loads((P/phase/'commands.json').read_text())
    orders = commands['orders']
    assert len(orders) == len(set(map(tuple, orders))) == 6
    for col in range(3):
        assert sorted(o[col] for o in orders) == sorted(V*2)
    assert len(commands['measurements']) == 18
    for i, row in enumerate(commands['measurements']):
        assert row['exit_code'] == 0 and row['round'] == i//3+1 and row['variant'] == orders[i//3][i%3]
        if i:
            assert commands['measurements'][i-1]['end_utc'] <= row['start_utc']
ignore = {'ns_per_frame', 'encode_fps', 'macroblocks_per_second', 'latency_ns', 'allocs_per_frame'}
w = {v: [json.loads((P/'fresh-whole'/('r%d-%s.json' % (r, v))).read_text()) for r in range(1, 7)] for v in V}
rows = [x for group in w.values() for x in group]
fields = set(rows[0])-ignore
assert len(fields) == 22
assert all({k: row[k] for k in fields} == {k: rows[0][k] for k in fields} and row['allocs_per_frame'] == 0 for row in rows)
med = {v: statistics.median(x['ns_per_frame'] for x in w[v]) for v in V}
assert med == json.loads((P/'fresh-whole/summary.json').read_text())['median_ns_per_frame']
pairs = {v: [100*(a['ns_per_frame']/b['ns_per_frame']-1) for a,b in zip(w['candidate'], w[v])] for v in ('parent', 'old691')}
profiles = {v: json.loads((P/'profiles600'/(v+'.json')).read_text()) for v in V}
pf = set(profiles['candidate'])-ignore
assert len(pf) == 22
assert all({k: row[k] for k in pf} == {k: profiles['candidate'][k] for k in pf} and row['allocs_per_frame'] == 0 for row in profiles.values())
def group(prefix):
    cases = [x for k,x in summary['cases'].items() if k.startswith(prefix)]
    return {'median_speedup_range': [min(x['parent_over_candidate'] for x in cases), max(x['parent_over_candidate'] for x in cases)], 'all_paired_speedup_range': [min(min(x['paired_parent_over_candidate']) for x in cases), max(max(x['paired_parent_over_candidate']) for x in cases)], 'old_over_current_range': [min(x['old_over_candidate'] for x in cases), max(x['old_over_candidate'] for x in cases)]}
out = {'baseline': 'exact frozen first saturated-byte TM winner on ac259fce; not merged45 directly', 'parent_production_Go_byte_verified':len(a), 'only_production_change':delta, 'fresh486_raw_samples_all27_medians_pairs_verified':True, 'all9_TM8_cases_all54_pairs_at_least1_5':True, 'TM8':group('8x8/'), 'TM16':group('16x16/'), 'TM4_controls':group('4x4/'), 'whole18_all22_semantics_zero_alloc_verified':True, 'whole_median_ns_per_input':med, 'whole_parent_time_difference_percent':100*(med['candidate']/med['parent']-1), 'whole_old_time_difference_percent':100*(med['candidate']/med['old691']-1), 'whole_all6_paired_time_difference_percent':pairs, 'profiles3_all22_semantics_zero_alloc_verified':True, 'scope':'TM8 independently qualifies for kernel retention. TM16 improves but does not meet1.5median. Keep separate from first-stage measurements; no multiplication/cross-window gap-attribution. Profiles diagnostic only.'}
Path('/private/tmp/govpx-root-tm-left-proof.json').write_text(json.dumps(out, indent=2)+'\n')
print(json.dumps(out, indent=2))
