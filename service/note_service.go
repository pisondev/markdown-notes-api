package service

import (
	"context"
	"io"
	"pisondev/markdown-notes-api/model/web"
)

type NoteService interface {
	UploadNote(ctx context.Context, req web.NoteRequest, file io.Reader) (web.NoteResponse, error)
	FindAll(ctx context.Context, userID int, page int, limit int) (web.PaginatedNoteResponse, error)
}
