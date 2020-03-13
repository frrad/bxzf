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
	index, err := r.parseIndexAt(indexOffset)
	if err != nil {
		return err
	}

	fmt.Println(index)

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
