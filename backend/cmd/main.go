package main

import (
	"context"
	backend "courses"
	"courses/internal/auth"
	"courses/internal/config"
	"courses/internal/database"
	"courses/internal/handlers"
	"courses/internal/mailer"
	"courses/internal/repository"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var frontendFS embed.FS

func main() {
	logger := log.New(os.Stdout, "\nCourses-", log.Default().Flags())

	// Initialize ENV
	// And Database and Redis
	config.Init(logger)
	pool, err := database.ConnectOnlineDatabase()
	if err != nil {
		logger.Fatalf("Unable to connect with Database %s", err)
	}
	defer pool.Close()

	// This is for local redis client
	// redisClient, err := database.NewRedisClient(config.REDIS_ADDR, "", 0)

	redisClient, err := database.NewRedisOnlineClient()
	if err != nil {
		logger.Fatalf("Unable to Connect With Redis %s", err)
	}

	defer redisClient.Close()

	// Initialize Frontend
	strippedFS, err := fs.Sub(backend.FrontendFS, "dist")
	if err != nil {
		log.Fatal("Failed to read embedded folder: ", err)
	}

	fileServer := http.FileServer(http.FS(strippedFS))

	// Initialize Core Services
	resendMailer := mailer.NewResendMailer()
	otpRepo := repository.NewRedisOTPRepo(redisClient)
	userRepo := repository.NewUserRepo(pool)
	refreshRepo := repository.NewRedisRefreshTokenRepo(redisClient)
	authService := auth.NewAuthService(userRepo, otpRepo, refreshRepo)

	// Set up a router
	mux := handlers.RegisterRoutes(fileServer, logger, authService, resendMailer)

	server := &http.Server{
		Handler:      mux,
		Addr:          ":" + config.Port,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Printf("Server Starting on Port %s\n", config.Port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatalf("Error Starting the Server")
		}
	}()

	SigChan := make(chan os.Signal, 1)
	signal.Notify(SigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-SigChan
	logger.Printf("Received %s Signal! Initiating Exit Routine", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(timeoutContext)
}
