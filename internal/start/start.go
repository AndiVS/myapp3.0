package start

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"myapp3.0/internal/handler"
	"net/http"
	"strings"
)

func New() *echo.Echo {
	e := echo.New()
	var configPath string
	Run(configPath ,e)

	return e
}

func initHandlers(pool *pgxpool.Pool ,ea *echo.Echo) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	hndlr:= handler.CatHandler{}
	hndlr.Pool = pool
	e.POST("/records",hndlr.Insert)
	e.GET("/records/:id", hndlr.Select)
	e.GET("/records", hndlr.SelectAll)
	e.PUT("/records/:id" ,hndlr.Update )
	e.DELETE("/records/:id",hndlr.Delete)

	return e
}

func initViper(configPath string) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("restexample")

	viper.SetDefault("loglevel", "debug")
	viper.SetDefault("listen", "10.1.0.1:8080")
	viper.SetDefault("db.url", "postgres://andeisaldyun:e3cr3t@10.1.0.1:5432/catsDB")

	if configPath != "" {
		log.Infof("Parsing config: %s", configPath)
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Unable to read config file: %s", err)
		}
	} else {
		log.Infof("Config file is not specified.")
	}
}

func Run(configPath string, e *echo.Echo)  {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	initViper(configPath)

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
	http.Handle("/", initHandlers(pool,e))
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}

	log.Info("HTTP server terminated")
}
