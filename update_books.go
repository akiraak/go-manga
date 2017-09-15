package main

import (
	"errors"
	"fmt"
	"github.com/DDRBoxman/go-amazon-product-api"
	"github.com/PuerkitoBio/goquery"
	"github.com/akiraak/go-manga/db"
	"github.com/akiraak/go-manga/elastic"
	. "github.com/akiraak/go-manga/model"
	"github.com/akiraak/go-manga/tools/userbooks"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const updateBookInterval = time.Duration(12) * time.Hour
var url2AsinReg = regexp.MustCompile(`/dp/(\w+)`)


func httpGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}
	userAgent := "Linux (wget)"
	req.Header.Set("User-Agent", userAgent)
	return client.Do(req)
}

func url2Asin(url string) (string, error) {
	result := url2AsinReg.FindAllStringSubmatch(url, 1)
	if len(result) > 0 {
		return result[0][1], nil
	}
	return "", errors.New("URL does not include asin. " + url)
}

type Asins []string

func (a Asins)Exist(checkAsin string) bool {
	for _, asin := range a {
		if asin == checkAsin {
			return true
		}
	}
	return false
}

type TitleAsin Asins

func (ts *TitleAsin)AddAsins(asins Asins) {
	*ts = append(*ts, asins...)
}

func getUrlAsins(url string) ([]*TitleAsin, error) {
	titleAsins := []*TitleAsin{}
	success := false
	for i := 0; i < 10; i++ {
		resp, err := httpGet(url)
		time.Sleep(10 * time.Second)
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			continue
		}
		doc.Find(".s-result-item").Each(func(i int, book *goquery.Selection) {
			asins := []string{}
			book.Find("a").Each(func(i int, a *goquery.Selection) {
				url, exists := a.Attr("href")
				if exists {
					asin, err := url2Asin(url)
					if err == nil {
						asins = append(asins, asin)
					}
				}
			})
			asins = unique(asins)
			titleAsin := &TitleAsin{}
			titleAsin.AddAsins(asins)
			titleAsins = append(titleAsins, titleAsin)
		})
		if len(titleAsins) == 0 {
			if doc.Find("#noResultsTitle").Length() != 0 {
				success = true
				break
			} else {
				log.Printf("Error. There aren't books. count:%d\n", i)
				continue
			}
		} else {
			success = true
			break
		}
	}
	if success {
		return titleAsins, nil
	}
	return titleAsins, errors.New("Cant't fetch url.")
}

func unique(values []string) []string {
	values_map := map[string]int{}
	for _, value := range values {
		_, exists := values_map[value]
		if !exists {
			values_map[value] = 0
		}
	}
	result := []string{}
	for value, _ := range values_map {
		result = append(result, value)
	}
	return result
}

func getUrl(page int) string {
	return fmt.Sprintf("%s&page=%d",
		"https://www.amazon.co.jp/s/ref=sr_nr_p_n_publication_date_3?fst=as%3Aoff&rh=n%3A465392%2Cn%3A%21465610%2Cn%3A466280%2Cn%3A2278488051%2Cp_n_publication_date%3A2315442051%7C2285539051&bbn=2278488051&ie=UTF8&qid=1495328200",
		page)
}

func validAsins(asins []string) []string {
	checkTime := time.Now()
	checkTime = checkTime.Add(-updateBookInterval)
	var excludeBooks []Book
	db.ORM.Where("asin in (?)", asins).Where("updated_at > ?", checkTime).Find(&excludeBooks)
	excludeAsins := map[string]int{}
	for _, book := range excludeBooks {
		excludeAsins[book.Asin] = 0
	}
	var validAsins []string
	for _, asin := range asins {
		_, exist := excludeAsins[asin]
		if !exist {
			validAsins = append(validAsins, asin)
		}
	}
	return validAsins
}

func getAsins(dummy bool) ([]*TitleAsin, Asins, int) {
	titleAsins := []*TitleAsin{}
	if dummy {
		asins := []string {"4799210300","4041055342","4088810716","4063882543","4091895565","478596006X","B06XPZ7VZ1","B06ZYBPLSB","4772811559","B06XPYHDTY","4829685905","4087925161","4758066507","4063959376","B06ZXZXJN1","B071V54538","B06ZY98NMC","B071H6KKBJ"}
		//asins := []string {"4820759728"}
		for _, asin := range asins {
			titleAsin := &TitleAsin{}
			titleAsin.AddAsins([]string{asin})
			titleAsins = append(titleAsins, titleAsin)
		}
	} else {
		page := GetOneLastUpdateBookPage() + 1
		log.Println("Start page:", page)
		for ; ; page++ {
			url := getUrl(page)
			urlTitleAsins, err := getUrlAsins(url)
			titleAsins = append(titleAsins, urlTitleAsins...)
			log.Printf("URL:%s Books:%d\n", url, len(urlTitleAsins))
			if len(urlTitleAsins) == 0 {
				if err == nil {
					SetOneLastUpdateBookPage(0)
				} else {
					SetOneLastUpdateBookPage(page - 1)
				}
				break
			}
		}
	}
	asins := Asins{}
	for _, titleAsin := range titleAsins {
		asins = append(asins, *titleAsin...)
	}
	asins = unique(asins)
	fetchAsinCount := len(asins)
	asins = validAsins(asins)
	log.Printf("Fetch asin:%d. Update asin:%d", fetchAsinCount, len(asins))
	return titleAsins, asins, fetchAsinCount
}

