package basicfile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// A FileInfo describes a file and is returned by Stat.
//
//	type FileInfo interface {
// 		Name() string       // base name of the file
// 		Size() int64        // length in bytes for regular files; system-dependent for others
// 		Mode() FileMode     // file mode bits
// 		ModTime() time.Time // modification time
// 		IsDir() bool        // abbreviation for Mode().IsDir()
// 		Sys() interface{}   // underlying data source (can return nil)
//	}
//
// Reference: standard library fs.go
type FileInfo = fs.FileInfo

// FileInfo returns the same information as
// Stat() but without the error. An error is
// implied if the return value is nil.
//
// Errors are logged if Err is active.
func (f *basicFile) FileInfo() fs.FileInfo {
	fi, err := f.Stat()
	if Err(err) != nil {
		return nil
	}
	return fi
}

// Stat returns a FileInfo describing the named file.
// If the FileInfo is cached, it will be returned
// directly unless IsDirty is true.
//
// If there is an error, it will be of type *PathError.
//
// It is a convenience wrapper for os.Stat() that traps
// and processes errors that may occur using the
// the ErrorLogger package.
//
// Errors are logged if Err is active.
func (f *basicFile) Stat() (fs.FileInfo, error) {
	if f.fi == nil || f.isDirty {
		fi, err := Stat(f.Name())
		if Err(err) != nil {
			return nil, err
		}
		f.fi = fi
	}
	return f.fi, nil
}

// Info is an alias of Stat() that satisfies
// the fs.DirEntry interface.
func (f *basicFile) Info() (fs.FileInfo, error) {
	return f.Stat()
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
//
// It is a convenience wrapper for os.Stat that traps
// and processes errors that may occur using the
// the ErrorLogger package.
//
// Errors are logged if Err is active.
func Stat(filename string) (os.FileInfo, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, Err(NewGoFileError("gofile.Stat()", filename, err))
	}
	return fi, nil
}

// RegularFileInfo returns file information
// (after symlink evaluation and path cleaning)
// using os.Stat().
//
// If the file does not exist, is a directory,
// is not a regular file, or if the user lacks
// adequate permissions, an error is logged and
// nil is returned.
//
// It is a convenience wrapper for os.Stat that traps
// and processes errors that may occur using the
// the ErrorLogger package.
//
// Errors are logged if Err is active.
func RegularFileInfo(filename string) os.FileInfo {

	// EvalSymlinks also calls Abs and Clean as well as
	// checking for existance of the file.
	filename, err := filepath.EvalSymlinks(filename)
	if err != nil {
		Err(err)
		return nil
	}

	fi, err := Stat(filename)
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
