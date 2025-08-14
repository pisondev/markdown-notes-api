package controller

import "github.com/gofiber/fiber/v2"

type NoteController interface {
	UploadNote(ctx *fiber.Ctx) error
}
