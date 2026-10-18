# Deterministic noisy-stream fixture correction

The original test inherited realtime CpuUsed=+8 from newSizedTestEncoder. At tiny frame sizes that enables actual wall-clock feedback and changes later picker dispatch. Thus an ARM hash cannot be assumed stable across native hosts or different test selections.

The native oracle-only diagnostic39162342059 ran the unchanged assertion on original baseline88fe29 and candidate112 in the same host: both produced exactly731832e65341b129394786ac2aa2e361f9c17461f8787cecb0e498955937d5d0 for the isolated partial Threads2 case. The earlier all-case current hashdd098e... was also kept. The diagnostic performed0 benchmark calls and0 whole-codec builds. No production bug or tolerance adjustment was inferred from these hashes.

The corrected fixture sets CpuUsed=-8 before encoder construction and asserts effective speed8 on every packet. Immutable new pins were generated first from unchanged PR110 runtime, then matched by candidate. All4 configurations repeated3 times on each ARM binary. This changes only the test workload from adaptive to fixed speed; runtime and timed codec binaries remain unchanged. All old adaptive pins, failures and measurements remain retained. Exact settings/pins/source/binary hashes and logs are in fixed-speed-stream-proof/. A new native x64 run will exercise the same corrected fixture on both refs.
