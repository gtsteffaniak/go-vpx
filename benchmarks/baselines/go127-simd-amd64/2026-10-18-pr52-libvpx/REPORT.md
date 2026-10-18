# PR52 optimized libvpx reference audit

Exact head5c5cf3977052368746325c4033a1909e7eacd26a, run39131584002/job117687687088 SUCCESS. Native AMD EPYC7763, separatejob fromfocused; no crossjobpooling. Pinned optimized libvpx1.16.0 reference1024874c5919305883187e2953de8fcb4c3d7fa6, runtimeCPUdetect/SSE2/AVX2config and objects verified; no separatelinkedCcapabilityprobe. ImmutableGitCLI/driver/workflow plus all74artifactfiles checked; no local timing.

Go/C TIME medians: fixedVP8720=2.15966349, adaptiveVP8=1.65434952, VP9320=2.87780382, VP9720=2.79790933. Values>1meanCFASTER. Every32helper invocation/control/timer/quality/allocation audited;24timed rows usevpxenc-stats timer,4warmups+4separatequalitypasses retained.

FixedCPU−8,minq2/maxq56 is not fixedquantizer. Bothsides58encoded/2drop,723593B,PSNR27.15761087/SSIM0.925927329 match measuredoutput/quality and everytimedsample, withoutbitstreamidentityclaim. Adaptive unequalwork: Go685957B/PSNR27.98395/SSIM.90850 vsCquality678205B/29.17664/.93309; Ctimedround4 differsfromitsseparatequalityoutput, allotheradaptiveCtimedrows match. Thus adaptive quality cannot be assigned indiscriminately to everytimedrow. VP9output/quality unequal, eachside stable against ownqualitypass.

Goallocations preserve baseline0/.6166667/1.25 exceptVP9720rounds1and6=1.2666667, eachoneextraallocation over60inputs. Callocations unmeasured. Allpaired losses and actualquality/output stats retained. No Go-beats-C, fixedworkforunequalcases, bitidentity, or crossCPUcausalclaim.


Publication: `raw-evidence.tar.gz` preserves all 5 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
