package server

import (
	"context"
	"falcon-seed/pkg/logger/zap"
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
}

// Config represents server specific config
type Config struct {
	Port                int
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

// New instantiates new Echo server
func New() *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger(), echoMiddleware.Recover(),
		middleware.CORS(), middleware.Headers())

	e.GET("/", healthCheck)

	return &Server{
		e,
	}
}

// Start starts echo server
func (server *Server) Start(cfg *Config) {
	logger := zap.GetLogger()
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
			logger.Info("shutting down the server")
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
		logger.Fatal(err)
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
