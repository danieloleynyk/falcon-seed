package middleware

import (
	"falcon-seed/pkg/logger"
	"github.com/labstack/echo/v4"
	"time"
)

// Logger adds basic logging for each request
func Logger(logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			logger.Info("request received",
				"remote_ip", c.RealIP(),
				"host", req.Host,
				"uri", req.RequestURI,
				"method", req.Method,
				"path", path,
				"status", res.Status,
				"latency", stop.Sub(start).String(),
			)

			return err
		}
	}
}
