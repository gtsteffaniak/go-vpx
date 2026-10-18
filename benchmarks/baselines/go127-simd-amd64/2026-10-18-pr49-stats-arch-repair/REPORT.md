# Stats architecture reduction alias trial

Private three-file overlay only; ROOT and all ARM50 snapshots remain unchanged. ARM uses `pixelStatsReduction = int32`, AMD uses `= int64`, with compile-time vector/scalar reduction selection retained. All accumulator conversion and odd-row addition use that alias.

All 13 command receipts pass: ARM ordinary and race/checkptr=2 Stats/extreme/variance tests across pixel, VP8 DSP and VP9 DSP; AMD v1 checkptr=2 selected SIMD and actual AVX/AVX2-disabled tests across the same three packages. Required AMD native selection passes; CPU-off native-only tests skip while public tests pass. No timing claim.

AMD all seven selected functions exactly match the 31123 baseline instruction sequences after address/relocation normalization. ARM six functions exactly match41596; Pairs8 returns to169 instructions with no extra sign extensions. Its sole remaining textual instruction difference is the commuted scalar odd-row addition `add x0,x8,x0` to `add x0,x0,x8`, which has identical inputs, destination and operation. No other opcode/register/branch/frame change occurs. A separate commutative-add normalization proves the full Pairs8 body equivalent. This is not presented as exact machine-byte equality.

The initial mechanical replacement accidentally changed two `Uint64` spellings and failed compilation; that script error and compiler output are preserved as `initial-transform-failure.log`. The corrected overlay and all successful outputs are retained. Baseline source overlays and binaries remain in the previous Stats package; full baseline pins are41596b4e0ad1087e4afe4020825d03cb38e2c4ce ARM and31123a9ecb7fa3b90afdb293beb06dadb2a6f857 AMD. All seven functions use the same ROOT link context and baseline overlays.


Publication: `raw-evidence.tar.gz` preserves all 7 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
