package router

import (
	"fmt"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/web"
	"github.com/akiraak/go-manga/elastic"
	. "github.com/akiraak/go-manga/model"
	"github.com/akiraak/go-manga/pagination"
	"github.com/labstack/echo"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type BaseParam struct {
	Nav			string
	SearchKey	string
}

func (BaseParam)NowUnix() int64 {
	return time.Now().Unix()
}

type Day struct {
	Date			time.Time
	PublisherBooks	map[int64]map[int64]*TitleBook
}

type DaysParam struct {
	BaseParam
	Days	[]Day
}

func titleGroupBooks(books []Book) map[int64]*TitleBook {
	resultBooks := map[int64]*TitleBook{}
	for _, book := range books {
		if !book.TitleID.Valid {
			continue
		}
		_, exist := resultBooks[book.TitleID.Int64]
		if exist {
			resultBooks[book.TitleID.Int64].AddBook(book)
		} else {
			tBook := &TitleBook{}
			tBook.AddBook(book)
			resultBooks[book.TitleID.Int64] = tBook
		}
	}
	return resultBooks
}

func publisherGroupBooks(titleBooks map[int64]*TitleBook) map[int64]map[int64]*TitleBook {
	resultBooks := map[int64]map[int64]*TitleBook{}
	for key, tBook := range titleBooks {
		publisherID := tBook.PublisherID()
		_, exist := resultBooks[publisherID]
		if exist {
			resultBooks[publisherID][key] = tBook
		} else {
			resultBooks[publisherID] = map[int64]*TitleBook{}
			resultBooks[publisherID][key] = tBook
		}
	}
	return resultBooks
}

func dateBooks(year int, month time.Month, day int, r18 bool) map[int64]map[int64]*TitleBook {
	books := []Book{}
	date := fmt.Sprintf("%d%02d%02d", year, month, day)
	r18Val := 0
	if r18 {
		r18Val = 1
	}
	db.ORM.
		Joins("left join publishers on publishers.id = books.publisher_id").
		Where("date_publish = ?", date).
		Where("publishers.r18 = ?", r18Val).
		Find(&books)
	tboos := titleGroupBooks(books)
	pboos := publisherGroupBooks(tboos)
	return pboos
}

func daysBooks(nav string, r18 bool) DaysParam {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	days := 5
	param := DaysParam{BaseParam: BaseParam{nav, ""}}
	param.Days = make([]Day, days)
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		param.Days[i].Date = date
		param.Days[i].PublisherBooks = dateBooks(date.Year(), date.Month(), date.Day(), r18)
	}
	return param
}

func GetIndexHandler(c echo.Context) error {
	param := daysBooks("index", false)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/index.html",
			"template/days_books.html"},
		param)
}

func GetR18Handler(c echo.Context) error {
	param := daysBooks("r18", true)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/r18.html",
			"template/days_books.html"},
		param)
}

func searchBooks(keyword string, offset int, limit int) ([]*TitleBook, int64, int) {
	if len(keyword) > 1 {
		asins, hitTotal := elastic.SearchAsins(keyword, offset, limit)
		books := booksWithAsins(asins)
		return books, hitTotal, len(asins)
	}
	return []*TitleBook{}, 0, 0
}

func GetSearchHandler(c echo.Context) error {
	keyword := c.QueryParam("key")
	page := PageQuery(c)
	type Param struct {
		BaseParam
		TitleBooks	[]*TitleBook
		Page		pagination.Page
	}
	bookMax := 10000
	limit := 200
	pageMax := bookMax / limit
	if page > pageMax {
		page = pageMax
	}
	offset := limit * (page - 1)
	param := Param{BaseParam: BaseParam{"search", keyword}}
	titleBooks, hitTotal, asinsCount := searchBooks(keyword, offset, limit)
	param.TitleBooks = titleBooks

	showPageMax := int(math.Ceil(float64(hitTotal) / float64(limit)))
	if showPageMax > pageMax {
		showPageMax = pageMax
	}

	param.Page = pagination.CreatePage(
		page,
		showPageMax,
		int(hitTotal),
		offset + 1,
		offset + asinsCount)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/search.html",
			"template/books.html",
			"template/pagination.html"},
		param)
}

