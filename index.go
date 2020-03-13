package bxzf

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
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
	CRC32      [4]byte
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

	offsetInIndex := int64(xzStreamIndicatorSize)
	numRecords, bytesInNumRep, err := r.parseNumberAt(offset + offsetInIndex)
	if err != nil {
		return ix, err
	}

	offsetInIndex += int64(bytesInNumRep)

	ix.NumRecords = numRecords
	for i := uint64(0); i < numRecords; i++ {
		unpaddedSize, offsetDiff, err := r.parseNumberAt(offset + offsetInIndex)
		if err != nil {
			return ix, err
		}
		offsetInIndex += int64(offsetDiff)

		uncompressedSize, offsetDiff, err := r.parseNumberAt(offset + offsetInIndex)
		if err != nil {
			return ix, err
		}
		offsetInIndex += int64(offsetDiff)

		ix.Records = append(ix.Records, xzStreamIndexRecord{
			UncompressedSize: uncompressedSize,
			UnpaddedSize:     unpaddedSize,
		})
	}

	offsetWithPadding := roundUpFour(offsetInIndex)
	if offsetInIndex != offsetWithPadding {
		paddingLen := uint64(offsetWithPadding - offsetInIndex)
		paddingBytes, err := readExactly(r.cprsReader, offset+offsetInIndex, paddingLen)
		if err != nil {
			return ix, err
		}
		err = nullBytes(paddingBytes)
		if err != nil {
			return ix, err
		}
	}

	checksum32, err := readExactly(r.cprsReader, offset+offsetWithPadding, 4)
	if err != nil {
		return ix, err
	}

	copy(ix.CRC32[:], checksum32)

	allBytes, err := readExactly(r.cprsReader, offset, uint64(offsetWithPadding))
	if err != nil {
		return ix, err
	}

	if crc32.ChecksumIEEE(allBytes) != ix.CRC32Uint() {
		return ix, fmt.Errorf("checksum fail")
	}

	return ix, nil
}

func (f xzStreamIndex) CRC32Uint() uint32 {
	return binary.LittleEndian.Uint32(f.CRC32[:])
}
