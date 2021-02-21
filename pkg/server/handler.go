package server

import (
	"github.com/labstack/echo/v4"
)

type HTTP interface {
	RegisterHandlers(*echo.Group)
}
