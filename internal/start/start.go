package start

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"myapp3.0/internal/handler"
	"strings"
)

func New() *echo.Echo {
	e := echo.New()

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	initViper()

	logLevelString := viper.GetString("loglevel")

	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		log.Fatalf("Unable to parse loglevel: %s", logLevelString)
	}

	log.SetLevel(logLevel)

	dbURL := viper.GetString("db.url")
	log.Infof("Using DB URL: %s", dbURL)

	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer pool.Close()
	log.Infof("Connected!")

	listenAddr := viper.GetString("listen")
	log.Infof("Starting HTTP server at %s...", listenAddr)

	initHandlers(pool, e)
	e.Start(listenAddr)

	log.Info("HTTP server terminated")

	return e
}

func initHandlers(pool *pgxpool.Pool, e *echo.Echo) *echo.Echo {

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.Recover())

	hndlr := handler.CatHandler{}
	hndlr.Pool = pool

	e.POST("/records", hndlr.Add)
	e.GET("/records/:id", hndlr.Get)
	e.GET("/records", hndlr.GetAll)
	e.PUT("/records/:id", hndlr.Update)
	e.DELETE("/records/:id", hndlr.Delete)
	return e
}

func initViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("restexample")

	viper.SetDefault("loglevel", "debug")
	viper.SetDefault("listen", "localhost:8080")
	viper.SetDefault("db.url", "postgres://andeisaldyun:e3cr3t@localhost:5432/catsDB")
}
