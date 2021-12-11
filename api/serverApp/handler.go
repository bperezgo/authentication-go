package serverApp

import (
	"context"
	"net/http"
)

type IHandler interface {
	// Handle method receive a context, and a request, and return a response of any interface
	Handle(ctx context.Context, req *http.Request) (*SuccessResponse, *ErrorResponse)
}

type iNextHandler interface {
	SetNext(handler *appHandler)
}
type iSetErrorHandler interface {
	SetErrorHandler(handler IHandler)
}

type IAppHandler interface {
	IHandler
	iNextHandler
}

type appHandler struct {
	// From IHandler is gotten the Handle Method, and that is decission of the user
	Handler      IHandler
	NextHandler  *appHandler
	ErrorHandler IHandler
}

// Method to set the next handler
func (ah *appHandler) SetNext(handler *appHandler) {
	ah.NextHandler = handler
}

// Handle method to verify if exists another Handler in the chain
func (ah *appHandler) Handle(ctx context.Context, req *http.Request) (*SuccessResponse, *ErrorResponse) {
	response, err := ah.Handler.Handle(ctx, req)
	// TODO: If the response is nil and err is nil, lets continue with the next handler
	// This is a provisional solution, but the problem is when the response is sent, it is possible
	// to the user to execute another handler, for example, to save tracking data
	if ah.NextHandler != nil && response == nil && err == nil {
		return ah.NextHandler.Handle(ctx, req)
	}
	return response, err
}
