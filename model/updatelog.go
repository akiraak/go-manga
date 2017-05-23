package model

import (
	"time"
)

type UpdateLog struct {
	ID					int64
	Date				time.Time
	FetchAsinCount		int
	UpdateAsinCount		int
	UpdatedBookCount	int
	CreatedAt			time.Time
	UpdatedAt			time.Time
}
