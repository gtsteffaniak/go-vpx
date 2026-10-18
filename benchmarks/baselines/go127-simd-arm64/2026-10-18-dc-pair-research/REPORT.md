# DC pair reconstruction research

Native Apple M4 Max, Go1.27.2 archsimd, PGO off, GOMAXPROCS1. Two balanced
parent/ASM/current diagnostic rounds pervariant,100ms per18 public fixtures.
Parent eb820478, assembly691cd051. Preserve every sample/source; these short
diagnostics are not fresh six-round acceptance or whole-codec evidence.

The new B.Loop public fixtures cover pitches8/32/1312, exact in-place and
disjoint unequal pitches, unaligned starts, three DC pairs including clipping
and identity, exact visible spans, independent scalar equations and full
backing canaries. In-place reset restores32 pixels inside the measured loop
identically in every reference; disjoint inputs remain immutable. Results are
observable; all samples have zero B/op/allocations. OldASM remains legitimate.

* diagnostic: 1.2020–1.2415x parent; 0.5455–0.6394x assembly.
* guards-diagnostic: 1.2211–1.3506x parent; 0.5607–0.6841x assembly.
* narrow-diagnostic: 1.3019–1.4216x parent; 0.6176–0.7011x assembly.

Variants remove closure/slice checks after original bounds/alias guards,
simplify stride/delta interval guards and exact-in-place selection, then use
one signed-delta pattern plus widened add and unsigned narrowing saturation.
All remain below1.5x parent and still trail assembly. Preserve their code and
data for continuing research without claiming a qualifying gain or changing
the production DC path in this checkpoint.

Guard proofs: uint(stride-8)<=maxPitch-8 accepts exactly8..floor((MaxInt-7)/3),
including rejection of MinInt wrapping. Then3*stride+7<len equals the original
stride<=(len-8)/3. uint16(delta+255)<=510 selects exactly[-255,255] over every
signedint16 value. Exact aliases have non-overlapping rows; other accepted
spans are fully disjoint. Thus individual row read/store preserves semantics.
Widened pixel+delta lies[-255,510]; signed add plus unsigned saturation equals
the original positive/negative byte saturation. No padding or retained pointer.

Existing256-byte×511-delta scalar clipping tests and alignment/alias cases pass.
The new451,584-case guard matrix checks original division-based acceptance,
independent delta pairs, negative/extreme strides and short lengths, full96-byte
canaries and immutable inputs. Native pixel/public DSP tests pass. AMD builds
and Rosetta correctness pass after a review caught the first signed narrow
intrinsic was ARM-only: the final version uses the existing architecture
packBytes helper. The original source is recovered exactly by its reviewed
SHA and retained alongside the correction. No emulated speed claim.

Independent numeric/bounds/alias review passed after that correction. No
concurrent local builds/profiles during timing. No README or dependency edits,
production cgo or handwritten assembly introduced.
