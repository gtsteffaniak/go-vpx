# ARM copy-kernel x64 control evidence

PR41 tested head ab3abd6b335242d5732fa1956a63b7591eaf781b merged as cf3117194f0590fa0af35951a32e8fadb6da9ab3 with the exact tested tree. Its predecessor34c4e4a1 and merged main2905 have identical trees. Full audit, independent root proof, merge/retarget/cleanup proof and all160 native/C raw files are preserved in raw-evidence.tar.gz. bundle-sha256.json verifies every frozen bundle member, including its original manifest.

Native EPYC7763: all338 cases have six zero-allocation samples per variant;23 exact source hashes and all22 public/profile semantic fields match. Every public/profile allocation sample matches. The eight copy cases are unchanged AMD controls; inherited AVX512 is unavailable on this host and its two dedicated tests correctly skip. No AVX512 or AMD copy speedup is claimed.

Whole controls versus parent / old Go assembly build: VP8 -0.097% / -19.573%; VP9 320p +0.512% / -40.472%; VP9 720p -0.386% / -49.732%. All179 median regressions and five all-six-slower controls are retained. This does not alter the separately measured retained M4 copy gains of1.581–1.840× or establish a new ARM whole-codec gain.

Separate optimized-libvpx EPYC7763 job: fixed VP8 Go/libvpx encode-time ratio2.355204× with matching measured quality, payload and frame counts. Adaptive/VP9 perform unequal coding work; full controls and quality evidence remain in the bundle. Matching statistics do not prove bitstream identity.
