/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors implements the functions, types, and interfaces for the module.
package errors

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Coder interface {
	String() string
	Descriptor() protoreflect.EnumDescriptor
	Type() protoreflect.EnumType
	Number() protoreflect.EnumNumber
	EnumDescriptor() ([]byte, []int)
}
