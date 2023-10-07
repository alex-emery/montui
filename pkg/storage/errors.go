package storage

import "errors"

var ErrUniqueConstraintFailed = errors.New("unique constraint failed")

// type ErrUniqueConstraintFailed struct {
// 	err error
// }

// func (err ErrUniqueConstraintFailed) Error() string {
// 	return err.Error()
// }
