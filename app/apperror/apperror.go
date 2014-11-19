package apperror

import (
	"net/http"
	"strconv"
)

type Apperror struct {
	id          int
	message     string
	sys_message string
	field       string
	http_error  int
}

func NewDBError(message string, err error) Apperror {
	if len(message) == 0 {
		message = "Some error occured while querying database. Please try later."
	}
	return Apperror{id: 1, message: message, sys_message: err.Error(), http_error: http.StatusInternalServerError}
}

func NewRequiredError(message string, field string) Apperror {
	if len(message) == 0 {
		message = "A required field is missing"
	}
	return Apperror{id: 2, message: message, field: field, http_error: http.StatusBadRequest}
}

func NewInvalidInputError(message string, field string) Apperror {
	if len(message) == 0 {
		message = "A field is invalid"
	}
	return Apperror{id: 3, message: message, field: field, http_error: http.StatusBadRequest}
}

func NewNotNumericInputError(message string, err error, field string) Apperror {
	if len(message) == 0 {
		message = "A field has non-numeric chars"
	}
	return Apperror{id: 4, message: message, sys_message: err.Error(), field: field, http_error: http.StatusBadRequest}
}

func NewTokenInvalidError(message string, err error, field string) Apperror {
	if len(message) == 0 {
		message = "Invalid token"
	}
	return Apperror{id: 5, message: message, sys_message: err.Error(), field: field, http_error: http.StatusBadRequest}
}

func NewTokenExpiredError(message string) Apperror {
	if len(message) == 0 {
		message = "Token expired"
	}
	return Apperror{id: 6, message: message, http_error: http.StatusBadRequest}
}

func NewUserNameExistsError(message string, field string) Apperror {
	if len(message) == 0 {
		message = "User Name exists"
	}
	return Apperror{id: 7, message: message, field: field, http_error: http.StatusBadRequest}
}
func NewInvalidUserNamePasswordError(message string) Apperror {
	if len(message) == 0 {
		message = "User Name/Password Invalid"
	}
	return Apperror{id: 8, message: message, http_error: http.StatusBadRequest}
}
func NewInvalidUserNameError(message string) Apperror {
	if len(message) == 0 {
		message = "User Name Invalid"
	}
	return Apperror{id: 8, message: message, http_error: http.StatusBadRequest}
}

func NewInvalidPasswordError(message string, field string) Apperror {
	if len(message) == 0 {
		message = "User Password Invalid"
	}
	return Apperror{id: 8, message: message, field: field, http_error: http.StatusBadRequest}
}

func (err *Apperror) GetId() int {
	return err.id
}

func (err *Apperror) GetIdString() string {
	return strconv.Itoa(err.id)
}

func (err *Apperror) GetMessage() string {
	return err.message
}

func (err *Apperror) GetSysMesage() string {
	return err.sys_message
}

func (err *Apperror) GetField() string {
	return err.field
}

func (err *Apperror) GetHttpStatusCode() int {
	return err.http_error
}
