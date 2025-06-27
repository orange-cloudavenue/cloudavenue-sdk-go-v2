package errors

type (
	ClientError struct {
		Message string
		APIErr  APIError
	}
)

// Error returns the error message for ClientError.
func (e *ClientError) Error() string {
	if e == nil {
		return "nil ClientError"
	}
	return e.Message
}
