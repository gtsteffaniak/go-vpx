# Accepted ARM odd-eight-pixel Stats packet

Baseline is exact `10caaa67ed53f8c3c1248f20070e5a4117245852`. Only internal/pixel/goav1_pixelstats_simd.go changes: a complete final8-byte row forms a low-half SIMD packet, with the existing square/sum helpers and one final reduction. Original paired rows/other widths are unchanged. Short final rows retain the original scalar fallback after evaluating source/reference slices in the original order.

Each delta is in[-255,255], each square<=65025. A public width8 chunk has<=4096 rows, SSE<=2130739200 and absolute signed sum<=8355840; both fit int32. Moving reduction preserves all eight squares and differences. Public tall batching preserves unsigned modulo2^32 SSE and full scalar sum. Exact checked8-byte loads read no extra byte, and inputs remain read-only with arbitrary mutual aliases. Existing stride/index arithmetic is unchanged; malformed windows preserve exact parent panic state.

84,514 independent/frozen-parent comparisons pass per mode:65,536 exhaustive byte pairs,2,346 odd/even/tall/minimal/alignment/alias states and16,632 invalid/short/spare-capacity/overflow/panic-order states. Native ordinary, native race/checkptr2, AMD selected/off checkptr2 all pass (338,056 candidate comparisons), plus84,514 parent self-checks. Valid inputs also compare against independent StatsScalar; tall heights include4095/4096/4097,8191/8192/8193 and65535/65536/65537. Existing Stats16 tests pass; the SIMD-selection test skips under actual CPU-off as intended. AMD runs establish correctness only. No test thresholds were weakened.

Native optimized LLVM pixelStatsPairs8 grows169→200 static instructions/frame16 unchanged, with no hot helper calls; only cold bounds/stack calls remain. Stats202 and the existing B.Loop fixture336 static instruction counts are unchanged. The existing18-case public fixture is identical on both references; it compares StatsScalar and full input state after the loop and consumes the output.

Six alternating AB/BA75ms rounds on Apple M4 Max, official Go1.27.2 SIMD, CGO1, PGOoff, GOMAXPROCS1 produce216 zero-allocation samples and36 excluded one-shot validation rows. Width8,height1 median is1.546155x parent, all six paired ratios>=1.509123x: this meets the user's mandatory1.5x local-component retention rule without rounding. Width8,height7 improves1.260546x, all six pairs>=1.240039x. This is a local ARM component win, with no whole-codec, native-x64 or old-assembly claim.

All controls and individual losses remain: width32,height1 median0.999603x (3/6 slower), width64,height1 median0.991250x (4/6 slower); other mixed slower pairs are retained in full. Unchanged source is not used to dismiss observed losses. No measurements from a different window are pooled.

All source/overlay/binary hashes and complete parent source set were independently checked against the exact baseline. Raw outputs/commands/timestamps, both excluded preflights, every ratio/pair/outlier, correctness logs, frozen-parent tests, build receipts, relevant source and LLVM are included. Executables/full snapshots are omitted with SHA256/size/private paths recorded. No ROOT/Git/README/dependency/assembly edits or new timing occurred during the audit.


Publication: `raw-evidence.tar.gz` preserves all 5 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
