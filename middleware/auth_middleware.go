package middleware

import (
	"os"
	"pisondev/markdown-notes-api/model/web"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Info("get auth header...")
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			errorResponse := web.ErrorResponse{
				Code:   fiber.StatusUnauthorized,
				Status: "UNAUTHORIZED",
				Data:   "missing jwt",
			}
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		logger.Info("trim authHeader prefix...")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &web.CustomClaims{}

		logger.Info("parse with claims...")
		token, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
			secretKey := os.Getenv("JWT_SECRET_KEY")
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			logger.Errorf("failed to parse: %v", err)
			return err
		}

		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
