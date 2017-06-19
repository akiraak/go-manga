package model

import (
	"time"
)

type User struct {
	ID				int64
	Name			string
	UpdatedAt		time.Time
	CreatedAt		time.Time
}
