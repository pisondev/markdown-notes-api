package repository

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"pisondev/markdown-notes-api/model/domain"

	"github.com/sirupsen/logrus"
)

type NoteRepositoryImpl struct {
	StoragePath string
	Log         *logrus.Logger
}

func NewNoteRepository(storagePath string, log *logrus.Logger) NoteRepository {
	return &NoteRepositoryImpl{
		StoragePath: storagePath,
		Log:         log,
	}
}

func (r *NoteRepositoryImpl) SaveFile(note domain.Note, file io.Reader) error {
	r.Log.Info("---REPOSITORY: SAVE FILE---")

	r.Log.Info("join filepath...")
	filePath := filepath.Join(r.StoragePath, note.StoredFilename)

	r.Log.Info("create filepath...")
	dst, err := os.Create(filePath)
	if err != nil {
		r.Log.Errorf("failed to create filepath: %v", err)
		return err
	}

	r.Log.Info("copy filepath...")
	_, err = io.Copy(dst, file)
	if err != nil {
		r.Log.Errorf("failed to copy file: %v", err)
		return err
	}
	return nil
}

func (r *NoteRepositoryImpl) SaveMetadata(ctx context.Context, tx *sql.Tx, note domain.Note) (domain.Note, error) {
	r.Log.Info("---REPOSITORY: SAVE METADATA---")

	SQL := "INSERT INTO notes(id, user_id, original_filename, stored_filename, created_at) VALUES (?,?,?,?,?)"

	r.Log.Info("exec context...")
	_, err := tx.ExecContext(ctx, SQL, note.ID, note.UserID, note.OriginalFilename, note.StoredFilename, note.CreatedAt)
	if err != nil {
		r.Log.Errorf("failed to insert note: %v", err)
		return domain.Note{}, err
	}
	return note, nil
}

func (r *NoteRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, userID int, limit int, offset int) ([]domain.Note, error) {
	r.Log.Info("REPOSITORY: FindAll")

	SQL := "SELECT * FROM notes WHERE user_id = ? LIMIT ? OFFSET ?"

	r.Log.Info("--query context...")
	rows, err := tx.QueryContext(ctx, SQL, userID, limit, offset)
	if err != nil {
		r.Log.Errorf("--failed to query context: %v", err)
		return []domain.Note{}, err
	}
	defer rows.Close()

	var notes []domain.Note

	r.Log.Info("--scan each note...")
	for rows.Next() {
		note := domain.Note{}
		err := rows.Scan(&note.ID, &note.UserID, &note.OriginalFilename, &note.StoredFilename, &note.CreatedAt)
		if err != nil {
			r.Log.Errorf("--failed to scan note: %v", err)
			return []domain.Note{}, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *NoteRepositoryImpl) CountAll(ctx context.Context, tx *sql.Tx) (int, error) {
	r.Log.Info("REPOSITORY: CountAll")
	SQL := "SELECT COUNT(user_id) FROM notes"
	r.Log.Info("--query context...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		r.Log.Errorf("--failed to query context...")
		return 0, err
	}
	defer rows.Close()

	var total int
	r.Log.Info("--scan rows integer count...")
	if rows.Next() {
		rows.Scan(&total)
		return total, nil
	} else {
		r.Log.Error("--failed to scan integer")
		return total, err
	}
}
