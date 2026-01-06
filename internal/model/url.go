package model

import "time"

type URL struct {
	Code      string
	LongURL   string
	CreatedAt time.Time
	ExpiresAt time.Time
}
