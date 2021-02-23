package user

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterHandlers(group *echo.Group) {
	userHandlers := group.Group("/users")

	userHandlers.GET("/:user", getUser)
}

func getUser(context echo.Context) error {
	context.Set("request_name", "get user")
	return context.JSON(http.StatusOK, context.Get("username"))
}
