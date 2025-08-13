package exception

import (
	"errors"
	"pisondev/markdown-notes-api/model/web"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	status := "INTERNAL SERVER ERROR"

	validationErr := err.(validator.ValidationErrors)
	if errors.As(err, &validationErr) {
		code = fiber.ErrBadRequest.Code
		status = "BAD REQUEST"
	}

	//409
	if errors.Is(err, ErrConflictUser) {
		code = fiber.StatusConflict
		status = "CONFLICT"
	}

	errorResponse := web.ErrorResponse{
		Code:   code,
		Status: status,
		Data:   err.Error(),
	}
	return c.Status(code).JSON(errorResponse)
}
