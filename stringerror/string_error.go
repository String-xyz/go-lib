package stringerror

import (
	"strings"

	"github.com/pkg/errors"
)

func Is(err error, targets ...error) bool {
	if err == nil {
		return false
	}

	for _, target := range targets {
		if strings.Contains(errors.Cause(err).Error(), target.Error()) {
			return true
		}
	}
	return false
}

var NOT_FOUND = errors.New("not found")
var FORBIDDEN = errors.New("invoking member lacks authority")
var INVALID_RESET_TOKEN = errors.New("invalid password reset token")
var INVALID_PASSWORD = errors.New("invalid password")
var ALREADY_IN_USE = errors.New("already in use")
var DEACTIVATED = errors.New("deactivated")
var INVALID_DATA = errors.New("invalid data")
var DB_ERROR = errors.New("database error")
var EXPIRED = errors.New("expired")
var UNKNOWN_DEVICE = errors.New("unknown device")
var FUNC_NOT_ALLOWED = errors.New("function is not allowed on this contract")
var CONTRACT_NOT_ALLOWED = errors.New("contract not allowed by platform on network")

/* Marlon's Proposal */

// // SError -> String Error
// type SError struct {
// 	Code        string
// 	Message     string
// 	NativeError error
// 	StackTrace  errors.StackTrace
// }

// func (e SError) Error() string {
// 	return e.Code
// }

// func (e SError) UnWrap() error {
// 	return e.NativeError
// }

// func New(code, message string) SError {
// 	return SError{Code: code, Message: message, NativeError: errors.New(code)}
// }

// var NOT_FOUND = SError{
// 	Code:        "not_found",
// 	Message:     "not found",
// 	NativeError: errors.New("not found"),
// }

// var FORBIDDEN = SError{
// 	Code:        "forbidden",
// 	Message:     "invoking member lacks authority",
// 	NativeError: errors.New("invoking member lacks authority"),
// }

// var INVALID_RESET_TOKEN = SError{
// 	Code:        "invalid_reset_token",
// 	Message:     "invalid password reset token",
// 	NativeError: errors.New("invalid password reset token"),
// }

// var INVALID_PASSWORD = SError{
// 	Code:        "invalid_password",
// 	Message:     "invalid password",
// 	NativeError: errors.New("invalid password"),
// }

// var ALREADY_IN_USE = SError{
// 	Code:        "already_in_use",
// 	Message:     "already in use",
// 	NativeError: errors.New("already in use"),
// }

// var DEACTIVATED = SError{
// 	Code:        "deactivated",
// 	Message:     "deactivated",
// 	NativeError: errors.New("deactivated"),
// }
