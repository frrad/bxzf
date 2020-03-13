package bxzf

import (
	"fmt"
	"io"
	"log"
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
	footerOffset := int64(r.size) - int64(xzStreamFooterSize)
	footer, err := r.parseFooterAt(footerOffset)

	indexSize := footer.RealBackwardSize()
	indexOffset := r.size - int64(xzStreamFooterSize) - int64(indexSize)
	indexBytes := make([]byte, indexSize)
	bytesRead, err := r.cprsReader.ReadAt(indexBytes, indexOffset)
	if uint64(bytesRead) != indexSize {
		return fmt.Errorf("couldn't read a complete index")
	}
	if err != nil {
		return err
	}

	log.Fatalf("%s", showBytes(indexBytes[:100]))

	return nil
}

func showBytes(in []byte) string {
	ans := "[ "
	for _, x := range in {
		ans += fmt.Sprintf("%x ", x)
	}
	return ans + "]"
}
