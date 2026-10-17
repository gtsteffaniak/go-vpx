# ARM variance and loop-filter transpose (2026-10-17)

Apple M4 Max, macOS/arm64, Go 1.27.2, GOEXPERIMENT=simd,
GOMAXPROCS=1, PGO=off. All Go commands use the root project wrapper/cache.
Native timings are serialized against all local Go builds/tests/profiles.
`public.json` records six balanced permutations of original assembly
691cd051, previous 611f47fe and the candidate. VP9 is a parent/current
control. Workload:1280×720, 60 frames, 30 fps, 2500 kbps, realtime, cpu8, one thread,
encode-only,no libvpx autodetection. Only four timing fields are excluded
from the deep semantic comparison; every other JSON field matches.

VP8 median: parent 8.621926 -> candidate 8.320135 ms/frame (-3.500%). Every
paired change is faster (3.208–4.365%). Assembly is 7.065953 ms/frame in this
same run: candidate still takes 17.750% more time. The parent is 22.021%
above assembly here. The older 26.590% ARM gap belongs to its separate run;
this optimization accounts for 3.500% less encoder time, not the entire
change between separately measured gaps. Every assembly/candidate pair is
still slower by 9.248–21.587%. ARM assembly parity remains open.

VP9 parent 12.650947 -> candidate 12.391071 ms/frame (-2.054%), with mixed
paired changes (+0.176, -1.389, +0.402, -2.638, -3.558, -4.727%). This control
does not establish a stable VP9 speedup. Allocations remain unchanged:
VP8 zero, VP9 1.25 per input frame. Byte counts, frame decisions, quantizer
statistics/histograms, quality fields and all options match exactly.

The 16-column variance path accumulates signed sums in int16 for1..64 rows:
2*255*64=32640 per lane. Squared errors use independent int32 accumulators;
NEON ADDV replaces four scalar extractions at reduction. A stride bound
proves final-row index multiplication cannot overflow, and the two bounds
checks prove all exact 16-byte loads. Zero strides/read-only overlapping rows,
minimal buffers, unaligned addresses and extreme signed differences pass
an independent scalar oracle, including the64/65 row transition. Race and
checkptr=1 pass. The same public B.Loop fixture is overlaid on the parent:
h1–h64 improve 24.9–45.5%; h16 is 28.230 -> 16.870 ns (-40.24%), zero allocations.
Unchanged h65/h720 controls have small mixed changes, medians+0.37/+0.40%.

The ARM 16×8 transpose pairs rows eight apart and uses 24 TRN instructions.
It replaces two separate 8×8 transposes plus eight joins, using only exact
8-byte row reads. Independent lane tests cover 1120 minimal/unaligned tiles;
existing sequential loop-filter oracle tests cover masks and overlap
fallbacks. No read or write window is widened. The selected x64 variance
bodies are instruction-identical. x64's column helper has only a reorder
of three argument-save instructions in its prologue; its instruction count
and arithmetic are unchanged. Fresh native x64 CI remains the check.

Design exploration retains raw results. Separate safe-indexed variance,
bounds-proven variance and accumulation-order variants were compared.
Integer zero-fill, integer broadcast and float-bit row pairings were
compared; the final candidate uses integer broadcast. The combined integer
candidate won all six exploration pairs and all six fresh final pairs.
A byte16 denoiser improved its direct kernel 76.605 -> 21.590 ns but public
encoder results were mixed, so it was rejected. Two/four-accumulator
row-major full-plane SSE prototypes also had mixed public results and
were rejected. Macro-only constant/load/store specializations were left
experimental because direct kernel gains did not repeat in the encoder.

No new assembly or module dependency. Other architectures, purego and
experiment-off use their existing code. Required-native Rosetta correctness
is recorded separately and never used for performance claims.

Final-source checks: full `scripts/go.sh test ./...`, fmt/PGO/noasm and all
six build selections pass. Strict libvpx byte-parity gate passes in 36.901 s;
regression report passes 18/18. PGO was refreshed after all production Go
source/comments were frozen; fingerprint 099eaf8d75a17d8b630d168e795433f6f9992dda.
