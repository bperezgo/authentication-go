package serverApp

import (
	"net/http"
)

type appHandler struct {
	// From IHandler is gotten the Handle Method, and that is decission of the user
	Handler     IHandler
	NextHandler *appHandler
}

// Method to set the next handler
func (ah *appHandler) SetNext(handler *appHandler) {
	ah.NextHandler = handler
}

// Handle method to verify if exists another Handler in the chain
func (ah *appHandler) Handle(res http.ResponseWriter, req *http.Request) (*SuccessResponse, *ErrorResponse) {
	response, err := ah.Handler.Handle(res, req)
	// TODO: If the response is nil and err is nil, lets continue with the next handler
	// This is a provisional solution, but the problem is when the response is sent, it is possible
	// to the user to execute another handler, for example, to save tracking data
	if ah.NextHandler != nil {
		return ah.NextHandler.Handle(res, req)
	}
	return response, err
}
