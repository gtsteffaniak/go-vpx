package dsp

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/gtsteffaniak/go-vpx/internal/vp9/tables"
)

func refHoriz(src []byte, ss int, dst []byte, ds int, f *[tables.SubpelShifts][tables.SubpelTaps]int16, x0, step, w, h, off int) {
	start := off - 3
	for y := 0; y < h; y++ {
		xq := x0
		for x := 0; x < w; x++ {
			b := start + y*ss + xq>>4
			sum := 0
			for k := 0; k < 8; k++ {
				sum += int(src[b+k]) * int(f[xq&15][k])
			}
			dst[y*ds+x] = clipPixel(roundPowerOfTwo(int32(sum), tables.FilterBits))
			xq += step
		}
	}
}

func refVert(src []byte, ss int, dst []byte, ds int, f *[tables.SubpelShifts][tables.SubpelTaps]int16, y0, step, w, h, off int) {
	start := off - 3*ss
	for x := 0; x < w; x++ {
		yq := y0
		for y := 0; y < h; y++ {
			b := start + (yq>>4)*ss + x
			sum := 0
			for k := 0; k < 8; k++ {
				sum += int(src[b+k*ss]) * int(f[yq&15][k])
			}
			dst[y*ds+x] = clipPixel(roundPowerOfTwo(int32(sum), tables.FilterBits))
			yq += step
		}
	}
}

// TestConvolveUnscaledMatchesReference checks the step-16 fast paths are
// bit-identical to the original per-pixel loops for every filter set, phase
// and block size, including extreme pixel values that hit the clip.
func TestConvolveUnscaledMatchesReference(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	sets := []*[tables.SubpelShifts][tables.SubpelTaps]int16{
		&tables.SubPelFilters8, &tables.SubPelFilters8s, &tables.SubPelFilters8lp,
		&tables.SubPelFilters4, &tables.BilinearFilters,
	}
	const ss = 160
	src := make([]byte, ss*160)
	for iter := 0; iter < 400; iter++ {
		for i := range src {
			switch rng.Intn(4) {
			case 0:
				src[i] = 0
			case 1:
				src[i] = 255
			default:
				src[i] = byte(rng.Intn(256))
			}
		}
		f := sets[iter%len(sets)]
		w := []int{2, 4, 8, 16, 32, 64}[rng.Intn(6)]
		h := []int{2, 4, 8, 16, 32, 64}[rng.Intn(6)]
		q := rng.Intn(16 * 4)
		off := 8*ss + 8
		ds := w + rng.Intn(8)
		got := make([]byte, ds*h)
		want := make([]byte, ds*h)
		convolveHoriz(src, ss, got, ds, f, q, 16, w, h, off)
		refHoriz(src, ss, want, ds, f, q, 16, w, h, off)
		if !bytes.Equal(got, want) {
			t.Fatalf("horiz mismatch iter %d w=%d h=%d q=%d", iter, w, h, q)
		}
		convolveVert(src, ss, got, ds, f, q, 16, w, h, off)
		refVert(src, ss, want, ds, f, q, 16, w, h, off)
		if !bytes.Equal(got, want) {
			t.Fatalf("vert mismatch iter %d w=%d h=%d q=%d", iter, w, h, q)
		}
		// Averaging variants: blend into an existing prediction.
		for i := range got {
			got[i] = byte(rng.Intn(256))
		}
		copy(want, got)
		tmp := make([]byte, len(want))
		convolveAvgHoriz(src, ss, got, ds, f, q, 16, w, h, off)
		refHoriz(src, ss, tmp, ds, f, q, 16, w, h, off)
		avgInto(want, tmp, ds, w, h)
		if !bytes.Equal(got, want) {
			t.Fatalf("avg horiz mismatch iter %d", iter)
		}
		copy(want, got)
		convolveAvgVert(src, ss, got, ds, f, q, 16, w, h, off)
		refVert(src, ss, tmp, ds, f, q, 16, w, h, off)
		avgInto(want, tmp, ds, w, h)
		if !bytes.Equal(got, want) {
			t.Fatalf("avg vert mismatch iter %d", iter)
		}
	}
}

func BenchmarkConvolveVert16(b *testing.B) {
	src := make([]byte, 160*160)
	dst := make([]byte, 64*64)
	for i := range src {
		src[i] = byte(i * 7)
	}
	b.Run("fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			convolveVert(src, 160, dst, 64, &tables.SubPelFilters8, 5, 16, 64, 64, 8*160+8)
		}
	})
	b.Run("ref", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refVert(src, 160, dst, 64, &tables.SubPelFilters8, 5, 16, 64, 64, 8*160+8)
		}
	})
	b.Run("horizfast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			convolveHoriz(src, 160, dst, 64, &tables.SubPelFilters8, 5, 16, 64, 64, 8*160+8)
		}
	})
	b.Run("horizref", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			refHoriz(src, 160, dst, 64, &tables.SubPelFilters8, 5, 16, 64, 64, 8*160+8)
		}
	})
}

func avgInto(dst, p []byte, ds, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := y*ds + x
			dst[i] = uint8((int(dst[i]) + int(p[i]) + 1) >> 1)
		}
	}
}
