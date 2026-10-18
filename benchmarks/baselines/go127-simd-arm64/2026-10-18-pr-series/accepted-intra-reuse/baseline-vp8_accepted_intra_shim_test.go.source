package govpx

import (
	vp8common "github.com/thesyncim/govpx/internal/vp8/common"
	vp8dec "github.com/thesyncim/govpx/internal/vp8/decoder"
	vp8enc "github.com/thesyncim/govpx/internal/vp8/encoder"
)

func addWholeIntraResidualToAnalysisMacroblockCoeffs(img *vp8common.Image, row, col int, mode *vp8dec.MacroblockMode, c *vp8enc.MacroblockCoefficients, d *vp8common.MacroblockDequant) bool {
	return acceptedIntraLegacy(img, row, col, mode, c, d)
}
