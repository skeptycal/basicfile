package basicfile

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"time"
)

func NewBasicFile(filename string) (BasicFile, error) {
	return &basicFile{providedName: filename}, nil
}

func newFileWithErr(providedName string) (BasicFile, error) {
	return nil, ErrNotImplemented
}

// newFile returns a new BasicFile, but no error,
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
func newFile(providedName string) BasicFile {
	f, err := newFileWithErr(providedName)
	if Err(err) != nil {
		return nil
	}
	return f
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
// Reference: standard library fs.go
type BasicFile interface {

	// A File provides access to a single file.
	// The File interface is the minimum
	// implementation required of the file.
	fsFile

	// GoFile implements the most common
	// file functionality in Go.
	// GoFile // TODO fix implementation ...

	// Dirty sets isDirty to true, forcing any
	// cached values to be recalculated.
	Dirty()
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
		providedName     string // original user input
		lock             bool
		isDirty          bool
		fi               os.FileInfo // cached file information
		mode             os.FileMode // cached file mode
		modTime          time.Time   // used to validate cache entries
		bufio.ReadWriter             // only allocated when needed.
		*os.File                     // only opened when needed.
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

	b := bufRWC{}
	b.ReadWriter = bufio.NewReadWriter(f.Reader().(*bufio.Reader), f.Writer().(*bufio.Writer))

	return &b
}

func (f *basicFile) Reader() io.Reader { return bufio.NewReader(f.File) }
func (f *basicFile) Writer() io.Writer { return bufio.NewWriter(f.File) }

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

// Abs returns an absolute representation of path. If the path is not absolute it will be joined with the current working directory to turn it into an absolute path. The absolute path name for a given file is not guaranteed to be unique. Abs calls Clean on the result.
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
func (f *basicFile) Split() (dir, file string)    { return filepath.Split(f.Abs()) }
func (f *basicFile) Link(newname string) error    { return os.Link(f.Abs(), newname) }
func (f *basicFile) Move(newname string) error    { return os.Rename(f.Abs(), newname) }
func (f *basicFile) Rename(newname string) error  { return os.Rename(f.Abs(), newname) }
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
