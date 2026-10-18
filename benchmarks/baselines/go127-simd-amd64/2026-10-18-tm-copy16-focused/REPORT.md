# PR46 focused SSE8 control audit

Frozen head `46a08f8e6eb03f6ee4549d975384bd6c952a81fa`; base `78ce74d39806d74ecf91bc3d082899a74d7bc85f`. This separate focused runner used **AMD EPYC 7763 64-Core Processor**. SSE8 production sources are unchanged. These are acceptance-relevant timing controls; source identity does not establish a cause or make regressions irrelevant.

All 720 measured B.Loop/legacy samples and 120 excluded one-shot samples are preserved. All measured allocations are zero. All exact Git hashes, dispatch/correctness/race checks and six balanced orders passed. Source inventories: head 274, base 261, old reference 147.

Across 36 public fixtures, 20 median times increase and five are slower in all six pairs. No fixture is faster in all six pairs. Four pointer fixtures regress in every pair; all nine pointer medians increase by 3.471–6.534%. The remaining all-six-slower public fixture is tight-pitch slice8x16 pattern1 (+0.124%).

| API | Parent/candidate median ratio range | All-six slower |
|---|---:|---:|
| slice8x4 | 0.999307–1.004034× | 0 / 9 |
| slice8x8 | 0.998170–1.000915× | 0 / 9 |
| slice8x16 | 0.998760–1.000931× | 1 / 9 |
| pointer8x8 | 0.938667–0.966455× | 4 / 9 |

The separate B.N pointer diagnostic increases 6.072261% in all six pairs. B.N slice4 improves 4.263%, while slice8 increases 0.394% and slice16 improves 0.102% with mixed pairs. These weak B.N diagnostics do not replace the 36 live B.Loop fixtures. All raw samples and old-reference comparisons are in the audit JSON. No new Copy16 gain, whole-codec gain, or cross-CPU claim follows from this focused job.

# Separate optimized libvpx audit

The optimized C job used **AMD EPYC 9V74 80-Core Processor**; its timings are not pooled with focused/native jobs. Upstream is v1.16.0 / `1024874c5919305883187e2953de8fcb4c3d7fa6`, with runtime CPU detection, assembly objects, SSE2 and AVX2 available. All 24 timers came from vpxenc statistics; all quality-output statistics match every corresponding timed output.

| Workload | Go/libvpx TIME ratio | Scope |
|---|---:|---|
| vp8-720p-fixed | 2.282511644× | Matched measured payload, counts, PSNR and SSIM; 58 quality frames |
| vp8-720p-adaptive | 1.661755936× | Unequal output / quality work |
| vp9-320p | 2.664478208× | Unequal output / quality work |
| vp9-720p | 2.576795879× | Unequal output / quality work |

A ratio above one means libvpx takes less time. Fixed VP8 means CPU-used −8 with min-q=2/max-q=56, not a fixed quantizer. Measured matching statistics do not establish bitstream identity. Adaptive VP8 and VP9 retain different output/quality work and cannot support equivalent-work comparisons.

PR46 full native/root audit and explicit rootPERFREADY accepted documented control regressions and Copy16 disjoint wins. Normal guarded merge `2df3eaf5fb547db07bd970cd65bf504ecb287590` preserves testedhead ancestry and exact testedtree. Final body/fourchecks/parents/tree proof live in the companion full bundle. CI success alone did not authorize merge.


Publication: `raw-evidence.tar.gz` preserves all 6 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
