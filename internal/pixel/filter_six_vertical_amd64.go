//go:build goexperiment.simd && amd64 && !purego

package pixel

func filterVerticalSix8(src []byte, ss int, dst []byte, ds, w, h int, taps []int16) bool {
	return false
}
func filterVerticalFour8(src []byte, ss int, dst []byte, ds, w, h int, taps []int16) bool {
	return false
}
