package service

import (
	"context"
	"database/sql"
	"io"
	"pisondev/markdown-notes-api/helper"
	"pisondev/markdown-notes-api/model/domain"
	"pisondev/markdown-notes-api/model/web"
	"pisondev/markdown-notes-api/repository"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type NoteServiceImpl struct {
	NoteRepository repository.NoteRepository
	DB             *sql.DB
	Log            *logrus.Logger
}

func NewNoteService(noteRepository repository.NoteRepository, DB *sql.DB, log *logrus.Logger) NoteService {
	return &NoteServiceImpl{
		NoteRepository: noteRepository,
		DB:             DB,
		Log:            log,
	}
}

func (s *NoteServiceImpl) UploadNote(ctx context.Context, req web.NoteRequest, file io.Reader) (web.NoteResponse, error) {
	s.Log.Info("---SERVICE: UPLOAD NOTE---")

	noteID := uuid.NewString()

	s.Log.Info("begin tx...")
	tx, err := s.DB.Begin()
	if err != nil {
		s.Log.Errorf("failed to begin tx: %v", err)
		return web.NoteResponse{}, err
	}

	note := domain.Note{
		ID:               noteID,
		UserID:           req.UserID,
		OriginalFilename: req.OriginalFilename,
		StoredFilename:   noteID + ".md",
		CreatedAt:        time.Now().UTC().Truncate(time.Second),
	}

	s.Log.Info("call save metadata repo...")
	savedMetadata, err := s.NoteRepository.SaveMetadata(ctx, tx, note)
	if err != nil {
		s.Log.Errorf("failed to use save metadata repo: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			s.Log.Errorf("failed to rollback tx: %v", err)
			return web.NoteResponse{}, errRollback
		}
		return web.NoteResponse{}, err
	}

	s.Log.Info("call save file repo...")
	err = s.NoteRepository.SaveFile(savedMetadata, file)
	if err != nil {
		s.Log.Errorf("failed to use save file repo: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			s.Log.Errorf("failed to rollback tx: %v", err)
			return web.NoteResponse{}, errRollback
		}
		return web.NoteResponse{}, err
	}

	s.Log.Info("commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		s.Log.Errorf("failed to commit tx: %v", err)
		return web.NoteResponse{}, errCommit
	}

	return helper.ToNoteResponse(savedMetadata), nil
}

func (s *NoteServiceImpl) FindAll(ctx context.Context, userID int, page int, limit int) (web.PaginatedNoteResponse, error) {
	s.Log.Info("SERVICE: FindAll")

	offset := (page - 1) * limit
	s.Log.Info("--begin tx...")
	tx, err := s.DB.Begin()
	if err != nil {
		s.Log.Errorf("--failed to begin tx: %v", err)
		return web.PaginatedNoteResponse{}, err
	}

	s.Log.Info("--call FindAll repository...")
	notes, err := s.NoteRepository.FindAll(ctx, tx, userID, limit, offset)
	if err != nil {
		s.Log.Errorf("--failed to use FindAll repository: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			s.Log.Errorf("--failed to rollback tx: %v", err)
			return web.PaginatedNoteResponse{}, errRollback
		}
		return web.PaginatedNoteResponse{}, err
	}

	s.Log.Info("--call CountAll repository...")
	total, err := s.NoteRepository.CountAll(ctx, tx)
	if err != nil {
		s.Log.Errorf("--failed to use CountAll repository: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			s.Log.Errorf("--failed to rollback tx: %v", err)
			return web.PaginatedNoteResponse{}, errRollback
		}
		return web.PaginatedNoteResponse{}, err
	}

	noteResponses := helper.ToNoteResponses(notes)
	totalPages := (total + limit - 1) / limit

	pagination := web.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}
	return web.PaginatedNoteResponse{
		Data:       noteResponses,
		Pagination: pagination,
	}, nil
}
