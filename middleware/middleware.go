package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
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
			span, ok := tracer.SpanFromContext(c.Request().Context())
			if !ok {
				logger.Warn().Msg("no span found in context")
			}
			logger.Info().
				Str("path", v.URI).
				Str("method", v.Method).
				Int("status_code", v.Status).
				Str("request_id", v.RequestID).
				Str("host", v.Host).
				Dur("latency", v.Latency).
				Uint64("dd.trace_id", span.Context().TraceID()).
				Uint64("dd.span_id", span.Context().SpanID()).
				Str("env", env).
				Err(v.Error).
				Msg("request")
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
