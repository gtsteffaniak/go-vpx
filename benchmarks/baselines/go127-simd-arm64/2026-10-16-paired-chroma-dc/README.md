# Paired chroma filters and DC reconstruction

Apple M4 Max, macOS ARM64, Go 1.27.2, `GOEXPERIMENT=simd`, `GOMAXPROCS=1`, PGO off. Builds and tests finished before each serialized timing run. The parent is `3eb51360`.

UV joins U/V eight-position edges into one sixteen-byte packet; aliased planes and the 255 boundary retain sequential calls. DC joins two adjacent 4x4 blocks. Both int16 and byte-saturation designs were tested; bytes won. DC retains safe slice loads/stores. Exact in-place and disjoint layouts run together; partial aliases or overlapping rows retain block order.

Six balanced DC triples, all zero allocations:

| Layout | Parent | int16 design | Byte design | Byte time change |
| --- | ---: | ---: | ---: | ---: |
| In-place pair | 13.32 ns | 7.559 ns | 7.131 ns | -46.46% |
| Disjoint pair | 12.695 ns | 7.003 ns | 6.559 ns | -48.33% |

Six fresh UV pairs, all zero allocations:

| Filter | Parent | Paired | Time change |
| --- | ---: | ---: | ---: |
| Horizontal | 24.58 ns | 16.49 ns | -32.91% |
| Vertical | 39.12 ns | 25.325 ns | -35.26% |
| Macroblock horizontal | 28.4 ns | 19.205 ns | -32.38% |
| Macroblock vertical | 44.755 ns | 31.235 ns | -30.21% |

The earlier UV version and its public pairs are retained separately. Final UV bounds prove the entire 8x8 tile before multiplying a stride. Independent libvpx scalar tests cover full allocations, aliases, alignments and threshold extremes. DC tests exhaust all 130,816 byte/delta combinations; other cases cover minimal spans, different strides, partial aliases, oversized dimensions and full int32 DC inputs. The UV profile confirms live paired dispatch (about 4% cumulative in this capture).

Ten public quintets compare five variants. Balanced cyclic/reversed execution puts each variant twice in each position and balances predecessor pairs. Workload: VP8 1280x720, 60 input frames, 30 fps, 2500 kbps, realtime, CPU-used 8, threads 1, encode-only, auto-libvpx=false.

| Variant | Median ns/frame |
| --- | ---: |
| `asm`: main691cd051 | 6,951,804 |
| `parent`: 3eb51360 | 10,407,510.5 |
| `uv`: final UV only | 10,245,667 |
| `dc`: byte DC only | 10,209,247 |
| `paired`: both | 9,999,657 |

Combined time falls 3.92% versus parent (all ten pairs -5.14% to -2.32%), but still exceeds original assembly by 43.84% on this ARM workload. All output/count/quantizer/allocation fields match across all 50 samples: 685,957 output bytes, 58 encoded/2 dropped frames, zero allocations/frame.

ARM and Rosetta AVX2 pixel, both DSP and both encoder suites pass. Native x64 measurement of these paired changes is pending. No cross-CPU performance inference is made.
