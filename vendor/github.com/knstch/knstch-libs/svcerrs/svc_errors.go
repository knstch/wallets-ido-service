package svcerrs

import "errors"

var (
	ErrInvalidData  = errors.New("invalid data")
	ErrDataNotFound = errors.New("data not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
	ErrForbidden    = errors.New("forbidden")
	ErrGone         = errors.New("gone")
)
