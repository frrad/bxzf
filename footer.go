package bxzf

import (
	"encoding/binary"
	"fmt"
)

const (
	xzStreamFooterSize = uint64(12) // length of a stream footer in bytes
	xzFooterMagicOne   = 0x59
	xzFooterMagicTwo   = 0x5a
)

// xzStreamFooter as described in 2.1.2 of the spec
type xzStreamFooter struct {
	CRC32              [4]byte
	StoredBackwardSize [4]byte
	StreamFlags        [2]byte
}

// RealBackwardSize gives the backwards size of the index in bytes. Computed as
// described in 2.1.2.2 of the spec.
func (f xzStreamFooter) RealBackwardSize() uint64 {
	backwardsSize := binary.LittleEndian.Uint32(f.StoredBackwardSize[:])
	return 4 * (1 + uint64(backwardsSize))
}

// Check checks the CRC32 checksum. See 2.1.2.1 of the spec.
func (f xzStreamFooter) Check() error {
	// todo implement
	return nil
}

// parseFooterAt attempts to parse an xzStreamFooter starting at the given
// offset.
func (r *ReaderAt) parseFooterAt(footerOffset int64) (xzStreamFooter, error) {
	footerBytes := make([]byte, xzStreamFooterSize)
	bytesRead, err := r.cprsReader.ReadAt(footerBytes, footerOffset)
	if uint64(bytesRead) != xzStreamFooterSize {
		return xzStreamFooter{}, fmt.Errorf("couldn't read a complete footer")
	}
	if err != nil {
		return xzStreamFooter{}, err
	}

	if footerBytes[10] != xzFooterMagicOne {
		return xzStreamFooter{}, fmt.Errorf("stream footer[%d] == %x != %x", 10, footerBytes[10], xzFooterMagicOne)
	}
	if footerBytes[11] != xzFooterMagicTwo {
		return xzStreamFooter{}, fmt.Errorf("stream footer[%d] == %x != %x", 11, footerBytes[11], xzFooterMagicTwo)
	}

	footer := xzStreamFooter{}
	copy(footer.CRC32[:], footerBytes[:4])
	copy(footer.StoredBackwardSize[:], footerBytes[4:8])
	copy(footer.StreamFlags[:], footerBytes[8:10])

	return footer, footer.Check()
}
