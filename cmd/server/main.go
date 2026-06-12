package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/yourname/user-api/config"
	"github.com/yourname/user-api/internal/logger"
	"github.com/yourname/user-api/internal/middleware"
	"github.com/yourname/user-api/internal/routes"
	"go.uber.org/zap"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		return
	}
	if logger.Logger != nil {
		defer logger.Logger.Sync()
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Logger.Fatal("failed to load config", zap.Error(err))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Logger.Fatal("failed to open database connection", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Logger.Fatal("failed to ping database", zap.Error(err))
	}
	logger.Logger.Info("database connection established")

	app := fiber.New()

	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())

	routes.RegisterRoutes(app, db)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	logger.Logger.Info("starting server", zap.String("addr", addr))

	if err := app.Listen(addr); err != nil {
		logger.Logger.Fatal("failed to start server", zap.Error(err))
	}
}
