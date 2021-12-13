package serverApp

import (
	"net/http"
)

type IHandler interface {
	// Handle method receive a context, and a request, and return a response of any interface
	Handle(res http.ResponseWriter, req *http.Request) (*SuccessResponse, *ErrorResponse)
}

type IErrorHandler interface {
	// The first parameter can be a struct
	Handle(err *ErrorResponse, res http.ResponseWriter, req *http.Request)
}

type iNextHandler interface {
	SetNext(handler *appHandler)
}

type IAppHandler interface {
	IHandler
	iNextHandler
}
