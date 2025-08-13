package domain

import "time"

type Note struct {
	ID               string
	OriginalFileName string
	StoredFileName   string
	CreatedAt        time.Time
}
