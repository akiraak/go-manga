package model

import (
	"time"
)

type BbsComment struct {
	ID				int64
	Name			string
	Comment			string
	IpHash			string
	CreatedAt		time.Time
	UpdatedAt		time.Time
}

func (c *BbsComment) UpdatedAtJp() time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return c.UpdatedAt.In(jst)
}