func getXml(asins []string) (xmlString string, err error)  {
	var api amazonproduct.AmazonProductAPI
	api.AccessKey = os.Getenv("MANGANOW_AMAZON_ACCESS_KEY")
	api.SecretKey = os.Getenv("MANGANOW_AMAZON_SECRET_KEY")
	api.AssociateTag = os.Getenv("MANGANOW_AMAZON_ASSOCIATE_TAG")
	api.Host = "ecs.amazonaws.jp"
	xmlString, err = api.ItemSearch("", map[string]string{
		"Operation": "ItemLookup",
		"IdType": "ASIN",
		"ItemId": strings.Join(asins, ","),
		"ResponseGroup": "Medium",
		"ItemPage": "1",
	})
	time.Sleep(2 * time.Second)
	return xmlString, err
}

func get10BooksInfo(asins Asins) ([]Book, int) {
	books := []Book{}
	const tryMax = 10
	countTry := 0
	for ; countTry < tryMax; countTry++ {
		xmlString, err := getXml(asins)
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(xmlString)))
		if err != nil {
			continue
		}
		items := doc.Find("Item")
		if items.Length() > 0 {
			items.Each(func(_ int, item *goquery.Selection) {
				book := Book{}
				book.Asin = item.Find("ASIN").Text()
				if asins.Exist(book.Asin) {
					book.Region = "JP"
					attributes := item.Find("ItemAttributes").First()
					if attributes.Length() > 0 {
						book.PublishType = attributes.Find("Binding").Text()
						book.Name = CleanName(attributes.Find("Title").Text())
						book.Publisher.Name = CleanName(attributes.Find("Publisher").Text())
						book.Author.Name = CleanName(attributes.Find("Author").Text())
						if book.Publisher.Name == "" || book.Author.Name == "" {
							return
						}
						dateStr := attributes.Find("PublicationDate").Text()
						jst, _ := time.LoadLocation("Asia/Tokyo")
						datePublish, _ := time.ParseInLocation("2006-01-02", dateStr, jst)
						book.DatePublish = fmt.Sprintf("%d%02d%02d", datePublish.Year(), datePublish.Month(), datePublish.Day())
					}
					smallImage := item.Find("SmallImage").First()
					book.ImageS_Url = smallImage.Find("URL").Text()
					book.ImageS_Width, _ = strconv.Atoi(smallImage.Find("Width").Text())
					book.ImageS_Height, _ = strconv.Atoi(smallImage.Find("Height").Text())
					mediumImage := item.Find("MediumImage").First()
					book.ImageM_Url = mediumImage.Find("URL").Text()
					book.ImageM_Width, _ = strconv.Atoi(mediumImage.Find("Width").Text())
					book.ImageM_Height, _ = strconv.Atoi(mediumImage.Find("Height").Text())
					largeImage := item.Find("LargeImage").First()
					book.ImageL_Url = largeImage.Find("URL").Text()
					book.ImageL_Width, _ = strconv.Atoi(largeImage.Find("Width").Text())
					book.ImageL_Height, _ = strconv.Atoi(largeImage.Find("Height").Text())
					xml, err := item.Html()
					if err == nil {
						book.Xml.String = "<item>" + xml + "</item>"
						book.Xml.Valid = true
					}
					books = append(books, book)
				}
			})
		} else {
			continue
		}
		break
	}
	return books, countTry + 1
}

func getBooksInfo(asins []string) []Book {
	const maxAsins = 10
	asinsCount := len(asins)
	updateCount := 0
	totalTryCount := 0
	allUpdatedBooks := []Book{}
	for i := 0; i < asinsCount; i += maxAsins {
		num := maxAsins
		if i + num > asinsCount {
			num = asinsCount - i
		}
		updatedBooksInfos, tryCount := get10BooksInfo(asins[i:i+num])
		allUpdatedBooks = append(allUpdatedBooks, updatedBooksInfos...)
		log.Printf("%d / %d: tryCount: %d\n", i, asinsCount, tryCount)
		updateCount++
		totalTryCount += tryCount
	}
	log.Printf("Getting try average: %f", (float32(totalTryCount) / float32(updateCount)))
	return allUpdatedBooks
}

