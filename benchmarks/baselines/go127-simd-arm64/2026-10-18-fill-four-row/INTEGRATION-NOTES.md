# Scope correction for the earlier frozen FillBlock archive

Keep /private/tmp/govpx-arm-fill-exact-ready byte-for-byte unchanged. Its
REPORT.md incorrectly describes current production FillBlock callers as
only n8/n16: internal/vp8/dsp/intra4x4_other.go:10 also calls
intraDCGo(dst,stride,above,left,4,true,true), which calls pixel.FillBlock.
The earlier component proof includes 4x4 FillBlock, while actual DSP caller
benchmarks cover only IntraDCPredict8x8/16x16. No actual 4x4 DSP caller
performance was measured or inferred. This correction changes neither the
verified 4x4 kernel data nor acceptance from 8x8/16x16 oracles/measurements.

The new unroll4 probe uses that exact locked FillBlock candidate as its strong
parent. Initial3-round screen results are provisional and cannot be combined
multiplicatively with prior medians. Fresh confirmation and separate whole
quantification require a new authorized timing window. The frozen earlier
winner remains independently eligible if this probe fails confirmation.

Old691's ARM intraFill8x8NEON and intraFill16x16NEON each store ONE row per
loop, then increment the pointer/decrement counter/branch. No old-NEON-unroll
explanation or speed claim is made. New Go source computes one row base per
four stores; actual compiler output is preserved separately.
