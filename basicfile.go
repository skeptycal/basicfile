package gofile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/skeptycal/goutil/errorlogger"
	"github.com/skeptycal/goutil/gofile"
)

var Err = errorlogger.Err

func NewFileWithErr(providedName string) (BasicFile, error) {
	return nil, gofile.ErrNotImplemented
}

// NewFile returns a new BasicFile, but no error,
// as is the custom in the standard library os
// package. Most often, if a file cannot be opened
// or created, we do not care why. It is often
// beyond the scope of the application to correct
// these types of errors. It simply means  that
// we cannot proceed.
//
// This allows for convenient inline usage with
// the standard pattern of Go nil checking.
// If errorlogger is active, any error is still
// recoreded in the log. This offloads error
// logging duties to errorlogger, or whichever
// standard library compatible logger function
// is assigned to the global logger function:
//  func Err(err error) error
//
// TLDR: Check for nil if you only want to know
// whether *any* error occurred.
// If you care about a *specific* error, use
//  NewFileWithErr() (f *os.File, err error)
// for a more os.Open()-ish way.
func NewFile(providedName string) BasicFile {
	f, err := NewFileWithErr(providedName)
	if Err(err) != nil {
		return nil
	}
	return f
}

// A BasicFile provides access to a single file as an in
// memory buffer.
//
// The BasicFile interface is the minimum implementation
// required of the file and may be extended to specific file
// types. (e.g. CSV, JSON, Esri Shapefile, config files, etc.)
//
// It may also be implemented as an abstract "file" interface
// that provides access to a single file that is too large to
// fit in memory at once.
//
// An implementation for large files should include a way to
// cache one section at a time, perhaps using a maxAlloc
// value or a mutex of file sections.
//
// Caching write requests will likely be the bottleneck and
// collecting multiple write requests and then writing the results
// of the most recent or most active areas of the file may be
// effective. However, performance profiling and some research
// into whether a database is more efficient is warranted.
//
// It could also be implemented as a way to access a database,
// API, buffer, or other storage.
//
// A file may implement additional interfaces, such as
// ReadDirFile, ReaderAt, or Seeker, to provide additional
// or optimized functionality.
//
//  type FileModer interface {
//  	String() string
//  	IsDir() bool
//  	IsRegular() bool
//  	Perm() FileMode
//  	Type() FileMode
//  }
//
// A FileInfo describes a file and is returned by Stat.
//
//  type FileInfo interface {
//      Name() string       // base name of the file
//      Size() int64        // length in bytes for regular files; system-dependent for others
//      Mode() FileMode     // file mode bits
//      ModTime() time.Time // modification time
//      IsDir() bool        // abbreviation for Mode().IsDir()
//      Sys() interface{}   // underlying data source (can return nil)
//  }
//
// Reference: standard library fs.go
// using File, FileInfo, and FileModer interfaces
//
// Minimum required to implement fs.File interface:
//  type File interface {
//     Stat() (fs.FileInfo, error)
//     Read([]byte) (int, error)
//     Close() error
//  }
//
// Implements fs.FileInfo interface:
// 	// A FileInfo describes a file and is returned by Stat.
//  type FileInfo interface {
//  	Name() string       // base name of the file
// 	Size() int64        // length in bytes for regular files; sy stem-dependent for others
//  	Mode() FileMode     // file mode bits
//  	ModTime() time.Time // modification time
//  	IsDir() bool        // abbreviation for Mode().IsDir()
//  	Sys() interface{}   // underlying data source (can return nil)
//  }
type BasicFile interface {
	// Handle returns the file handle, *os.File.
	// The minimum interface that is implemented
	// by a File is:
	GoFile
	//  Stat()
	Handle() *os.File

	// Stat returns the FileInfo portion of the
	// BasicFile interface.
	Stat() (fs.FileInfo, error)
	// io.ReadCloser

	FileModer
	FileInfo

	// Additional basic methods:
	Abs() string // absolute path of the file
	// IsRegular() bool // is a regular file?
	// String() string

}

