package secondary

import "errors"

var ErrUserNotFound = errors.New("user not found")

var ErrComicNotFound = errors.New("comic not found")
var ErrInternal = errors.New("internal error")
