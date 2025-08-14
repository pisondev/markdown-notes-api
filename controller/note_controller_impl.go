package controller

import (
	"pisondev/markdown-notes-api/model/web"
	"pisondev/markdown-notes-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type NoteControllerImpl struct {
	NoteService service.NoteService
	Log         *logrus.Logger
}

func NewNoteController(noteService service.NoteService, log *logrus.Logger) NoteController {
	return &NoteControllerImpl{
		NoteService: noteService,
		Log:         log,
	}
}

func (controller *NoteControllerImpl) UploadNote(ctx *fiber.Ctx) error {
	controller.Log.Info("---CONTROLLER: UPLOAD NOTE---")

	file, err := ctx.FormFile("note")
	if err != nil {
		controller.Log.Errorf("failed to return file from multipartform: %v", err)
		return err
	}

	fileContent, err := file.Open()
	if err != nil {
		controller.Log.Errorf("failed to open fileContent: %v", err)
		return err
	}

	userIDInterface := ctx.Locals("userID")
	userID, ok := userIDInterface.(int)
	if !ok {
		controller.Log.Errorf("failed to convert type interface of userID")
		return err
	}
	noteRequest := web.NoteRequest{
		UserID:           userID,
		OriginalFilename: file.Filename,
	}

	controller.Log.Infof("call upload note service...")
	noteResponse, err := controller.NoteService.UploadNote(ctx.Context(), noteRequest, fileContent)
	if err != nil {
		controller.Log.Errorf("failed to use upload note service: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(noteResponse)
}
