package web

import "time"

type NoteResponse struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"originalFileName"`
	CreatedAt        time.Time `json:"createdAt"`
}

type PaginatedNoteResponse struct {
	Data       []NoteResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	Limit       int `json:"limit"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
}
