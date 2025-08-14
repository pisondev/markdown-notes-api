package helper

import (
	"pisondev/markdown-notes-api/model/domain"
	"pisondev/markdown-notes-api/model/web"
)

func ToUserRegisterResponse(user domain.User) web.UserRegisterResponse {
	return web.UserRegisterResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}

func ToNoteResponse(note domain.Note) web.NoteResponse {
	return web.NoteResponse{
		ID:               note.ID,
		OriginalFilename: note.OriginalFilename,
		CreatedAt:        note.CreatedAt,
	}
}

func ToNoteResponses(notes []domain.Note) []web.NoteResponse {
	var noteResponses []web.NoteResponse
	for _, note := range notes {
		noteResponses = append(noteResponses, ToNoteResponse(note))
	}
	return noteResponses
}
