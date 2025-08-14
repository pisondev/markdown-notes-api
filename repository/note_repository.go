package repository

import (
	"context"
	"database/sql"
	"io"
	"pisondev/markdown-notes-api/model/domain"
)

type NoteRepository interface {
	SaveFile(note domain.Note, file io.Reader) error
	SaveMetadata(ctx context.Context, tx *sql.Tx, note domain.Note) (domain.Note, error)
}
