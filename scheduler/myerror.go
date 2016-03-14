package scheduler

import ()

const (
	ErrorNotValidEntity = 422
)

type RespError struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NotFoundError(msg string) RespError {
	resp := RespError{
		Type:    "error",
		Status:  404,
		Code:    "404 not found",
		Message: msg,
	}

	return resp
}

func NotValidEntityError(msg string) RespError {
	resp := RespError{
		Type:    "error",
		Status:  422,
		Code:    "Unprocessable Entity",
		Message: msg,
	}
	return resp
}

type UnloginUserError struct {
}

func (e UnloginUserError) Error() string {
	return "not logined user"
}
