# PR111 native VP9 setup and actual codec feedback

https://github.com/thesyncim/govpx/actions/runs/39161547354

AMD EPYC9V74;actualAVX2/AVX512/Enabled=true. Exacthead35231892/base88fe29. Two720p60input VP9CPU8pairs,2500kbps,threads1,PGOoff; fouractualprepareframecontext componentcontrols, two40msAB/BApairs. No C/profiles/fullCI.

| Round | Parent ms/input | Current ms/input | Throughputdelta | Timedelta |
|---|---:|---:|---:|---:|
| 1 | 18.631151 | 18.644900 | -0.0737% | +0.0738% |
| 2 | 18.931176 | 18.587690 | +1.8479% | -1.8144% |

Codec result is mixed/nearflat, not a robust x64whole gain. All22semantic/output fields match; codec allocations unchanged75/60=1.25/frame. Setupcomponent allocations alone are zero.

| Setupcase | Median speedup | Pairedspeedups |
|---|---:|---|
| BenchmarkVP9NonrdCoeffCostPrepare/compound-inter | 0.998277x | 0.998410x/0.998144x |
| BenchmarkVP9NonrdCoeffCostPrepare/fullrd-inter | 0.995232x | 0.997621x/0.992834x |
| BenchmarkVP9NonrdCoeffCostPrepare/key | 1.004643x | 1.008121x/1.001161x |
| BenchmarkVP9NonrdCoeffCostPrepare/nonrd-inter | 6.909802x | 6.903315x/6.916283x |

All requestedordinaryoraclesPASS bothrefs. Normalmerge e9d4d08b3208c63254f251aa363e8b2b1c66abbf after exactheadfastCI+customcorrectness/no materialregression. No force/delete/rootcheckout changes.
