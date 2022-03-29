package basicfile

import (
	"bufio"
	"bytes"
	"strings"
)

type TextFile interface {
	BasicFile
	Text() string
	Lines() (retval []string, err error)
	Sep(c byte)
}

// textfile is a basicfile type that is
// specialized for utf-8 string data
type textfile struct {
	basicFile
	linesep   byte `default:"\n"`
	recordsep byte `default:"\t"`
	wordsep   byte `default:" "`
	data      string
	dirty     bool
	lines     []string // only used JIT
	records   []string // only used JIT
	words     []string // only used JIT
}

func (d *textfile) Data() string        { return d.data }
func (d *textfile) Dirty() bool         { return d.dirty }
func (d *textfile) Sep() byte           { return d.linesep }
func (d *textfile) RecordSep() byte     { return d.recordsep }
func (d *textfile) WordSep() byte       { return d.wordsep }
func (d *textfile) SetSep(c byte)       { d.linesep = c }
func (d *textfile) SetRecordSep(c byte) { d.recordsep = c }
func (d *textfile) String() string      { return d.Data() }

func (d *textfile) Lines() ([]string, error) {
	if len(d.lines) == 0 || d.dirty {

		// allocate lines
		count := strings.Count(d.data, string(d.linesep))
		size := int(len(d.data) / count)
		retval := make([]string, size)

		// allocate scanner buffer
		// size := len(d.data)
		// buf := bytes.NewBuffer(make([]byte, 0, size))
		buf := bytes.NewBufferString(d.data)
		defer buf.Reset()

		// scan ...
		scanner := bufio.NewScanner(buf)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			retval = append(retval, scanner.Text())
		}
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		d.lines = append(d.lines, retval...)
	}
	return d.lines, nil
}
