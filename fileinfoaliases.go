package basicfile

import (
	"os"
	"time"
)

//////////////////////////// FileInfo Aliases

// IsDir returns true if the file is a directory.
//
// Alias of FileInfo().IsDir()
func (f *basicFile) IsDir() bool {
	return f.FileInfo().IsDir()
}

// Type returns type bits in m (m & ModeType).
//
// Alias of FileInfo().Mode().Type().
func (f *basicFile) Type() os.FileMode {
	return f.Mode().Type()
}

// Size returns the length in bytes for
// regular files; system-dependent for others.
//
// Alias of FileInfo().Size()
func (f *basicFile) Size() int64 {
	return f.FileInfo().Size()
}

// ModTime returns the modification time.
//
// Alias of FileInfo().ModTime()
func (f *basicFile) ModTime() time.Time {
	return f.FileInfo().ModTime()
}

// Sys returns the underlying data source
// (can return nil).
//
// Alias of FileInfo().Sys()
func (f *basicFile) Sys() interface{} {
	return f.FileInfo().Sys()
}

// Mode returns the file mode bits.
//
// Alias of FileInfo().Mode()
func (f *basicFile) Mode() os.FileMode {
	return f.FileInfo().Mode()
}

// IsRegular reports whether the file a regular file.
//
// Alias of FileInfo().Mode().IsRegular()
func (f *basicFile) IsRegular() bool {
	return f.FileInfo().Mode().IsRegular()
}

// Perm returns the Unix permission bits
//
// Alias of FileInfo().Mode().Perm()
func (f *basicFile) Perm() os.FileMode {
	return f.FileInfo().Mode().Perm()
}
