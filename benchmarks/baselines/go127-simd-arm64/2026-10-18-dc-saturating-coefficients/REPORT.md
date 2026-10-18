# ARM DC signed-vector saturation evidence

Apple M4 Max, official Go1.27.2 SIMD, CGO1, PGOoff, runtimeGOMAXPROCS1. Final candidate includes accepted SEPARATEpacket plus vector coefficient packing/clamp. Root/PR49 source remained unchanged during this private research. Baseline49 canonical source is41596b4e0ad1087e4afe4020825d03cb38e2c4ce; component also directly measures exactPR48 and old691. No cross-window speedup multiplication.

Final source has1,947Go files: onlytwo ARM production deltas andone new coefficient-boundary test relative to acceptedpacket. Coefficients wrap at signed32 before arithmeticshift; SQXTN packs signed16 coefficients and selfZIP×2 broadcasts each tofourcolumns; SQADD thenSQXTUN preserves final byte clipping. Independent review exhausts16,777,216signed16delta×byte cases, proves remaining signed32 tails, and passes132,096native full-backing comparisons against TWO unchanged public single-block transforms under race/checkptr2. Original55,296/156,000full-state/panic suites, byte/delta, disabled/purego/nosimd and rootboundary tests also pass. Candidate/root files were not changed by independent overlaytest.

The final CLI in-place owner is152instructions/frame16, first-return prefix66 versus83 for acceptedpacket, nohotcalls/vectorstack. GenericDSP owner remains220/frame16 and prefix122, with real SQXTN/fourSQADD/fourSQXTUN. Prefix counts are static code-layout diagnostics, not dynamic instruction or speed claims. Actualtwo retained decoder CLI caller bodies normalize identically. Generic owner was reviewed in the publicDSP binary because CLI dead-strips it. Early basic/globals codegen/build variants are untimed diagnostics, not pooled performance results; complete sources/logs/static comparisons preserved.

## Five-way component window

All1,800samples=36cases×10rounds×5references pass independent public outputs with0B/op and0allocs/op. Ten cyclic/reverse orders balancepositions but notallorderedpairadjacencies. All raw samples and losses are retained. No preflight samples are timing denominators.

Nine actual in-place target medians are1.618931–1.645445×PR48, with90/90pairs faster but minimum1.247932×: this is not an all-paired1.5floor. Versus final49 they are1.103039–1.114863×,88/90pairs faster; versus acceptedpacket1.058967–1.079962×,87/90faster. Allnine medians beat the direct-call old assembly reference by1.012420–1.033995×, while70/90pairs are faster,19slower andone tie. No universal or every-pair assembly win is claimed.

All36medians beat49. Originalgeneric disjoint medians1.139802–1.161448×49(85/90faster); originalgenericinplace1.054676–1.079576×(79/90); newdisjointcontrols1.158198–1.191693×(89/90). Old remains faster across all27generic/disjoint medians. All360paired samples beatPR48, but several ratios remain below1.5 and every loss/outlier against stronger references is preserved.

## Separate whole-codec window

Six balanced three-way permutations,18JSON, VP8 realtimeCPU8 one thread1280×720,2500kbps,300inputframes. Median49/current/old=7,135,325/7,089,421.5/6,807,178ns/input. Current is0.64332739%faster median than49, four ofsixpairs faster andtwo slower: +0.05402%,−1.33258%,−0.07720%,−1.14194%,+1.10822%,−0.85879%. This small mixed whole result does not establish a stable whole-codec gain.

Current remains4.14626296%slower than old median; five paired gaps are+3.20905–5.35567%,one pair favorscurrent0.78562%. Goal of beating old assembly wholeARM remains open. All18recorded outputs match22non-timing semantic fields and0allocations/input. Quality was skipped; equal recorded fields do not establish bitstream identity or quality equality. Component and whole windows are separate and no gap-closure attribution/subtraction is made.

All source/binary SHA bindings were independently rechecked after the windows. Actual1947candidate source and finalCLI/DSP receipts are preserved, as are the priorpacket source/binary references, allmath/state/codegen tests, exact drivers/orders/commands/timestamps/stderr, and root independentaudits. Executables alone are omitted with SHA,size,privatepath receipts; full private snapshots/binaries remain unchanged. No entire snapshots are copied. Review-script syntax/scope failures are preserved; no codec output failure was discarded. No README,dependency,assembly,Git or root source edits were made by this evidence work.


Publication: `raw-evidence.tar.gz` preserves all 9 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
