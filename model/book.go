package model

import (
	"database/sql"
	"fmt"
	"time"
	"sort"
	"strconv"
	"strings"
)

type Book struct {
	ID				int64
	Asin			string
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
	TitleID			sql.NullInt64
	PublisherID		sql.NullInt64
	AuthorID		sql.NullInt64
	CreatedAt		time.Time
	UpdatedAt		time.Time

	Title			Title
	Publisher		Publisher
	Author			Author
}

func (b Book) Url() string {
	return fmt.Sprintf("http://amazon.jp/o/ASIN/%s", b.Asin)
}

type TitleBook []Book

func (tbs *TitleBook) AddBook(book Book) {
	*tbs = append(*tbs, book)
	tbs.sorte()
}

func (tbs *TitleBook) PublisherID() int64 {
	if (*tbs)[0].PublisherID.Valid {
		return (*tbs)[0].PublisherID.Int64
	}
	return 0
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

func (tbs *TitleBook) DatePublish() string {
	return (*tbs)[0].DatePublish
}

func (tbs *TitleBook) DatePublishTime() time.Time {
	timeStr := tbs.DatePublish()
	year, _ := strconv.Atoi(timeStr[0:4])
	month, _ := strconv.Atoi(timeStr[4:6])
	day, _ := strconv.Atoi(timeStr[6:8])
	time := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return time
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

func CleanName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "\"", "\\u0022", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}
