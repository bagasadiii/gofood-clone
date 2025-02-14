package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bagasadiii/gofood-clone/app"
	"github.com/bagasadiii/gofood-clone/config"
	"github.com/bagasadiii/gofood-clone/handler"
	"github.com/bagasadiii/gofood-clone/middleware"
	"github.com/bagasadiii/gofood-clone/repository"
	"github.com/bagasadiii/gofood-clone/service"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load(".env")
	secretKey := os.Getenv("SECRETKEY")
	db := config.InitDB()
	logger := config.NewLogger()

	jwtService := middleware.NewJWTService([]byte(secretKey), logger)
	userRepo := repository.NewUserRepo(db, logger)
	userService := service.NewUserService(userRepo, logger, jwtService)
	userHandler := handler.NewUserHandler(userService, logger)

	merchantRepo := repository.NewMerchantRepo(db, logger)
	merchantService := service.NewMerchantService(merchantRepo, logger)
	merchantHandler := handler.NewMerchantHandler(merchantService, logger)

	driverRepo := repository.NewDriverRepo(db, logger)
	driverService := service.NewDriverService(driverRepo, logger)
	driverHandler := handler.NewDriverHandler(driverService, logger)

	dependencies := app.HandlerDependencies{
		UserEndpoint:     userHandler,
		MerchantEndpoint: merchantHandler,
		DriverEndpoint:   driverHandler,
		Middleware:       jwtService,
	}

	app := app.NewRouter(dependencies)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler

	server := &http.Server{
		Addr:    ":8080",
		Handler: cors(app.Route()),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Server is starting on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Could not listen on :8080: %v\n", zap.Error(err))
		}
	}()

	<-done
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
