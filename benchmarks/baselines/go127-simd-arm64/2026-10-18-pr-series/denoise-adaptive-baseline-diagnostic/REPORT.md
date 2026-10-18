# PR112 narrow native baseline comparison

[Oracle-only run39162342059](https://github.com/thesyncim/govpx/actions/runs/39162342059) failed the original assertion on both references. It did not measure performance.

- Baseline88fe29f7 and candidatec3815d05 produced the identical isolated65x49/threads2 stream hash`731832e65341b129394786ac2aa2e361f9c17461f8787cecb0e498955937d5d0`.
- OriginalARMexpected hash`737cefc7e128dd3a1bf750b7cdcf385684517a95f845b70bcfa53b01c973fd8b` remains unchanged.
- Prior fulltest candidate hash was`dd098e...`; this context difference is retained and unproven.
- AMD EPYC7763, actualAVX2/Enabled=true,AVX512=false; identical sharedfixturebytes, both binary/source bindings saved.
- Zero benchmark invocations and zero whole-codec builds. No C/profiles/fullCI/pin changes.

This proves baseline/candidate parity for the isolated case on this host. It does not establish a universal architecture-specific stream pin. PR112 remains on hold pending root diagnosis.
