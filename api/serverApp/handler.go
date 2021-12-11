package serverApp

import (
	"context"
	"net/http"
)

type IHandler interface {
	// Handle method receive a context, and a request, and return a response of any interface
	Handle(ctx context.Context, req *http.Request) (*SuccessResponse, *ErrorResponse)
}

type INextHandler interface {
	SetNext(handler IHandler)
	SetErrorMiddleware(handler IHandler)
}

type IAppHandler interface {
	IHandler
	INextHandler
}
