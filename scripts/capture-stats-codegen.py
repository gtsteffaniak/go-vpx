#!/usr/bin/env python3
"""Preserve actual native Stats code/layout only after optional timing ends."""
from pathlib import Path
import argparse, hashlib, json, platform, re, subprocess

OWNERS = (
    'Stats', 'StatsScalar', 'statsTallSIMD', 'pixelStatsSIMD',
    'pixelStatsPairs8', 'pixelStatsRows16', 'pixelStatsQuads4',
    'pixelStats16Narrow', 'pixelStatsPureGo', 'reducePixelAcc',
    'BenchmarkStatsReductionNative', 'BenchmarkStatsReductionNative.func1',
)
PREFIX = 'github.com/thesyncim/govpx/internal/pixel.'
FIXTURES = ('fp32_optional_native_test.go', 'stats_reduction_native_test.go')

def digest(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()

def symbols(text):
    found = {}
    for line in text.splitlines():
        m = re.fullmatch(r'\s*([0-9a-fA-F]+)\s+(\d+)\s+([Tt])\s+(\S+)\s*', line)
        if not m:
            continue
        address, size, kind, name = m.groups()
        if name in {PREFIX + n for n in OWNERS}:
            assert name not in found, ('duplicate symbol', name)
            value = int(address, 16)
            found[name] = {'address_hex': hex(value), 'size_bytes': int(size),
                           'address_mod16': value % 16, 'address_mod32': value % 32,
                           'address_mod64': value % 64, 'kind': kind}
    return found

def bind_protocol(text, binary):
    binary = Path(binary).resolve()
    entries = []
    for line in text.splitlines():
        m = re.fullmatch(r'([0-9a-f]{64})\s+(.+)', line)
        if m and Path(m.group(2)).resolve() == binary:
            entries.append(m.group(1))
    assert entries and set(entries) == {digest(binary)}, 'binary/protocol SHA mismatch'
    feature = re.findall(r'FP32PublicHost AVX512=(true|false) AVX2=(true|false) Enabled=(true|false)', text)
    assert feature, 'actual SDK feature probe missing'
    orders = re.findall(r'^# round (\d+) order (parent current|current parent)$', text, re.M)
    assert orders == [(str(r), 'parent current' if r % 2 else 'current parent') for r in range(1, 7)], 'order mismatch'
    assert 'optional pixel inventory PASS:' in text, 'strict optional parser did not finish'
    return {'actual_feature_probes': feature, 'orders': orders,
            'protocol_lines': [line for line in text.splitlines() if line.startswith(
                ('# head ', '# base ', '# driver checkout ', '# SDK=', '# SELECTED_',
                 '# round ', 'Model name:', 'Flags:', 'go version '))]}

def main():
    parser = argparse.ArgumentParser()
    for name in ('variant', 'binary', 'output-dir', 'source-dir', 'common-source-dir',
                 'source-commit', 'head-commit', 'base-commit', 'protocol-log', 'go-wrapper'):
        parser.add_argument('--' + name, required=True)
    parser.add_argument('--objdump', default='objdump')
    args = parser.parse_args()
    assert platform.system() == 'Linux' and platform.machine() == 'x86_64', 'native Linux/x64 capture required'
    assert args.variant in ('parent', 'current')
    out = Path(args.output_dir); out.mkdir(parents=True, exist_ok=True)
    commands = []
    def run(argv, filename):
        proc = subprocess.run(argv, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        (out / filename).write_text(proc.stdout)
        commands.append({'argv': argv, 'exit_code': proc.returncode, 'output': filename})
        (out / 'commands.json').write_text(json.dumps(commands, indent=2) + '\n')
        assert proc.returncode == 0, argv
        return proc.stdout
    binary = Path(args.binary).resolve(); log = Path(args.protocol_log).resolve()
    protocol = bind_protocol(log.read_text(), binary)
    tool = run([args.objdump, '--version'], 'objdump-version.txt')
    assert tool.startswith('GNU objdump'), 'GNU objdump required'
    run([args.go_wrapper, 'version', '-m', str(binary)], 'binary-buildinfo.txt')
    nm = run([args.go_wrapper, 'tool', 'nm', '-size', '-sort', 'address', str(binary)], 'symbols-all.txt')
    active = symbols(nm)
    assert PREFIX + 'BenchmarkStatsReductionNative.func1' in active, 'actual B.Loop caller missing'
    (out / 'symbols.json').write_text(json.dumps(active, indent=2) + '\n')
    for name in active:
        # Raw bytes and absolute branch/call addresses are retained, not normalized.
        run([args.objdump, '-d', '--disassemble=' + name, str(binary)],
            name.removeprefix(PREFIX) + '.asm.txt')
    source = Path(args.source_dir)
    actual = {str(p.relative_to(source)): digest(p) for p in sorted((source / 'internal/pixel').glob('*.go'))}
    common = Path(args.common_source_dir)
    fixture_sha = {}
    for name in FIXTURES:
        rel = Path('internal/pixel') / name
        assert (source / rel).read_bytes() == (common / rel).read_bytes(), 'fixture mismatch'
        fixture_sha[str(rel)] = digest(source / rel)
    metadata = {'scope': 'Actual native executable layout; captured after ALL optional timers/parser ended. No binary execution.',
                'variant': args.variant, 'binary': str(binary), 'binary_sha256': digest(binary),
                'binary_size': binary.stat().st_size, 'source_commit': args.source_commit,
                'head_commit': args.head_commit, 'base_commit': args.base_commit,
                'actual_source_SHA256': actual, 'identical_common_fixture_SHA256': fixture_sha,
                'protocol_log': str(log), 'protocol_SHA256': digest(log), 'ISA_and_orders': protocol,
                'capture_helper_SHA256': digest(__file__), 'live_symbols': active,
                'absent_may_be_inlined': [PREFIX + n for n in OWNERS if PREFIX + n not in active]}
    (out / 'metadata.json').write_text(json.dumps(metadata, indent=2) + '\n')

if __name__ == '__main__':
    main()
