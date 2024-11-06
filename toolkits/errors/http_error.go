package errors

import (
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/origadmin/toolkits/errors/httperr"
	"github.com/origadmin/toolkits/errors/rpcerr"

	"origadmin/basic-layout/api/v1/services/helloworld"
)

// ErrorHTTP returns an error with the given reason, code, and message.
// It is also used id for display the error message at the client with i18n support.
func ErrorHTTP(reason helloworld.HELLO_WORLD_ERROR_REASON, code int32, msg string) *httperr.Error {
	id := "http.response.status." + strings.ToLower(reason.String())
	return &httperr.Error{
		ID:     id,
		Code:   code,
		Detail: msg,
	}
}

func ErrorGRPC(reason helloworld.HELLO_WORLD_ERROR_REASON) *rpcerr.Error {
	id := "grpc.response.status." + strings.ToLower(reason.String())
	return &rpcerr.Error{
		Id:     id,
		Code:   int32(reason),
		Detail: reason.String(),
	}
}

func ToHttpError(err error) *httperr.Error {
	if err == nil {
		return nil
	}

	var httpErr *httperr.Error
	if errors.As(err, &httpErr) {
		return httpErr
	}

	var rpcErr *rpcerr.Error
	if errors.As(err, &rpcErr) {
		id := rpcErr.Id
		if strings.HasPrefix(id, "grpc.response.status.") {
			id = strings.Replace(id, "grpc.response.status.", "http.response.status.", 1)
		}
		return &httperr.Error{
			ID:     id,
			Code:   rpcErr.Code,
			Detail: rpcErr.Detail,
		}
	}

	kerr := errors.FromError(err)
	return &httperr.Error{
		ID:     "http.response.status." + strings.ToLower(kerr.Reason),
		Code:   kerr.Code,
		Detail: kerr.Message,
	}
}
