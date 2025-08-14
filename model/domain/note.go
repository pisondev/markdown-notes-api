package domain

import "time"

type Note struct {
	ID               string
	UserID           int
	OriginalFilename string
	StoredFilename   string
	CreatedAt        time.Time
}
