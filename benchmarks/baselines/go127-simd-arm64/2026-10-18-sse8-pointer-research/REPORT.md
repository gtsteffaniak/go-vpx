# Rejected ARM trusted-pointer SSE8 probes

All three private ARM pointer designs are rejected. None reaches the user’s 1.5× component threshold. The final balanced unsigned-byte-square interaction gives 1.320181–1.330271× the parent and is 0.544–1.784% slower than the direct-return control in the same run. The old NEON build remains faster. No six-round whole-matrix confirmation or whole-codec timing was started for these nonqualifying designs. ARM SSE8 research stops here.

## Pinned scope

Parent is `ab3abd6b335242d5732fa1956a63b7591eaf781b`; old assembly build is `691cd0512c4829d9ee40d483205dd223f9eb028d`. Native hardware is Apple M4 Max, compiler Go 1.27.2 darwin/arm64, SIMD enabled and PGO disabled. Builds use the root project wrapper with private `-C` snapshots, modest build concurrency and distinct agent IDs. Timing is serial with GOMAXPROCS=1 and test.cpu=1. Parent sources and all rejected snapshots remain outside the root checkout; no root, Git, dependency or README changes were made.

The public trusted-window operation is DSP.SSE8x8PtrFast, used by macroblockChromaSSE in the VP8 encoder. Inputs must already describe valid 8×8 windows. Read-only aliases, unaligned exact eight-byte rows and repeated/overlapping rows are supported. Valid source/reference strides are nonnegative and below 2^30. All other cases retain the original unsafe.Slice construction and pixel.SSE8 call, in its original argument order. Invalid tiny pointers with otherwise eligible large strides violate that trusted contract and were not dereferenced as guard tests.

## Three isolated designs

1. **Bool helper:** fixed eight rows as four explicit paired-row packets; original byte Max-Min, signed widening and square arithmetic. New ARM pixel.SSE8Ptr returns (uint32,bool), and the DSP wrapper uses the original fallback after a decline. The nine-case screen improves 1.318640–1.343107× the parent. Optimized helper has no bounds checks/calls; the wrapper saves four arguments around the bool helper call.
2. **Direct-return helper:** same packet arithmetic and row loads, returning uint32 only. Unsupported cases perform the exact original fallback internally; DSP forwarding inlines away. The initial nine-case screen improves 1.328139–1.356500× the parent. Removing the bool ABI did not materially reduce measured public latency, so the ABI is not established as the limiting cost.
3. **Unsigned-square interaction:** same direct API, guard, loads and unroll, replacing four signed widening multiplies per packet with two unsigned byte-to-uint16 products followed by zero widening and uint32 accumulation. Each square is at most 65025; zero widening is essential because that value does not fit signed int16. The entire 8×8 sum is at most 4,161,600, safely within int32. Codegen has eight UMULL and no SMULL, versus sixteen SMULL in the direct kernel, but both complete helpers have 169 instructions including cold fallback paths. That reduction does not establish a speedup.

The bool and direct initial screens each contain one native sample per variant/case and are screens only. The final interaction uses four balanced Williams orders, comparing parent/direct/unsigned/oldbuild in every order position. It preserves four samples per case per variant. No favorable-only reruns or cross-CPU subtraction are used.

## Final balanced public pointer screen

| Case | Parent ns | Direct ns | Unsigned ns | Old NEON ns | Parent / unsigned | Unsigned vs direct |
|---|---:|---:|---:|---:|---:|---:|
| pointer8x8/pitch8-8/pattern0 | 5.9980 | 4.4535 | 4.5155 | 4.2215 | 1.328314× | +1.392163% |
| pointer8x8/pitch8-8/pattern1 | 5.9915 | 4.4590 | 4.5260 | 4.2145 | 1.323796× | +1.502579% |
| pointer8x8/pitch8-8/pattern2 | 5.9890 | 4.4570 | 4.5365 | 4.2180 | 1.320181× | +1.783711% |
| pointer8x8/pitch32-37/pattern0 | 6.0120 | 4.5020 | 4.5265 | 4.2120 | 1.328179× | +0.544203% |
| pointer8x8/pitch32-37/pattern1 | 5.9835 | 4.4560 | 4.5090 | 4.2260 | 1.327013× | +1.189408% |
| pointer8x8/pitch32-37/pattern2 | 6.0100 | 4.4745 | 4.5295 | 4.2075 | 1.326857× | +1.229188% |
| pointer8x8/pitch1312-1320/pattern0 | 6.0120 | 4.4795 | 4.5490 | 4.2170 | 1.321609× | +1.551512% |
| pointer8x8/pitch1312-1320/pattern1 | 6.0095 | 4.4765 | 4.5175 | 4.2270 | 1.330271× | +0.915894% |
| pointer8x8/pitch1312-1320/pattern2 | 5.9845 | 4.4735 | 4.5115 | 4.2245 | 1.326499× | +0.849447% |

All 198 timed pointer samples across the three screens report 0 B/op and 0 allocs/op. The existing 36-case public fixture is copied identically across refs and passes its one-iteration scalar-result, complete input immutability and backing-canary checks. Its other 27 slice cases are correctness controls only in this experiment; they were not performance-measured after the pointer screens failed acceptance. The complete public-call fixture uses testing.B.Loop and consumes the result. Oldbuild dispatch selects its NEON or DotProd kernel according to its existing CPU capability logic; no forced feature override was introduced.

## Correctness and fallback evidence

Bool/direct ordinary and race+checkptr tests cover minimal allocations, all 16×16 source/reference alignment combinations, zero/overlapping/unequal/padded strides, maximum errors and read-only aliases. Bool decline tests check nil/negative/huge inputs without dereferencing invalid windows. Direct-return tests verify nine nil/negative/overflow fallback failures in child processes, retaining ordinary unsafe.Slice panic and unrecoverable checkptr diagnostics. The unchanged parent also passes that checkptr failure test. Purego and nosimd SSE/variance/public-row tests pass for both APIs.

The unsigned interaction retains all direct-return bounds/alias/fallback tests and passes ordinary plus race/checkptr. An additional exhaustive 65,536-byte-pair test covers every byte-square result through the direct pointer entry. Only the ARM arithmetic and ARM test differ from the direct snapshot; generic fallback source is identical, so its experiment-off results are reused rather than rerun. No correctness or allocation gate was weakened.

## Evidence

Each variant directory preserves exact relevant sources, raw screen logs, one-iteration fixture logs, correctness/checkptr/fallback logs and generated-code evidence. `compact.json` retains every sample, all medians, pair ranges and regressions. `baselines/` records relevant parent/oldbuild code and fixture context; `unsigned-square-prior-source.go.txt` records the copied prior arithmetic source without introducing a dependency. `binary-sha256.json` pins the compiled binaries without copying them. The reproduction driver reconstructs pinned source snapshots, applies these private overlays and runs the same screening protocols through the root wrapper.

`verification.json` independently counts/validates all timed rows, allocation fields, copied source bytes and fixture checks. `SHA256SUMS` covers all other package files. The package is frozen after verification. No new Go jobs or timings were run during packaging.
