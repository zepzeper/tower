package clients

// ErrorResponse is the deafult error response data structure.
type ErrorResponse struct {
	Code    string
	Message string
	Data    ErrorData
}

// ErrorData is used to store extra data from error responses.
type ErrorData struct {
	Status int
}
