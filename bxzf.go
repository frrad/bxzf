package bxzf

import (
	"fmt"
	"io"
	"os"
)

type ReaderAt struct {
	cprsReader io.ReaderAt
	size       int64
}

// Open takes a io.ReaderAt for a xz compressed file, along with the file's
// length and returns a bxzf.ReaderAt for its contents.
func Open(compressedReader io.ReaderAt, size int64) (ReaderAt, error) {
	opened := ReaderAt{
		cprsReader: compressedReader,
		size:       size,
	}

	err := opened.init()
	return opened, err
}

// OpenFile wraps Open for the common case of opening a compressed file.
func OpenFile(filename string) (ReaderAt, error) {
	f, err := os.Open(filename)
	if err != nil {
		return ReaderAt{}, err
	}

	info, err := f.Stat()
	if err != nil {
		return ReaderAt{}, err
	}

	return Open(f, info.Size())
}

// init builds the index to keep track of offsets for separately compressed
// blocks within the compressed file
func (r *ReaderAt) init() error {
	footerOffset := r.size - int64(xzStreamFooterSize)
	footer, err := r.parseFooterAt(footerOffset)
	if err != nil {
		return err
	}

	indexOffset := footerOffset - int64(footer.RealBackwardSize())
	_, err = r.parseIndexAt(indexOffset)
	if err != nil {
		return err
	}

	//	fmt.Println(index)

	return nil
}

func showBytes(in []byte) string {
	ans := "[ "
	for _, x := range in {
		ans += fmt.Sprintf("%x ", x)
	}
	return ans + "]"
}

func readExactly(r io.ReaderAt, off int64, nBytes uint64) ([]byte, error) {
	outBytes := make([]byte, nBytes)
	n, err := r.ReadAt(outBytes, off)
	if n != int(nBytes) {
		return outBytes, fmt.Errorf("read %d != %d bytes", n, nBytes)
	}
	return outBytes, err
}

func readExactlyOne(r io.ReaderAt, off int64) (byte, error) {
	bs, err := readExactly(r, off, 1)
	if err != nil {
		return 0x0, err
	}
	return bs[0], nil
}

// parseNumberAt attempts to parse a variable-length multibyte integer as
// described in 1.2 of the spec at the given offset. It returns the number, and
// the length of its encoded representation in bytes.
func (r *ReaderAt) parseNumberAt(offset int64) (uint64, uint32, error) {

	sizeMax := uint32(9)

	currentByte, err := readExactlyOne(r.cprsReader, offset)
	if err != nil {
		return 0, 0, err
	}

	num := uint64(currentByte & 0x7f)

	i := uint32(1)

	for currentByte&0x80 != 0 {

		currentByte, err = readExactlyOne(r.cprsReader, offset+int64(i))
		if err != nil {
			return 0, 0, err
		}

		if i >= sizeMax || currentByte == 0x00 {
			return num, 0, nil
		}

		num |= uint64(currentByte&0x7F) << (i * 7)

		i++
	}

	return num, i, nil
}

func roundUpFour(x int64) int64 {
	lastTwoBits := x & 3

	if lastTwoBits != 0 {
		return x - lastTwoBits + 4
	}
	return x
}

func nullBytes(in []byte) error {
	for _, b := range in {
		if b != 0x0 {
			return fmt.Errorf("expected null byte")
		}
	}
	return nil
}
