#!/usr/bin/env python3
"""Post-timing, current/base/old native code provenance; never measures throughput."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

variant, binary_arg, output_arg = sys.argv[1:]
assert variant in {"head", "base", "assembly"}
binary = Path(binary_arg).resolve()
output = Path(output_arg).resolve()
root = Path.cwd()
wrapper = str(root / "scripts/go.sh")

def run(args):
    p = subprocess.run(args, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if p.returncode:
        raise RuntimeError(f"command failed ({p.returncode}): {args!r}\n{p.stderr}")
    return p.stdout

gnu_version = run(["objdump", "--version"])
assert gnu_version.startswith("GNU objdump"), "GNU binutils objdump required; LLVM/Go decoders are not used"
expected = (output / f"{variant}-binary-sha256.txt").read_text().split()[0]
actual = hashlib.sha256(binary.read_bytes()).hexdigest()
assert actual == expected, "executable differs from the measured binary"
probe = subprocess.run([wrapper, "tool", "-n", "nm"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
nm_cmd = [wrapper, "tool", "nm"] if probe.returncode == 0 else [wrapper, "run", "cmd/nm"]
nm_args = nm_cmd + ["-size", "-sort", "address", str(binary)]
full_map = run(nm_args)
selected = []
for line in full_map.splitlines():
    m = re.match(r"^\s*([0-9a-fA-F]+)\s+(\d+)\s+([Tt])\s+(.+)$", line)
    if not m:
        continue
    address, size, kind, name = m.groups()
    keep = False
    for package, expression in [
        ("github.com/thesyncim/govpx/internal/pixel.", r"(?i)^sse8"),
        ("github.com/thesyncim/govpx/internal/vp8/dsp.", r"(?i)^(sse8|varianceblock8|benchmarksse8)"),
    ]:
        if name.startswith(package) and re.match(expression, name[len(package):]):
            keep = True
    if keep:
        selected.append({"address": address, "size": int(size), "kind": kind, "name": name, "nm_line": line})
assert selected, "no relevant Go text symbols selected"
names = {entry["name"] for entry in selected}
caller = "github.com/thesyncim/govpx/internal/vp8/dsp.BenchmarkSSE8RowsNative.func1"
assert any(name.startswith(caller) for name in names), "actual live B.Loop caller missing"
if variant != "assembly":
    assert "github.com/thesyncim/govpx/internal/pixel.SSE8Ptr" in names, "current/base pointer helper missing"
else:
    assert any("varianceBlock8xNSSE2" in name for name in names), "old pointer assembly kernel missing"
assert sum(entry["size"] for entry in selected) < 1_000_000, "unexpectedly broad code selection"
symbols = output / f"{variant}-codegen-symbols.txt"
symbols.write_text("\n".join(entry["nm_line"] for entry in selected) + "\n")
code = output / f"{variant}-codegen.txt"
with code.open("w") as dest:
    dest.write(f"# measured executable SHA256 {actual}\n# Go nm command {nm_args!r}\n")
    for entry in selected:
        text = run(["objdump", "-d", "--disassemble=" + entry["name"], str(binary)])
        assert f"<{entry['name']}>:" in text, f"GNU decoder did not find {entry['name']}"
        assert re.search(r"^\s*[0-9a-f]+:\s+[0-9a-f]{2}\s", text, re.M), "empty machine-code body"
        entry["gnu_disassembly_sha256"] = hashlib.sha256(text.encode()).hexdigest()
        dest.write("\n# selected symbol " + entry["name"] + "\n" + text)
optional = ["github.com/thesyncim/govpx/internal/vp8/dsp.SSE8x8PtrFast", "github.com/thesyncim/govpx/internal/vp8/dsp.sse8x8PtrFast"]
metadata = {
    "variant": variant,
    "measured_executable_sha256": actual,
    "nm_command": nm_args,
    "symbols": selected,
    "optional_wrapper_status": {name: "emitted" if name in names else "not emitted (inlined or dead-stripped); actual B.Loop caller and kernel captured" for name in optional},
    "filtered_nm_sha256": hashlib.sha256(symbols.read_bytes()).hexdigest(),
    "codegen_sha256": hashlib.sha256(code.read_bytes()).hexdigest(),
    "post_timing_only": True,
    "no_performance_threshold": True,
}
(output / f"{variant}-codegen.json").write_text(json.dumps(metadata, indent=2) + "\n")
tools = output / "codegen-tools.txt"
if not tools.exists():
    tools.write_text(gnu_version + "\n" + run([wrapper, "version"]) + "\n")
with tools.open("a") as dest:
    dest.write(f"{variant} Go nm command: {nm_args!r}\n")
assert hashlib.sha256(binary.read_bytes()).hexdigest() == actual, "measured binary mutated during capture"
