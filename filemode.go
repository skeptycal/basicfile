package basicfile

import (
	"os"
)

func (f *basicFile) FileMode() os.FileMode {
	if f.mode == 0 {
		f.mode = f.FileInfo().Mode()
	}
	return f.mode
}

// Mode returns the filemode of file.
// If an error occurs, it is logged
// and 0 is returned.
func Mode(file string) os.FileMode {
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
// type FileMode = fs.FileMode