func updateDb(books []Book) {
	records := []elastic.AsinRecord{}
	for _, book := range books {
		var publisher Publisher
		if db.ORM.Where(&Publisher{Name: book.Publisher.Name}).First(&publisher).RecordNotFound() {
			publisher.Name = book.Publisher.Name
			db.ORM.Create(&publisher)
		}
		book.PublisherID.Int64 = publisher.ID
		book.PublisherID.Valid = true

		var author Author
		if db.ORM.Where(&Author{Name: book.Author.Name}).First(&author).RecordNotFound() {
			author.Name = book.Author.Name
			db.ORM.Create(&author)
		}
		book.AuthorID.Int64 = author.ID
		book.AuthorID.Valid = true

		var existBook Book
		if db.ORM.Where(&Book{Asin: book.Asin}).First(&existBook).RecordNotFound() {
			db.ORM.Set("gorm:save_associations", false).Create(&book)
		} else {
			book.ID = existBook.ID
			book.CreatedAt = existBook.CreatedAt
			db.ORM.Set("gorm:save_associations", false).Save(&book)
		}

		records = append(
			records,
			elastic.AsinRecord{
				elastic.AsinParam{
					book.Name,
					book.Publisher.Name,
					book.Author.Name,
					"",
					book.DatePublishTime()},
				book.Asin})
	}

	log.Println("Elasticsearch Added index:", len(records))
	updatedIndex := elastic.BulkAsinIndex(records)
	log.Println("Elasticsearch Updated index:", updatedIndex)
}

func initLog() *os.File {
	filePath := os.Getenv("MANGANOW_UPDATING_BOOK_LOG_FILE")
	f, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Error opening:%v", err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	return f
}

func updateTitleID(titleAsins []*TitleAsin) {
	for _, titleAsin := range titleAsins {
		if len(*titleAsin) == 0 {
			//log.Panic("updateTitleID(): len(*titleAsin) == 0")
			continue
		}

		// Search existing title_id
		var titleID int64 = 0
		books := []Book{}
		if db.ORM.Where("asin IN (?)", *titleAsin).Find(&books).RecordNotFound() {
			log.Panicf("updateTitleID(): Book does not exist. Asins:%v", *titleAsin)
		} else {
			if len(books) == 0 {
				// Fetch asin in web, but doesn't fetch xml by API
				continue
			}
			if books[0].TitleID.Valid {
				titleID = books[0].TitleID.Int64
			}
		}
		if titleID == 0 {
			// Create title
			title := Title{}
			db.ORM.Create(&title)
			titleID = title.ID
		}

		// Remove old book's title_id
		db.ORM.Table("books").Where("title_id = ?", titleID).UpdateColumn("title_id", nil)
		// Set title_id
		db.ORM.Table("books").Where("asin IN (?)", *titleAsin).UpdateColumn("title_id", titleID)
	}
}

func updateLog(fetchAsinCount, fetchTitleCount, updateAsinCount, updatedBookCount int) {
	var log UpdateLog
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if db.ORM.Where("date = ?", date).Find(&log).RecordNotFound() {
		log.Date = date
		log.FetchAsinCount = fetchAsinCount
		log.FetchTitleCount = fetchTitleCount
		log.UpdateAsinCount = updateAsinCount
		log.UpdatedBookCount = updatedBookCount
		db.ORM.Create(&log)
	} else {
		log.FetchAsinCount += fetchAsinCount
		log.FetchTitleCount += fetchTitleCount
		log.UpdateAsinCount += updateAsinCount
		log.UpdatedBookCount += updatedBookCount
		db.ORM.Save(&log)
	}
}

func main() {
	logFile := initLog()
	defer logFile.Close()

	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	titleAsins, asins, fetchAsinCount := getAsins(false)
	log.Printf("TitleAsins:%d asins:%d", len(titleAsins), len(asins))
	books := getBooksInfo(asins)
	log.Printf("Book Info:%d", len(books))
	updateDb(books)
	updateTitleID(titleAsins)

	userbooks.Run()

	updateLog(fetchAsinCount, len(titleAsins), len(asins), len(books))

	log.Printf("fetchAsinCount: %d\n", fetchAsinCount)
	log.Printf("updateAsinCount: %d\n", len(asins))
	log.Printf("updatedBookCount: %d\n", len(books))
}
