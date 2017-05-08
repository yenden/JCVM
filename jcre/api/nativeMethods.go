package api

import (
	"errors"
)

func CheckArrayArgs(bArray []byte, offset int16, length int16) error {
	ln := int16(len(bArray) - 1)
	if offset < 0 && offset > ln && length <= ln-offset {
		return nil
	}
	return errors.New("Array arguments exception")
}
