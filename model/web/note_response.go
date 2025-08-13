package web

import "time"

type NoteResponse struct {
	ID               string    `json:"id"`
	OriginalFileName string    `json:"originalFileName"`
	CreatedAt        time.Time `json:"createdAt"`
}
