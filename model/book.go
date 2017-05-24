package model

import (
	"fmt"
	"time"
	"strings"
)

type Book struct {
	ID				int64
	TreeType		string	`gorm:"type: enum('main', 'sub'); default: 'main'"`
	Asin			string
	SubAsinsCol		string	`gorm:"column:sub_asins"`
	PublishType		string
	Title			string
	Region			string
	DatePublish		time.Time
	ImageS_Url		string	`gorm:"column:image_s_url"`
	ImageS_Width	int		`gorm:"column:image_s_width"`
	ImageS_Height	int		`gorm:"column:image_s_height"`
	ImageM_Url		string	`gorm:"column:image_m_url"`
	ImageM_Width	int		`gorm:"column:image_m_width"`
	ImageM_Height	int		`gorm:"column:image_m_height"`
	ImageL_Url		string	`gorm:"column:image_l_url"`
	ImageL_Width	int		`gorm:"column:image_l_width"`
	ImageL_Height	int		`gorm:"column:image_l_height"`
	PublisherID		int64
	AuthorID		int64
	CreatedAt		time.Time
	UpdatedAt		time.Time

	Publisher		Publisher
	Author			Author
}

func (b Book) Url() string {
	return fmt.Sprintf("http://amazon.jp/o/ASIN/%s", b.Asin)
}

func (b Book) SubAsins() []string {
	return strings.Split(b.SubAsinsCol, ",")
}
