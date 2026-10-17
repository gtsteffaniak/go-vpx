# Native x64 optimized libvpx comparison at 9a442f66

Job: https://github.com/thesyncim/govpx/actions/runs/39031352624/job/117392274827
Artifact: https://github.com/thesyncim/govpx/actions/runs/39031352624/artifacts/12100701439
The libvpx-reference job passes on exact head 9a442f6635fa1d430dbb1f58d95c038bf21de83f.
All 74 raw files are copied unchanged, including environment, complete build
log/configuration, object list, source/binary hashes, command manifest and
individual results. This README is the only added file.

Native EPYC 7763, official Go 1.27.2, SIMD/AVX2, PGO off, one Go execution
thread and one encoder thread. Pristine libvpx v1.16.0 is pinned to immutable
commit 1024874c5919305883187e2953de8fcb4c3d7fa6. GCC 13.3.0 and NASM 2.16.01
build optimized upstream assembly with runtime CPU dispatch. The script
checks actual SSE2/AVX2 selection and assembly object presence.

| Workload | Go ms/input | libvpx ms/input | Go time delta | Encoded/dropped, both |
| --- | ---: | ---: | ---: | ---: |
| VP9 320x180, cpu8 | 4.353493 | 1.428450 | +204.77% | 117/3 |
| VP9 720p, cpu8 | 32.668384 | 10.101667 | +223.40% | 48/12 |
| VP8 720p, adaptive cpu8 | 22.183473 | 10.259817 | +116.22% | 58/2 |
| VP8 720p, fixed cpu-8 | 18.320404 | 6.781175 | +170.17% | 58/2 |

Six samples per case alternate Go-first/reference-first. Discarded warmups
precede all timing rounds; decoded quality is a separate final pass. Every
paired Go/C time delta is positive. The fixed-speed VP8 case is the strongest
matched-quality evidence: both codecs give PSNR 27.157611 dB, SSIM 0.925927,
723593 payload bytes and 58/2 encoded/dropped frames in the quality pass.
Its payload/count statistics match every timing sample. This does not prove
bitstream identity. Actual bitrate is 2894.372 kbps for both, 15.775% above
the 2500 kbps target; equivalent target settings do not imply exact target rate.

Adaptive VP8 performs different work. C delivers higher measured PSNR/SSIM
using fewer bytes: Go/C 27.983951/29.176639 dB, 0.908504/0.933088 SSIM,
685957/678205 bytes and 2743.828/2712.820 kbps. This ratio is reported as a
matched-input/control comparison, not identical-quality throughput. Adaptive
and fixed speed are distinct workloads and are not compared as speedups.

VP9 also differs in output/quality:

| Workload | PSNR Go/C dB | SSIM Go/C | Payload bytes Go/C | Actual kbps Go/C |
| --- | --- | --- | --- | --- |
| VP9 320p | 28.910974/28.969583 | 0.966373405/0.967011486 | 246893/246912 | 493.786/493.824 |
| VP9 720p | 27.516434/27.546117 | 0.949347795/0.949095623 | 642654/643059 | 2570.616/2572.236 |

All quality-pass payload, bitrate, keyframe/interframe-size and count fields
match corresponding timing samples within each codec/workload. Encode-only
timing samples skip quality decode; quality metrics come from the separate
pass rather than being asserted for unretained timed bitstreams. The command
manifest and raw parity flags preserve CBR, buffering, q-range, GOP, drop,
noise, threading and speed controls on the exact same generated I420 inputs.

The internal timer measures vpx_codec_encode calls. Reference wall medians
are separately 1.583220 / 11.042451 / 11.073850 / 7.500706 ms per input in
table order, including subprocess startup and its file I/O. Input generation
and raw-input serialization precede these timers. Go timing sums EncodeInto
calls; neither codec time includes its setup or decoded quality work. Mean
latency placeholders from the legacy helper are not used as quantiles.

The result establishes a substantial remaining gap against optimized C even
where Go SIMD beats earlier Go assembly. Absolute times across other jobs or
CPU models do not isolate implementation changes. No wall-timer fallback,
Rosetta timing, quality-equivalence assumption or discarded pair is used.
