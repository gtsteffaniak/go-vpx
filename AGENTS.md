# govpx agent instructions

## SIMD performance acceptance

- Keep correctness-verified kernel or component optimizations that achieve at least a 1.5× speedup in representative native benchmarks, including local measurements, even when overall codec performance is unchanged.
- Report kernel and end-to-end results separately, with hardware and workload scope. Preserve raw measurements and correctness evidence.
- Keep experiment reports and benchmark logs out of README files.

## Fast performance iteration

- During tuning, run affected correctness tests and a short native benchmark screen. Use broader tests and longer measurements once a candidate is ready for its PR, or when a failure requires them.
- Reuse unchanged source, oracle and benchmark evidence. Do not repeat full audits, profile collections or evidence packaging for each small kernel change.
- Report short screens as preliminary. Let CI run broad validation; do not block publishing a draft PR on redundant local validation.

## Build caches

- Treat all branches and Codex worktrees of this Git project as one cache scope, identified by `git rev-parse --git-common-dir`.
- Run Go builds and tests through `scripts/go.sh`, which invokes the installed project environment wrapper. In this workspace the wrapper is `/Users/thesyncim/.codex/bin/project-env`.
- Set a distinct `CODEX_AGENT_ID` for concurrent agents. A cache path containing `.codex/worktrees` is invalid; check `go env GOCACHE` through the wrapper before a large build if uncertain.
- Build wrappers must propagate the same project environment to child processes.
