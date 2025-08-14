package main

import (
	"os"
	"pisondev/markdown-notes-api/app"
	"pisondev/markdown-notes-api/controller"
	"pisondev/markdown-notes-api/exception"
	"pisondev/markdown-notes-api/middleware"
	"pisondev/markdown-notes-api/repository"
	"pisondev/markdown-notes-api/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	serverPort := os.Getenv("SERVER_PORT")
	storagePath := os.Getenv("STORAGE_PATH")

	db := app.NewDB()
	validate := validator.New()

	userRepository := repository.NewUserRepository(log)
	userService := service.NewUserService(userRepository, db, validate, log)
	userController := controller.NewUserController(userService, log)

	noteRepository := repository.NewNoteRepository(storagePath, log)
	noteService := service.NewNoteService(noteRepository, db, log)
	noteController := controller.NewNoteController(noteService, log)

	server := fiber.New(fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	})
	log.Infof("server starting on port %s...", serverPort)

	server.Post("/register", userController.Register)
	server.Post("/login", userController.Login)

	noteRoutes := server.Group("/api/notes", middleware.AuthMiddleware(log))

	noteRoutes.Post("", noteController.UploadNote)

	err = server.Listen(serverPort)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
