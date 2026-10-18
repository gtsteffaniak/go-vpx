# Focused native amd64 SSE8 CI overlay

Delivery is private and reviewable. Apply only the two files under `overlay/`:

- `.github/workflows/sse8-amd64-bench.yml`
- `scripts/benchmark-sse8-amd64.sh`

No root repository files, runtime kernels, existing full benchmark workflows,
existing full 338-case benchmark scripts, CI heads, or PRs were edited by this lane.
The full native+C jobs remain the final broader evidence. This focused lane has
independent concurrency and does not cancel active measurements.

## Runtime protocol

The workflow checks out the exact PR head (or workflow commit for manual runs).
The driver resolves the supplied head/base refs and old assembly 691cd051 to full
commit IDs once, records them, and creates three temporary git archives. It
copies only the existing head fixture
`internal/vp8/dsp/sse8_rows_native_test.go` to base and assembly. All production
source remains pinned. Builds run from the root checkout through its
`scripts/go.sh -C archive` wrapper, preserving the project's common-dir cache
scope. Go 1.27.2 is required; SIMD is enabled, GOAMD64=v1, GOMAXPROCS=1, PGO is off,
and module files are readonly.

Before timing, head AVX2 dispatch must actually execute and pass, targeted head
SSE8 correctness/race/checkptr2 tests run, and all three binaries execute every
fixture once with independent result/input/canary checks. These one-shot timing
values are excluded from performance summaries. Every measured benchmark also
checks its result and immutable backing storage before/after its B.Loop.

Six rounds use all six variant permutations, so every variant appears twice in
each order position. Each sample runs 80ms. The measured inventory is 36 public
B.Loop cases × 6 rounds × 3 variants = 648 samples, plus 4 existing legacy B.N diagnostics
× 6 rounds × 3 variants = 72 samples. All 720 samples and all 120 one-shot records are
retained. The four diagnostics are accurately named SSE8x4, SSE8x8, SSE8x16 and
SSE8x8PtrFast; all have width 8. No SSE16 fixture is included or relabeled.

The script rejects missing/extra/duplicate cases, invalid measured timings,
reported measured allocations, missing final PASS, and unbalanced variant order.
Performance ratios never fail the job. Every regression, all six paired ratios,
all raw allocation fields, base/head and assembly/head medians are preserved.
The Markdown summary separates live public cases from legacy diagnostics.
Legacy one-shot setup allocations and absent zero-resolution one-shot ns/op are
diagnostic only; measured samples still require positive finite timings and
reported 0 B/op, 0 allocs/op.

Artifacts include exact commits, CPU metadata/flags, toolchain/cache settings,
Git blob identities plus SHA256 of archived pixel/DSP/module sources, the exact
copied fixture and its hash, driver/workflow/wrapper hashes, binary hashes,
all 18 raw measured logs, all 3 raw one-shot logs, dispatch/correctness logs,
orders, complete JSON/Markdown summaries, status/phase and SHA256SUMS. All three
archives/binaries are deleted by a checked exit trap; partial evidence survives
failures and is hashed. No timing retries/outlier removal or timing threshold
failure is added. Workflow upload runs even on failure.

## Private validation

Pinned validation inputs:

- head c638b7c76eac857affa0780820e86c211cdd84d3
- base ab3abd6b335242d5732fa1956a63b7591eaf781b
- old assembly 691cd0512c4829d9ee40d483205dd223f9eb028d

Checks passed:

- POSIX shell syntax and YAML syntax; all embedded Python programs compile.
- Linux/amd64 and Darwin/amd64 DSP crossbuilds for all three variants.
- Final readonly Linux/amd64 builds for all three variants.
- Rosetta 1x correctness/canary execution of all 36 public and 4 legacy cases in every variant.
- Actual advertised-AVX2 head dispatch on Rosetta.
- Parser tests exercise 720 synthetic samples, every actual case name, retained
  all-case regressions, scientific numeric formatting, omitted one-shot timings,
  legacy one-shot setup allocations, and intentional malformed inventories,
  duplicates, allocation/timing/PASS failures and unbalanced order.
- Exit 0/exit 7 shell finalization preserves status, deletes temporary archives,
  retains partial evidence and verifies the final manifest.
- The actual driver refuses timing on Darwin/arm64 before any Go command or
  output creation. No platform bypass was added.

The parser's synthetic values are software test inputs, not measurements.
Rosetta 1x values are correctness-only and support no speed conclusion. No native
Linux timing or end-to-end GitHub Actions run was performed privately. The
workflow's first native run is needed for actual performance feedback.

`validate.py`, `parser_tests.py`, and `finalization_tests.py` reproduce these
checks using private snapshots/temporary directories. Validation binaries/full
snapshots remain outside this delivery. Evidence and source files are hashed in
SHA256SUMS; no experiment content is written to any README.
