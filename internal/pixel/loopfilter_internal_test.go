//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import (
	"bytes"
	"testing"
)

func TestInternalLoopFilterMinimalWindows(t *testing.T) {
	if !Enabled {
		t.Skip("native SIMD unavailable")
	}
	for _, stride := range []int{16, 17, 32} {
		for _, vertical := range []bool{false, true} {
			step, outer := stride, 1
			if vertical {
				step, outer = 1, stride
			}
			for align := range 8 {
				for trial := range 12 {
					original := make([]byte, align+15*stride+16)
					for i := range original {
						original[i] = byte(i*37 + trial*19)
						if trial%2 == 0 {
							original[i] = byte(112 + (i/stride+i%stride+trial)%17)
						}
					}
					for _, blimit := range []int{0, 64, 254, 255} {
						got, want := bytes.Clone(original), bytes.Clone(original)
						filtered := LoopFilterInternal16(got[align:], step, outer, blimit, 32, 0)
						if blimit == 255 {
							if filtered || !bytes.Equal(got, original) {
								t.Fatal("255 boundary must fall back without writes")
							}
							continue
						}
						if !filtered {
							t.Fatal("native dispatch not selected")
						}
						for edge := 1; edge <= 3; edge++ {
							LoopFilter(want[align:], 4*edge*step, step, outer, 16, blimit, 32, 0, 4)
						}
						if !bytes.Equal(got, want) {
							t.Fatalf("stride%d vertical%v align%d trial%d blimit%d", stride, vertical, align, trial, blimit)
						}
					}
				}
			}
		}
	}
}
