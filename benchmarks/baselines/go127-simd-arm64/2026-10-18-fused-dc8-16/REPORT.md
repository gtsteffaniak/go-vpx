Publication scope: all raw measurements, profiles, sources, oracle results and proofs are preserved in `raw-evidence.tar.gz`. Executable blobs remain private; every original binary hash, size and private path is retained in `OMITTED-BINARIES.json`. The private full package is unchanged.

# ARM fused8/16 DC prediction retention package

Exact strongest baseline PR44 `46b89f5ab4c6c469755323a61c77fe3a766b4081`; old reference `691cd0512c4829d9ee40d483205dd223f9eb028d`, actual NEON assembly. M4 Max, official Go1.27.2 SIMD, PGO off, native GOMAXPROCS=1. No root repository/PR changes. This package contains complete production Go/assembly source for all three fresh variants, source/Git proofs, relevant independent tests, all benchmark/CLI binaries, commands, raw measurements, full LLVM, failed inline trial, and profiles.

Production overlay: exact ARM fused sum/fill helper and other-architecture false stub; ARM constant selects an exact no-reference FillBlock route in existing n8/n16 wrappers. Both available sources are fully read before stores, preserving aliases. Complete positive bounded destination/ref guards precede SIMD writes; failure returns to the original flow, preserving panic diagnostics and partial writes. The no-ref route invokes the exact original FillBlock on every geometry. n4 production caller exists and is unchanged; its public performance controls are explicitly measured here. AMD normalized n4/generic/public-fixture codegen and ARM n4/generic codegen are identical to44. The initial helper112-entry body is unchanged in the final repair.

Correctness: exhaustive24488 available-sum public calls; full backing/alias/unaligned/minimal/short+sparecap/clippedcap/extreme strides and exact original panic/partial-state comparisons; race/checkptr2; purego; SIMD-disabled fallback; all3 fresh27-case one-shot scalar/full-input/output/canary oracles. All486 fresh native component samples0B/0alloc. The identical15-case supplement adds12 one-neighbor cases plus3 actual n4 controls; its full-state alias tests also pass on all3 references with race/checkptr2.

Initial candidate3round3way screen108samples and no-ref repair3round3way screen108samples are preserved separately and never pooled. Initial16none regression motivated the repaired no-ref route. Failed inline-wrapper cost110 exceeded budget80, added a call; no timing run. Repair's additional branch makes initial both-neighbor path a few percent slower but removes none regressions. No production variant beyond this final repair was tested.

Fresh acceptance:27 public B.Loop cases ×6 balanced permutations ×3 variants,100ms,486 samples. Both8 median1.770–1.811×44, both16 median1.590–1.620×44; every both-neighbor paired ratio>1.5. Both also beat old NEON. None and one-neighbor medians improve44; some8one-neighbor paired samples trail oldNEON. n4 controls remain approximatelyflat vs44 and ~3.2×slowerold; unchanged codegen and this gap are disclosed, not attributed as a win.

| Public case |44 ns|candidate ns|oldNEON ns|44/candidate|old/candidate|
|---|---:|---:|---:|---:|---:|
|n16/both/pitch1312|11.5350|7.1195|7.8175|1.6202×|1.0980×|
|n16/both/pitch16|11.1050|6.9835|7.9500|1.5902×|1.1384×|
|n16/both/pitch37|11.2800|7.0945|8.0055|1.5900×|1.1284×|
|n16/none/pitch1312|6.2610|6.0575|6.9850|1.0336×|1.1531×|
|n16/none/pitch16|5.2470|4.9690|7.1100|1.0559×|1.4309×|
|n16/none/pitch37|6.0995|5.8795|6.9335|1.0374×|1.1793×|
|n8/both/pitch1312|8.1880|4.5210|6.1535|1.8111×|1.3611×|
|n8/both/pitch37|8.1945|4.5605|6.1925|1.7968×|1.3579×|
|n8/both/pitch8|8.2040|4.6350|6.1365|1.7700×|1.3239×|
|n8/none/pitch1312|4.2115|3.5265|4.2600|1.1942×|1.2080×|
|n8/none/pitch37|4.2405|3.6660|4.2505|1.1567×|1.1594×|
|n8/none/pitch8|4.3695|3.7125|4.2780|1.1770×|1.1523×|
|n16/left-only/pitch1312|7.9285|7.0650|7.5320|1.1222×|1.0661×|
|n16/left-only/pitch16|7.7805|6.9170|7.5090|1.1248×|1.0856×|
|n16/left-only/pitch37|7.7915|6.9810|7.3885|1.1161×|1.0584×|
|n16/up-only/pitch1312|7.8890|7.1180|7.7465|1.1083×|1.0883×|
|n16/up-only/pitch16|7.7865|6.8935|7.6115|1.1295×|1.1042×|
|n16/up-only/pitch37|7.9010|7.0455|7.5905|1.1214×|1.0774×|
|n4/both/pitch1312|9.8940|9.8615|3.1145|1.0033×|0.3158×|
|n4/both/pitch37|10.0020|9.8445|3.1035|1.0160×|0.3153×|
|n4/both/pitch4|9.8420|10.1310|3.1035|0.9715×|0.3063×|
|n8/left-only/pitch1312|6.2995|4.7715|4.7710|1.3202×|0.9999×|
|n8/left-only/pitch37|6.4730|4.7450|4.8635|1.3642×|1.0250×|
|n8/left-only/pitch8|6.4255|4.7220|4.9040|1.3608×|1.0385×|
|n8/up-only/pitch1312|6.3485|4.8450|5.0305|1.3103×|1.0383×|
|n8/up-only/pitch37|6.4220|4.7165|5.1460|1.3616×|1.0911×|
|n8/up-only/pitch8|6.4960|5.0335|5.0510|1.2906×|1.0035×|

All6 paired ratios, including every slower control, are preserved in `fresh/fresh-dsp/summary.json`; all18 raw logs and full commands accompany it.

Whole encoder: fresh six balanced300-input VP8 realtime720p cpu8/30fps/2500kbps/1thread runs for all3. Medians44 7.330647ms, candidate7.3022785ms, old6.901636ms. Candidate−0.386985%vs44 (5/6 faster; paired changes+0.17919,−0.24515,−1.45844,−0.40656,−0.59636,−2.09108%). Candidate remains+5.80504%slowerold (all6 slower; paired+2.78545,+3.45116,+9.02898,+2.91534,+8.21586,+8.61933%). All22 semantic fields and zero allocations match across18 JSONs. This is a retained component win; whole-encoder improvement is small/noisy and the assembly gap remains open.

Separate600-input profiles: all3 same binaries/options with raw pgo,22semantic/zeroalloc equality, flat/cumulative top50 text and commands. They are outside timed pairs and are diagnostic only. Candidate LF fusedV7.11%flat/12.32%cum; fusedH6.87%flat/9.24%cum; pickLFlevel29.15%cum; tokenwrite5.45%flat; memmove4.03%flat. Old's LFvertical12.44%, DCpair8.85%, LFhorizontal4.07%, SSE16dot3.83%flat. Profile percentages are sampling estimates and not comparative timing denominators.

The original frozen FillBlock archives remain unchanged. The new fused helper uses those exact proven stores; no blanket ARM or x64 encoder completion claim is made. All binaries are retained only as private evidence, not intended for repository inclusion.
