package errjson

import (
	"net/http"
)

const (
	ErrorNotValidEntity = 422
)

type RespError struct {
	Type   string `json:"type"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Data   string `json:"data"`
}

func (e RespError) Error() string {
	return e.Data
}

//404 error
type NotFoundError struct {
	Resp RespError
}

func NewNotFoundError(msg string) NotFoundError {

	e := NotFoundError{
		Resp: RespError{
			Type:   "error",
			Status: http.StatusNotFound,
			Code:   "404 not found",
			Data:   msg,
		},
	}
	return e

}
func (e NotFoundError) Error() string {
	return e.Resp.Data
}

//403
type ErrForbidden struct {
	Resp RespError
}

func NewErrForbidden(msg string) ErrForbidden {
	e := ErrForbidden{
		Resp: RespError{
			Type:   "error",
			Status: http.StatusForbidden,
			Code:   "Forbidden",
			Data:   msg,
		},
	}
	return e
}

func (e ErrForbidden) Error() string {
	return e.Resp.Data
}

//422错误
type NotValidEntityError struct {
	Resp RespError
}

func NewNotValidEntityError(msg string) NotValidEntityError {
	e := NotValidEntityError{
		Resp: RespError{
			Type:   "error",
			Status: 422,
			Code:   "Unprocessable Entity",
			Data:   msg,
		},
	}
	return e
}
func (e NotValidEntityError) Error() string {
	return e.Resp.Data
}

//401
type UnauthorizedError struct {
	Resp RespError
}

func NewUnauthorizedError(msg string) UnauthorizedError {
	e := UnauthorizedError{
		Resp: RespError{
			Type:   "error",
			Status: http.StatusUnauthorized,
			Code:   "Unauthorized",
			Data:   msg,
		},
	}
	return e
}

func (e UnauthorizedError) Error() string {
	return e.Resp.Data
}

//500
type InternalServerError struct {
	Resp RespError
}

func NewInternalServerError(msg string) InternalServerError {
	e := InternalServerError{
		Resp: RespError{
			Type:   "error",
			Status: http.StatusInternalServerError,
			Code:   "internel server error",
			Data:   msg,
		},
	}
	return e
}

func (e InternalServerError) Error() string {
	return e.Resp.Data
}
