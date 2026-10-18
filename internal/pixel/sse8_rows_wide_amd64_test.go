//go:build goexperiment.simd && amd64 && !purego

package pixel

import (
	"simd/archsimd"
	"testing"
	"unsafe"
)

func TestSSE8RowsWideDispatch(t *testing.T) {
	t.Logf("AVX2=%v Enabled=%v available=%v", archsimd.X86.AVX2(), Enabled, sse8WideAvailable)
	if !Enabled {
		t.Skip("wide helper requires public Enabled guard; scalar fallback covered")
	}
	a, b := make([]byte, 64), make([]byte, 64)
	for i := range a {
		a[i] = 255
	}
	got := sse8RowsWide(unsafe.Pointer(unsafe.SliceData(a)), 8, unsafe.Pointer(unsafe.SliceData(b)), 8, 8)
	if got != 64*255*255 {
		t.Fatalf("wide sum%d", got)
	}
}
