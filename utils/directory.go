package utils

import (
	"errors"
	"os"
)

// @author: [piexlmax](https://github.com/piexlmax)
// @function: PathExists
// @description: directory is exists
// @param: path string
// @return: bool, error

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("directory with same name already exists")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
