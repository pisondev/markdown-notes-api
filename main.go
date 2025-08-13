package main

import (
	"os"

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

	app := fiber.New()
	log.Infof("server starting on port %s...", serverPort)

	err = app.Listen(serverPort)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
