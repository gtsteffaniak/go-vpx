# Native AMD 611f47fe FP class-index follow-up

Run: https://github.com/thesyncim/govpx/actions/runs/39028739721
Raw artifact: https://github.com/thesyncim/govpx/actions/runs/39028739721/artifacts/12099567466
Hosted CI: https://github.com/thesyncim/govpx/actions/runs/39028739727
Both workflows pass on exact head 611f47fea29c3e024833dbf5dcbd43c751f59282.

All 85 raw files are preserved unchanged in this directory. Added summary.json,
complete-parity.json and checks.json retain analysis and exact-head check status.
Environment: native EPYC 7763, official Go 1.27.2, AVX2 without AVX512,
GOAMD64=v1, GOMAXPROCS=1, GOEXPERIMENT=simd, PGO off.
Variants are original assembly 691cd051, parent 63c2f188 and current 611f47fe.
Six rounds cover every variant permutation. No Rosetta timing is used.

| Public workload | Assembly ms/input | Parent ms/input | Current ms/input | Current/parent | Current/assembly |
| --- | ---: | ---: | ---: | ---: | ---: |
| VP9 320x180 | 7.023516 | 4.427351 | 4.360532 | -1.509% | -37.915% |
| VP8 720p | 23.307879 | 22.410145 | 22.601558 | +0.854% | -3.030% |
| VP9 720p | 60.774650 | 33.563238 | 32.827483 | -2.192% | -45.985% |

VP9 320p is faster than parent in 5/6 pairs: -2.510..+0.155% time.
Its paired mean is -1.270%; a descriptive paired t95 interval is
[-2.253,-0.288]%, treating six balanced rounds as independent observations.
This limited experiment supports a gain; it is not a universal CPU/layout fix.
VP9 720p improves in 4/6 pairs (-3.388..+1.634%); its paired interval crosses zero.
VP8 control improves in 2/6 pairs (-1.089..+2.472%); its interval also crosses zero.
No VP8 change is attributed to the VP9-only arithmetic change.
Current is faster than assembly in all 6 pairs for each public workload.

All 22 non-timing semantic/configuration fields agree across all 18 samples
per workload, and across corresponding profile runs. Allocations also agree
in every sample: VP9 320p 0.616667, VP8 0, VP9 720p 1.25 per input; profile allocation
samples 0.625 / 0 for VP9/VP8 agree. There is no allocation outlier in this run.
Exact bitstreams are independently checked by hosted oracle gates.

All 75 kernel cases have 6 samples per variant and 0 B/op / 0 allocs/op.
Nine equal-full-class-output composite fixtures improve 2.043–4.258% over
parent medians. They beat original quantization plus its scalar full raster
class walk by 48.483–79.174%. This defined output contract can do more class
work than production callers that stop at EOB/zero; public encoder timing is
therefore retained as the acceptance evidence. The identical fixture SHA256
is recorded in environment.txt.

The no-class 64-coefficient control changes 62.15 → 61.79 ns (-0.579%), with mixed
pairs. It still takes 53.477% more time than original bare assembly 40.26 ns.
The full-class composite does not imply all quantizer contracts beat assembly.
Unchanged scalar DCT16 control has a +5.705% median change; production dispatch
is -0.364%. The full run is not treated as globally noise-free.

The quantizer stays live in separate 720p profiles: parent 0.26 s flat / 0.34 s
cumulative, current 0.21 / 0.27 s, each 4.56 s total CPU. These short profiles support
reachability and provide no precise isolated causality or high-confidence
ratio against the assembly path, which performs different class work.

The previous 63c2 run's +0.479% VP9 320p change versus a6c1f016 is preserved in
its linked artifact. Cross-job absolute times do not isolate code changes.
The EPYC 9V45 VP8 gap and ARM assembly gaps remain separate open evidence.
