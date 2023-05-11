package httperror

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/validator"
	"github.com/labstack/echo/v4"
)

type JSONError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details *Details `json:"details,omitempty"`
}

type Details struct {
	Params *validator.InvalidParams `json:"invalidParams,omitempty"`
}

func InvalidPayload400(c echo.Context, err error) error {
	errorParams := validator.ExtractErrorParams(err)
	message := errorParams[0].Message
	return c.JSON(http.StatusBadRequest, JSONError{Message: message, Code: "BAD_REQUEST", Details: &Details{Params: &errorParams}})
}

func BadRequest400(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusBadRequest, JSONError{Message: strings.Join(message, " "), Code: "BAD_REQUEST"})
	}
	return c.JSON(http.StatusBadRequest, JSONError{Message: "Bad Request", Code: "BAD_REQUEST"})
}

func Conflict409(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusConflict, JSONError{Message: strings.Join(message, " "), Code: "CONFLICT"})
	}
	return c.JSON(http.StatusConflict, JSONError{Message: "Conflict", Code: "CONFLICT"})
}

func Forbidden403(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusForbidden, JSONError{Message: strings.Join(message, " "), Code: "FORBIDDEN"})
	}
	return c.JSON(http.StatusForbidden, JSONError{Message: "not enough permissions", Code: "FORBIDDEN"})
}

func NotFound404(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusNotFound, JSONError{Message: strings.Join(message, " "), Code: "NOT_FOUND"})
	}
	return c.JSON(http.StatusNotFound, JSONError{Message: "Resource Not Found", Code: "NOT_FOUND"})
}

func Unauthorized401(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusUnauthorized, JSONError{Message: strings.Join(message, " "), Code: "UNAUTHORIZED"})
	}
	return c.JSON(http.StatusUnauthorized, JSONError{Message: "This action requires authentication", Code: "UNAUTHORIZED"})
}

func Internal500(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: strings.Join(message, " "), Code: "INTERNAL_SERVER"})
	}
	return c.JSON(http.StatusInternalServerError, JSONError{Message: "Something went wrong", Code: "INTERNAL_SERVER"})
}

func Unprocessable422(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusUnprocessableEntity, JSONError{Message: strings.Join(message, " "), Code: "UNPROCESSABLE_ENTITY"})
	}
	return c.JSON(http.StatusUnprocessableEntity, JSONError{Message: "Unable to process entity", Code: "UNPROCESSABLE_ENTITY"})
}

func NotAllowed405(c echo.Context, message ...string) error {
	if len(message) > 0 {
		return c.JSON(http.StatusMethodNotAllowed, JSONError{Message: strings.Join(message, " "), Code: "NOT_ALLOWED"})
	}
	return c.JSON(http.StatusMethodNotAllowed, JSONError{Message: "Not Allowed", Code: "NOT_ALLOWED"})
}
