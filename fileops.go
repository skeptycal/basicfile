package basicfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return nil, Err(NewGoFileError("gofile.Stat()", filename, err))
	}
	return fi, nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, os.ErrNotExist)
}

func NotExists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, os.ErrNotExist)
}

// FileInfo returns file information (after symlink evaluation
// and path cleaning) using os.Stat().
//
// If the file does not exist, is not a regular file,
// or if the user lacks adequate permissions, an error is logged and nil is returned.
//
// It is a convenience wrapper for os.Stat that traps
// and processes errors that may occur using the
// the ErrorLogger package.
//
// If the file does not exist, nil is returned.
// Errors are logged if Err is active.
func FileInfo(filename string) os.FileInfo {

	// EvalSymlinks also calls Abs and Clean as well as
	// checking for existance of the file.
	filename, err := filepath.EvalSymlinks(filename)
	if err != nil {
		Err(err)
		return nil
	}

	fi, err := os.Stat(filename)
	if err != nil {
		Err(err)
		return nil
	}

	//Check 'others' permission
	m := fi.Mode()
	if m&(1<<2) == 0 {
		Err(fmt.Errorf("insufficient permissions: %v", filename))
		return nil
	}

	if fi.IsDir() {
		Err(fmt.Errorf("the filename %s refers to a directory", filename))
		return nil
	}

	if !fi.Mode().IsRegular() {
		Err(fmt.Errorf("the filename %s is not a regular file", filename))
		return nil
	}

	return fi
}

// Mode returns the filemode of file.
// If an error occurs, it is logged
// and 0 is returned.
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
	if Err(err) != nil {
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
func Create(name string) (BasicFile, error) {
	b := &basicFile{providedName: name}

	err := b.Create()
	if err != nil {
		return nil, err
	}

	defer b.timeStamp()

	// standard library: OpenFile is the generalized open call; most users
	// will use Open or Create instead. It opens the named file with specified
	// flag (O_RDONLY etc.). If the file does not exist, and the O_CREATE flag
	// is passed, it is created with mode perm (before umask). If successful,
	// methods on the returned File can be used for I/O. If there is an error,
	// it will be of type *PathError.
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, NormalMode)
	if err != nil {
		return nil, Err(NewGoFileError("gofile.Create", name, err))
	}

	b.File = f

	return b, nil
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
		return nil, NewGoFileError("gofile.CreateSafe", name, err)
	}
	return f, nil
}
