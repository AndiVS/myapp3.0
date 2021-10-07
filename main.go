package main

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"myapp3.0/internal/config"
	"myapp3.0/internal/handler"
	"myapp3.0/internal/repository"
	"myapp3.0/internal/service"
)

func main() {
	e := echo.New()

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	cfg := config.Config{}
	config.New(&cfg)

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Unable to parse loglevel: %s", cfg.LogLevel)
	}

	log.SetLevel(logLevel)

	log.Infof("Using DB URL: %s", cfg.DBURL)

	pool, err := pgxpool.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer pool.Close()
	log.Infof("Connected!")

	log.Infof("Starting HTTP server at %s...", cfg.Port)

	initHandlers(pool, e)
	err = e.Start(cfg.Port)
	if err != nil {
		log.Error("Start server error")
	}

	log.Info("HTTP server terminated")
}

func initHandlers(pool *pgxpool.Pool, e *echo.Echo) *echo.Echo {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.Recover())

	recordRepository := repository.New(pool)
	recordService := service.New(recordRepository)
	recordHandler := handler.New(recordService)

	e.POST("/records", recordHandler.Add)
	e.GET("/records/:id", recordHandler.Get)
	e.GET("/records", recordHandler.GetAll)
	e.PUT("/records/:id", recordHandler.Update)
	e.DELETE("/records/:id", recordHandler.Delete)

	return e
}
