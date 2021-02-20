package logger

import (
	"github.com/labstack/echo/v4"
)

const (
	Source = "source"
	User   = "user"
	Id     = "id"
	Error  = "error"
)

// Logger represents logging interface
type Logger interface {
	Info(string)
	Debug(string)
	Error(error)
	LogRequest(echo.Context, string, string, error, map[string]interface{})
}
