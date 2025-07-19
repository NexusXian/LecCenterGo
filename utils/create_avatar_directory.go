package utils

import "os"

func CreateAvatarDirectory() error {
	dir := "./images/avatars"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
