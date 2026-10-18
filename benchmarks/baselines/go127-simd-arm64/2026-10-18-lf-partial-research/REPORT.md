# Rejected ARM partial C/D filter-packet cache

Four balanced native public Picker rounds find no qualifying 1.5× component gain. Flat/ramp/noop cases improve only 1.013666–1.022311× the parent. Textured smooth/mixed cases range 0.992246–1.001948×, with mixed paired signs and a pitch32/smooth regression in three of four pairs. Reject this private design; no fresh six-round or whole-codec timing was run.

Parent: `ab3abd6b335242d5732fa1956a63b7591eaf781b`. Native Apple M4 Max, Go 1.27.2 darwin/arm64, SIMD enabled, PGO off, GOMAXPROCS=1. Two 150ms benchmark runs per order direction were balanced AB/BA/BA/AB. Both variants use the exact existing ten-case BenchmarkVP8PickerFilterNative fixture, with testing.B.Loop, full input-window resets, independent public separate-edge reference and complete backing-buffer comparison. Every one-shot fixture result passes. All eighty timed samples report zero B/op and allocs/op.

## Source and bounds

This is a distinct partial packet cache, not the previously rejected full20-column input tile. The existing wide initial eight-column loader and finalized-output stores are unchanged. For incoming words B/C/D, the candidate loads the p+4 sixteen-byte row windows once. B uses odd-word interleave plus low-half zip; D shares the same odd intermediates and uses high-half zip; C uses even-word interleave and high-half zip. Only eight packed C/D vectors remain live, and their final four-column transposes are deferred until consumption. Sequential corrected-neighbor carries and all scalar byte arithmetic remain unchanged.

The established wrapper proves sixteen disjoint twenty-byte rows. Every wide load remains columns4..19. Cached C (12..15) is consumed at edge12 before its own correction, after edge8 modifies only6..9; cached D (16..19) is consumed at edge16 after edge12 modifies only10..13. Corrected q vectors continue through the original carry assignments. No eager image stores, widened window, alias restriction or new terminal pointer is introduced. Horizontal source prefix and other production files are byte-identical. The private generated body was edited directly for screening; scripts/gensimd.go was not changed, and this is not a regeneration-ready production patch.

## Diagnostics and correctness

Existing pixel/DSP scalar loop-filter oracles, 7200 macro/internal cases, skipped/later-edge cases and alias/decline guards pass ordinary and race+checkptr=2 tests. The optimized vertical body changes698→668 instructions, nonstack vector loads57→27 including constants (thirty-two incoming data reloads removed, with changed constant scheduling). Stack vector loads0/stores19 and total stack-memory instruction count45 remain unchanged; frame336→352 bytes. Calls remain initial-load/final-store helpers and cold stack growth. The intended load savings were real and not canceled by additional spill instructions, but those diagnostics did not translate to an accepted public speedup.

| Public case | Parent ns | Candidate ns | Parent / candidate | Faster pairs |
|---|---:|---:|---:|---:|
| pitch32/flat | 46.0500 | 45.0450 | 1.022311× | 4/4 |
| pitch32/ramp | 45.5800 | 44.8650 | 1.015937× | 4/4 |
| pitch32/smooth | 87.0200 | 87.7000 | 0.992246× | 1/4 |
| pitch32/mixed | 87.0000 | 86.9950 | 1.000057× | 3/4 |
| pitch32/noop | 40.4250 | 39.8800 | 1.013666× | 4/4 |
| pitch1312/flat | 45.7450 | 44.9250 | 1.018253× | 4/4 |
| pitch1312/ramp | 45.5750 | 44.8950 | 1.015146× | 4/4 |
| pitch1312/smooth | 87.4200 | 87.2500 | 1.001948× | 2/4 |
| pitch1312/mixed | 87.2850 | 87.3300 | 0.999485× | 3/4 |
| pitch1312/noop | 40.4700 | 39.6500 | 1.020681× | 4/4 |

All samples, slower controls and outliers are retained. The old NEON implementation was not measured in this screen; no assembly or whole-codec timing claim is made. Shared desktop timings and instruction counts do not establish a general bottleneck. Prior full-tile evidence was read to avoid repeating it, rather than pooled or subtracted across runs.

## Preserved evidence

Exact relevant parent/candidate sources, one-file patch, original public fixture, correctness/race logs, build commands/logs, both disassemblies, process-idle proof, every raw screen sample and compact summary are included. Binary hashes pin the already-built private test executables without copying them. The reproduction driver reconstructs parent snapshots, overlays only the changed body, validates and reruns the recorded public protocol through the root project wrapper. No root/index/branch/README edits were made, and all jobs were stopped before packaging. `SHA256SUMS` covers all other package members; `verification.json` recomputes sample counts, allocations and medians.
