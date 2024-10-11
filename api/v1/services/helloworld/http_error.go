package helloworld

import (
	"strings"

	"github.com/origadmin/toolkits/errors/httperr"
	"github.com/origadmin/toolkits/errors/rpcerr"
)

func ErrorHTTP(reason HelloWorldErrorReason, code int32, msg string) *httperr.Error {
	id := "http.response.status." + strings.ToLower(reason.String())
	return &httperr.Error{
		ID:     id,
		Code:   code,
		Detail: msg,
	}
}

func ErrorGRPC(reason HelloWorldErrorReason) *rpcerr.Error {
	id := "grpc.response.status." + strings.ToLower(reason.String())
	return &rpcerr.Error{
		Id:     id,
		Code:   int32(reason),
		Detail: reason.String(),
	}
}
