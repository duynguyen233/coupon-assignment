package app

import (
	"coupon-be/config"
	"coupon-be/internal/controller"
	"coupon-be/internal/repositories"
	router "coupon-be/internal/router/http"
	"coupon-be/internal/services"
	"coupon-be/pkg/httpserver"
	"coupon-be/pkg/logger"
	"io"
	"time"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	env := os.Getenv("ENV")
	if env == "PROD" {
		err := os.MkdirAll("./.log", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(
			"./.log/server.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)

	}
	l := logger.New(cfg.Log.Level)

	//connect postgres with gorm
	db, err := gorm.Open(mysql.Open(cfg.MYSQL.URL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Caching
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.RedisAddress,
		Password: cfg.Redis.RedisPassword,
		DB:       cfg.Redis.RedisDB,
	})

	// Repositories
	couponRepo := repositories.NewCouponRepository(db)
	// middleware

	// Services
	couponServices := services.NewCouponService(l)

	// Controllers
	couponController := controller.NewCouponController(l, couponServices, couponRepo, redisClient)
	orderController := controller.NewOrderController(l, couponRepo, couponServices)
	// HTTP Server
	handler := gin.New()
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.NewRouter(handler, l, couponController, orderController)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
