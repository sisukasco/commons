package utils

import (
	"os"
	"strings"
)

func FileExists(path string) bool {
	//Source : https://stackoverflow.com/a/12518877
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		return true
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

func IsImage(filetype string) bool {
	p := strings.Split(filetype, "/")
	if len(p) > 0 {
		tt := strings.ToLower(p[0])
		if tt == "image" {
			return true
		}
	}
	return false
}
