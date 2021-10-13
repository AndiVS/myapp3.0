package main

import (
	"context"

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

	log.Infof("Using DB URL: %s", cfg.DBURLMONGO)

	pool, err := pgxpool.Connect(context.Background(), cfg.DBURLMONGO)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	defer pool.Close()
	log.Infof("Connected!")

	log.Infof("Starting HTTP server at %s...", cfg.Port)

	mongoClient, mongoDB := getMongo(&cfg)
	defer mongoClient.Disconnect(context.Background()) //nolint:errcheck,gocritic

	//recordRepository := repository.New(pool)
	recordRepository := repository.NewRepositoryMongo(mongoDB)
	recordService := service.New(recordRepository)
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
	// Configure middleware with the custom claims type
	configur := middleware.JWTConfig{
		Claims:     &model.Claims{},
		SigningKey: []byte(cfg.AuthenticationKey),
	}

	admin.Use(middleware.JWTWithConfig(configur))
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

	user.Use(middleware.JWTWithConfig(configur))
	user.Use(userHandler.Service.TokenRefresherMiddleware)
	user.GET("/records/:id", recordHandler.GetC)
	user.GET("/records", recordHandler.GetAllC)

	return e
}

func getMongo(cfg *config.Config) (*mongo.Client, *mongo.Database) {
	mongoURI, dbName := "mongodb://andeisaldyun:e3cr3t@localhost:5432", "catsDB"

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	check(err)

	db := mongoClient.Database(dbName)
	return mongoClient, db
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
