package basicfile

import "io/fs"

func (f *basicFile) FileMode() FileMode {
	if f.mode == 0 {
		f.mode = f.FileInfo().Mode()
	}
	return f.mode
}

// Mode returns the filemode of file.
// If an error occurs, it is logged
// and 0 is returned.
func Mode(file string) FileMode {
	fi, err := Stat(file)
	if err != nil {
		Err(err)
		return 0
	}
	return fi.Mode()
}

// A FileMode represents a file's mode
// and permission bits. The bits have the
// same definition on all systems, so that
// information about files can be moved
// from one system to another portably. Not
// all bits apply to all systems. The only
// required bit is ModeDir for directories.
//
//  type FileMode uint32
//
// Reference: standard library fs.go
type FileMode = fs.FileMode

// // FileModer implements fs.FileMode methods
// //
// // A FileMode represents a file's mode and permission bits.
// // The bits have the same definition on all systems, so that
// // information about files can be moved from one system
// // to another portably. Not all bits apply to all systems.
// // The only required bit is ModeDir for directories.
// type FileModer interface {

// 	// human-readable representation of the file
// 	String() string

// 	// IsDir reports whether m describes a directory.
// 	// That is, it tests for the ModeDir bit being set in m.
// 	// IsDir() bool // duplicated in FileInfo interface

// 	// IsRegular reports whether m describes a regular file.
// 	// That is, it tests that no mode type bits are set.
// 	IsRegular() bool

// 	// Perm returns the Unix permission bits in m (m & ModePerm).
// 	Perm() FileMode

// 	// Type returns type bits in m (m & ModeType).
// 	Type() FileMode
// }
