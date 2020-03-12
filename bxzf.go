package bxzf

import (
	"io"
	"os"
)

type ReaderAt struct{}

// Open takes a io.ReaderAt for a xz compressed file, along with the file's
// length and returns a bxzf.ReaderAt for its contents.
func Open(compressedReader io.ReaderAt, size int64) (ReaderAt, error) {

	return ReaderAt{}, nil

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
