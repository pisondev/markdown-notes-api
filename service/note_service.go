package service

import (
	"context"
	"io"
	"pisondev/markdown-notes-api/model/web"
)

type NoteService interface {
	UploadNote(ctx context.Context, req web.NoteRequest, file io.Reader) (web.NoteResponse, error)
}
