package basicfile

import "io/fs"

// A DirEntry is an entry read from a directory
// (using the ReadDir function or a ReadDirFile's
// ReadDir method).
//
//	 type DirEntry interface {
//	 	// Name returns the name of the file (or
//	 	// subdirectory) described by the entry.
//	 	// This name is only the final element of
//	 	// the path (the base name), not the entire
//	 	// path.
//	 	// For example, Name would return "hello.go"
//	 	// not "home/gopher/hello.go".
//	 	Name() string
//
// 		// IsDir reports whether the entry describes
//	 	// a directory.
//	 	IsDir() bool
//
//	 	// Type returns the type bits for the entry.
//	 	// The type bits are a subset of the usual
//	 	// FileMode bits, those returned by the
//	 	// FileMode.Type method.
//	 	Type() fs.FileMode
//
//	 	// Info returns the FileInfo for the file or
// 		// subdirectory described by the entry.
//	 	// The returned FileInfo may be from the
//	 	// time of the original directory read
//	 	// or from the time of the call to Info.
//	 	// If the file has been removed or renamed
//	 	// since the directory read, Info may return
//	 	// an error satisfying errors.Is(err, ErrNotExist).
//	 	// If the entry denotes a symbolic link, Info
//	 	// reports the information about the link itself,
//	 	// not the link's target.
//	 	Info() (FileInfo, error)
//	 }
type DirEntry = fs.DirEntry

// DirEntry returns an entry read from a directory
// (using the ReadDir function or a ReadDirFile's
// ReadDir method).
func (f *basicFile) DirEntry() fs.DirEntry {
	return fs.FileInfoToDirEntry(f.FileInfo())
}
