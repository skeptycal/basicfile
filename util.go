package basicfile

import "os"

// PWD returns the current working directory.
//
// If an error is encountered, "" is returned
// and the error is logged.
func PWD() string {
	path, err := os.Getwd()
	if Err(err) != nil {
		path = ""
	}
	return path
}
