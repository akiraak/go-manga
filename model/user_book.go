package model

import (
	"time"
)

type UserBook struct {
	ID			int64
	UserID		int64
	BookID		int64
	CreatedAt	time.Time

	User		User
	Book		Book
}
