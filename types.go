package basicfile

import (
	"io/fs"
)

// Reference types copied from io package
type (
	Closer interface {
		Close() error
	}

	GoDir interface {
		// GoFile
		fsFile

		// Readdir(count int) ([]os.FileInfo, error)
		ReadDir(n int) ([]fs.DirEntry, error)

		Chdir() error

		Readdirnames(dir string) (n int, err error)
	}
)

//******************* Reference: standard library fs.go

type (

	// A ReadDirFile is a directory file whose entries can be read with the ReadDir method.
	// Every directory file should implement this interface.
	// (It is permissible for any file to implement this interface,
	// but if so ReadDir should return an error for non-directories.)
	//   type ReadDirFile interface {
	//   	File
	// 		// ReadDir reads the contents of the directory and returns
	//		// a slice of up to n DirEntry values in directory order.
	//		// Subsequent calls on the same file will yield further DirEntry values.
	//		//
	//		// If n > 0, ReadDir returns at most n DirEntry structures.
	//		// In this case, if ReadDir returns an empty slice, it will return
	//		// a non-nil error explaining why.
	//		// At the end of a directory, the error is io.EOF.
	//		//
	//		// If n <= 0, ReadDir returns all the DirEntry values from the directory
	//		// in a single slice. In this case, if ReadDir succeeds (reads all the way
	//		// to the end of the directory), it returns the slice and a nil error.
	//		// If it encounters an error before the end of the directory,
	//		// ReadDir returns the DirEntry list read until that point and a non-nil error.
	//   	ReadDir(n int) ([]DirEntry, error)
	//   }
	//
	// Reference: standard library fs.go
	ReadDirFile = fs.ReadDirFile

	// An FS provides access to a hierarchical file system.
	//
	// The FS interface is the minimum implementation required of the file system.
	// A file system may implement additional interfaces,
	// such as ReadFileFS, to provide additional or optimized functionality.
	//
	//  type FS interface {
	//  	// Open opens the named file.
	//  	//
	//  	// When Open returns an error, it should be of type *PathError
	//  	// with the Op field set to "open", the Path field set to name,
	//  	// and the Err field describing the problem.
	//  	//
	//  	// Open should reject attempts to open names that do not satisfy
	//  	// ValidPath(name), returning a *PathError with Err set to
	//  	// ErrInvalid or ErrNotExist.
	//  	Open(name string) (File, error)
	//  }
	//
	// Reference: standard library fs.go
	FS = fs.FS
)
