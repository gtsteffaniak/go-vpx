# PR51 additive fixture and harness validation

This supplements the immutable integrated ARM51 performance archive without changing its source, measurements or manifests. No new throughput measurements or acceptance reruns are included.

The exact shared token fixture passes against the old691 native ARM source through a Go overlay:12,288full-state/error/panic/backing oracle cases and12excluded one-shot B.Loop cases, all0B/op and0allocs/op. The12one-shot durations are compatibility checks, never throughput denominators. Full1730Go source-map equality proves the old snapshot remains unchanged. Fixture SHA, overlay, complete output, exact officialSDK/CGO1/PGOoff wrapper command/environment are retained. The command receipt was recorded after execution from the actual invocation; it contains no invented execution timestamps. This establishes native ARM compatibility only, without native-x64 or clean-assembly performance inference.

The independent harness review preserves inventory.json, validation.json and capture-schema-proof.json plus their validation source and output. All four reviewed ROOT scripts match the recorded SHA/size. Static shell/embedded-Python/helper checks pass. All549previous benchmark names/timing-protocol blocks remain preserved; the prospective full matrix is561cases/10098timings,35common fixtures,132source hashes. These counts describe a prepared protocol, not completed new CI results.

Nine capture-schema tests pass with explicitly SYNTHETIC command/tool stubs. They are not native execution, GNU Linux disassembly or throughput evidence. Actual native script builds/parsers/feature selection remain future CI gates. The original optional/focused timing and parser blocks retain their reviewed semantics; captured metadata is additive after timing. Complete scripts, proof outputs and SHA bindings are retained, with no fake native addresses or source/layout causal claim.

No ROOT production/README/Git edits, source mutation or benchmark rerun occurred during this publication. All bytes are covered by MANIFEST.json and the published SHA256SUMS.


Publication: `raw-evidence.tar.gz` preserves all 17 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
