# VP9 non-RD large-model duplicate pass — preliminary screen

Exact baseline: `7e0e37b7bbbcda65523f5581e7879c7d2d93827c`.
Private production patch: `model-once.patch`. No root production/Git edits.

The filter sweep and ordinary non-RD candidate previously ran predictor +
whole-block variance/SSE + `ModelRdForSbY`, then ran `ModelRdForSbYLarge` and
replaced every small-model output. Pinned libvpx `vp9_pickmode.c:1527–1538`
and `2336–2348` select one of these models. The large model already gathers
whole-block variance/SSE while producing 8×8 transform-skip statistics.

V2 defers that first statistics pass only for prepared, fully visible source
blocks. The existing large-model call stays in the picker. Clipped and
unprepared cases retain original prediction/statistics; rejected large-model
windows recover small-model statistics without repeating prediction. Chroma
skip/early-termination order is unchanged.

Native Apple M4 Max, Darwin arm64, Go 1.27.2, SIMD, CGO=1, GOMAXPROCS=1.
Two balanced short A/B and B/A samples. These are preliminary, not a broad
performance conclusion. No C comparison, profile, full suite or matrix run.

| V2 component | Pair 1 speedup | Pair 2 speedup |
|---|---:|---:|
| 32×32 fullpel | 1.190× | 1.203× |
| 64×64 fullpel | 1.164× | 1.168× |
| 32×32 subpel | 1.112× | 1.079× |
| 64×64 subpel | 1.106× | 1.054× |
| 32×32 clipped | 0.984× | 0.985× |
| 64×64 clipped | 1.003× | 1.020× |
| 32×32 legacy | 0.982× | 0.984× |
| 64×64 legacy | 1.019× | 0.955× |

All component cases reported zero bytes and zero allocations per operation.

Canonical Go-only VP9 CPU8 CBR, one thread, encode-only, PGO off:

| Workload | Pair 1 speedup | Pair 2 speedup |
|---|---:|---:|
| 320×180, 120 frames, 500 kbps | 0.989× | 0.996× |
| 1280×720, 60 frames, 2500 kbps | 0.987× | 0.990× |

The component clears the user's 15% criterion, but this is **not an ARM codec
win**. Retained independently for requested x64 screening; losses are included.
All non-timing CLI fields, including bytes, frame/drop counts, quantizers and
allocations, matched within each A/B pair.

Correctness: frozen old-order model outputs and predictor bytes matched across
prepared/fullpel/subpel/clipped/legacy/invalid-model-window/short-buffer/invalid-
filter cases; existing prepared-predictor oracle passed. Exact 96 VP8 and
72 VP9 ML packet stream hashes matched all nine baseline configurations.
`candidate-v2-tests.log`, `baseline-stream.log`, and `binding.json` retain proof.

V1 moved the large-model argument/result through another helper call and was
rejected. Most component controls slowed; codec results were inconsistent.
Its exact sources and raw logs are in `rejected-v1/`, with all results retained
in `summary.json`. V2 raw results are in `native-v2/` and commands in its
`commands.json`. Source and binary hashes are in `binding.json`.

`publish/manifest.json` + `publish/overlays.tar.gz` are an independent archive
bundle for native x64 variant 201. Its shared fixture compiles against both
baseline and candidate through an optional internal method interface. Both
shared-fixture correctness checks passed; this adapter is separately measured
on x64, and is not substituted for the exact local V2 fixture in the binding.
