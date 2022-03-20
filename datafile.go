package basicfile

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrBadCount   = errors.New("datafile: bad read count")
	ErrNotRegular = errors.New("datafile: not regular file")
)

type DataFile interface {
	BF
	Data() ([]byte, error)
}

// datafile is a file type that is specialized for binary data
type Datafile struct {
	Basicfile
	data []byte
}

func NewDataFile(filename string) (DataFile, error) {
	src, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	if !src.Mode().IsRegular() {
		return nil, ErrNotRegular
	}

	name, err := filepath.Abs(src.Name())
	if err != nil {
		return nil, &os.PathError{Op: "abs", Path: src.Name(), Err: err}
	}

	df := Datafile{}

	df.ProvidedName = filename
	df.name = name
	df.size = src.Size()
	// df.info = src
	df.FileInfo = src

	return &df, nil
}
