package basicfile

import (
	"io"
	"io/fs"
	"os"
	"syscall"
	"time"
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

		fs.File

		io.Writer
		io.StringWriter

		// ReaderFrom is the interface that wraps
		// the ReadFrom method.
		//
		// ReadFrom reads data from r until EOF or error.
		// The return value n is the number of bytes read.
		// Any error except EOF encountered during the
		// read is also returned.
		//
		// The Copy function uses ReaderFrom if available.
		io.ReaderFrom

		// WriterTo is the interface that wraps the
		// WriteTo method:
		//     WriteTo(w Writer) (n int64, err error)
		//
		// WriteTo writes data to w until there's no
		// more data to write or when an error occurs.
		// The return value n is the number of bytes
		// written. Any error encountered during the
		// write is also returned.
		//
		// The Copy function uses WriterTo if available.
		//
		io.WriterTo

		// ReaderAt is the interface that wraps the basic ReadAt method.
		// 	ReadAt(p []byte, off int64) (n int, err error)
		//
		// ReadAt reads len(p) bytes into p starting at offset off in the underlying input source. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
		//
		// When ReadAt returns n < len(p), it returns a non-nil error explaining why more bytes were not returned. In this respect, ReadAt is stricter than Read.
		//
		// Even if ReadAt returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, ReadAt blocks until either all the data is available or an error occurs. In this respect ReadAt is different from Read.
		//
		// If the n = len(p) bytes returned by ReadAt are at the end of the input source, ReadAt may return either err == EOF or err == nil.
		//
		// If ReadAt is reading from an input source with a seek offset, ReadAt should not affect nor be affected by the underlying seek offset.
		//
		// Clients of ReadAt can execute parallel ReadAt calls on the same input source.
		//
		// Implementations must not retain p.
		io.ReaderAt

		// WriterAt is the interface that wraps the basic WriteAt method.
		// 	WriteAt(p []byte, off int64) (n int, err error)
		//
		// WriteAt writes len(p) bytes from p to the underlying data stream at offset off. It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early. WriteAt must return a non-nil error if it returns n < len(p).
		//
		// If WriteAt is writing to a destination with a seek offset, WriteAt should not affect nor be affected by the underlying seek offset.
		io.WriterAt

		// FileInfo methods
		// Name() string       // base name of the file
		// Size() int64        // length in bytes for regular files; system-dependent for others
		// Mode() FileMode     // file mode bits
		// ModTime() time.Time // modification time
		// IsDir() bool        // abbreviation for Mode().IsDir()
		// Sys() interface{}   // underlying data source (can return nil)
		fs.FileInfo

		// FileMode methods
		// String() string 	// human-readable representation of the file
		// IsDir() bool 	// abbreviation for Mode().IsDir()
		// IsRegular() bool // IsRegular reports whether m is a regular file.
		// Perm() FileMode	// Perm returns the Unix permission bits
		// Type() FileMode
		FileModer

		// FileOps methods
		// Abs() (string, error)
		// Base(path string) string
		// Chmod(mode os.FileMode) error
		// Dir(path string) string
		// Ext(path string) string
		// Move(path string) error
		// Split(path string) (dir, file string)
		FileOps

		// Unix File Operations
		// 	Fd() uintptr
		// 	Link(newname string) error
		// 	Readlink() (string, error)
		// 	Remove() error
		// 	Symlink(newname string) error
		// 	Truncate(size int64) error
		FileUnix
	}

	GoDir interface {
		// GoFile
		File

		// Readdir(count int) ([]os.FileInfo, error)
		ReadDir(n int) ([]fs.DirEntry, error)

		FileOps
		Chdir() error

		Readdirnames(dir string) (n int, err error)
	}

	FileOps interface {
		Abs() string
		Base() string
		Dir() string
		Ext() string
		Split() (dir, file string)

		Chmod(mode os.FileMode) error
		Chown(uid int, gid int) error
		Move(newpath string) error
		Sync() error

		SetDeadline(t time.Time) error
		SetReadDeadline(t time.Time) error
		SetWriteDeadline(t time.Time) error

		SyscallConn() (syscall.RawConn, error)
	}

	FileUnix interface {
		Fd() uintptr
		Link(newname string) error
		Readlink() (string, error)
		Remove() error
		Symlink(newname string) error
		Truncate(size int64) error
	}
)

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
