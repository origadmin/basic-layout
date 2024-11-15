// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package errors implements the functions, types, and interfaces for the module.
package errors

import (
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/toolkits/runtime/transport/gins"
)

func GinErrorEncoder(ctx *gins.Context, err error) {
	se := ToHttpError(err)
	codec, _ := http.CodecForRequest(ctx.Request, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		ctx.Writer.WriteHeader(nethttp.StatusInternalServerError)
		return
	}
	ctx.Writer.Header().Set("Content-Type", ctx.Request.Header.Get("Content-Type"))
	ctx.Writer.WriteHeader(int(se.Code))
	_, _ = ctx.Writer.Write(body)
}

func HttpErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := ToHttpError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}
