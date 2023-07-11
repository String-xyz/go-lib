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

// LogRequest is a middleware that logs the request and extracts the span Id and trace Id from the context
// it give us the ability to trace and correlate the request in the logs as well as the starting
// point of the root span, which is use to trace the request through the system
// the span id and trace id are logged as dd.span_id and dd.trace_id for datadog to pick them up.
// Keep in mind that the tracer middleware must be added before this middleware.
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
			// we need the span to get the traceId and spanId
			// and log them, datadog needs these two values to be logged to correlate the logs with the traces
			span, ok := tracer.SpanFromContext(c.Request().Context())
			if !ok {
				logger.Warn().Msg("no span found in context")
			}
			var logEnvent *zerolog.Event
			if v.Error != nil {
				logEnvent = logger.Error().Err(v.Error)
			} else {
				logEnvent = logger.Info()
			}

			logEnvent.
				Str("path", v.URI).
				Str("method", v.Method).
				Int("status_code", v.Status).
				Str("request_id", v.RequestID).
				Str("host", v.Host).
				Dur("latency", v.Latency).
				Uint64("dd.trace_id", span.Context().TraceID()).
				Uint64("dd.span_id", span.Context().SpanID()).
				Str("env", env).
				Msg("request")
			return nil
		},
	})
}

// RequestID generates a unique request ID
func RequestId() echo.MiddlewareFunc {
	return echomiddleware.RequestID()
}

// Tracer is a middleware that traces the request and adds the span to the context
// for any other middleware to use it, or to extract the span later on.
func Tracer(service string) echo.MiddlewareFunc {
	return ddmiddleware.Middleware(ddmiddleware.WithServiceName(service))
}
