package internal

import "errors"

var ErrRecordNotFound = errors.New("cannot find record")
var ErrUniqueConstraintViolation = errors.New("unique record already exists")
var ErrCtxKeyNotFound = errors.New("key not found in context")
