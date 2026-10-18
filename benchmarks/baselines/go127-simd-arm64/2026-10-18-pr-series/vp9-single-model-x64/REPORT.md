# PR113 actual native VP9 feedback

https://github.com/thesyncim/govpx/actions/runs/39161954456

AMD EPYC7763, actual AVX2/Enabled=true,AVX512=false. Exacthead52ea8879/basee9d4. Two720p60inputVP9CPU8pairs at2500kbps,one thread,PGOoff; all22semanticfields match and allocations unchanged1.25/frame. No C/profile/fullCI.

| Round | Parent ms/input | Current ms/input | Throughputdelta | Timedelta |
|---|---:|---:|---:|---:|
| 1 | 24.666288 | 24.524666 | +0.5775% | -0.5742% |
| 2 | 24.614503 | 24.422758 | +0.7851% | -0.7790% |

Both short codec pairs improved modestly; preliminary, not a robust largegain claim. The unchanged framecontextsetup control ranged.9972–1.0059×,all16componenttimings0alloc; no controlgainclaim. Allrequestedpreparedpredictor/costtransition/96VP8/72MLoraclespassedbothrefs.

| Unchangedcontrol | Median speedup | Pairedspeedups |
|---|---:|---|
| BenchmarkVP9NonrdCoeffCostPrepare/compound-inter | 1.002016x | 0.997686x/1.006366x |
| BenchmarkVP9NonrdCoeffCostPrepare/fullrd-inter | 1.005859x | 1.004544x/1.007175x |
| BenchmarkVP9NonrdCoeffCostPrepare/key | 0.997219x | 0.993237x/1.001227x |
| BenchmarkVP9NonrdCoeffCostPrepare/nonrd-inter | 0.999398x | 1.001206x/0.997596x |

ExactheadfastCI+customnativepassed;normalmerge 9cdd3b7f0c15ed5ea49bd81b8b24ac5026231f8a at 2026-10-18T20:12:45Z. No force/delete/rootcheckoutchanges.
