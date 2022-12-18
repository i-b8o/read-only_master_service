package model

import (
	"fmt"
)

func ErrDuplicate(err error) error {
	return fmt.Errorf("record already exists: %v", err)
}

func ErrNotExists(err error) error {
	return fmt.Errorf("row not exists: %v", err)
}

func ErrUpdateFailed(err error) error {
	return fmt.Errorf("update failed: %v", err)
}

func ErrDeleteFailed(err error) error {
	return fmt.Errorf("delete failed: %v", err)
}
