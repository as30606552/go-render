package fsutils

import "os"

// Checks for the presence of a directory and creates a new one if it is not present.
// Returns an error communicating with the file system, if it occurred.
func MakeDirIfNotExists(name string) error {
	var err error
	if _, err = os.Stat(name); os.IsNotExist(err) {
		err = os.Mkdir(name, os.ModePerm)
	}
	return err
}
