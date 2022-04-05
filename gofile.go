package basicfile

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Handle returns the underlying io.ReadWriteCloser.
// For convenience, the following are implemented:
// 	io.ReadWriteCloser
// 	io.StringWriter
// 	io.ReaderFrom
// 	io.WriterTo
// not included:
// 	io.ReaderAt
// 	io.WriterAt
type Handle interface {
	io.ReadWriteCloser
	io.StringWriter
	RWToFrom
}

// Dirty sets isDirty to true.
func (f *basicFile) Dirty() {
	f.isDirty = true
}

// OsFile returns the underlying
// open file descriptor (*os.File).
func (f *basicFile) OsFile() *os.File {
	return f.File
}

// Handle returns a file 'handle' that is
// ready for most go file-type operations,
// such as io.Reader and io.Writer.
//  io.ReadWriteCloser
//  io.StringWriter
//  io.ReaderFrom
//  io.WriterTo
func (f *basicFile) Handle() Handle {
	return f.rwc()
}

// String is a custom default string
// representation of the object.
func (f *basicFile) String() string {
	return fmt.Sprintf("%8s %15s", f.Mode(), f.Name())
}

// ModeString returns a human-readable
// representation of the file
//
// Alias of FileInfo().Mode().String()
func (f *basicFile) ModeString() string {
	return f.FileInfo().Mode().String()
}

// Name returns the base name of the file
// after processing with Abs().
func (f *basicFile) Name() string {
	return filepath.Base(f.Abs())
}

type (
	GoFile interface {
		fsFile
		osFile
		fs.DirEntry
		fs.FileInfo
		fs.FileMode

		// OsFile returns the file descriptor, *os.File.
		OsFile() *os.File

		// Handle returns a buffered io.ReadWriteCloser
		// (with io.WriteString, io.WriterTo, and io.ReaderFrom)
		Handle() Handle

		// A DirEntry is an entry read from a directory
		// (using the ReadDir function or a ReadDirFile's
		// ReadDir method).
		DirEntry() DirEntry

		// FileOps provides access to common file
		// operations of *os.File.
		FileOps() FileOps

		// FileUnix provides access to Unix file operations.
		FileUnix() FileUnix
	}

	// A File provides access to a single file.
	// The File interface is the minimum
	// implementation required of the file.
	// A file may implement additional interfaces,
	// such as ReadDirFile, ReaderAt, or Seeker,
	// to provide additional or optimized
	// functionality.
	//
	//  type File interface {
	//  	Stat() (fs.FileInfo, error)
	//  	Read([]byte) (int, error)
	//  	Close() error
	//  }
	//
	// Reference: standard library fs.go
	fsFile = fs.File

	// minimal os.File interface
	osFile interface {
		Seek(offset int64, whence int) (int64, error)
		Open() error
		Create() error
	}
)
