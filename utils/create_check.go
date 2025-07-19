package utils

import "os"

func CreateCheckDirectory() error {
	dir := "./images/clock"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
