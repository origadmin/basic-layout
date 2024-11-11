/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors implements the functions, types, and interfaces for the module.
package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
)

type (
	Error  = errors.Error
	Status = errors.Status
)

var (
	UnknownCode              = errors.UnknownCode
	UnknownReason            = errors.UnknownReason
	SupportPackageIsVersion1 = errors.SupportPackageIsVersion1
	E_DefaultCode            = errors.E_DefaultCode
	E_Code                   = errors.E_Code
	File_errors_errors_proto = errors.File_errors_errors_proto
	BadRequest               = errors.BadRequest
	IsBadRequest             = errors.IsBadRequest
	Unauthorized             = errors.Unauthorized
	IsUnauthorized           = errors.IsUnauthorized
	Forbidden                = errors.Forbidden
	IsForbidden              = errors.IsForbidden
	NotFound                 = errors.NotFound
	IsNotFound               = errors.IsNotFound
	Conflict                 = errors.Conflict
	IsConflict               = errors.IsConflict
	InternalServer           = errors.InternalServer
	IsInternalServer         = errors.IsInternalServer
	ServiceUnavailable       = errors.ServiceUnavailable
	IsServiceUnavailable     = errors.IsServiceUnavailable
	GatewayTimeout           = errors.GatewayTimeout
	IsGatewayTimeout         = errors.IsGatewayTimeout
	ClientClosed             = errors.ClientClosed
	IsClientClosed           = errors.IsClientClosed
	Is                       = errors.Is
	As                       = errors.As
	Unwrap                   = errors.Unwrap
	New                      = errors.New
	Newf                     = errors.Newf
	Errorf                   = errors.Errorf
	Code                     = errors.Code
	Reason                   = errors.Reason
	Clone                    = errors.Clone
	FromError                = errors.FromError
)
