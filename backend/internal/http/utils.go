package http

import (
	"errors"
	"fmt"
)

func wraps(err1 error, err2 error) error {
	return fmt.Errorf("%w: %w", err1, err2)
}

func bytesToErr(raw []byte) error {
	return errors.New(string(raw))
}
