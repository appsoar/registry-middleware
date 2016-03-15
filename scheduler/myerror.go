package scheduler

import (
	"net/http"
)

const (
	ErrorNotValidEntity = 422
)

type RespError struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e RespError) Error() string {
	return e.Message
}

//404 error
type NotFoundError struct {
	resp RespError
}

func NewNotFoundError(msg string) NotFoundError {

	e := NotFoundError{
		resp: RespError{
			Type:    "error",
			Status:  http.StatusNotFound,
			Code:    "404 not found",
			Message: msg,
		},
	}
	return e

}

func (e NotFoundError) Error() string {
	return e.resp.Message
}

//422错误
type NotValidEntityError struct {
	resp RespError
}

func NewNotValidEntityError(msg string) NotValidEntityError {
	e := NotValidEntityError{
		resp: RespError{
			Type:    "error",
			Status:  422,
			Code:    "Unprocessable Entity",
			Message: msg,
		},
	}
	return e
}
func (e NotValidEntityError) Error() string {
	return e.resp.Message
}

//401
type UnauthorizedError struct {
	resp RespError
}

func NewUnauthorizedError(msg string) UnauthorizedError {
	e := UnauthorizedError{
		resp: RespError{
			Type:    "error",
			Status:  http.StatusUnauthorized,
			Code:    "Unauthorized",
			Message: msg,
		},
	}
	return e
}

func (e UnauthorizedError) Error() string {
	return e.resp.Message
}

//500
type InternalServerError struct {
	resp RespError
}

func NewInternalServerError(msg string) InternalServerError {
	e := InternalServerError{
		resp: RespError{
			Type:    "error",
			Status:  http.StatusInternalServerError,
			Code:    "internel server error",
			Message: msg,
		},
	}
	return e
}

func (e InternalServerError) Error() string {
	return e.resp.Message
}
