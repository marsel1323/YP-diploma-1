package main

import (
	"flag"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/config"
	"github.com/marsel1323/YP-diploma-1/internal/handlers"
	"github.com/marsel1323/YP-diploma-1/internal/middlewares"
	"github.com/marsel1323/YP-diploma-1/internal/repository"
	"log"
	"os"
)

func main() {
	serverAddressFlag := flag.String("a", "127.0.0.1:8080", "Listen to address:port")
	dbDsnFlag := flag.String("d", "", "Database URI")
	accrualSystemAddressFlag := flag.String("r", "", "Accrual System Address")
	keyFlag := flag.String("k", "", "Hashing key")
	flag.Parse()

	serverAddress := GetEnv("RUN_ADDRESS", *serverAddressFlag)
	dbDsn := GetEnv("DATABASE_URI", *dbDsnFlag)
	accrualSystemAddress := GetEnv("ACCRUAL_SYSTEM_ADDRESS", *accrualSystemAddressFlag)
	key := GetEnv("KEY", *keyFlag)

	cfg := config.Config{
		Address:        serverAddress,
		DSN:            dbDsn,
		AccrualAddress: accrualSystemAddress,
		Secret:         key,
	}

	dbStorage, err := repository.NewPostgresStorage(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	app := &config.Application{
		Config:   cfg,
		Sessions: make(map[string]string),
	}

	repo := handlers.NewRepo(app, dbStorage)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.POST("/api/user/register", repo.RegisterUser)
	router.POST("/api/user/login", repo.LoginUser)
	// Middleware IsAuthed

	authorized := router.Group("/")
	authorized.Use(middlewares.AuthRequired(app.Sessions))
	{
		authorized.POST("/api/user/orders", repo.CreateOrder)
		authorized.GET("/api/user/orders", repo.GetAllOrders)
		authorized.GET("/api/user/balance", nil)
		authorized.POST("/api/user/balance/withdraw", nil)
		authorized.GET("/api/user/balance/withdrawals", nil)
	}

	log.Fatal(router.Run(cfg.Address))
}

func GetEnv(key string, defaultValue string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	}
	return defaultValue
}
