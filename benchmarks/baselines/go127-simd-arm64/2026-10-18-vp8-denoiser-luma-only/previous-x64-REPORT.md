# Independent VP8 work-elision native screens

https://github.com/thesyncim/govpx/actions/runs/39155920162

IntelXeonPlatinum8370C;actualAVX2/AVX512/Enabled=true. Both independent overlays screened and ordinary/MLstream oracles passed. Two40ms AB/BApairs,25cases100zeroallocsamples. No whole/C/profiles/fullrun; kernel/component scope only.

## PR107

Omit discarded accepted-pass dequantized coefficients and EOB staging.

| Case | Baseline ns | Candidate ns | Pairedspeedups | Median |
|---|---|---|---|---:|
| BenchmarkVP8AcceptedQuantPipeline/dc | 613.5,619.5 | 542.3,543.8 | 1.131293x,1.139206x | 1.135249x |
| BenchmarkVP8AcceptedQuantPipeline/dense | 620.7,614.1 | 542.1,541.7 | 1.144992x,1.133653x | 1.139323x |
| BenchmarkVP8AcceptedQuantPipeline/zero | 446.0,444.4 | 368.5,367.1 | 1.210312x,1.210569x | 1.210441x |

## PR108

Stop full-pel neighborhoods once initial cost reaches the attainable global minimum.

| Case | Baseline ns | Candidate ns | Pairedspeedups | Median |
|---|---|---|---|---:|
| BenchmarkAstraVP8MinimumMotionCost/fast/hex/dense | 1235.0,1234.0 | 1237.0,1238.0 | 0.998383x,0.996769x | 0.997576x |
| BenchmarkAstraVP8MinimumMotionCost/fast/hex/exact | 1026.0,1026.0 | 828.9,830.8 | 1.237785x,1.234954x | 1.236370x |
| BenchmarkAstraVP8MinimumMotionCost/fast/hex/flat | 1059.0,1060.0 | 830.2,843.5 | 1.275596x,1.256669x | 1.266132x |
| BenchmarkAstraVP8MinimumMotionCost/fast/hex/improved | 1084.0,1083.0 | 1080.0,1078.0 | 1.003704x,1.004638x | 1.004171x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method1/dense | 571.6,570.9 | 570.4,572.5 | 1.002104x,0.997205x | 0.999655x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method1/exact | 411.2,411.0 | 178.8,179.7 | 2.299776x,2.287145x | 2.293461x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method1/flat | 448.1,444.8 | 179.3,188.5 | 2.499163x,2.359682x | 2.429423x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method1/improved | 317.1,317.7 | 315.5,315.7 | 1.005071x,1.006335x | 1.005703x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method1/partial | 770.0,766.5 | 478.4,476.8 | 1.609532x,1.607592x | 1.608562x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method1/shifted | 409.2,406.9 | 179.1,179.4 | 2.284757x,2.268116x | 2.276437x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method2/dense | 412.1,412.7 | 417.5,416.1 | 0.987066x,0.991829x | 0.989447x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method2/exact | 307.8,306.6 | 179.0,178.5 | 1.719553x,1.717647x | 1.718600x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method2/flat | 325.6,324.7 | 178.8,179.0 | 1.821029x,1.813966x | 1.817498x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method2/improved | 263.6,263.2 | 261.9,261.5 | 1.006491x,1.006501x | 1.006496x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method2/partial | 644.3,646.0 | 479.5,476.5 | 1.343691x,1.355719x | 1.349705x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method2/shifted | 306.3,306.6 | 179.1,178.7 | 1.710218x,1.715725x | 1.712971x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method3/dense | 409.8,409.9 | 405.9,407.0 | 1.009608x,1.007125x | 1.008367x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method3/exact | 289.3,288.8 | 134.3,134.9 | 2.154133x,2.140845x | 2.147489x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method3/flat | 317.2,317.4 | 134.4,135.2 | 2.360119x,2.347633x | 2.353876x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method3/improved | 293.7,292.4 | 293.8,293.2 | 0.999660x,0.997271x | 0.998466x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method3/partial | 636.8,634.0 | 429.2,427.4 | 1.483691x,1.483388x | 1.483539x |
| BenchmarkAstraVP8MinimumMotionCost/fullpel/method3/shifted | 288.7,287.2 | 134.8,134.4 | 2.141691x,2.136905x | 2.139298x |

