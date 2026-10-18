# VP8 accepted coefficient work elimination

Runtime commit5274ea78b8cc; parent715645fd. Apple M4 Max / Go1.27.2 / simd / CGO1 / PGOoff. Short native measurements are preliminary. No full CI, new C comparison or profile audit.

The accepted fast-quant pass previously generated DQCoeff for 24 luma/chroma blocks although only RD distortion scoring reads those arrays. Reconstruction later dequantizes QCoeff again. The new ARM/x64 SIMD entry omits this discarded multiplication and stores; generic/signed-table/alias fallbacks retain existing per-block semantics. RD scoring still produces DQCoeff. EOBs are written directly into the authoritative macroblock side channel, removing staging and an otherwise empty block/context walk. No coding-control or quality-setting change.

## Component results

Two short samples per case, zero allocations. Shared input and output consumption keep the work live. Eight-block quantization is the representative chroma batch; 16-block luma gains are smaller and retained explicitly.

| Batch/pattern | Old ns samples | Quant-only ns samples | Median speedup |
|---|---|---|---:|
| 8/zero | [17.04, 14.45] | [11.98, 11.87] | 1.320× |
| 8/sparse | [14.97, 14.94] | [11.93, 11.93] | 1.254× |
| 8/dense | [15.08, 14.6] | [12.27, 12.15] | 1.215× |
| 16/zero | [25.78, 25.41] | [23.29, 23.42] | 1.096× |
| 16/sparse | [25.94, 25.71] | [23.4, 23.35] | 1.105× |
| 16/dense | [25.75, 25.54] | [23.76, 23.38] | 1.088× |

Portable accepted-MB pipeline fixture (gather + FDCT + quantization + EOB, identical source fixture on parent/candidate), two balanced80ms pairs:

| Pattern | Parent ns samples | Candidate ns samples | Paired speedups |
|---|---|---|---|
| zero | [167.0, 144.5] | [144.8, 125.3] | 1.153×, 1.153× |
| dc | [162.9, 159.1] | [140.1, 141.2] | 1.163×, 1.127× |
| dense | [217.6, 217.9] | [199.2, 200.2] | 1.092×, 1.088× |

## End-to-end result and limits

Go-only VP8 1280×720,60inputs,30fps,2500kbps,realtime,threads1/noise4; fixed-speed cpu=-8 and adaptive cpu=8. Two balanced samples per host/workload/reference. All22 non-performance fields and zero-allocation observations match. Same-host immutable fixed baseline c4a8bd9f is measured independently of latest-main parent. Warmups and the earlier exploratory screen remain in raw evidence and are not pooled.

| Mode | Reference | Paired throughput speedups | Median |
|---|---|---|---:|
| cpu-8 | fixed-c4 | 1.0239×, 1.0381× | 1.0310× |
| cpu-8 | main | 0.9953×, 1.0150× | 1.0052× |
| cpu8 | fixed-c4 | 1.0168×, 1.0362× | 1.0265× |
| cpu8 | main | 1.0155×, 1.0134× | 1.0145× |

Latest-main fixed-speed throughput remains flat; adaptive improvement in the final short screen is about1.4%, below the2% codec target. The change is retained for the21–32% eight-block kernel gain and15.3% zero-residual accepted-pipeline gain. These figures do not establish closure of the libvpx gap. x64 native correctness/performance feedback is pending at publication; no x64 speed claim is made. The1.5× fixed-baseline codec gate remains unmet.

## Correctness and iteration ledger

- Full signed coefficients (including−32768), signed table fallback, counts1/2/8/16/25, zero/sparse cases, untouched tails and coefficient/output overlaps compare against the independent scalar quantizer or the unchanged overlap/store-order entry. Native and pure-Go focused checks pass.
- Existing full/chroma RD equivalence check passes;96 VP8 packet/flag hashes and72 root ML packet hashes match their immutable expectations.
- An initial pipeline fixture omitted required coefficient probabilities and measured the early-return guard. Those4ns samples are preserved as `invalid-nil-probs-*` and excluded. The corrected shared fixture supplies probabilities and asserts live nonzero coefficient output before timing.
- First Go-only exploratory screen: fixed paired1.006/0.989×, adaptive1.028/1.017×. Final source-bound screen above supersedes it; no pooled2% claim. All regressions/noise remain visible.
- README unchanged. Only affected correctness, compact component and Go-only codec screens ran. No automatic full validation was dispatched.

Raw logs, command manifests, invalid experiments, warmups, final whole outputs, source/binary bindings and tests are preserved beside this report. PR fast CI can additionally execute the `TestVP8WorkElisionQuantOnly*` tests on x64 through the existing affected-test selector.
