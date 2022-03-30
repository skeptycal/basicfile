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

// SameFile reports whether fi1 and fi2 describe the same file.
// For example, on Unix this means that the device and inode fields
// of the two underlying structures are identical; on other systems
// the decision may be based on the path names.
// SameFile only applies to results returned by this package gofile
// It returns false in other cases.
var SameFile = os.SameFile
