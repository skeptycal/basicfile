package basicfile

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	NormalMode os.FileMode = 0644
	DirMode    os.FileMode = 0755
)

func NewFileWithErr(providedName string) (BasicFile, error) {
	return nil, ErrNotImplemented
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
	Handle() *os.File

	// The minimum interface that is implemented
	// by a File is:
	Seek(offset int64, whence int) (int64, error)
	Open() error
	Create() error

	// GoFile interface {
	// 	Seek(offset int64, whence int) (int64, error)
	// 	Open() error
	// 	Create() error
	//
	// 	fs.File
	//
	// 	io.Writer
	// 	io.StringWriter
	//
	// 	// ReaderFrom is the interface that wraps
	// 	// the ReadFrom method.
	// 	//
	// 	// ReadFrom reads data from r until EOF or error.
	// 	// The return value n is the number of bytes read.
	// 	// Any error except EOF encountered during the
	// 	// read is also returned.
	// 	//
	// 	// The Copy function uses ReaderFrom if available.
	// 	io.ReaderFrom
	//
	// 	// WriterTo is the interface that wraps the
	// 	// WriteTo method:
	// 	//     WriteTo(w Writer) (n int64, err error)
	// 	//
	// 	// WriteTo writes data to w until there's no
	// 	// more data to write or when an error occurs.
	// 	// The return value n is the number of bytes
	// 	// written. Any error encountered during the
	// 	// write is also returned.
	// 	//
	// 	// The Copy function uses WriterTo if available.
	// 	//
	// 	io.WriterTo
	//
	// 	// ReaderAt is the interface that wraps the basic ReadAt method.
	// 	// 	ReadAt(p []byte, off int64) (n int, err error)
	// 	//
	// 	// ReadAt reads len(p) bytes into p starting at offset off in the underlying input source. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
	// 	//
	// 	// When ReadAt returns n < len(p), it returns a non-nil error explaining why more bytes were not returned. In this respect, ReadAt is stricter than Read.
	// 	//
	// 	// Even if ReadAt returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, ReadAt blocks until either all the data is available or an error occurs. In this respect ReadAt is different from Read.
	// 	//
	// 	// If the n = len(p) bytes returned by ReadAt are at the end of the input source, ReadAt may return either err == EOF or err == nil.
	// 	//
	// 	// If ReadAt is reading from an input source with a seek offset, ReadAt should not affect nor be affected by the underlying seek offset.
	// 	//
	// 	// Clients of ReadAt can execute parallel ReadAt calls on the same input source.
	// 	//
	// 	// Implementations must not retain p.
	// 	io.ReaderAt
	//
	// 	// WriterAt is the interface that wraps the basic WriteAt method.
	// 	// 	WriteAt(p []byte, off int64) (n int, err error)
	// 	//
	// 	// WriteAt writes len(p) bytes from p to the underlying data stream at offset off. It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early. WriteAt must return a non-nil error if it returns n < len(p).
	// 	//
	// 	// If WriteAt is writing to a destination with a seek offset, WriteAt should not affect nor be affected by the underlying seek offset.
	// 	io.WriterAt
	//
	// 	// FileInfo methods
	// 	// Name() string       // base name of the file
	// 	// Size() int64        // length in bytes for regular files; system-dependent for others
	// 	// Mode() FileMode     // file mode bits
	// 	// ModTime() time.Time // modification time
	// 	// IsDir() bool        // abbreviation for Mode().IsDir()
	// 	// Sys() interface{}   // underlying data source (can return nil)
	// 	fs.FileInfo
	//
	// 	// FileMode methods
	// 	// String() string 	// human-readable representation of the file
	// 	// IsDir() bool 	// abbreviation for Mode().IsDir()
	// 	// IsRegular() bool // IsRegular reports whether m is a regular file.
	// 	// Perm() FileMode	// Perm returns the Unix permission bits
	// 	// Type() FileMode
	// 	FileModer
	//
	// 	// FileOps methods
	// 	// Abs() (string, error)
	// 	// Base(path string) string
	// 	// Chmod(mode os.FileMode) error
	// 	// Dir(path string) string
	// 	// Ext(path string) string
	// 	// Move(path string) error
	// 	// Split(path string) (dir, file string)
	// 	FileOps
	//
	// 	// Unix File Operations
	// 	// 	Fd() uintptr
	// 	// 	Link(newname string) error
	// 	// 	Readlink() (string, error)
	// 	// 	Remove() error
	// 	// 	Symlink(newname string) error
	// 	// 	Truncate(size int64) error
	// 	FileUnix
	// }
	GoFile

	// fs.File
	Stat() (fs.FileInfo, error)
	Read([]byte) (int, error)
	Close() error

	// FileModer
	// fs.FileInfo

	// Additional basic methods:
	// Abs() string // absolute path of the file
	// IsRegular() bo,ol // is a regular file?
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
		providedName string // original user input
		lock         bool
		isDirty      bool
		fi           os.FileInfo // cached file information
		mode         FileMode    // cached file mode
		modTime      time.Time   // used to validate cache entries

		bufio.ReadWriter

		*os.File // underlying file handle
	}
)

