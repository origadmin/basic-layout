package dto

import (
	"context"
	"net/http"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = helloworld.ErrorHTTP(helloworld.HelloWorldErrorReason_USER_NOT_FOUND, http.StatusNotFound, "user not found")
)

// Greeter is a Greeter model.
type Greeter = helloworld.GreeterData

// GreeterQueryParam defines the query parameters for the `Menu` struct.
type GreeterQueryParam struct {
	CodePath         string   `form:"code" json:"code,omitempty"`                           // Code path (like xxx.xxx.xxx)
	Name             string   `form:"name" json:"name,omitempty"`                           // Display name of menu
	IncludeResources bool     `form:"include_resources" json:"include_resources,omitempty"` // Include resources
	InIDs            []string `form:"-" json:"-"`                                           // Include menu IDs
	Status           string   `form:"-" json:"-"`                                           // Status of menu (disabled, enabled)
	ParentID         string   `form:"-" json:"-"`                                           // Parent ID (From Menu.ID)
	ParentPathPrefix string   `form:"-" json:"-"`                                           // Parent path (split by .)
	UserID           string   `form:"-" json:"-"`                                           // User ID
	RoleID           string   `form:"-" json:"-"`                                           // Role ID
}

// GreeterDao is a Greater dao.
type GreeterDao interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, string, *GreeterQueryParam) (*Greeter, error)
	ListByHello(context.Context, string, *GreeterQueryParam) ([]*Greeter, error)
	ListAll(context.Context, *GreeterQueryParam) ([]*Greeter, error)
}