/*
Chdir
Chmod
Chown
Close
Fd
Name
Read
ReadAt
ReadDir
ReadFrom
Readdir
Readdirnames
Seek
SetDeadline
SetReadDeadline
SetWriteDeadline
Stat
Sync
SyscallConn
Truncate
Write
WriteAt
WriteString
Chdir(). Error
Close). Error
Sync(). Error
*/

type (
	basicFile struct {
		providedName string      // original user input
		fi           os.FileInfo // cached file information
		modTime      time.Time   // used to validate cache entries
		*os.File                 // underlying file handle
		lock         bool
	}
)

////////////// Return component interfaces

// Handle returns the file handle, *os.File.
// The minimum interface that is implemented
// by a File is:
//  io.ReadCloser
//  Stat()
func (f *basicFile) Handle() *os.File {
	return f.file()
}

func (f *basicFile) file() *os.File {
	if f.File == nil {
		ff, err := os.Open(f.providedName)
		if Err(err) != nil {
			return nil
		}
		f.File = ff
	}
	return f.File
}
func (f *basicFile) Stat() (FileInfo, error) {
	if f.fi == nil {
		fi, err := os.Stat(f.Name())
		if Err(err) != nil {
			return nil, NewGoFileError("Gofile.Stat()", f.providedName, err)
		}
		f.fi = fi
	}
	return f.fi, nil
}

func (f *basicFile) FileInfo() FileInfo {
	fi, _ := f.Stat()
	return fi
}

// Flush flushes any in-memory copy of recent changes,
// closes the underlying file, and resets the file
// pointer / fileinfo to nil.
// This includes running os.File.Sync(), which commits
// the current contents of the file to stable storage.
//
// The BasicFile object remains available and the
// underlying will be reopened and used as needed.
// During the flushing and closing process, any new
// concurrent read or write operations will block and
// be unavailable.
func (f *basicFile) Flush() error {
	if f.Locked() {
		return errFileLocked
	}
	f.Lock()
	defer f.Unlock()
	err := Err(f.File.Sync())
	if err != nil {
		// TODO: retry ... could get stuck here ...
		return f.Flush()
	}
	err = f.File.Close()
	if err != nil {
		Err(err)
	}
	f.fi = nil
	f.File = nil
	f.timeStamp()
	return nil
}

func (f *basicFile) Locked() bool {
	return f.lock
}

func (f *basicFile) Lock() {
	f.lock = true
}

func (f *basicFile) Unlock() {
	f.lock = false
}

// timeStamp sets the most recent mod time in
// the basicFile struct and returns that time.
// This is separate and unrelated from the
// underlying file modTime, which is at:
//  (*basicFile).ModTime() time.Time
func (f *basicFile) timeStamp() time.Time {
	f.modTime = time.Now()
	return f.modTime
}

// Mode - returns the file mode bits
func (f *basicFile) Mode() FileMode {
	return f.FileInfo().Mode()
}

// Name - returns the base name of the file
func (f *basicFile) Name() string {
	return filepath.Base(f.Abs())
}

func (f *basicFile) Abs() string {
	s, err := filepath.Abs(f.providedName)
	if err != nil {
		return ""
	}
	return s
}

// Size - returns the length in bytes for regular files; system-dependent for others
func (f *basicFile) Size() int64 {
	return f.FileInfo().Size()
}

// ModTime - returns the modification time
func (f *basicFile) ModTime() time.Time {
	return f.FileInfo().ModTime()
}

// IsDir - returns true if the file is a directory
func (f *basicFile) IsDir() bool {
	return f.FileInfo().IsDir()
}

// Sys - returns the underlying data source (can return nil)
func (f *basicFile) Sys() interface{} {
	return f.FileInfo().Sys()
}

func (f *basicFile) IsRegular() bool {
	return f.FileInfo().Mode().IsRegular()
}

func (f *basicFile) Perm() FileMode {
	return f.FileInfo().Mode().Perm()
}

func (f *basicFile) Type() FileMode {
	return f.FileInfo().Mode().Type()
}

func (f *basicFile) String() string {
	return fmt.Sprintf("%8s %15s", f.Mode(), f.Name())
}
