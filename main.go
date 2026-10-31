// Command govpx is a small demo CLI for the govpx VP8/VP9 codec.
//
// Library import: github.com/gtsteffaniak/go-vpx/pkg/govpx
//
// Usage:
//
//	govpx info   <file.ivf>
//	govpx decode <file.ivf> <out.yuv>
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gtsteffaniak/go-vpx/internal/vpx/ivf"
	"github.com/gtsteffaniak/go-vpx/pkg/govpx"
)

type decoder interface {
	Decode(packet []byte) error
	NextFrame() (govpx.Image, bool)
	Close() error
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		fatal(err)
	}
	hdr, err := ivf.ParseHeader(data)
	if err != nil {
		fatal(err)
	}
	frames, err := ivf.Frames(data)
	if err != nil {
		fatal(err)
	}
	switch os.Args[1] {
	case "info":
		fmt.Printf("fourcc=%s size=%dx%d timebase=%d/%d frames=%d\n",
			hdr.FourCC[:], hdr.Width, hdr.Height,
			hdr.TimebaseNumerator, hdr.TimebaseDenominator, len(frames))
	case "decode":
		if len(os.Args) < 4 {
			usage()
		}
		if err := decode(hdr, frames, os.Args[3]); err != nil {
			fatal(err)
		}
	default:
		usage()
	}
}

func decode(hdr ivf.Header, frames []ivf.Frame, outPath string) error {
	var dec decoder
	var err error
	switch hdr.FourCC {
	case ivf.FourCCVP9:
		dec, err = govpx.NewVP9Decoder(govpx.VP9DecoderOptions{})
	default:
		dec, err = govpx.NewVP8Decoder(govpx.DecoderOptions{})
	}
	if err != nil {
		return err
	}
	defer dec.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	w := bufio.NewWriterSize(out, 1<<20)
	shown := 0
	for _, f := range frames {
		if err := dec.Decode(f.Data); err != nil {
			out.Close()
			return fmt.Errorf("frame %d: %w", f.Index, err)
		}
		for {
			img, ok := dec.NextFrame()
			if !ok {
				break
			}
			if err := writeI420(w, img); err != nil {
				out.Close()
				return err
			}
			shown++
		}
	}
	if err := w.Flush(); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	fmt.Printf("decoded %d frames to %s\n", shown, outPath)
	return nil
}

// writeI420 writes the visible planes of img as tightly packed I420.
func writeI420(w *bufio.Writer, img govpx.Image) error {
	cw, ch := (img.Width+1)/2, (img.Height+1)/2
	planes := []struct {
		data         []byte
		stride, w, h int
	}{
		{img.Y, img.YStride, img.Width, img.Height},
		{img.U, img.UStride, cw, ch},
		{img.V, img.VStride, cw, ch},
	}
	for _, p := range planes {
		for y := 0; y < p.h; y++ {
			if _, err := w.Write(p.data[y*p.stride : y*p.stride+p.w]); err != nil {
				return err
			}
		}
	}
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:\n  govpx info   <file.ivf>\n  govpx decode <file.ivf> <out.yuv>")
	os.Exit(2)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "govpx:", err)
	os.Exit(1)
}
