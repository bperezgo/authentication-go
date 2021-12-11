package serverApp

import (
	"context"
	"log"
	"net/http"
)

type DefaultErrorHandler struct{}

// TODO: Handle method of the error handler
func (h *DefaultErrorHandler) Handle(ctx context.Context, req *http.Request) (*SuccessResponse, *ErrorResponse) {
	log.Println("[INFO] TODO Error Handler")
	return nil, nil
}
