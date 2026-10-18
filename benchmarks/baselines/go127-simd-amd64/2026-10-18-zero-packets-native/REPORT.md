# Packet native compact audit: PASS

Run39136145443/job117700941976, exacthead2366eb72882de11ea9758270fac9524b27b88123, parent5c5cf3977052368746325c4033a1909e7eacd26a. NativeAMD EPYC9V45, actualAVX512/AVX2/Enabled=true. Selected ordinary/race/checkptr,7568numeric/450alias/30panic/60wrap and post-timing packet oracles pass. Three source deltas match exactGit; onlytwoAMD files compile onx64, ARMshared delta excluded.

1008component measurements/168excluded1x;84cases; sixbalancedAB/BAorders; allcomponentallocationszero. Zero28medians1.033–1.311x; sparse20/28losses(range0.947–1.024x);dense20/28losses(range0.951–1.015x). Zero cases have3slower individualpairs; sparse111/168pairs anddense115/168pairs slower. No1.5xqualifyingmedian. Everyloss/rawpair retained.

VP9320time−0.635709%, VP9720−0.953557%, each4/6pairs faster. All22semanticfields and timedallocationrates equal:32074/120;72075/60. All48excluded320diagnosticprofileallocation samples preserved. Parent720r3=19.158857ms andcurrent320r3=2.461351ms retained; nooutlierdiscarding.

Excluded packettrace zero fractions79.15699%320 and75.49782%720. This is anisolatedGo-vs-parent comparison, not libvpx throughput or historical-cause proof. Source pin2366 is an available preliminarybaseline option; finalPR53hasadditionalARMkeepers and mustbe measured separately.

Raw462files, producerMANIFEST, source/binarybindings, profiles/GNUdumps and compactproof retained. No executable/fullsnapshot included. Compact audit intentionally avoids repeating the driver’s full sourceinventory hashing.


Publication: `raw-evidence.tar.gz` preserves all 468 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
