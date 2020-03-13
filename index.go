package bxzf

import (
	"fmt"
	"log"
)

const (
	// xzStreamIndicatorSize is the length of a stream index indicator in bytes
	xzStreamIndicatorSize  = uint64(1)
	xzStreamIndicatorMagic = 0x0
)

// xzStreamIndex as described in Section 4 the spec
type xzStreamIndex struct {
	NumRecords uint64
	Records    []xzStreamIndexRecord
}

type xzStreamIndexRecord struct {
	UnpaddedSize     uint64
	UncompressedSize uint64
}

// parseIndexAt attempts to parse an xzStreamIndex starting at the given offset.
func (r *ReaderAt) parseIndexAt(offset int64) (xzStreamIndex, error) {
	ix := xzStreamIndex{}

	indexIndicatorBytes, err := readExactly(r.cprsReader, offset, xzStreamIndicatorSize)
	if err != nil {
		return ix, err
	}
	if indexIndicatorBytes[0] != xzStreamIndicatorMagic {
		return ix, fmt.Errorf("index indicator byte %x != %x", indexIndicatorBytes[0], xzStreamIndicatorMagic)
	}

	offset += int64(xzStreamIndicatorSize)
	numRecords, bytesInNumRep, err := r.parseNumberAt(offset)
	if err != nil {
		return ix, err
	}

	offset += int64(bytesInNumRep)

	ix.NumRecords = numRecords
	for i := uint64(0); i < numRecords; i++ {
		unpaddedSize, offsetDiff, err := r.parseNumberAt(offset)
		if err != nil {
			return ix, err
		}
		offset += int64(offsetDiff)

		uncompressedSize, offsetDiff, err := r.parseNumberAt(offset)
		if err != nil {
			return ix, err
		}
		offset += int64(offsetDiff)

		ix.Records = append(ix.Records, xzStreamIndexRecord{
			UncompressedSize: uncompressedSize,
			UnpaddedSize:     unpaddedSize,
		})
	}

	log.Println(len(ix.Records[:3]))

	return ix, nil
}
