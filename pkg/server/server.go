package server

import (
	"context"
	"falcon-seed/pkg/logger"
	"falcon-seed/pkg/server/middleware"
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	*echo.Echo
	logger.Logger
}

// Config represents server specific config
type Config struct {
	Port                int
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

// New instantiates new Echo server
func New(logger logger.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger(logger), echoMiddleware.Recover(),
		middleware.CORS(), middleware.Headers())

	e.GET("/", healthCheck)

	return &Server{
		e,
		logger,
	}
}

// Start starts echo server
func (server *Server) Start(cfg *Config) {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	server.Echo.Debug = cfg.Debug

	// Start server
	go func() {
		server.Logger.Info("starting server...")
		if err := server.Echo.StartServer(httpServer); err != nil {
			server.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Echo.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
