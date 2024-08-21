package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct{}

func Start() (*echo.Echo, error) {
	// setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// use zerolog
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// bind routes
	apiGroup := e.Group("/api")
	handler := &Handler{}
	apiGroup.GET("/auth/twitch", handler.TwitchLogin)
	apiGroup.POST("/auth/twitch", handler.TwitchCallback)
	apiGroup.GET("/auth/twitch/valid", handler.TwitchValidateAuth)
	apiGroup.POST("/twitch/link-filtering", handler.TwitchLinkFiltering)

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		address := fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Print("debug2")

	return e, nil
}
