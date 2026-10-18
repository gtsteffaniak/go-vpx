//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"bytes"
	"os"
	"os/exec"
	"simd/archsimd"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

func TestSSE8PtrExactWindows(t *testing.T) {
	t.Logf("AVX2=%v Enabled=%v", archsimd.X86.AVX2(), Enabled)
	a, b := make([]byte, 64), make([]byte, 64)
	for _, strides := range [][2]int{{0, 0}, {1, 3}, {8, 8}, {9, 17}, {31, 37}, {1312, 1320}} {
		ss, rs := strides[0], strides[1]
		n, m := 7*ss+8, 7*rs+8
		for align := range 8 {
			for mode := range 3 {
				a, b := make([]byte, n+align)[align:], make([]byte, m+align+1)[align+1:]
				for i := range a {
					a[i] = byte(i*37 + 11)
				}
				for i := range b {
					b[i] = byte(i*19 + 53)
				}
				if mode == 1 {
					for i := range a {
						a[i] = 255
					}
					clear(b)
				} else if mode == 2 {
					z := make([]byte, max(n, m)+5)
					for i := range z {
						z[i] = byte(i*17 + 29)
					}
					a, b = z[:n:n], z[5:5+m:5+m]
				}
				beforeA, beforeB := bytes.Clone(a), bytes.Clone(b)
				want := SSE8Scalar(a, ss, b, rs, 8)
				got := SSE8Ptr(&a[0], ss, &b[0], rs)
				if got != want || !bytes.Equal(a, beforeA) || !bytes.Equal(b, beforeB) {
					t.Fatalf("strides%v align%d mode%d got%d want%d", strides, align, mode, got, want)
				}
			}
		}
	}
	for base := 0; base < 1<<16; base += 64 {
		for i := range a {
			v := base + i
			a[i], b[i] = byte(v>>8), byte(v)
		}
		got := SSE8Ptr(&a[0], 8, &b[0], 8)
		if got != SSE8Scalar(a, 8, b, 8, 8) {
			t.Fatalf("byte pairs base%d got%d", base, got)
		}
	}
	old := Enabled
	Enabled = false
	defer func() { Enabled = old }()
	if got, want := SSE8Ptr(&a[0], 8, &b[0], 8), SSE8Scalar(a, 8, b, 8, 8); got != want {
		t.Fatalf("backend disabled fallback got%d want%d", got, want)
	}
}

// Invalid unsafe.Slice geometry throws under checkptr and cannot be recovered.
// Compare original and new fallback in child processes under the same flags.
func TestSSE8PtrFallbackPanics(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	cases := [][3]int{{-maxInt - 1, 8, 0}, {8, -maxInt - 1, 0}, {-2, 8, 0}, {8, -2, 0}, {maxInt, 8, 0}, {8, maxInt, 0}, {8, 8, 1}, {8, 8, 2}}
	if raw := os.Getenv("GOVPX_SSE8_INVALID_CASE"); raw != "" {
		i, err := strconv.Atoi(raw)
		if err != nil {
			t.Fatal(err)
		}
		tc := cases[i]
		a, b := make([]byte, 64), make([]byte, 64)
		ap, bp := &a[0], &b[0]
		if tc[2] == 1 {
			ap = nil
		} else if tc[2] == 2 {
			bp = nil
		}
		if os.Getenv("GOVPX_SSE8_INVALID_API") == "old" {
			SSE8(unsafe.Slice(ap, 7*tc[0]+8), tc[0], unsafe.Slice(bp, 7*tc[1]+8), tc[1], 8)
		} else {
			SSE8Ptr(ap, tc[0], bp, tc[1])
		}
		t.Fatal("invalid geometry unexpectedly returned")
	}
	for i := range cases {
		kind := ""
		for _, api := range []string{"old", "new"} {
			cmd := exec.Command(os.Args[0], "-test.run=^TestSSE8PtrFallbackPanics$")
			cmd.Env = append(os.Environ(), "GOVPX_SSE8_INVALID_CASE="+strconv.Itoa(i), "GOVPX_SSE8_INVALID_API="+api)
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("case%d %s unexpectedly exited successfully", i, api)
			}
			got := "len"
			if strings.Contains(string(out), "ptr is nil") {
				got = "nil"
			} else if !strings.Contains(string(out), "unsafe.Slice: len out of range") {
				t.Fatalf("case%d %s unexpected failure: %s", i, api, out)
			}
			if kind != "" && got != kind {
				t.Fatalf("case%d fallback panic kind changed %s to%s", i, kind, got)
			}
			kind = got
		}
	}
}
