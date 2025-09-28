package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/handler"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/infrastructure/postgres"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://app:app@localhost:5432/appdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Infrastructure layer
	userRepository := postgres.NewUserRepository(db)
	hasher := service.NewPasswordHasher(bcrypt.DefaultCost)

	// Service layer (business logic)
	userService := service.NewUserService(userRepository, hasher, logger)

	// Handler layer (presentation)
	userHandler := handler.NewUserHandler(userService, logger)
	r := handler.NewRouter(userHandler)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
