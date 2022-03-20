package basicfile

import (
	"bufio"
	"bytes"
)

type TextFile interface {
	BF
	Text() string
	Lines() (retval []string, err error)
	Sep(c byte)
}

// textfile is a file type that is specialized for utf-8 text
type textfile struct {
	Basicfile
	linesep   byte
	recordsep byte
	data      string
}

func (d *textfile) Sep(c byte) {
	d.linesep = c
}

func (d *textfile) RecordSep(c byte) {
	d.recordsep = c
}

func (d *textfile) String() string {
	return d.data
}

func (d *textfile) Lines() (retval []string, err error) {

	buf := &bytes.Buffer{}
	defer buf.Reset()

	buf.WriteString(d.data)
	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		retval = append(retval, scanner.Text())
	}

	return retval, scanner.Err()
}
