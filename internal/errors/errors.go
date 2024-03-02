package errors

import "github.com/jackvonhouse/enrichment/pkg/errors"

var (
	ErrCantEnrichment = errors.NewType("can't user")
	ErrInternal       = errors.NewType("internal error")
	ErrAlreadyExists  = errors.NewType("already exists")
	ErrNotFound       = errors.NewType("not found")
	ErrEmptyField     = errors.NewType("empty field")
	ErrInvalidValue   = errors.NewType("invalid value")
)
