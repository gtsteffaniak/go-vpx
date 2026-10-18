# PR110 actual-codec x64 feedback

https://github.com/thesyncim/govpx/actions/runs/39160666370

Source 2bbc6431b7e1c69e09ab086fa030f2ab7d6dc24d; baseline 7e0e37b7bbbcda65523f5581e7879c7d2d93827c. AMD EPYC7763, actualAVX2/Enabled=true/AVX512=false. VP8 1280x720,60inputs,30fps,2500kbps,realtime,threads1,GOMAXPROCS1,PGOoff. Two balanced AB/BA rounds per CPU setting. No C/profile/fullCI.

| CPU-used | Round | Parent ms/input | Current ms/input | Throughputdelta | Timedelta |
|---|---:|---:|---:|---:|---:|
| -8 | 1 | 14.545363 | 14.226403 | +2.2420% | -2.1929% |
| -8 | 2 | 14.635821 | 14.329881 | +2.1350% | -2.0904% |
| 8 | 1 | 17.152662 | 16.506257 | +3.9161% | -3.7685% |
| 8 | 2 | 16.747204 | 16.765289 | -0.1079% | +0.1080% |

Fixed−8 median throughput +2.1883%; both pairs exceed2%. Adaptive+8 is preliminary/mixed: +3.9161% and−0.1079% throughput. All four paired reports match22semantic/output fields; all8codecallocation samples0. Ordinary stream/ML/threaded-reference/new-helper selectors passed. Unrelatedsubpelcontrol median.99821× (mixedpairs), no componentgain claim.

Normalmerge 88fe29f77914a38e96a23adb28bd0275b827ce87 at 2026-10-18T19:55:05Z afterexactheadfastCI+customcodec checks passed. No force/delete/localcheckout edits.
