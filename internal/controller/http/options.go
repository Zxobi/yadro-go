package http

import (
	"time"
)

type Option func(*router)

func ScanTimeout(timeout time.Duration) Option {
	return func(r *router) {
		r.scanTimeout = timeout
	}
}

func ScanLimit(limit int) Option {
	return func(r *router) {
		r.scanLimit = limit
	}
}
