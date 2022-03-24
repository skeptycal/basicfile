package gofile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skeptycal/goutil/gofile"
)

const (
	smallBufferSize = 64
	maxInt          = int(^uint(0) >> 1)
	minRead         = bytes.MinRead
)

// Stat returns the os.FileInfo for file if it exists.
//
// It is a convenience wrapper for os.Stat that traps
// and processes errors that may occur using the
// the ErrorLogger package.
//
// If the file does not exist, nil is returned.
// Errors are logged if Err is active.
func Stat(filename string) (os.FileInfo, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, Err(gofile.NewGoFileError("gofile.Stat()", filename, err))
	}
	return fi, nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, ErrNotExist)
}

func NotExists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, os.ErrNotExist)
}

// StatCheck returns file information (after symlink evaluation)
// using os.Stat(). If the file does not exist, is not a regular file,
// or if the user lacks adequate permissions, an error is returned.
// StatCheck returns file information (after symlink evaluation
// and path cleaning) using os.Stat().
//
// If the file does not exist, is not a regular file,
// or if the user lacks adequate permissions, an error is
// returned.
//
// It is a convenience wrapper for os.Stat that traps
// and processes errors that may occur using the
// the ErrorLogger package.
//
// If the file does not exist, nil is returned.
// Errors are logged if Err is active.
func StatCheck(filename string) (os.FileInfo, error) {

	// EvalSymlinks also calls Abs and Clean as well as
	// checking for existance of the file.
	filename, err := filepath.EvalSymlinks(filename)
	if err != nil {
		return nil, Err(err)
	}

	fi, err := os.Stat(filename)
	if err != nil {
		return nil, Err(err)
	}

	//Check 'others' permission
	m := fi.Mode()
	if m&(1<<2) == 0 {
		return nil, fmt.Errorf("insufficient permissions: %v", filename)
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("the filename %s refers to a directory", filename)
	}

	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("the filename %s is not a regular file", filename)
	}

	return fi, err
}

// Mode returns the filemode of file.
func Mode(file string) os.FileMode {
	fi, err := Stat(file)
	if err != nil {
		Err(err)
		return 0
	}
	return fi.Mode()
}

func NewBasicFile(filename string) (BasicFile, error) {

	return &basicFile{providedName: filename}, nil
}

// Open opens the named file for reading as an in memory object.
// If successful, methods on the returned file can be used for
// reading; the associated file descriptor has mode O_RDONLY.
// If there is an error, it will be of type *os.PathError.
func Open(name string) (BasicFile, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, NewGoFileError("gofile.Open", name, err)
	}

	return NewBasicFile(f.Name())
}

// Create creates or truncates the named file and returns an
// opened file as io.ReadWriteCloser.
//
// If the file already exists, it is truncated. If the file
// does not exist, it is created with mode 0644 (before umask).
// If successful, methods on the returned File can be used
// for I/O; the associated file descriptor has mode O_RDWR.
//
// If there is an error, it will be of type *PathError.
func Create(name string) (io.ReadWriteCloser, error) {

	// standard library: OpenFile is the generalized open call; most users
	// will use Open or Create instead. It opens the named file with specified
	// flag (O_RDONLY etc.). If the file does not exist, and the O_CREATE flag
	// is passed, it is created with mode perm (before umask). If successful,
	// methods on the returned File can be used for I/O. If there is an error,
	// it will be of type *PathError.
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, NormalMode)
	if err != nil {
		return nil, &PathError{Op: "gofile.Create", Path: name, Err: err}
	}

	return f, nil
}

// CreateSafe creates the named file and returns an
// opened file as io.ReadWriteCloser.
//
// If the file already exists, an error is returned. If the file
// does not exist, it is created with mode 0644 (before umask).
// If successful, methods on the returned File can be used
// for I/O; the associated file descriptor has mode O_RDWR.
//
// If there is an error, it will be of type *PathError.
func CreateSafe(name string) (io.ReadWriteCloser, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, NormalMode)
	if err != nil {
		return nil, &PathError{Op: "gofile.CreateSafe", Path: name, Err: err}
	}
	return f, nil
}
