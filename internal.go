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

	// Implements ReadWriteCloser for bufio.ReadWriter
	bufRWC struct {
		*bufio.ReadWriter
		Closer
	}

	Handle interface {
		io.ReadWriteCloser
		io.StringWriter
		RWToFrom
	}

	RWToFrom interface {
		// ReaderFrom is the interface that wraps
		// the ReadFrom method.
		//
		// ReadFrom reads data from r until EOF or error.
		// The return value n is the number of bytes read.
		// Any error except EOF encountered during the
		// read is also returned.
		//
		// The Copy function uses ReaderFrom if available.
		io.ReaderFrom

		// WriterTo is the interface that wraps the
		// WriteTo method:
		//     WriteTo(w Writer) (n int64, err error)
		//
		// WriteTo writes data to w until there's no
		// more data to write or when an error occurs.
		// The return value n is the number of bytes
		// written. Any error encountered during the
		// write is also returned.
		//
		// The Copy function uses WriterTo if available.
		//
		io.WriterTo
	}

	RWAt interface {
		// ReaderAt is the interface that wraps the basic ReadAt method.
		// 	ReadAt(p []byte, off int64) (n int, err error)
		//
		// ReadAt reads len(p) bytes into p starting at offset off in the underlying input source. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
		//
		// When ReadAt returns n < len(p), it returns a non-nil error explaining why more bytes were not returned. In this respect, ReadAt is stricter than Read.
		//
		// Even if ReadAt returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, ReadAt blocks until either all the data is available or an error occurs. In this respect ReadAt is different from Read.
		//
		// If the n = len(p) bytes returned by ReadAt are at the end of the input source, ReadAt may return either err == EOF or err == nil.
		//
		// If ReadAt is reading from an input source with a seek offset, ReadAt should not affect nor be affected by the underlying seek offset.
		//
		// Clients of ReadAt can execute parallel ReadAt calls on the same input source.
		//
		// Implementations must not retain p.
		// io.ReaderAt

		// WriterAt is the interface that wraps the basic WriteAt method.
		// 	WriteAt(p []byte, off int64) (n int, err error)
		//
		// WriteAt writes len(p) bytes from p to the underlying data stream at offset off. It returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered that caused the write to stop early. WriteAt must return a non-nil error if it returns n < len(p).
		//
		// If WriteAt is writing to a destination with a seek offset, WriteAt should not affect nor be affected by the underlying seek offset.
		// io.WriterAt
	}

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
