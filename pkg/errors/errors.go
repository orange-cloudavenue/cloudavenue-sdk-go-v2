package errors

import (
	"errors"
)

func parseErrorType[errType any](err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(errType)
	return ok
}

func IsAPIError(err error) bool {
	return parseErrorType[*APIError](err)
}

func IsClientError(err error) bool {
	return parseErrorType[*ClientError](err)
}

var New = errors.New
