# PR52 native full audit: PASS

Exacthead5c5cf3977052368746325c4033a1909e7eacd26a, base6bda90f7eeee6f407d6ed53e299eab505afa440e, old691cd0512c4829d9ee40d483205dd223f9eb028d. Run39131584002/job117687686848 completedSUCCESS14:10:00UTC. Exhaustive39131584200SUCCESS13:40:50; focused/C green. No merge action bythisauditor.

NativeIntelXeon6973P-C; Go1.27.2/Linuxamd64/CGO1/PGOoff/GOMAXPROCS1. ActualAVX512/AVX2/Enabled=true; selected32ordinary+race and optional execution proven. These results do not imply anotherCPU shares this behavior.

Independentlyvalidated131rawfiles(99top+32actualStatscapture),132recordedsourceSHA,35sharedfixtures,561cases×18=10098zero-allocation main measurements,527excludedchecks,guard396+66,optional792+132. All sixvariantpermutations and preflights pass. Every22semanticfield and timed/profile allocation sample matches; no allocation exceptions. ActualStatscapture machine-code symbols/binaryreceipts/nativebuildinfo independentlypass.

Coverage caveat: the132pathlist is not fullrepository coverage. Three changedARMpaths were outside it: pixel/loopfilter_columns_wide20_arm64.go, pixel/loopfilter_columns_reuse_arm64_test.go, DSP/loopfilter_macro_axis_native_test.go. Supplemental exactGitblobproof binds allfivechangedGofiles, with135frozen source copies. All changedcodecGo and bothnewtestfiles explicitlyrequirearm64 and are excludedfromLinuxamd64; gensimd is ageneratorcommand. The unchanged132list is preserved truthfully.

Whole results vs51/old assembly:
- VP9 320x180: -0.174706%time vsparent (4/6pairs faster); -45.416696%vsold. All sixpairedchanges remain in compactsummary/fullproof/raw.
- VP8 720p: -0.273250%time vsparent (6/6pairs faster); -34.279041%vsold. All sixpairedchanges remain in compactsummary/fullproof/raw.
- VP9 720p: +0.202691%time vsparent (4/6pairs faster); -56.869134%vsold. All sixpairedchanges remain in compactsummary/fullproof/raw.

Losses/outliers:267/561kernel medians slower and3all6pairs slower; allrows andsamples remain in fullproof/CSV/raw. VP9320parentr3=3.015424ms andoldr4=5.575729ms remain unchanged. VP9720currentr1=18.521168ms is retained. No samples removed, no causal performance claim from unchangedx64production source. NativeStats9/9normalized bodies matchparent after onlypadding/relocation normalization; rawdiffs retained.

This lean package includes allrawprofiles/codegen/oracles, strictauditors, sourcebindings and exact135sourcecopies; no executablebinary/fullcheckout. GNUcapture binarySHA/size/address receipts remain. Kernel and whole results are separate. ExistingARMcorrectness/nativeperformance evidence belongs totheparent review; this x64run alone does not measure newARMcode.


Publication: `raw-evidence.tar.gz` preserves all 303 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
