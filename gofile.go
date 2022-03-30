package basicfile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func (f *basicFile) Handle() *os.File {
	return f.File
}

func (f *basicFile) Io() Handle {
	return f.rwc()
}

func (f *basicFile) DirEntry() fs.DirEntry {
	return fs.FileInfoToDirEntry(f.FileInfo())
}

// Name - returns the base name of the file
func (f *basicFile) Name() string {
	return filepath.Base(f.Abs())
}

// IsDir - returns true if the file is a directory
func (f *basicFile) IsDir() bool {
	return f.FileInfo().IsDir()
}

// Type returns type bits in m (m & ModeType).
func (f *basicFile) Type() FileMode { return f.Mode().Type() }

// Info is an alias of Stat() that satisfies
// the fs.DirEntry interface.
func (f *basicFile) Info() (fs.FileInfo, error) {
	return f.Stat()
}

// Size - returns the length in bytes for
// regular files; system-dependent for others.
// Alias of FileInfo().Size()
func (f *basicFile) Size() int64 {
	return f.FileInfo().Size()
}

// ModTime - returns the modification time.
// Alias of FileInfo().ModTime()
func (f *basicFile) ModTime() time.Time {
	return f.FileInfo().ModTime()
}

// Sys - returns the underlying data source
// (can return nil).
// Alias of FileInfo().Sys()
func (f *basicFile) Sys() interface{} {
	return f.FileInfo().Sys()
}

// Mode - returns the file mode bits
func (f *basicFile) Mode() fs.FileMode {
	return f.FileInfo().Mode()
}

func (f *basicFile) String() string {
	return fmt.Sprintf("%8s %15s", f.Mode(), f.Name())
}

// human-readable representation of the file
func (f *basicFile) ModeString() string {
	return f.FileInfo().Mode().String()
}

// IsRegular reports whether the file a regular file
func (f *basicFile) IsRegular() bool {
	return f.FileInfo().Mode().IsRegular()
}

// Perm returns the Unix permission bits
func (f *basicFile) Perm() FileMode {
	return f.FileInfo().Mode().Perm()
}

type (
	GoFile interface {

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
		fsFile

		Seek(offset int64, whence int) (int64, error)
		Open() error
		Create() error

		// Io returns the underlying io.ReadWriteCloser.
		// For convenience, the following are implemented:
		// 	io.ReadWriteCloser
		// 	io.StringWriter
		// 	io.ReaderFrom
		// 	io.WriterTo
		// not included:
		// 	io.ReaderAt
		// 	io.WriterAt
		Io() Handle

		// DirEntry returns the file's directory entry.
		// A DirEntry is an entry read from a directory
		// (using the ReadDir function or a ReadDirFile's
		// ReadDir method).
		// 	type DirEntry interface {
		//		Name() string
		//		IsDir() bool
		//		Type() FileMode
		//		Info() (FileInfo, error)
		//	}
		DirEntry() fs.DirEntry
		Name() string
		IsDir() bool
		Type() FileMode
		Info() (FileInfo, error)

		// FileInfo interface ... also implements fs.DirEntry
		// 	Name() string       // base name of the file
		// 	Size() int64        // length in bytes for regular files; system-dependent for others
		// 	Mode() FileMode     // file mode bits
		// 	ModTime() time.Time // modification time
		// 	IsDir() bool        // abbreviation for Mode().IsDir()
		// 	Sys() interface{}   // underlying data source (can return nil)
		Size() int64        // Info().Size()
		ModTime() time.Time // Info().ModTime()
		Sys() interface{}   // Info().Sys()

		// Mode returns the type bits for the entry
		// and implements fs.DirEntry
		// 	String() string 	// human-readable representation of the file
		// 	IsDir() bool 		// abbreviation for Mode().IsDir()
		// 	IsRegular() bool 	// IsRegular reports whether m is a regular file.
		// 	Perm() FileMode		// Perm returns the Unix permission bits
		// 	Type() FileMode
		Mode() FileMode  // Info().Mode()
		String() string  // human-readable representation of the file
		IsRegular() bool // IsRegular reports whether the file a regular file
		Perm() FileMode  // Perm returns the Unix permission bits

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
)
