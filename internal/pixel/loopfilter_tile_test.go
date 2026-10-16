//go:build goexperiment.simd && (arm64 || amd64) && !purego

package pixel

import (
	"bytes"
	"fmt"
	"simd/archsimd"
	"testing"
)

func TestLoopFilterEdgeTileBoundsAndWrites(t *testing.T) {
	if !Enabled {
		t.Skip("native SIMD unavailable")
	}
	for _, width := range []int{4, 8} {
		for _, pitch := range []int{width, width + 3, 31} {
			for align := range 4 {
				t.Run(fmt.Sprintf("width%d/pitch%d/align%d", width, pitch, align), func(t *testing.T) {
					original := make([]byte, align+7*pitch+width)
					for i := range original {
						original[i] = byte(i*37 + 11)
					}
					var a [16]archsimd.Int16x8
					const first = 4
					loadEdgeTile(original, align, pitch, width, &a, first)
					for x := range width {
						var lanes [8]int16
						a[first+x].Store(lanes[:])
						for y, v := range lanes {
							if want := original[align+y*pitch+x]; v != int16(want) {
								t.Fatalf("column %d row %d: %d != %d", x, y, v, want)
							}
						}
						// The expected rows differ in every column, so swapped transpose
						// lanes and writes outside the filter window cannot pass silently.
						a[first+x] = a[first+x].Add(archsimd.BroadcastInt16x8(13)).Min(archsimd.BroadcastInt16x8(255))
					}
					for _, count := range []int{2, 4, 6} {
						if count > width {
							continue
						}
						offset := (width - count) / 2
						got, want := bytes.Clone(original), bytes.Clone(original)
						for y := range 8 {
							for x := offset; x < offset+count; x++ {
								i := align + y*pitch + x
								want[i] = byte(min(int(want[i])+13, 255))
							}
						}
						storeEdgeTile(got, align, pitch, first, width, first+offset, first+offset+count, &a)
						if !bytes.Equal(got, want) {
							t.Fatalf("store count %d changed wrong bytes", count)
						}
					}
				})
			}
		}
	}
}

// Compare the wide byte/AVX2 dispatch against the int16 kernel with exact
// final-row allocations, including the 255 block-limit fallback boundary.
func TestLoopFilterMinimalWindows(t *testing.T) {
	if !Enabled {
		t.Skip("native SIMD unavailable")
	}
	seed := uint32(0x531728)
	next := func() byte { seed = seed*1664525 + 1013904223; return byte(seed >> 24) }
	for _, kind := range []int{2, 4, 6} {
		left, right := 4, 3
		if kind == 2 {
			left, right = 2, 1
		}
		for _, n := range []int{8, 16} {
			for _, layout := range [][2]int{{n + 3, 1}, {1, left + right + 4}, {1, 1}, {3, 1}, {1, 3}} {
				step, outer := layout[0], layout[1]
				for align := range 4 {
					off := align + left*step
					for trial := range 12 {
						original := make([]byte, off+right*step+(n-1)*outer+1)
						for i := range original {
							if trial%3 == 0 {
								original[i] = 120 + next()%9
							} else {
								original[i] = next()
							}
						}
						for _, blimit := range []int{0, 1, 64, 254, 255} {
							limit, thresh := 255, 8
							if trial%2 == 0 {
								limit, thresh = 32, 0
							}
							got, want := bytes.Clone(original), bytes.Clone(original)
							if !LoopFilter(got, off, step, outer, n, blimit, limit, thresh, kind) {
								t.Fatal("native dispatch unavailable")
							}
							for base := 0; base < n; base += 8 {
								loopFilterNarrow8(want, off+base*outer, step, outer, blimit, limit, thresh, kind)
							}
							if !bytes.Equal(got, want) {
								t.Fatalf("kind%d n%d step%d outer%d align%d trial%d blimit%d", kind, n, step, outer, align, trial, blimit)
							}
						}
					}
				}
			}
		}
	}
}
