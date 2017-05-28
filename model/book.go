package model

import (
	"database/sql"
	"fmt"
	"time"
	"sort"
	"strings"
)

type Book struct {
	ID				int64
	TreeType		string	`gorm:"type: enum('main', 'sub'); default: 'main'"`
	Asin			string
	SubAsinsCol		string	`gorm:"column:sub_asins"`
	PublishType		string
	Name			string
	Region			string
	DatePublish		string
	ImageS_Url		string	`gorm:"column:image_s_url"`
	ImageS_Width	int		`gorm:"column:image_s_width"`
	ImageS_Height	int		`gorm:"column:image_s_height"`
	ImageM_Url		string	`gorm:"column:image_m_url"`
	ImageM_Width	int		`gorm:"column:image_m_width"`
	ImageM_Height	int		`gorm:"column:image_m_height"`
	ImageL_Url		string	`gorm:"column:image_l_url"`
	ImageL_Width	int		`gorm:"column:image_l_width"`
	ImageL_Height	int		`gorm:"column:image_l_height"`
	Xml				sql.NullString
	TitleID			int64
	PublisherID		int64
	AuthorID		int64
	CreatedAt		time.Time
	UpdatedAt		time.Time

	Title			Title
	Publisher		Publisher
	Author			Author
}

func (b Book) Url() string {
	return fmt.Sprintf("http://amazon.jp/o/ASIN/%s", b.Asin)
}

func (b Book) SubAsins() []string {
	return strings.Split(b.SubAsinsCol, ",")
}

type TitleBook []Book

func (tbs *TitleBook) AddBook(book Book) {
	*tbs = append(*tbs, book)
	tbs.sorte()
}

func (tbs *TitleBook) PublisherID() int64 {
	return (*tbs)[0].PublisherID
}

func (tbs *TitleBook) Url() string {
	return (*tbs)[0].Url()
}

func (tbs *TitleBook) Name() string {
	return (*tbs)[0].Name
}

func (tbs *TitleBook) ImageURL() string {
	for _, book := range *tbs {
		if len(book.ImageL_Url) > 0 {
			return book.ImageL_Url
		}
	}
	return ""
}

func prio(publishType string) int {
	switch publishType {
	case "Kindle版":
		return 0
	case "コミック":
		return 1
	default:
		return 2
	}
}

func (tbs *TitleBook) sorte() {
	books := *tbs
	sort.Slice(books, func(i, j int) bool {
		return prio(books[i].PublishType) < prio(books[j].PublishType)
	})
}
