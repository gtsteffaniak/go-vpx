# Short actual-work pass ledger

Fixed full-validation baseline remains c4a8bd9f66805fc9b7b4f99e26752cc68505e572. All samples are preliminary; no1.5x codec threshold claim or full-validation dispatch. No new C comparison/profile/full matrix. README untouched.

| Probe | Native codec/component result | Decision |
|---|---|---|
| Encoder full LF via fused luma+chroma | flat/slightly slower; first fixed parent outlier | rejected; raw+patch in PR110 evidence |
| ONE token unrolling | mostly1–3% component; one slow parent outlier | rejected below threshold |
| Accepted VP8 whole-intra predictor reuse | ARM fixed1.0383x, x64EPYC7763 fixed1.0219x; DC1.91x/skip2.89x ARM | merged PR110,88fe29f7 |
| Subpel MV-rate floor, all searches | fixed0.9525x, adaptive1.0491x | rejected; mixed codec loss |
| Subpel MV-rate floor, iterative only | fixed0.9813x, adaptive0.9908x; first pair had compile overlap | rejected |
| VP9 large model v1 | extra structs/helper overhead; mixed/slower | rejected; raw retained |
| VP9 large model v2 | ARM fullpel32/64 component1.19–1.20x/1.164–1.168x; codec0.987–0.996x | winning component code preserved in archive for scoped native follow-up |
| VP9 unused coefficient-cost setup | ARM720p1.0373/1.0162x;320p outlier excluded | merged PR111/e9d4d08b; x64 non-RD setup6.91x, codec0.9993/1.0185x |
| VP8 average-backed denoiser | ARM fixed1.0404/1.0478x, adaptive1.0955/1.0156x | PR112 held: x64 partial Thread2 noisy hash mismatch; baseline-only diagnosis pending |
| VP9 small-model-only fallback v3 | ARM720p1.0373/1.0074x;320p outlier excluded | merged PR113/9cdd3b7f; x64 codec1.0058/1.0079x |

Codec and kernel results are separate; unchanged VP9 codec allocation count is1.25/frame on this720p workload. Setup components and VP8 codec allocate zero. Raw results retain all losses. The V2 and V3 picker variants overlap and are alternatives; the cost-table gate is independent.
