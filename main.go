package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"myapp3.0/internal/config"
	"myapp3.0/internal/handler"
	"myapp3.0/internal/middlewar"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"myapp3.0/internal/service"
	"time"
)

func main() {
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

	cfg.DBURL = getURL(&cfg)
	log.Infof("Using DB URL: %s", cfg.DBURL)

	var recordRepository repository.Cats
	switch cfg.System {
	case "mongodb":
		mongoClient, mongoDB := getMongo(cfg.DBURL, cfg.DBName)
		defer mongoClient.Disconnect(context.Background())
		recordRepository = repository.NewRepository(mongoDB)
	case "postgres":
		pool := getPostgres(cfg.DBURL)
		if err != nil {
			log.Fatalf("Unable to connection to database: %v", err)
		}
		defer pool.Close()
		recordRepository = repository.NewRepository(pool)
	}

	log.Infof("Connected!")

	log.Infof("Starting HTTP server at %s...", cfg.Port)

	recordService := service.NewService(recordRepository)
	recordHandler := handler.NewC(recordService)

	userService := service.NewAuthorizer(
		recordRepository,
		cfg.HashSalt,
		[]byte(cfg.AuthenticationKey),
		[]byte(cfg.RefreshKey),
		time.Duration(cfg.AuthenticationTokenDuration)*time.Second,
		time.Duration(cfg.RefreshTokenDuration)*time.Second,
	)
	userHandler := handler.NewU(userService)

	e := echo.New()
	initHandlers(recordHandler, userHandler, e, &cfg)
	err = e.Start(cfg.Port)
	if err != nil {
		log.Error("Start server error")
	}

	log.Info("HTTP server terminated")
}

func initHandlers(recordHandler *handler.CatHandler, userHandler *handler.UserHandler, e *echo.Echo, cfg *config.Config) *echo.Echo {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.Recover())

	e.POST("/auth/sign-up", userHandler.AddU)
	e.POST("/auth/sign-in", userHandler.SignIn)

	admin := e.Group("/admin")
	config := middleware.JWTConfig{
		Claims:     &model.Claims{},
		SigningKey: []byte(cfg.AuthenticationKey),
	}

	admin.Use(middleware.JWTWithConfig(config))
	admin.Use(middlewar.Check)
	admin.Use(userHandler.Service.TokenRefresherMiddleware)

	admin.POST("/records", recordHandler.AddC)
	admin.GET("/records/:id", recordHandler.GetC)
	admin.GET("/records", recordHandler.GetAllC)
	admin.PUT("/records/:id", recordHandler.UpdateC)
	admin.DELETE("/records/:id", recordHandler.DeleteC)

	admin.POST("/user", userHandler.AddU)
	admin.GET("/user", userHandler.GetAllU)
	admin.PUT("/user/:username", userHandler.UpdateU)
	admin.DELETE("/user/:username", userHandler.DeleteU)

	user := e.Group("/user")

	user.Use(middleware.JWTWithConfig(config))
	user.Use(userHandler.Service.TokenRefresherMiddleware)

	user.POST("/records", recordHandler.AddC)
	user.GET("/user", userHandler.GetAllU)

	user.GET("/records/:id", recordHandler.GetC)
	user.GET("/records", recordHandler.GetAllC)

	return e
}

func getPostgres(url string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	return pool
}

func getMongo(url, dbname string) (*mongo.Client, *mongo.Database) {

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	//mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}

	db := mongoClient.Database(dbname)
	return mongoClient, db
}

func getURL(cfg *config.Config) (URL string) {
	var str string
	if cfg.System == "mongodb" {
		str = fmt.Sprintf("%s://%s:%s@%s:%d",
			cfg.System,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
		)
	} else {
		str = fmt.Sprintf("%s://%s:%s@%s:%d/%s",
			cfg.System,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
	}
	return str
}