func booksWithAsins(asins []elastic.AsinRecord) []*TitleBook {
	ids := make([]string, len(asins))
	for _, asin := range asins {
		ids = append(ids, asin.Asin)
	}

	books := []Book{}
	db.ORM.Where("asin IN (?)", ids).Find(&books)
	tbooks := titleGroupBooks(books)
	sortedBooks := []*TitleBook{}
	for _, tbook := range tbooks {
		sortedBooks = append(sortedBooks, tbook)
	}
	sort.Slice(sortedBooks, func(i, j int) bool {
		int1, _ := strconv.Atoi(sortedBooks[i].DatePublish())
		int2, _ := strconv.Atoi(sortedBooks[j].DatePublish())
		return int1 > int2
	})
	return sortedBooks
}

func searchUserBooks(keywords []string, offset int, limit int) ([]*TitleBook, int64, int) {
	if len(keywords) > 1 {
		asins, hitTotal := elastic.SearchUserAsins(keywords, offset, limit)
		books := booksWithAsins(asins)
		return books, hitTotal, len(asins)
	}
	return []*TitleBook{}, 0, 0
}

func GetDeveloperHandler(c echo.Context) error {
	keywords := []string{
		"HUNTER×HUNTER",
		"ヴィンランド・サガ",
		"落日のパトス",
		"狼と香辛料",
		"食戟のソーマ",
		"小説家になる方法",
		"重版出来",
		"山と食欲と私",
		"釣り船御前丸",
		"木根さんの1人でキネマ",
		"インベスターZ",
		"波よ聞いてくれ",
		"食戟のソーマ",
		"アルキメデスの大戦",
		"ふらいんぐうぃっち",
		"亜人",
		"宇宙兄弟",
		"BLUE GIANT",
		"ベイビーステップ",
		"ハイスコアガール",
		"蛇蔵",
		"僕らはみんな河合荘",
		"のの湯",
		"百姓貴族",
		"ドロヘドロ",
		"からかい上手の高木さん",
		"乙嫁語り",
		"ばらかもん",
		"君に届け",
		"のんのんびより",
		"海街diary",
		"後遺症ラジオ",
		"ワンパンマン",
		"いぶり暮らし",
		"ヒストリエ",
		"つれづれダイアリー",
		"ダンジョン飯",
		"メイドインアビス",
		"ドメスティックな彼女",
		"東京喰種",
		"進撃の巨人",
		"アオバ自転車店",
		"ちはやふる",
		"甘々と稲妻",
		"ろんぐらいだぁす",
		"はたらく細胞",
		"猫のお寺の知恩さん",
	}
	type Param struct {
		BaseParam
		TitleBooks	[]*TitleBook
		Tags		[]string
		Page		pagination.Page
	}
	param := Param{BaseParam: BaseParam{"dev", ""}}
	page := PageQuery(c)

	bookMax := 10000
	limit := 200
	pageMax := bookMax / limit
	if page > pageMax {
		page = pageMax
	}
	offset := limit * (page - 1)
	titleBooks, hitTotal, asinsCount := searchUserBooks(keywords, offset, limit)
	param.TitleBooks = titleBooks
	param.Tags = keywords

	showPageMax := int(math.Ceil(float64(hitTotal) / float64(limit)))
	if showPageMax > pageMax {
		showPageMax = pageMax
	}

	param.Page = pagination.CreatePage(
		page,
		showPageMax,
		int(hitTotal),
		offset + 1,
		offset + asinsCount)

	return web.RenderTemplate(
		c,
		http.StatusOK,
		[]string{
			"template/user_books.html",
			"template/books.html",
			"template/pagination.html"},
		param)
}

func PageQuery(c echo.Context) int {
	page, err := strconv.Atoi(c.QueryParam("p"))
	if page == 0 || err != nil {
		page = 1
	}
	return page
}
