package controller

import (
	"pisondev/markdown-notes-api/model/web"
	"pisondev/markdown-notes-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserControllerImpl struct {
	UserService service.UserService
	Log         *logrus.Logger
}

func NewUserController(userService service.UserService, log *logrus.Logger) UserController {
	return &UserControllerImpl{
		UserService: userService,
		Log:         log,
	}
}

func (controller *UserControllerImpl) Register(ctx *fiber.Ctx) error {
	controller.Log.Info("parsing request body...")
	userAuthRequest := web.UserAuthRequest{}
	err := ctx.BodyParser(&userAuthRequest)
	if err != nil {
		controller.Log.Errorf("failed to parse body: %v", err)

		errorResponse := web.ErrorResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	controller.Log.Info("calling register service...")
	userResponse, err := controller.UserService.Register(ctx.Context(), userAuthRequest)
	if err != nil {
		controller.Log.Errorf("failed to use register service in controller layer: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(userResponse)
}

func (controller *UserControllerImpl) Login(ctx *fiber.Ctx) error {
	controller.Log.Info("parsing request body...")
	userAuthRequest := web.UserAuthRequest{}
	err := ctx.BodyParser(&userAuthRequest)
	if err != nil {
		controller.Log.Errorf("failed to parse body: %v", err)

		errorResponse := web.ErrorResponse{
			Code:   fiber.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	controller.Log.Info("calling login service...")
	userResponse, err := controller.UserService.Login(ctx.Context(), userAuthRequest)
	if err != nil {
		controller.Log.Errorf("failed to use login service in controller layer: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(userResponse)
}
