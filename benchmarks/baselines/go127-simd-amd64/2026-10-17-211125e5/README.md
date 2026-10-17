# Native AMD 211125e5 ARM-checkpoint control

Hosted CI: https://github.com/thesyncim/govpx/actions/runs/39029805350
Native run: https://github.com/thesyncim/govpx/actions/runs/39029805450
Raw artifact: https://github.com/thesyncim/govpx/actions/runs/39029805450/artifacts/12100565978
Both workflows pass on exact head 211125e5a4987a351b39f2e54a9ec77f9373f34b.

All 85 raw files are retained unchanged. summary.json, complete-parity.json
and checks.json add analysis and exact-head status. Environment: native
EPYC 7763, official Go 1.27.2, SIMD / AVX2, GOAMD64=v1, GOMAXPROCS=1,
PGO off. Variants are original ASM 691cd051, parent 611f47fe and current
211125e5. Six rounds cover every variant permutation.

| Workload | ASM ms/input | Parent ms/input | Current ms/input | Current vs parent | Current vs ASM |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 7.050882 | 4.344040 | 4.337680 | -0.146% | -38.480% |
| VP8 720p | 23.293487 | 22.197201 | 22.124264 | -0.329% | -5.020% |
| VP9 720p | 59.969839 | 32.569408 | 32.311997 | -0.790% | -46.120% |

All parent/current paired ranges cross zero. VP9 320p improves in 2/6
pairs (-0.717..+0.378% time), VP8 improves in 4/6 (-1.679..+2.129%),
and VP9 720p improves in 4/6 (-4.372..+3.167%). These are mixed controls,
not evidence of a new x64 gain from an ARM-targeted change. Every pair
against original ASM is faster for each public workload. Absolute times
across runner jobs do not isolate code changes.

All 22 non-timing semantic/configuration fields agree across all 18 samples
per workload and across corresponding profile runs. Allocation samples also
agree throughout: VP9 320p 0.616667, VP8 zero and VP9 720p 1.25 per input;
profile allocations are 0.625 / 0 for VP9 / VP8. All 75 kernel cases have
six samples per variant and zero bytes/allocations per call. No allocation
outlier is discarded. Independent hosted byte-parity/oracle gates pass.

The preceding 611f47fe arithmetic experiment retains its own measurements.
The EPYC 9V45 VP8 limit and ARM ASM/libvpx gaps remain separate evidence;
this control establishes no universal CPU or codec-performance claim.
