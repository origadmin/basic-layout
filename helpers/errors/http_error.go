/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors provides helper functions for working with Kratos errors.
package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
)

// New is a helper to create a new Kratos error.
// It's a simple wrapper around errors.New.
func New(code int, reason, message string) *errors.Error {
	return errors.New(code, reason, message)
}

// FromError attempts to cast a generic error to a Kratos error.
// If it fails, it returns a new Kratos error with an unknown reason.
func FromError(err error) *errors.Error {
	if err == nil {
		return nil
	}
	if se := new(errors.Error); errors.As(err, &se) {
		return se
	}
	return errors.New(500, "UNKNOWN_ERROR", err.Error())
}

// Is matches a Kratos error by reason.
func Is(err error, reason string) bool {
	if se := new(errors.Error); errors.As(err, &se) {
		return se.Reason == reason
	}
	return false
}