////////////// Return component interfaces

// file returns the file handle, *os.File.
// The minimum interface that is implemented
// by a File is:
//  io.ReadCloser
//  Stat()
//
// This implementation also has:
//  io.Writer, io.StringWriter, io.ReaderFrom, io.WriterTo, io.ReaderAt, io.WriterAt
func (f *basicFile) file() *os.File {
	if f.File == nil || f.isDirty {
		ff, err := os.OpenFile(f.providedName, os.O_RDWR, NormalMode)
		if Err(err) != nil {
			return nil
		}
		f.File = ff
	}

	return f.File
}

func (f *basicFile) rwc() Handle {
	f.isDirty = true
	ff := f.file()
	if ff == nil {
		return nil
	}

	b := bufio.NewReadWriter(f.Reader().(*bufio.Reader), f.Writer().(*bufio.Writer))

	return b
}

func (f *basicFile) Reader() io.Reader { return bufio.NewReader(f.File) }
func (f *basicFile) Writer() io.Writer { return bufio.NewWriter(f.File) }

func (f *basicFile) Stat() (fs.FileInfo, error) {
	if f.fi == nil {
		fi, err := os.Stat(f.Name())
		if Err(err) != nil {
			return nil, NewGoFileError("Gofile.Stat()", f.providedName, err)
		}
		f.fi = fi
	}
	return f.fi, nil
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
		return ErrFileLocked
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
func (f *basicFile) Mode() fs.FileMode {
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

// Base returns the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns ".". If the path consists entirely of separators, Base returns a single separator.
func (f *basicFile) Base() string { return filepath.Base(f.Abs()) }

// Dir returns all but the last element of path, typically the path's directory. After dropping the final element, Dir calls Clean on the path and trailing slashes are removed. If the path is empty, Dir returns ".". If the path consists entirely of separators, Dir returns a single separator. The returned path does not end in a separator unless it is the root directory.
func (f *basicFile) Dir() string { return filepath.Dir(f.Abs()) }

// Ext returns the file name extension used by path. The extension is the suffix beginning at the final dot in the final element of path; it is empty if there is no dot.
func (f *basicFile) Ext() string { return filepath.Ext(f.Abs()) }

// Split splits path immediately following the final Separator, separating it into a directory and file name component. If there is no Separator in path, Split returns an empty dir and file set to path. The returned values have the property that path = dir+file.
func (f *basicFile) Split() (dir, file string) { return filepath.Split(f.Abs()) }

func (f *basicFile) Link(newname string) error { return os.Link(f.Abs(), newname) }

func (f *basicFile) Move(newname string) error   { return os.Rename(f.Abs(), newname) }
func (f *basicFile) Rename(newname string) error { return os.Rename(f.Abs(), newname) }

func (f *basicFile) Readlink() (string, error)    { return os.Readlink(f.Abs()) }
func (f *basicFile) Symlink(newname string) error { return os.Symlink(f.Abs(), newname) }

func (f *basicFile) WriteTo(w io.Writer) (n int64, err error) {
	return bufio.NewReader(f.File).WriteTo(w)
}

func (f *basicFile) Remove() error {
	err := os.Remove(f.Abs())
	if err != nil {
		return err
	}
	return f.Close()
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

func (f *basicFile) Perm() fs.FileMode {
	return f.FileInfo().Mode().Perm()
}

func (f *basicFile) Type() fs.FileMode {
	return f.FileInfo().Mode().Type()
}

func (f *basicFile) String() string {
	return fmt.Sprintf("%8s %15s", f.Mode(), f.Name())
}

func (bf *basicFile) Create() error {
	f, err := os.OpenFile(bf.providedName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, NormalMode)
	if err != nil {
		return Err(NewGoFileError("gofile.create", bf.providedName, err))
	}

	bf.File = f
	return nil
}

func (bf *basicFile) Open() error {
	f, err := os.OpenFile(bf.providedName, os.O_RDONLY, NormalMode)
	if err != nil {
		return Err(NewGoFileError("gofile.open", bf.providedName, err))
	}

	bf.File = f
	return nil
}
