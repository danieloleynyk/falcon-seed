package middleware

import (
	"falcon-seed/pkg/logger"
	"github.com/labstack/echo/v4"
	"strconv"
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

			requestName := c.Get("request_name")

			if requestName == nil {
				requestName = "Unknown request"
			}

			logger.Info(requestName.(string)+" request was received",
				"remote_ip", c.RealIP(),
				"host", req.Host,
				"uri", req.RequestURI,
				"method", req.Method,
				"path", path,
				"status", res.Status,
				"latency", strconv.FormatInt(int64(stop.Sub(start)), 10),
			)

			return err
		}
	}
}
