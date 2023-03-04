package stringerror

import (
	"github.com/pkg/errors"
)

func ErrorIs(err, target error) bool {
	return errors.Cause(err).Error() == target.Error()
}

var NOT_FOUND = errors.New("not found")
var FORBIDDEN = errors.New("invoking member lacks authority")
var INVALID_RESET_TOKEN = errors.New("invalid password reset token")
var INVALID_PASSWORD = errors.New("invalid password")
var ALREADY_IN_USE = errors.New("already in use")
var DEACTIVATED = errors.New("deactivated")

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
