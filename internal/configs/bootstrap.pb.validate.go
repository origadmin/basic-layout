// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: configs/bootstrap.proto

package configs

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on EntrySelectorConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *EntrySelectorConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EntrySelectorConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// EntrySelectorConfigMultiError, or nil if none found.
func (m *EntrySelectorConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *EntrySelectorConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Global

	// no validation rules for Name

	// no validation rules for Version

	if len(errors) > 0 {
		return EntrySelectorConfigMultiError(errors)
	}

	return nil
}

// EntrySelectorConfigMultiError is an error wrapping multiple validation
// errors returned by EntrySelectorConfig.ValidateAll() if the designated
// constraints aren't met.
type EntrySelectorConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EntrySelectorConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EntrySelectorConfigMultiError) AllErrors() []error { return m }

// EntrySelectorConfigValidationError is the validation error returned by
// EntrySelectorConfig.Validate if the designated constraints aren't met.
type EntrySelectorConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EntrySelectorConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EntrySelectorConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EntrySelectorConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EntrySelectorConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EntrySelectorConfigValidationError) ErrorName() string {
	return "EntrySelectorConfigValidationError"
}

// Error satisfies the builtin error interface
func (e EntrySelectorConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEntrySelectorConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EntrySelectorConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EntrySelectorConfigValidationError{}

// Validate checks the field values on Bootstrap with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Bootstrap) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Bootstrap with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in BootstrapMultiError, or nil
// if none found.
func (m *Bootstrap) ValidateAll() error {
	return m.validate(true)
}

func (m *Bootstrap) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Mode

	// no validation rules for ServiceName

	// no validation rules for CryptoType

	// no validation rules for Version

	if all {
		switch v := interface{}(m.GetEntry()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Entry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Entry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEntry()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BootstrapValidationError{
				field:  "Entry",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetRegistry()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Registry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Registry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRegistry()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BootstrapValidationError{
				field:  "Registry",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetData()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BootstrapValidationError{
				field:  "Data",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetSettings()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Settings",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Settings",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSettings()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BootstrapValidationError{
				field:  "Settings",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetService()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Service",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Service",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetService()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BootstrapValidationError{
				field:  "Service",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMiddlewares()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Middlewares",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, BootstrapValidationError{
					field:  "Middlewares",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMiddlewares()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return BootstrapValidationError{
				field:  "Middlewares",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Id

	if len(errors) > 0 {
		return BootstrapMultiError(errors)
	}

	return nil
}

// BootstrapMultiError is an error wrapping multiple validation errors returned
// by Bootstrap.ValidateAll() if the designated constraints aren't met.
type BootstrapMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BootstrapMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BootstrapMultiError) AllErrors() []error { return m }

// BootstrapValidationError is the validation error returned by
// Bootstrap.Validate if the designated constraints aren't met.
type BootstrapValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BootstrapValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BootstrapValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BootstrapValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BootstrapValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BootstrapValidationError) ErrorName() string { return "BootstrapValidationError" }

// Error satisfies the builtin error interface
func (e BootstrapValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBootstrap.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BootstrapValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BootstrapValidationError{}

// Validate checks the field values on Middlewares with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Middlewares) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Middlewares with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MiddlewaresMultiError, or
// nil if none found.
func (m *Middlewares) ValidateAll() error {
	return m.validate(true)
}

func (m *Middlewares) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RegisterAsGlobal

	if all {
		switch v := interface{}(m.GetLogger()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Logger",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Logger",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLogger()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewaresValidationError{
				field:  "Logger",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCors()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Cors",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Cors",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCors()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewaresValidationError{
				field:  "Cors",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMetrics()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Metrics",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Metrics",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetrics()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewaresValidationError{
				field:  "Metrics",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetSecurity()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Security",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewaresValidationError{
					field:  "Security",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSecurity()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewaresValidationError{
				field:  "Security",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return MiddlewaresMultiError(errors)
	}

	return nil
}

// MiddlewaresMultiError is an error wrapping multiple validation errors
// returned by Middlewares.ValidateAll() if the designated constraints aren't met.
type MiddlewaresMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MiddlewaresMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MiddlewaresMultiError) AllErrors() []error { return m }

// MiddlewaresValidationError is the validation error returned by
// Middlewares.Validate if the designated constraints aren't met.
type MiddlewaresValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MiddlewaresValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MiddlewaresValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MiddlewaresValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MiddlewaresValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MiddlewaresValidationError) ErrorName() string { return "MiddlewaresValidationError" }

// Error satisfies the builtin error interface
func (e MiddlewaresValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMiddlewares.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MiddlewaresValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MiddlewaresValidationError{}

// Validate checks the field values on Settings with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Settings) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Settings with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SettingsMultiError, or nil
// if none found.
func (m *Settings) ValidateAll() error {
	return m.validate(true)
}

func (m *Settings) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CryptoType

	if len(errors) > 0 {
		return SettingsMultiError(errors)
	}

	return nil
}

// SettingsMultiError is an error wrapping multiple validation errors returned
// by Settings.ValidateAll() if the designated constraints aren't met.
type SettingsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SettingsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SettingsMultiError) AllErrors() []error { return m }

// SettingsValidationError is the validation error returned by
// Settings.Validate if the designated constraints aren't met.
type SettingsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SettingsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SettingsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SettingsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SettingsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SettingsValidationError) ErrorName() string { return "SettingsValidationError" }

// Error satisfies the builtin error interface
func (e SettingsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSettings.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SettingsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SettingsValidationError{}

// Validate checks the field values on Bootstrap_Entry with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *Bootstrap_Entry) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Bootstrap_Entry with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// Bootstrap_EntryMultiError, or nil if none found.
func (m *Bootstrap_Entry) ValidateAll() error {
	return m.validate(true)
}

func (m *Bootstrap_Entry) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	// no validation rules for Scheme

	// no validation rules for Addr

	// no validation rules for Network

	// no validation rules for Weight

	// no validation rules for EnableSwagger

	// no validation rules for EnablePprof

	if all {
		switch v := interface{}(m.GetSelector()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, Bootstrap_EntryValidationError{
					field:  "Selector",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, Bootstrap_EntryValidationError{
					field:  "Selector",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSelector()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return Bootstrap_EntryValidationError{
				field:  "Selector",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Timeout != nil {

		if all {
			switch v := interface{}(m.GetTimeout()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, Bootstrap_EntryValidationError{
						field:  "Timeout",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, Bootstrap_EntryValidationError{
						field:  "Timeout",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetTimeout()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return Bootstrap_EntryValidationError{
					field:  "Timeout",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return Bootstrap_EntryMultiError(errors)
	}

	return nil
}

// Bootstrap_EntryMultiError is an error wrapping multiple validation errors
// returned by Bootstrap_Entry.ValidateAll() if the designated constraints
// aren't met.
type Bootstrap_EntryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Bootstrap_EntryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Bootstrap_EntryMultiError) AllErrors() []error { return m }

// Bootstrap_EntryValidationError is the validation error returned by
// Bootstrap_Entry.Validate if the designated constraints aren't met.
type Bootstrap_EntryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Bootstrap_EntryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Bootstrap_EntryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Bootstrap_EntryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Bootstrap_EntryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Bootstrap_EntryValidationError) ErrorName() string { return "Bootstrap_EntryValidationError" }

// Error satisfies the builtin error interface
func (e Bootstrap_EntryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBootstrap_Entry.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Bootstrap_EntryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Bootstrap_EntryValidationError{}
