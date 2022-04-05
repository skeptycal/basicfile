package basicfile

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pkg/errors"
)

const goFileErrorPrefix = "gofile: "

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// GoFileError records an error from a specific
// GoFile function call. If there is a relevant
// path, the path will be set and the error will
// be of type *os.PathError
type GoFileError struct {
	Op   string
	Path string
	Err  error
}

// NewGoFileError returns, as an error, a new
// GoFileError which implements the gofile.Errer
// interface.
//
// If there is a relevant path given, the path
// will be set and the error will wrap *os.PathError.
//
// As a convenience, if err is nil, NewGoFileError
// returns nil.
func NewGoFileError(op, path string, err error) *GoFileError {
	if err == nil {
		return nil
	}

	switch v := err.(type) {
	case *GoFileError:
		return v
	case *os.PathError:
		if path == "" {
			path = PWD()
		}
	case *os.LinkError:
		path = v.Old
	default:
		if path != "" {
			err = &os.PathError{Op: op, Path: path, Err: err}
		}
	}

	gfe := &GoFileError{
		Op:   prependGoFilePrefix(op),
		Path: path,
		Err:  &os.PathError{Op: op, Path: path, Err: err},
	}
	gfe.Wrap("gofile stack trace")
	return gfe
}

// prependGoFilePrefix adds the default prefix
// for annotating GoFileErrors.
func prependGoFilePrefix(msg string) string {
	return goFileErrorPrefix + msg
}

// checkGoFilePath checks if the path is not
// provided, and adds PWD() as a general path.
func checkGoFilePath(path string) string {
	if path == "" {
		return PWD()
	}
	if NotExists(path) {
		return "invalid path"
	}
	return path
}

// SetError sets the details of an error, Op
// and Path. As a convenience, a pointer to the
// error is returned.
func (e *GoFileError) SetError(op, path string) *GoFileError {
	if op != "" {
		e.Op = op
	}
	if path != "" {
		e.Path = path
	}
	return e
}

func (e *GoFileError) Error() string {
	return goFileErrorPrefix + e.Op + " " + e.Path + "): " + e.Err.Error()
}

// Wrap replaces the underlying error (err) with a
// wrapper annotating it with a stack trace at the
// point Wrap is called, and the supplied message.
//
// A pointer to the new, wrapped error is returned.
// If err is nil, Wrap returns nil and performs no
// other operations.
func (e *GoFileError) Wrap(message string) *GoFileError {
	if e.Err == nil {
		return nil
	}

	e.Err = errors.Wrap(e.Err, goFileErrorPrefix+message)

	return e
}

// Wrap returns an error annotating err with a stack
// trace at the point Wrap is called, and the supplied,
// formatted message. If err is nil, Wrap returns nil.
func (e *GoFileError) Wrapf(format string, args ...interface{}) *GoFileError {
	return e.Wrap(fmt.Sprintf(format, args...))
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func (e *GoFileError) WithMessage(msg string) *GoFileError {
	e.Err = errors.WithMessage(e.Err, msg)
	return e
}

// Unwrap returns the result of calling the Unwrap
// method on err, if err's type contains an Unwrap
// method. Otherwise, Unwrap returns nil.
func (e *GoFileError) Unwrap() error {
	return errors.Unwrap(e.Err)
}

// Is reports whether any error in err's chain
// matches target.
//
// The chain consists of err itself followed by
// the sequence of errors obtained by repeatedly
// calling Unwrap.
//
// An error is considered to match a target if
// it is equal to that target or if it implements
// a method Is(error) bool such that Is(target)
// returns true.
func (e *GoFileError) Is(target error) bool {
	return errors.Is(e, target)
}

// As finds the first error in err's chain that
// matches target, and if one is found, sets
// target to that error value and returns true.
// Otherwise, it returns false.
//
// The chain consists of err itself followed by
// the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete
// value is assignable to the value pointed to by
// target, or if the error has a method
//  As(interface{}) bool
// such that As(target) returns true. In the latter
// case, the As method is responsible for setting
// target.
//
// An error type might provide an As method so it
// can be treated as if it were a different errors
//type.
//
// As panics if target is not a non-nil pointer
// to either a type that implements error, or to
// any interface type.
func (e *GoFileError) As(target any) bool {
	return errors.As(e, target)
}

// Timeout reports whether this error represents
// a timeout.
func (e *GoFileError) Timeout() bool {
	t, ok := e.Err.(timeout)
	return ok && t.Timeout()
}
