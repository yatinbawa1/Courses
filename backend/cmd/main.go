package main

import (
	"context"
	"courses/internal/auth"
	"courses/internal/config"
	"courses/internal/database"
	"courses/internal/handlers"
	"courses/internal/mailer"
	"courses/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "\nCourses-", log.Default().Flags())

	// Initialize ENV
	// And Database and Redis
	config.Init(logger)
	pool, err := database.ConnectDataBase()
	if err != nil {
		logger.Fatalf("Unable to connect with Database %w", err)
	}
	defer pool.Close()

	redisClient, err := database.NewRedisClient(config.REDIS_ADDR, "", 0)
	if err != nil {
		logger.Fatalf("Unable to Connect With Redis %w", err)
	}
	defer redisClient.Close()

	// Initialize Core Services
	mailClient := mailer.MailConfig()
	otpRepo := repository.NewRedisOTPRepo(redisClient)
	userRepo := repository.NewUserRepo(pool)
	authService := auth.NewAuthService(userRepo, otpRepo)

	// Set up a router
	mux := handlers.RegisterRoutes(logger, authService, mailClient)

	server := &http.Server{
		Handler:      mux,
		Addr:         config.Port,
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
