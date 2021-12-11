package serverApp

type ErrorResponse struct {
	// Status code of the http request
	StatusCode int
	// This message can be returned to the user
	Message string
	// Another information useful to the user
	Body interface{}
	// Stack can be many information of the error and can be used to save
	// this information can not be returned to the user, it is only for internal use
	Stack string
}

type SuccessResponse struct {
	StatusCode int
	Body       interface{}
}
