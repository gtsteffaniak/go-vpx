# ARM64 fused edges and full-row residual gather

Apple M4 Max, macOS ARM64, Go 1.27.2, GOEXPERIMENT=simd, GOMAXPROCS=1, PGO off. Repeated measurements ran serially after builds/tests completed.

Six alternating kernel pairs compare bd25795a with fused internal filter edges, then fused edges with full-row gather. Horizontal internal edges improve roughly 36 to 23 ns, vertical 99 to 77 ns; residual gather16x16 improves 61 to 41 ns. All allocate zero.

Six public quartets use VP8 1280x720, 60 input frames, 30 fps, 2500 kbps, realtime CPU-used8, threads1, encode-only, auto-libvpx=false. Variants: asm=main691cd051, perf20=bd25795a, perf21=plus fused edges, perf22=plus full-row gather and amd64 high-byte widening (inactive on ARM). Median ns/frame: 7,330,123.5; 10,757,896.5; 10,596,859.5; 10,444,096. Combined change improves 2.92% against bd25795a, but remains 42.48% slower than original assembly. Every output byte/count, quantizer statistic and allocation field matches across all samples; 58 frames encoded, 2 dropped, 685957 output bytes, zero allocations/frame.
