package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/user-api/internal/handler"
	"github.com/yourname/user-api/internal/repository"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	repo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(repo)

	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/:id", userHandler.GetUser)
	app.Get("/users", userHandler.GetAllUsers)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}
