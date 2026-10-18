# PR112 corrected-fixture native feedback

https://github.com/thesyncim/govpx/actions/runs/39162855642

AMD EPYC7763;actualAVX2/Enabled=true,AVX512=false. Exactheade760c8f8/base88fe29;runtimeunchangedfroma8c8398c. Sharedfixed−8fixturebaseline-firstpins pass on both variants, includingpartial-thread2. All50requestedordinaryoraclepasses/ref retained.

| CPU | Round | Parentms/input | Currentms/input | Throughputdelta | Timedelta |
|---|---:|---:|---:|---:|---:|
| -8 | 1 | 14.379937 | 14.478177 | -0.6785% | +0.6832% |
| -8 | 2 | 14.450711 | 14.553973 | -0.7095% | +0.7146% |
| 8 | 1 | 16.597083 | 16.621064 | -0.1443% | +0.1445% |
| 8 | 2 | 16.565692 | 16.882311 | -1.8754% | +1.9113% |

All four codec pairs are slower. Preserve this result; component gains are separate. All22output/semantic fields match,all8codec allocations0. Two balanced720p60input rounds perCPU at2500kbps/threads1/PGOoff; no robust codecgain claim.

| Component | Median | Pairedspeedups |
|---|---:|---|
| BenchmarkAstraDenoiserLuma/fullpel | 1.220423x | 1.225013x/1.215850x |
| BenchmarkAstraDenoiserLuma/odd-integer | 1.211853x | 1.216710x/1.207039x |
| BenchmarkAstraDenoiserLuma/rejected | 0.956385x | 0.958125x/0.954647x |
| BenchmarkAstraDenoiserLuma/split-fallback | 1.075087x | 1.069682x/1.080522x |
| BenchmarkAstraDenoiserLuma/state-fallback | 1.213534x | 1.224311x/1.202927x |
| BenchmarkAstraDenoiserLuma/subpel | 1.016555x | 1.019145x/1.013973x |
| BenchmarkAstraDenoiserLuma/y-only | 1.222145x | 1.221992x/1.222299x |
| BenchmarkAstraDenoiserLuma/zero-yuv | 1.121837x | 1.123281x/1.120392x |

All32componenttimings0alloc. Originaladaptivefixturefailsbeforetimings because tinyframe+8useswall-clock feedback; isolatednativebaseline/currentequalreceipt andoriginalfailedpinlogs remainunchanged in siblingdirectories. Correctedfixedpins were capturedfrombaselinefirst. No C/profiles/fullCI/pinrelaxation.

Explicitly accepted tradeoff: retain the measured ARM codec gains and x64 component gains while preserving all four x64 codec regressions. Normal merge 159dce5b931af4f7c600f157d64c2a813092ec57 after exact-head fast checks passed; no extra timing or validation. Receipt: merge-receipt.json.
