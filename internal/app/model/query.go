package model

import "time"

type Query struct {
	ID           int       `json:"-"`
	Order        string    `json:"order"`
	ProcessingAt time.Time `json:"processing_at"`
}
