package service

import "errors"

var (
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrBadToken         = errors.New("bad token")
	ErrUpdateInProgress = errors.New("update already in progress")
	ErrInternal         = errors.New("internal error")
)
