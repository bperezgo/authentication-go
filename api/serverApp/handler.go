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
	// TODO: implement handlers when the response is sent to the user, it can be used to track data etc
	if ah.NextHandler != nil {
		return ah.NextHandler.Handle(res, req)
	}
	return response, err
}
