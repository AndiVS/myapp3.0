package main

import (
	"context"
	"fmt"
	"github.com/AndiVS/myapp3.0/internal/broker"
	"github.com/AndiVS/myapp3.0/internal/config"
	"github.com/AndiVS/myapp3.0/internal/handler"
	"github.com/AndiVS/myapp3.0/internal/middlewares"
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/AndiVS/myapp3.0/internal/server"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/AndiVS/myapp3.0/protocol"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"net"
	"os"
	"time"
)

const mongodatabase = "mongodb"
const postgresdatabase = "postgres"
const secho = "echo"
const sgrpc = "grpc"
const rediska = "redis"
const kafka = "kafka"
const rabbit = "rabbit"

func main() {
	setLog()

	cfg := config.Config{}
	config.New(&cfg)

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Unable to parse loglevel: %s", cfg.LogLevel)
	}
	log.SetLevel(logLevel)

	cfg.DBURL = getURL(&cfg)
	log.Infof("Using DB URL: %s", cfg.DBURL)

	access := service.NewJWTManager([]byte(cfg.AuthenticationKey), time.Duration(cfg.AuthenticationTokenDuration)*time.Second)
	refresh := service.NewJWTManager([]byte(cfg.RefreshKey), time.Duration(cfg.RefreshTokenDuration)*time.Second)

	var recordRepository repository.Cats
	switch cfg.System {
	case mongodatabase:
		mongoClient, mongoDatabase := getMongo(cfg.DBURL, cfg.DBName)
		defer mongoClient.Disconnect(context.Background()) //nolint:errorcheck,critic
		recordRepository = repository.NewRepository(mongoDatabase)
	case postgresdatabase:
		pool := getPostgres(cfg.DBURL)
		if err != nil {
			log.Errorf("Unable to connection to database: %v", err)
		}
		defer pool.Close()
		recordRepository = repository.NewRepository(pool)
	}

	log.Infof("Connected!")

	log.Infof("Starting HTTP server at %s...", cfg.Port)

	var recordService service.Cats
	switch cfg.Broker {
	case rediska:
		cons := runRedis()
		recordService = service.NewServiceCat(recordRepository, cons)
	case kafka:
		cons := runKafka()
		recordService = service.NewServiceCat(recordRepository, cons)
	case rabbit:
		cons := runRabbitMQ()
		recordService = service.NewServiceCat(recordRepository, cons)
	}

	//recordService := service.NewServiceCat(recordRepository, cons)
	userService := service.NewServiceUser(recordRepository)
	authenticationService := service.NewServiceAuthentication(recordRepository, access, refresh, cfg.HashSalt)

	switch cfg.Server {
	case secho:
		catHandler := handler.NewHandlerCat(recordService)
		userHandler := handler.NewHandlerUser(userService)
		authenticationHandler := handler.NewHandlerAuthentication(authenticationService)
		err = runEchoServer(catHandler, userHandler, authenticationHandler, &cfg)
	case sgrpc:
		catServer := server.NewCatServer(recordService)
		userServer := server.NewUserServer(userService)
		authenticationServer := server.NewAuthenticationServer(authenticationService)
		err = runGRPCServer(catServer, userServer, authenticationServer, &cfg, access, refresh)
	}

	log.Info("HTTP server terminated", err)
}

func setLog() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}

func getURL(cfg *config.Config) (URL string) {
	var str string
	switch cfg.System {
	case mongodatabase:
		str = fmt.Sprintf("%s://%s:%d",
			cfg.System,
			cfg.DBHost,
			cfg.DBPort,
		)
	case postgresdatabase:
		str = fmt.Sprintf("%s://%s:%s@%s:%d/%s",
			cfg.System,
			cfg.DBUser,
			// os.Getenv("DB_PASSWORD"),
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
	}
	return str
}

func getPostgres(url string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	return pool
}

func getMongo(url, dbname string) (*mongo.Client, *mongo.Database) {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}

	db := mongoClient.Database(dbname)
	return mongoClient, db
}

