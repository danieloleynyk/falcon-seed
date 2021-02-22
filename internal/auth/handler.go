package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents auth http service
type HTTP struct {
	service Service
}

func (httpHandler *HTTP) RegisterHandlers(
	group *echo.Group,
) {
	// swagger:route POST /login auth login
	// Logs in user by username and password.
	// responses:
	//  200: loginResp
	//  400: errMsg
	//  401: errMsg
	// 	403: err
	//  404: errMsg
	//  500: err
	group.POST("/login", httpHandler.login)

	// swagger:operation GET /refresh/{token} auth refresh
	// ---
	// summary: Refreshes jwt token.
	// description: Refreshes jwt token by checking at database whether refresh token exists.
	// parameters:
	// - name: token
	//   in: path
	//   description: refresh token
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/refreshResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	group.GET("/refresh/:token", httpHandler.refresh)
}

// NewHTTP creates new auth http service
func NewHTTP(service Service, e *echo.Echo, mw echo.MiddlewareFunc) {
	h := HTTP{service}
	h.RegisterHandlers(e.Group("/auth"))
}

type credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (httpHandler *HTTP) login(context echo.Context) error {
	context.Set("request_name", "login")

	creds := new(credentials)
	if err := context.Bind(creds); err != nil {
		return err
	}
	token, err := httpHandler.service.Authenticate(context, creds.Username, creds.Password)
	if err != nil {
		return err
	}
	return context.JSON(http.StatusOK, token)
}

func (httpHandler *HTTP) refresh(context echo.Context) error {
	context.Set("request_name", "refresh")

	token, err := httpHandler.service.Refresh(context, context.Param("token"))
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "error generating new token")
	}

	return context.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
