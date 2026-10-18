//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"
	"unsafe"
)

func fp16GuardParentProbe(coeff, iscan, qcoeff, dqcoeff []int16, classes []byte, n int) bool {
	type span struct{ p, bytes uintptr }
	spans := [5]span{
		{uintptr(unsafe.Pointer(unsafe.SliceData(coeff))), uintptr(n) * 2},
		{uintptr(unsafe.Pointer(unsafe.SliceData(iscan))), uintptr(n) * 2},
		{uintptr(unsafe.Pointer(unsafe.SliceData(dqcoeff))), uintptr(n) * 2},
	}
	count := 3
	if qcoeff != nil {
		spans[count] = span{uintptr(unsafe.Pointer(unsafe.SliceData(qcoeff))), uintptr(n) * 2}
		count++
	}
	if classes != nil {
		spans[count] = span{uintptr(unsafe.Pointer(unsafe.SliceData(classes))), uintptr(n)}
		count++
	}
	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			a, b := spans[i], spans[j]
			if a.p <= b.p && b.p-a.p < a.bytes || b.p < a.p && a.p-b.p < b.bytes {
				return false
			}
		}
	}
	return true
}

func fp16GuardSingleProbe(coeff, iscan, qcoeff, dqcoeff []int16, classes []byte, n int) bool {
	type span struct{ p, bytes uintptr }
	spans := [5]span{
		{uintptr(unsafe.Pointer(unsafe.SliceData(coeff))), uintptr(n) * 2},
		{uintptr(unsafe.Pointer(unsafe.SliceData(iscan))), uintptr(n) * 2},
		{uintptr(unsafe.Pointer(unsafe.SliceData(dqcoeff))), uintptr(n) * 2},
	}
	count := 3
	if qcoeff != nil {
		spans[count] = span{uintptr(unsafe.Pointer(unsafe.SliceData(qcoeff))), uintptr(n) * 2}
		count++
	}
	if classes != nil {
		spans[count] = span{uintptr(unsafe.Pointer(unsafe.SliceData(classes))), uintptr(n)}
		count++
	}
	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			a, b := spans[i], spans[j]
			if a.p-b.p+a.bytes-1 < a.bytes+b.bytes-1 {
				return false
			}
		}
	}
	return true
}

func fp16GuardUnrolledProbe(coeff, iscan, qcoeff, dqcoeff []int16, classes []byte, n int) bool {
	c := uintptr(unsafe.Pointer(unsafe.SliceData(coeff)))
	i := uintptr(unsafe.Pointer(unsafe.SliceData(iscan)))
	d := uintptr(unsafe.Pointer(unsafe.SliceData(dqcoeff)))
	bytes := uintptr(n) * 2
	sameLimit := bytes + bytes - 1
	if i-c+bytes-1 < sameLimit || d-c+bytes-1 < sameLimit || d-i+bytes-1 < sameLimit {
		return false
	}
	if qcoeff != nil {
		q := uintptr(unsafe.Pointer(unsafe.SliceData(qcoeff)))
		if q-c+bytes-1 < sameLimit || q-i+bytes-1 < sameLimit || q-d+bytes-1 < sameLimit {
			return false
		}
	}
	if classes != nil {
		cl := uintptr(unsafe.Pointer(unsafe.SliceData(classes)))
		classBytes := uintptr(n)
		classLimit := classBytes + bytes - 1
		if cl-c+classBytes-1 < classLimit || cl-i+classBytes-1 < classLimit || cl-d+classBytes-1 < classLimit {
			return false
		}
		if qcoeff != nil && cl-uintptr(unsafe.Pointer(unsafe.SliceData(qcoeff)))+classBytes-1 < classLimit {
			return false
		}
	}
	return true
}

// These exact helper bodies are frozen diagnostic candidates, not production dispatch.
var fp16GuardProbeSink bool

type fp16GuardProbeCase struct {
	name         string
	n            int
	optional     int
	pairA, pairB int
}

func fp16GuardProbeCases() []fp16GuardProbeCase {
	var cases []fp16GuardProbeCase
	for _, n := range []int{64, 256, 1024} {
		for optional, label := range []string{"full", "q-nil", "classes-nil", "both-nil"} {
			cases = append(cases, fp16GuardProbeCase{fmt.Sprintf("disjoint/n%d/%s", n, label), n, optional, -1, -1})
		}
	}
	names := []string{"coeff", "iscan", "qcoeff", "dqcoeff", "classes"}
	for a := 0; a < 5; a++ {
		for b := a + 1; b < 5; b++ {
			cases = append(cases, fp16GuardProbeCase{fmt.Sprintf("decline/n256/%s-%s", names[a], names[b]), 256, 0, a, b})
		}
	}
	return cases
}

type fp16GuardProbeInputs struct {
	backing                       []byte
	coeff, iscan, qcoeff, dqcoeff []int16
	classes                       []byte
}

