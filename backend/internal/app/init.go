package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/av-ugolkov/bitly-copy/internal/config"
	"github.com/av-ugolkov/bitly-copy/internal/db/redis"
	"github.com/av-ugolkov/bitly-copy/internal/handler"
	"github.com/av-ugolkov/bitly-copy/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(cfg *config.Config) {
	redis := redis.New(cfg)

	router := echo.New()
	router.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods:     []string{http.MethodOptions, http.MethodPost, http.MethodGet, http.MethodDelete},
			AllowCredentials: true,
		},
	))
	createHandlers(router, redis)

	srv := http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      router,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("start server")
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("fail run server: %v", err)
		}
	}()

	<-ctx.Done()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("fail shutdown server: %v", err)
	}
}

func createHandlers(r *echo.Echo, redis *redis.Redis) {
	slog.Info("init handlers")

	svc := service.New(redis)
	handler.Create(r, svc)
}
