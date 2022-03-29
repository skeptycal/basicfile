package basicfile

import (
	"io/fs"
	"os"
)

// Reference types copied from io package
type (
	Closer interface {
		Close() error
	}

	GoFile interface {
		Seek(offset int64, whence int) (int64, error)
		Open() error
		Create() error

		// io.ReadWriteCloser
		// io.StringWriter
		// io.ReaderFrom
		// io.WriterTo
		// io.ReaderAt
		// io.WriterAt
		Io() Handle

		// FileInfo methods
		// Name() string       // base name of the file
		// Size() int64        // length in bytes for regular files; system-dependent for others
		// Mode() FileMode     // file mode bits
		// ModTime() time.Time // modification time
		// IsDir() bool        // abbreviation for Mode().IsDir()
		// Sys() interface{}   // underlying data source (can return nil)
		FileInfo() fs.FileInfo

		// FileMode methods
		// String() string 	// human-readable representation of the file
		// IsDir() bool 	// abbreviation for Mode().IsDir()
		// IsRegular() bool // IsRegular reports whether m is a regular file.
		// Perm() FileMode	// Perm returns the Unix permission bits
		// Type() FileMode
		FileMode() FileModer

		// FileOps methods
		// Abs() (string, error)
		// Base(path string) string
		// Chmod(mode os.FileMode) error
		// Dir(path string) string
		// Ext(path string) string
		// Move(path string) error
		// Split(path string) (dir, file string)
		FileOps() FileOps

		// Unix File Operations
		// 	Fd() uintptr
		// 	Link(newname string) error
		// 	Readlink() (string, error)
		// 	Remove() error
		// 	Symlink(newname string) error
		// 	Truncate(size int64) error
		FileUnix() FileUnix
	}

	GoDir interface {
		// GoFile
		File

		// Readdir(count int) ([]os.FileInfo, error)
		ReadDir(n int) ([]fs.DirEntry, error)

		Chdir() error

		Readdirnames(dir string) (n int, err error)
	}
)

func (f *basicFile) Handle() Handle {}

//******************* Reference: standard library fs.go

type (

	// A File provides access to a single file.
	// The File interface is the minimum implementation required
	// of the file.
	// A file may implement additional interfaces, such as
	// ReadDirFile, ReaderAt, or Seeker, to provide additional
	// or optimized functionality.
	//
	//  type File interface {
	//  	Stat() (fs.FileInfo, error)
	//  	Read([]byte) (int, error)
	//  	Close() error
	//  }
	//
	// Reference: standard library fs.go
	File = fs.File

	// A FileInfo describes a file and is returned by Stat.
	//
	// Implements fs.FileInfo interface:
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
	// FileInfo = fs.FileInfo

	// A FileMode represents a file's mode and permission bits.
	// The bits have the same definition on all systems, so that
	// information about files can be moved from one system
	// to another portably. Not all bits apply to all systems.
	// The only required bit is ModeDir for directories.
	//
	//  type FileMode uint32
	//
	// Reference: standard library fs.go
	FileMode = fs.FileMode

	// FileModer implements fs.FileMode methods
	//
	// A FileMode represents a file's mode and permission bits.
	// The bits have the same definition on all systems, so that
	// information about files can be moved from one system
	// to another portably. Not all bits apply to all systems.
	// The only required bit is ModeDir for directories.
	FileModer interface {

		// human-readable representation of the file
		String() string

		// IsDir reports whether m describes a directory.
		// That is, it tests for the ModeDir bit being set in m.
		// IsDir() bool // duplicated in FileInfo interface

		// IsRegular reports whether m describes a regular file.
		// That is, it tests that no mode type bits are set.
		IsRegular() bool

		// Perm returns the Unix permission bits in m (m & ModePerm).
		Perm() FileMode

		// Type returns type bits in m (m & ModeType).
		Type() FileMode
	}

	// A DirEntry is an entry read from a directory
	// (using the ReadDir function or a ReadDirFile's ReadDir method).
	//
	//  type DirEntry interface {
	//  	// Name returns the name of the file (or subdirectory) described by the entry.
	//  	// This name is only the final element of the path (the base name), not the entire path.
	//  	// For example, Name would return "hello.go" not "/home/gopher/hello.go".
	//  	Name() string
	//
	//   	// IsDir reports whether the entry describes a directory.
	//  	IsDir() bool
	//
	//   	// Type returns the type bits for the entry.
	//  	// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
	//  	Type() FileMode
	//
	//   	// Info returns the FileInfo for the file or subdirectory described by the entry.
	//  	// The returned FileInfo may be from the time of the original directory read
	//  	// or from the time of the call to Info. If the file has been removed or renamed
	//  	// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
	//  	// If the entry denotes a symbolic link, Info reports the information about the link itself,
	//  	// not the link's target.
	//  	Info() (FileInfo, error)
	//  }
	//
	// Reference: standard library fs.go
	// DirEntry = fs.DirEntry

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
	// ReadDirFile = fs.ReadDirFile

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
	// FS = fs.FS
)

// SameFile reports whether fi1 and fi2 describe the same file.
// For example, on Unix this means that the device and inode fields
// of the two underlying structures are identical; on other systems
// the decision may be based on the path names.
// SameFile only applies to results returned by this package gofile
// It returns false in other cases.
var SameFile = os.SameFile
