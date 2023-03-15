package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	echoDatadog "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

func CORS() echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true, // allow cookie auth
	})
}

func Recover() echo.MiddlewareFunc {
	return echoMiddleware.Recover()
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
	return echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogMethod:    true,
		LogHost:      true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			env := os.Getenv("ENV")
			logger := c.Get("logger").(*zerolog.Logger)
			logger.Info().
				Str("path", v.URI).
				Str("method", v.Method).
				Int("status_code", v.Status).
				Str("request_id", v.RequestID).
				Str("host", v.Host).
				Dur("latency", v.Latency).
				Str("env", env).
				Err(v.Error).
				Msg("request")

			return nil
		},
	})
}

// RequestID generates a unique request ID
func RequestId() echo.MiddlewareFunc {
	return echoMiddleware.RequestID()
}

func Tracer() echo.MiddlewareFunc {
	return echoDatadog.Middleware()
}
