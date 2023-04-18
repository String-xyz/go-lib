package validator

import (
	"fmt"
	netmail "net/mail"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var tagsMessage = map[string]string{
	"required": "is required",
	"email":    "must be a valid email",
	"gte":      "must be greater or equal to",
	"gt":       "must be at least",
	"numeric":  "must be a valid numeric value",
}

type InvalidParamError struct {
	Param        string `json:"param"`
	Value        any    `json:"value"`
	ExpectedType string `json:"expectedType"`
	Message      string `json:"message"`
}

type InvalidParams []InvalidParamError

type Validator struct {
	validator *validator.Validate
}

// Validate runs validation on structs as default
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// New Returns an API Validator with the underlying struct validator
func New() *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{v}
}

// ExtractErrorParams loops over the errors returned by a validation
// this is the simplest validation, we at some point will want to extend it
func ExtractErrorParams(err error) InvalidParams {
	params := InvalidParams{}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return params
	}

	for _, err := range err.(validator.ValidationErrors) {
		p := InvalidParamError{
			Param:        err.Field(),
			Value:        err.Value(),
			ExpectedType: err.Type().String(),
			Message:      message(err),
		}

		params = append(params, p)
	}

	return params
}

func message(f validator.FieldError) string {
	message := tagsMessage[f.Tag()]
	if strings.HasPrefix(f.Tag(), "g") {
		return fmt.Sprintf("%v %s %s", f.Value(), message, f.Param())
	}
	if message == "" {
		return "Some fields are missing or invalid, please provide all required data"
	}
	return fmt.Sprintf("%s %s", f.Field(), message)
}

func IsUUID(id ...string) bool {
	validate := validator.New()

	for _, i := range id {
		if err := validate.Var(i, "uuid"); err != nil {
			return false
		}
	}

	return true
}

func ValidEmail(email string) bool {
	_, err := netmail.ParseAddress(email)
	return err == nil
}
