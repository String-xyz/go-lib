package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	ddmiddleware "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func CORS() echo.MiddlewareFunc {
	return echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowCredentials: true, // allow cookie auth
	})
}

func Recover() echo.MiddlewareFunc {
	return echomiddleware.Recover()
}

func Logger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			return next(c)
		}
	}
}

func LogRequest() echo.MiddlewareFunc {
	return echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogMethod:    true,
		LogHost:      true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			env := os.Getenv("ENV")
			logger := c.Get("logger").(*zerolog.Logger)
			span, _ := tracer.SpanFromContext(c.Request().Context())
			fmt.Println("logger", span)
			log.Info().Uint64("trace_id", span.Context().TraceID()).Msg("trace_id")
			logger.Info().
				Str("path", v.URI).
				Str("method", v.Method).
				Int("status_code", v.Status).
				Str("request_id", v.RequestID).
				Str("host", v.Host).
				Dur("latency", v.Latency).
				Str("env", env).
				Err(v.Error).
				Send()
			return nil
		},
	})
}

// RequestID generates a unique request ID
func RequestId() echo.MiddlewareFunc {
	return echomiddleware.RequestID()
}

func Tracer(service string) echo.MiddlewareFunc {
	return ddmiddleware.Middleware(ddmiddleware.WithServiceName(service))
}
