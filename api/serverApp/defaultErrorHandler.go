package serverApp

import (
	"net/http"
)

type DefaultErrorHandler struct{}

// TODO: Handle method of the error handler
func (h *DefaultErrorHandler) Handle(err *ErrorResponse, res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(err.StatusCode)
	jsonResponse(res, err.Body)
}
