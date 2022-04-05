package basicfile

import (
	"bufio"
	"io"
	"os"
	"syscall"
	"time"
)

// // NewBufferedSectionWriter converts incoming Write() requests into
// // buffered, asynchronous WriteAt()'s in a section of a file.
// // Reference: https://github.com/couchbase/moss
// func NewBufferedSectionWriter(w io.WriterAt, begPos, maxBytes int64,
// 	bufSize int) *bufferedSectionWriter {
// 	stopCh := make(chan struct{})
// 	doneCh := make(chan struct{})
// 	reqCh := make(chan ioBuf)
// 	resCh := make(chan ioBuf)

// 	go func() {
// 		defer close(doneCh)
// 		defer close(resCh)

// 		buf := make([]byte, bufSize)
// 		var pos int64
// 		var err error

// 		for {
// 			select {
// 			case <-stopCh:
// 				return
// 			case resCh <- ioBuf{buf: buf, pos: pos, err: err}:
// 			}

// 			req, ok := <-reqCh
// 			if ok {
// 				buf, pos = req.buf, req.pos
// 				if len(buf) > 0 {
// 					_, err = w.WriteAt(buf, pos)
// 				}
// 			}
// 		}
// 	}()

// 	return &bufferedSectionWriter{
// 		w:   w,
// 		beg: begPos,
// 		cur: begPos,
// 		max: maxBytes,
// 		buf: make([]byte, bufSize),

// 		stopCh: stopCh,
// 		doneCh: doneCh,
// 		reqCh:  reqCh,
// 		resCh:  resCh,
// 	}
// }

// replicate is a function found at
// Reference: https://github.com/maxymania/metaclusterfs
func replicate(dst io.WriterAt, src io.ReaderAt) (err error) {
	buf := make([]byte, 1<<12)
	p := int64(0)
	for {
		n, e := src.ReadAt(buf, p)
		err = e
		if n > 0 {
			dst.WriteAt(buf[:n], p)
		}
		if err != nil {
			break
		}
	}
	return
}

type (

	// Implements Handle by implementing
	// io.ReadWriteCloser (with io.StringWriter)
	// for bufio.ReadWriter
	bufRWC struct {
		*bufio.ReadWriter
		Closer
	}

	// RWToFrom implements io.ReaderFrom and io.WriterTo
	RWToFrom interface {
		io.ReaderFrom
		io.WriterTo
	}

	// RWAt implements io.ReaderAt and io.WriterAt
	RWAt interface {
		io.ReaderAt
		io.WriterAt
	}

	// FileOps implements common file operations
	// of *os.File.
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

	// FileUnix implements additional common file
	// operations (complementing FileOps)
	// for Unix files.
	FileUnix interface {
		Fd() uintptr
		Link(newname string) error
		Readlink() (string, error)
		Remove() error
		Symlink(newname string) error
		Truncate(size int64) error
	}
)

func (b *bufRWC) Close() error {
	// TODO: implement this ...
	return nil
}

// SameFile reports whether fi1 and fi2 describe the same file.
// For example, on Unix this means that the device and inode fields
// of the two underlying structures are identical; on other systems
// the decision may be based on the path names.
// SameFile only applies to results returned by this package gofile
// It returns false in other cases.
var SameFile = os.SameFile
