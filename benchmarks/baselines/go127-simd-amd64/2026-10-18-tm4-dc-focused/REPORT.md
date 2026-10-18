# PR48 focused SSE8 controls

Exact head `31123a9ecb7fa3b90afdb293beb06dadb2a6f857`, parent `d2610cd0747edd72f7656068dfd4604d29e45307`; separate EPYC 9V45 job. All 55 raw files, original checksum/disassembly/symbol records, 720 timed samples and 120 excluded one-shot samples are retained. All timed samples are zero B/alloc. Head/base/old source inventories contain 294/282/147 Go files; fixture identity and dispatch are verified.

All 36 public SSE8 medians are controls because PR48 changes ARM production only. Parent/current median ratios span 0.991905–1.030084×: seven median regressions, zero all-six regressions, one all-six win. The minimum paired ratio is 0.799643× and remains archived. Pointer medians span 1.002013–1.020612× with no universal six-pair win. The four separate legacy B.N diagnostics and every paired result are preserved. Actual measured GNU pointer loops have 28 instructions for both parent and head. Source identity does not remove observed drift; no cause or new performance claim is inferred. Old-reference results remain CPU/protocol scoped and are not codec-wide claims.

Original own/root audits and trace proof are archived alongside all source hashes and allocation samples. Merge receipt proves normal merge and exact tested-tree preservation. This job is separate from EPYC 7763 full/native-C jobs and M4 Max ARM evidence. No samples are pooled.


Publication: `raw-evidence.tar.gz` preserves all 9 original frozen files byte for byte, including nested raw archives, manifests and source/oracle evidence. Extract outside the repository to inspect historical Go/assembly references. Outer files are covered by `SHA256SUMS`.
