package stringerror

import (
	"fmt"
	"os"
	"strings"

	"github.com/String-xyz/go-lib/common"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func LogError(c echo.Context, err error, handlerMsg string) {
	lg := c.Get("logger").(*zerolog.Logger)
	sp, _ := tracer.SpanFromContext(c.Request().Context())
	lg.Error().Stack().Err(err).Uint64("trace_id", sp.Context().TraceID()).
		Uint64("span_id", sp.Context().SpanID()).Msg(handlerMsg)
}

func LogStringError(c echo.Context, err error, handlerMsg string) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	tracer, ok := errors.Cause(err).(stackTracer)
	if !ok {
		log.Warn().Str("error", err.Error()).Msg("error does not implement stack trace")
		return
	}

	cause := errors.Cause(err)
	st := tracer.StackTrace()

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		log.Warn().Msg("service name is missing from env")
	}

	if common.IsLocalEnv() {
		st2 := fmt.Sprintf("\nSTACK TRACE:\n%+v: [%+v ]\n\n", cause.Error(), st[0:5])
		// delete the string_api docker path from the stack trace
		st2 = strings.ReplaceAll(st2, "/"+serviceName+"/", "")
		fmt.Print(st2)
		return
	}

	LogError(c, err, handlerMsg)
}

func StringError(err error, optionalMsg ...string) error {
	if err == nil {
		return nil
	}

	concat := ""

	for _, msgs := range optionalMsg {
		concat += msgs + " "
	}

	if errors.Cause(err) == nil || errors.Cause(err) == err {
		return errors.Wrap(errors.New(err.Error()), concat)
	}

	return errors.Wrap(err, concat)
}

func ErrorIs(err, target error) bool {
	return errors.Cause(err).Error() == target.Error()
}

var ERR_NOT_FOUND = errors.New("not found")
var ERR_FORBIDDEN = errors.New("invoking member lacks authority")
var ERR_INVALID_RESET_TOKEN = errors.New("invalid password reset token")
var ERR_INVALID_PASSWORD = errors.New("invalid password")
