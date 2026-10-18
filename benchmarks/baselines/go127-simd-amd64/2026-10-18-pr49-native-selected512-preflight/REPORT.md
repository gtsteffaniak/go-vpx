# Failed native validation at5c37b7d7

The full-native job117651676914/run39119385796 failed before timing on AMD EPYC9V74 with actual SDK AVX512=true. Selected512 ordinary and race/checkptr2 oracle runs passed. The scalar-only lane disabled AVX and AVX2 but left AVX512 enabled, contradicting GOVPX_REQUIRE_FP32_OFF=1; the strict feature assertion correctly failed. Production Enabled=false prevented SIMD dispatch in that mode.

All4partial raw files and118exact source hashes are preserved. No main/kernel/optional throughput measurements completed. No speedup or performance acceptance is claimed. Selected512 correctness is predecessor evidence only for later source-identical heads.

Follow-up cb9bbdcd adds explicit cpu.avx512f=off to the two scalar-only FP32 commands in each of the full and focused scripts. Tests, feature assertions and1946Go source files are unchanged. This archive is private for later publication; no evidence-only push should restart that head.

No executables or repository snapshots are included. The118tracked source files are the exact failed head, distinct from subsequent repairs.


Publication: `raw-evidence.tar.gz` preserves all 129 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
