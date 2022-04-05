package basicfile

import (
	"bytes"
	"errors"
	"os"
)

const (
	NormalMode      os.FileMode = 0644
	DirMode         os.FileMode = 0755
	smallBufferSize             = 64
	maxInt                      = int(^uint(0) >> 1)
	minRead                     = bytes.MinRead
)

func (f *basicFile) FileOps() FileOps {
	return f
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, os.ErrNotExist)
}

func NotExists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, os.ErrNotExist)
}
