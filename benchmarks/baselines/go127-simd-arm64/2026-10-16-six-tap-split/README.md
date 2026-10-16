# VP8 six-tap arithmetic and split-block subpixel dispatch

Parent08c2b96e; original assembly691cd051. Apple M4 Max, official Go1.27.2,
GOEXPERIMENT=simd, PGO off, GOMAXPROCS=1. All timings are native ARM;
local Go validation/builds and other experiments were serialized.

VP8 six-tap prediction previously used the general int32 accumulator.
The new kernel checks coefficient signs and bounds, accumulates both
negative products first in int16, then saturates positive additions.
Negative weight is at most128, so the minimum is64-128*255=-32576;
each individual positive product is at most128*255=32640. Once a positive
sum saturates, rounding/clipping already yields255. This preserves exact
results while doubling arithmetic lanes. Odd phases omit the two zero
end taps through a dedicated four-tap kernel. All memory accesses use
safe slices; unsupported/custom coefficients retain the int32 fallback.

Three designs were compared, with patches and raw six-round schedules
retained: a common six-tap kernel, generic zero-end specialization, and
separate four/six-tap helpers. The first six triples compare the parent
against the first two designs. Six further triples choose the final
separate helpers against both candidates; the generic design adds cost
on dense4x4 cases. In the initial comparison, dense16x16 two-axis
prediction improves473.1→229.75ns (-51.44%). In the second comparison,
the final hybrid retains the dense design's result and reduces odd-phase
16x16 235.4→183.0ns (-22.26%);8x8 79.2→66.37ns (-16.20%). These are
within-run comparisons; absolute timings between sessions do not isolate
changes. Every kernel sample has zero bytes/allocations.

Both-fractional-axis split-block variance now dispatches to the existing
fused SIMD statistics kernel, eliminating uint16/byte staging. Integer
and one-axis cases retain their prior paths. Fresh six-pair B.Loop
benchmarks store both returned values in observable sinks:

| Shape | Parent | Fused | Time change |
| --- | ---: | ---: | ---: |
|16x8|127.60ns|20.945ns|-83.59%|
|8x16|129.25ns|25.280ns|-80.44%|
|8x8|72.89ns|13.700ns|-81.20%|
|8x4|45.015ns|9.850ns|-78.12%|
|4x8|47.19ns|13.615ns|-71.15%|
|4x4|31.405ns|9.531ns|-69.65%|

All are allocation-free. All bilinear phases on exact/misaligned input
windows match an independent scalar two-stage oracle, without changing
inputs. Six-tap tests include all codec tables, byte extrema, alternating
extrema, random inputs, exact row bounds and maximal signed weights.
Five kernel/DSP/encoder packages pass on ARM and Rosetta AVX2.

Ten balanced public quintets compare original assembly, parent, each
isolated change and the combination. Each variant occupies every position
twice; reversed orders balance predecessors. These are realtime CPU-used8,
1280x720,60inputs,2500kbps, one thread, encode-only:

| Assembly | Parent | Six-tap only | Split only | Combined |
| ---: | ---: | ---: | ---: | ---: |
|7.067ms|9.969ms|9.708ms|10.055ms|9.704ms|

Combined time improves2.661% versus parent; all ten pairs are faster
(-5.11..-1.65%). ARM still trails original assembly by37.31%.
Split-only public timing does not establish an end-to-end speedup; the
large split-block improvements above are DSP measurements. The fresh
720p CPU profile proves the six-tap kernel live (4.23% cumulative).
The CPU-used4 profile does not establish substantial split-block traffic.

Six alternating control pairs use VP9 720p CPU-used8 and VP8 320x180
CPU-used4. VP9 -0.447% has a paired range crossing zero (-1.52..+0.68%):
no established effect. VP8 CPU-used4 improves1.987%; all six pairs are
faster (-3.86..-1.15%). Every non-timing JSON field matches across the50
main and24control samples, including output byte/frame counts, quantizer
histograms and allocations. Bitstream parity remains a separate CI gate.
Native x64 performance of this checkpoint is pending.

After formatting and integration, six fresh balanced triples compare the
parent, prototype and final PGO-off binaries: 9.932/9.592/9.616ms. The final
build improves 3.184% versus parent, with every pair faster 2.32–4.12%; its
+0.243% median difference from the prototype is within noise. All 18 non-
timing/allocation samples agree. Final native/Rosetta five-package suites,
PGO refresh and pre-commit pass. ARM disassembly records the live vector
multiply, saturating-add and narrowing bodies in the retained code.
