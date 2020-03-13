package bxzf

import (
	"fmt"
)

const (
	// xzStreamIndicatorSize is the length of a stream index indicator in bytes
	xzStreamIndicatorSize  = uint64(1)
	xzStreamIndicatorMagic = 0x0
)

// xzStreamIndex as described in Section 4 the spec
type xzStreamIndex struct{}

// parseIndexAt attempts to parse an xzStreamIndex starting at the given offset.
func (r *ReaderAt) parseIndexAt(indexOffset int64) (xzStreamIndex, error) {
	indexIndicatorBytes, err := readExactly(r.cprsReader, indexOffset, xzStreamIndicatorSize)
	if err != nil {
		return xzStreamIndex{}, err
	}
	if indexIndicatorBytes[0] != xzStreamIndicatorMagic {
		return xzStreamIndex{}, fmt.Errorf("index indicator byte %x != %x", indexIndicatorBytes[0], xzStreamIndicatorMagic)
	}

	return xzStreamIndex{}, nil
}
