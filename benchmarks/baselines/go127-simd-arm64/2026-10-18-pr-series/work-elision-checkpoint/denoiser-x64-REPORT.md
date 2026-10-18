# PR109 native denoiser result

https://github.com/thesyncim/govpx/actions/runs/39156757003

Head 78f1c59b0e9fe289f6a4f159c4eb3a6db4fc645b; baseline 499ef24ee5f30dbb44890e3fb19b88bbe2e0407b. AMD EPYC7763,actualAVX2/Enabled=true,AVX512=false. Sharedcomplete-denoiserfixture,affectedfrozen/MLoraclesPASS,8cases32zeroallocationtimings. Two40msAB/BApairs; component evidence only, no x64whole/C/profile/fullrun.

| Case | Baseline ns | Candidate ns | Pairedspeedups | Median |
|---|---|---|---|---:|
| BenchmarkAstraDenoiserLuma/fullpel | 256.0,255.7 | 244.7,245.9 | 1.046179x,1.039854x | 1.043016x |
| BenchmarkAstraDenoiserLuma/odd-integer | 354.0,353.2 | 244.3,243.3 | 1.449038x,1.451706x | 1.450372x |
| BenchmarkAstraDenoiserLuma/rejected | 1015.0,1014.0 | 755.7,756.9 | 1.343126x,1.339675x | 1.341400x |
| BenchmarkAstraDenoiserLuma/split-fallback | 474.0,474.7 | 486.8,487.7 | 0.973706x,0.973344x | 0.973525x |
| BenchmarkAstraDenoiserLuma/state-fallback | 1094.0,1096.0 | 1112.0,1119.0 | 0.983813x,0.979446x | 0.981629x |
| BenchmarkAstraDenoiserLuma/subpel | 1055.0,1040.0 | 770.6,772.9 | 1.369063x,1.345582x | 1.357322x |
| BenchmarkAstraDenoiserLuma/y-only | 210.1,210.0 | 179.3,179.8 | 1.171779x,1.167964x | 1.169872x |
| BenchmarkAstraDenoiserLuma/zero-yuv | 448.7,447.1 | 473.0,473.0 | 0.948626x,0.945243x | 0.946934x |

ExactheadroutineCI+automaticquickSUCCESS; normalmerge 7e0e37b7bbbcda65523f5581e7879c7d2d93827c at 2026-10-18T19:03:31Z. No force/delete/localcheckout changes.