func runGRPCServer(recServer protocol.CatServiceServer, userServer protocol.UserServiceServer,
	serviceServer protocol.AuthServiceServer, cfg *config.Config, access, refresh *service.JWTManager) error {
	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := server.NewAuthInterceptor(access, refresh, nil)
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	grpcServer := grpc.NewServer(serverOptions...)
	protocol.RegisterUserServiceServer(grpcServer, userServer)
	protocol.RegisterCatServiceServer(grpcServer, recServer)
	protocol.RegisterAuthServiceServer(grpcServer, serviceServer)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return grpcServer.Serve(listener)
}

func runEchoServer(catHandler *handler.CatHandler, userHandler *handler.UserHandler,
	authenticationHandler *handler.AuthenticationHandler, cfg *config.Config) error {
	e := echo.New()
	initHandlers(catHandler, userHandler, authenticationHandler, e, cfg)
	err := e.Start(cfg.Port)
	if err != nil {
		log.Error("Start server error")
	}
	return err
}

func initHandlers(catHandler *handler.CatHandler, userHandler *handler.UserHandler,
	authenticationHandler *handler.AuthenticationHandler, e *echo.Echo, cfg *config.Config) *echo.Echo {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.Recover())

	/*e.POST("/records", catHandler.AddCat)
	e.GET("/records/:_id", catHandler.GetCat)
	e.GET("/records", catHandler.GetAllCat)
	e.PUT("/records/:_id", catHandler.UpdateCat)
	e.DELETE("/records/:_id", catHandler.DeleteCat)*/
	e.POST("/auth/sign-up", authenticationHandler.SignUp)
	e.POST("/auth/sign-in", authenticationHandler.SignIn)
	admin := e.Group("/admin")

	configuration := middleware.JWTConfig{
		Claims:     &model.Claims{},
		SigningKey: cfg.AuthenticationKey,
	}
	access := service.NewJWTManager([]byte(cfg.AuthenticationKey), time.Duration(cfg.AuthenticationTokenDuration)*time.Second)
	refresh := service.NewJWTManager([]byte(cfg.RefreshKey), time.Duration(cfg.RefreshTokenDuration)*time.Second)
	admin.Use(middleware.JWTWithConfig(configuration))
	admin.Use(middlewares.Check)
	admin.Use(middlewares.TokenRefresherMiddleware(access, refresh))

	admin.POST("/records", catHandler.AddCat)
	admin.GET("/records/:_id", catHandler.GetCat)
	admin.GET("/records", catHandler.GetAllCat)
	admin.PUT("/records/:_id", catHandler.UpdateCat)
	admin.DELETE("/records/:_id", catHandler.DeleteCat)

	admin.GET("/user", userHandler.GetAllUser)
	admin.PUT("/user/:username", userHandler.UpdateUser)
	admin.DELETE("/user/:username", userHandler.DeleteUser)

	user := e.Group("/user")

	user.Use(middleware.JWTWithConfig(configuration))
	// user.Use(middlewares.TokenRefresherMiddleware(authenticationHandler.Service.Access,authenticationHandler.Service.Refresh))

	user.GET("/records", catHandler.GetAllCat)
	user.GET("/records/:_id", catHandler.GetCat)

	return e
}

func runRedis() broker.Broker {
	adr := fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST"))
	client := redis.NewClient(&redis.Options{
		Addr:     "172.28.1.4:6379",
		Password: "",
		DB:       0, // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Errorf("ping %v__ error __ %v", adr, err)
	}

	redisStruct := broker.NewRedisClient(client, "STREAM")

	return redisStruct
}

func runKafka() broker.Broker {
	consumer := broker.StartKafkaConsumer()
	//producer := broker.StartKafkaProducer()

	//kafkaStruct := broker.NewKafka(consumer, producer, "Topic")
	kafkaStruct := broker.NewKafka(consumer, nil, "Topic")
	return kafkaStruct
}

func runRabbitMQ() broker.Broker {
	rabbitStruct := broker.NewRabbitMQ("Que")

	return rabbitStruct
}
