package model

import (
	"time"
)

type Publisher struct {
	ID			int64
	Name		string
	CreatedAt	time.Time
}
