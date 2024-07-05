package apperror

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Code is a custom apperror code that has the following format:
//
//	example: 101010
//	  1: service number; 01: Model number; 010: Error number
type Code int32

const (
	CODE_UNSPECIFIED Code = 0
	// Common apperror
	CODE_UNAUTHORIZED              Code = 100001
	CODE_NO_PERMISSION             Code = 100002
	CODE_INVALID_PARAMS            Code = 100003
	CODE_NOT_FOUND                 Code = 100004
	CODE_GET_FAILED                Code = 100005
	CODE_CREATE_FAILED             Code = 100006
	CODE_UPDATE_FAILED             Code = 100007
	CODE_DELETE_FAILED             Code = 100008
	CODE_CALL_THIRD_PARTY_FAILED   Code = 100009
	CODE_OTHER_INTERNAL_SERVER_ERR Code = 100010
)

// Error apperror implement apperror built in go.
type Error struct {
	Code     Code        `json:"code"`
	Message  string      `json:"message"`
	HTTPCode int         `json:"http_code"`
	Raw      error       `json:"raw"`
	Info     interface{} `json:"info"`
}

func (e *Error) Error() string {
	if e.Raw != nil {
		return errors.Wrap(e.Raw, e.Message).Error()
	}

	msg := fmt.Sprintf("%d", e.Code)
	if len(e.Message) > 0 {
		msg += ": " + e.Message
	}
	return msg
}

func (e *Error) WithInfo(info interface{}) *Error {
	e.Info = info
	return e
}

//ErrorAs finds the first apperror in err's chain that matches Error type and returns it.
func ErrorAs(err error) (*Error, bool) {
	pe := new(Error)
	if errors.As(err, &pe) {
		return pe, true
	}
	return nil, false
}

func ErrUnauthorized(err error) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusUnauthorized,
		Code:     CODE_UNAUTHORIZED,
		Message:  "UnAuthorized",
	}
}

func ErrNoPermission() *Error {
	return &Error{
		Raw:      nil,
		HTTPCode: http.StatusForbidden,
		Code:     CODE_NO_PERMISSION,
		Message:  "No permission.",
	}
}

func ErrInvalidParams(err error) *Error {
	if e, ok := ErrorAs(err); ok {
		return &Error{
			HTTPCode: http.StatusBadRequest,
			Code:     e.Code,
			Message:  e.Message,
		}
	}

	return &Error{
		Raw:      err,
		HTTPCode: http.StatusBadRequest,
		Code:     CODE_INVALID_PARAMS,
		Message:  "invalid params",
	}
}

func ErrGet(err error, msg string) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusInternalServerError,
		Code:     CODE_GET_FAILED,
		Message:  msg,
	}
}

func ErrCreate(err error, msg string) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusInternalServerError,
		Code:     CODE_CREATE_FAILED,
		Message:  msg,
	}
}

func ErrUpdate(err error, msg string) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusInternalServerError,
		Code:     CODE_UPDATE_FAILED,
		Message:  msg,
	}
}

func ErrDelete(err error, msg string) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusInternalServerError,
		Code:     CODE_DELETE_FAILED,
		Message:  msg,
	}
}

func ErrThirdParty(err error, msg string) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusInternalServerError,
		Code:     CODE_CALL_THIRD_PARTY_FAILED,
		Message:  msg,
	}
}

func ErrOtherInternalServerError(err error, msg string) *Error {
	return &Error{
		Raw:      err,
		HTTPCode: http.StatusInternalServerError,
		Code:     CODE_OTHER_INTERNAL_SERVER_ERR,
		Message:  msg,
	}
}