func fp16GuardProbeMake(c fp16GuardProbeCase) fp16GuardProbeInputs {
	span := 2*c.n + 64
	backing := make([]byte, 5*span+128)
	for i := range backing {
		backing[i] = byte(i*37 + i/7 + 11)
	}
	offsets := [5]int{64, 64 + span, 64 + 2*span, 64 + 3*span, 64 + 4*span}
	if c.pairA >= 0 {
		offsets[c.pairB] = offsets[c.pairA] + 1
	}
	view := func(off int) []int16 { return unsafe.Slice((*int16)(unsafe.Pointer(&backing[off])), c.n) }
	v := fp16GuardProbeInputs{backing, view(offsets[0]), view(offsets[1]), view(offsets[2]), view(offsets[3]), backing[offsets[4] : offsets[4]+c.n : offsets[4]+c.n]}
	if c.optional&1 != 0 {
		v.qcoeff = nil
	}
	if c.optional&2 != 0 {
		v.classes = nil
	}
	return v
}

// Independent interval oracle uses ordered starts/ends, not candidate modular math.
func fp16GuardProbeOracle(v fp16GuardProbeInputs, n int) bool {
	type span struct{ start, end uintptr }
	windows := []span{}
	add := func(p unsafe.Pointer, size uintptr) {
		start := uintptr(p)
		end := start + size
		if end < start {
			panic("oracle span wrapped")
		}
		windows = append(windows, span{start, end})
	}
	add(unsafe.Pointer(unsafe.SliceData(v.coeff)), uintptr(n)*2)
	add(unsafe.Pointer(unsafe.SliceData(v.iscan)), uintptr(n)*2)
	add(unsafe.Pointer(unsafe.SliceData(v.dqcoeff)), uintptr(n)*2)
	if v.qcoeff != nil {
		add(unsafe.Pointer(unsafe.SliceData(v.qcoeff)), uintptr(n)*2)
	}
	if v.classes != nil {
		add(unsafe.Pointer(unsafe.SliceData(v.classes)), uintptr(n))
	}
	for i := range windows {
		for j := 0; j < i; j++ {
			a, b := windows[i], windows[j]
			if a.start < b.end && b.start < a.end {
				return false
			}
		}
	}
	runtime.KeepAlive(v)
	return true
}
func TestFP16GuardComponentOracle(t *testing.T) {
	for _, c := range fp16GuardProbeCases() {
		t.Run(c.name, func(t *testing.T) {
			v := fp16GuardProbeMake(c)
			original := bytes.Clone(v.backing)
			want := fp16GuardProbeOracle(v, c.n)
			if want != (c.pairA < 0) {
				t.Fatal("fixture did not exercise declared policy")
			}
			if got := fp16GuardParentProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n); got != want {
				t.Fatal("parent", got, want)
			}
			if got := fp16GuardSingleProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n); got != want {
				t.Fatal("single", got, want)
			}
			if got := fp16GuardUnrolledProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n); got != want {
				t.Fatal("unrolled", got, want)
			}
			if got := fp16WindowsDisjoint(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n); got != want || got != fp16GuardUnrolledProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n) {
				t.Fatal("trial selected vs frozen unrolled/oracle mismatch")
			}
			if !bytes.Equal(v.backing, original) {
				t.Fatal("guard modified input backing")
			}
		})
	}
}

// Native Linux x64 only. Rosetta runs may execute the oracle test, never this timer.
// Every timed iteration calls one statically named guard; switch dispatch is outside b.Loop.
func BenchmarkFP16GuardComponent(b *testing.B) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" || os.Getenv("GOVPX_REQUIRE_NATIVE_GUARD") != "1" {
		b.Skip("native Linux x64 diagnostic only")
	}
	for _, variant := range []string{"parent", "single", "unrolled"} {
		for _, c := range fp16GuardProbeCases() {
			b.Run(variant+"/"+c.name, func(b *testing.B) {
				v := fp16GuardProbeMake(c)
				original := bytes.Clone(v.backing)
				want := fp16GuardProbeOracle(v, c.n)
				var got bool
				// Preflight all candidate results before timing; all backing checks remain outside timer.
				if fp16GuardParentProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n) != want || fp16GuardSingleProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n) != want || fp16GuardUnrolledProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n) != want {
					b.Fatal("independent oracle mismatch")
				}
				if fp16WindowsDisjoint(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n) != want {
					b.Fatal("actual trial guard mismatch")
				}
				b.ReportAllocs()
				switch variant {
				case "parent":
					for b.Loop() {
						got = fp16GuardParentProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n)
					}
				case "single":
					for b.Loop() {
						got = fp16GuardSingleProbe(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n)
					}
				case "unrolled":
					for b.Loop() {
						got = fp16WindowsDisjoint(v.coeff, v.iscan, v.qcoeff, v.dqcoeff, v.classes, c.n)
					}
				default:
					b.Fatal("unknown variant")
				}
				fp16GuardProbeSink = got
				if got != want || !bytes.Equal(v.backing, original) {
					b.Fatal("result/input backing mismatch")
				}
				runtime.KeepAlive(v)
			})
		}
	}
}
