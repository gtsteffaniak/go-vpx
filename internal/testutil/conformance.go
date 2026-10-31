package testutil

import "github.com/gtsteffaniak/govpx/internal/vpx/conformance"

var ErrInvalidOracleOutput = conformance.ErrInvalidOracleOutput

type FrameChecksum = conformance.FrameChecksum

func SameFrameChecksum(a FrameChecksum, b FrameChecksum) bool {
	return conformance.SameFrameChecksum(a, b)
}

func ParseFrameChecksumJSONLines(data []byte) ([]FrameChecksum, error) {
	return conformance.ParseFrameChecksumJSONLines(data)
}
