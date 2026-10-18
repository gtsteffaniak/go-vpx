Publication scope: all raw measurements, profiles, sources, oracle results and proofs are preserved in `raw-evidence.tar.gz`. Executable blobs remain private; every original binary hash, size and private path is retained in `OMITTED-BINARIES.json`. The private full package is unchanged.

# ARM4x4 DC fused sum/fill retention evidence

Strongest baseline is the exact frozen D8 winner in `/private/tmp/govpx-arm-fused-dc-ready` (source derived from PR44 `46b89f5ab4c6c469755323a61c77fe3a766b4081`, plus measured D8 overlay). This n4 candidate adds exactly3 production files: ARM helper, other-architecture false stub, and hook in existing4x4-only dispatch. It leaves all D8 production bytes unchanged. Baseline source exactness is proved against the immutable earlier package, with its Git source proof retained. Old actualNEON reference is `691cd0512c4829d9ee40d483205dd223f9eb028d`.

`overlay/` provides5 integration files:3 production and2 independent state/oracle tests. All27 public benchmark fixture bytes remain identical across current/baseline/old; they are retained under source/ for all3. The frozen prior D8 package, initial243sample n4 screen, failed inline attempt, and fresh acceptance samples remain separate.

Algorithm: guard Enabled, unsigned stride<2^30, complete4-row destination length and both4-byte reference lengths before any unsafe reads/stores. Exact uint32 reference loads occupy two packet lanes, widen8bytes to unsigned16, single reduction max2040, round(sum+4)>>3, broadcast and4 exact uint32 row stores. All references are consumed before first write, preserving destination/ref aliases and overlapping rows. Any decline returns to exact previous flow. Public above/left/destination index checks remain unchanged and in original order, preserving panic/partial-state behavior and capacity. Other-architecture stub folds away.

Correctness:2041 attainable sums×16align directpixel; full backing/alias/canaries/minimal windows/short+sparecap/clippedcap/extreme strides, independently short above/left and disabledstate; actual public frozen original prechecks+generic path as panic/partial-state oracle; race/checkptr2, purego,nosimd. Sixteen initial build/validation commands pass. New fresh all3 binaries rebuilt with officialGo1.27.2SIMD and PGOoff; all27 one-shot full-output/input/canary oracles pass. Optimized ARM helper37entries, no calls/vector spills, oneADDV. AMD n4/generic and both public closures normalized identical. ARM D8 helper/wrappers/generic and both public closures normalized identical too. All source and full LLVM proofs retained.

One authorized inline-wrapper attempt failed: exported guardwrappercost93>budget80, uncheckedworker96>80. It was stopped after compile, no timing run. Its source, build command, full compiler output, and binary are retained; initial37-entry measured source was never altered. No other n4 production variant.

Native fresh acceptance: M4Max, officialGo1.27.2, SIMD, PGOoff, GOMAXPROCS1,27 actual public B.Loop cases100ms ×6balancedpermutations×3variants=486samples, all0B/0alloc with output/input/canary oracles. Separate provisional3round75ms screen has243samples; never pooled.

Fresh n4 medians2.4694–2.5065×strongD8baseline; every18n4pairedratio>1.5 (2.3890–2.6886). The new kernel still takes27.3–28.2%more time than oldNEON in these public n4 cases. All24 D8 controls retained, including8up-onlypackedmedian0.9568×baseline, with paired0.9181–1.0390. Equal source/codegen does not prove the cause of control timing variation.

|Public case|D8baseline ns|n4winner ns|oldNEON ns|baseline/new|old/new|
|---|---:|---:|---:|---:|---:|
|n16/both/pitch1312|7.0860|7.1170|7.8980|0.9956×|1.1097×|
|n16/both/pitch16|6.9155|6.9635|7.8835|0.9931×|1.1321×|
|n16/both/pitch37|7.0295|7.0385|7.8525|0.9987×|1.1156×|
|n16/none/pitch1312|6.0325|6.0480|7.0695|0.9974×|1.1689×|
|n16/none/pitch16|4.9630|4.9830|6.9150|0.9960×|1.3877×|
|n16/none/pitch37|5.8755|5.8595|6.9680|1.0027×|1.1892×|
|n8/both/pitch1312|4.5570|4.4935|6.1400|1.0141×|1.3664×|
|n8/both/pitch37|4.5770|4.5720|6.1550|1.0011×|1.3462×|
|n8/both/pitch8|4.6855|4.7005|6.1195|0.9968×|1.3019×|
|n8/none/pitch1312|3.5785|3.6180|4.2160|0.9891×|1.1653×|
|n8/none/pitch37|3.6470|3.5960|4.2405|1.0142×|1.1792×|
|n8/none/pitch8|3.6875|3.5215|4.1790|1.0471×|1.1867×|
|n16/left-only/pitch1312|7.0200|7.0315|7.4490|0.9984×|1.0594×|
|n16/left-only/pitch16|6.8760|6.8865|7.4345|0.9985×|1.0796×|
|n16/left-only/pitch37|6.9675|6.9220|7.4160|1.0066×|1.0714×|
|n16/up-only/pitch1312|7.0500|7.0880|7.6320|0.9946×|1.0767×|
|n16/up-only/pitch16|6.8775|6.8660|7.5930|1.0017×|1.1059×|
|n16/up-only/pitch37|6.9675|6.9590|7.6785|1.0012×|1.1034×|
|n4/both/pitch1312|9.7900|3.9645|3.1060|2.4694×|0.7835×|
|n4/both/pitch37|9.8765|3.9715|3.1140|2.4868×|0.7841×|
|n4/both/pitch4|9.9510|3.9700|3.0970|2.5065×|0.7801×|
|n8/left-only/pitch1312|4.7095|4.7250|4.8705|0.9967×|1.0308×|
|n8/left-only/pitch37|4.7245|4.7500|4.8125|0.9946×|1.0132×|
|n8/left-only/pitch8|4.7245|4.7140|4.9280|1.0022×|1.0454×|
|n8/up-only/pitch1312|5.0295|4.8780|4.9795|1.0311×|1.0208×|
|n8/up-only/pitch37|4.8420|4.8270|5.0115|1.0031×|1.0382×|
|n8/up-only/pitch8|4.7890|5.0050|5.0185|0.9568×|1.0027×|

All6 paired ratios and every slower control are in fresh-dsp/summary.json, alongside all18rawlogs, orders, command timestamps and binary hashes.

Whole300-input VP8 realtime720p cpu8/30fps/2500kbps/1thread, sixbalanced3way runs: D8parent7.3279565ms, winner7.2998625ms, old6.7012025ms. Winner−0.383381%vsD8,4/6faster (paired+2.35241,−0.67308,−2.02767,+0.32316,−0.75795,−1.14069%). Winner+8.93362%slowerold,0/6faster (paired+5.55240,+9.05257,+9.86946,+9.67727,+8.64257,+7.75046%). All22 semantic fields and zeroalloc match across18JSONs. This is a retained2.5×component win; whole improvement is small/noisy and oldassembly gap remains open.

Separate600-input profiles for all3 use the exact same CLI binaries/options, outside timed pairs. Rawpgo, flat/cumulative top50, allcommands/JSON/stderr are retained. All22semantic and zeroalloc equality pass across3profileJSONs. These profiles are diagnostic only, not timing denominators. Binary blobs are private evidence; lean publication copy omits only executable blobs and retains their exact hashes/sizes/privatepaths. Full private package stays immutable.
