package main

import (
	"errors"
	"os"
	"io"
	"log"
	"strings"
	"strconv"
	"fmt"
	"time"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"github.com/DDRBoxman/go-amazon-product-api"
	"github.com/akiraak/go-manga/db"
	. "github.com/akiraak/go-manga/model"
)

const updateBookInterval = time.Duration(12) * time.Hour
//const updateBookInterval = time.Duration(1) * time.Minute
var url2AsinReg = regexp.MustCompile(`/dp/(\w+)`)

func url2Asin(url string) (string, error) {
	result := url2AsinReg.FindAllStringSubmatch(url, 1)
	if len(result) > 0 {
		return result[0][1], nil
	}
	return "", errors.New("URL does not include asin.")
}

type Asin struct {
	Type string
	Asin string
	SubAsins []string
}

func (a Asin) SubAsinsColString() string {
	return strings.Join(a.SubAsins, ",")
}

func getUrlAsins(url string) []Asin {
	var bookAsins []Asin
	//bookUrls := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return bookAsins
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

		// Make main and sub asin
		if len(asins) > 0 {
			// Main book
			subAsins := []string{}
			if len(asins) > 1 {
				subAsins = asins[1:]
			}
			bookAsins = append(bookAsins, Asin{Type:"main", Asin:asins[0], SubAsins:subAsins})

			// Sub book
			if len(asins) > 1 {
				for _, asin := range asins[1:] {
					bookAsins = append(bookAsins, Asin{Type:"sub", Asin:asin})
				}
			}
		}
	})
	return bookAsins
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
		"https://www.amazon.co.jp/s/ref=sr_pg_2?fst=as%3Aoff&rh=n%3A465392%2Cn%3A%21465610%2Cn%3A466280%2Cn%3A2278488051%2Cp_n_publication_date%3A2285539051%7C2315443051&bbn=2278488051&ie=UTF8&qid=1491933854",
		page)
}

func validAsins(asins []Asin) []Asin {
	checkTime := time.Now()
	checkTime = checkTime.Add(-updateBookInterval)
	checkAsins := []string{}
	for _, asin := range asins {
		checkAsins = append(checkAsins, asin.Asin)
	}
	var excludeBooks []Book
	db.ORM.Where("asin in (?)", checkAsins).Where("updated_at > ?", checkTime).Find(&excludeBooks)
	excludeAsins := []string{}
	for _, book := range excludeBooks {
		excludeAsins = append(excludeAsins, book.Asin)
	}
	var validAsins []Asin
	for _, asin := range asins {
		valid := true
		for _, excludeAsin := range excludeAsins {
			if asin.Asin == excludeAsin {
				valid = false
				//fmt.Printf("exclude: %s\n", excludeAsin)
				break
			}
		}
		if valid {
			validAsins = append(validAsins, asin)
		}
	}
	return validAsins
}

func getAsins(dummy bool) []Asin {
	booksAsins := []Asin{}
	if dummy {
		asins := []string {"4799210300","4041055342","4088810716","4063882543","4091895565","478596006X","B06XPZ7VZ1","B06ZYBPLSB","4772811559","B06XPYHDTY","4829685905","4087925161","4758066507","4063959376","B06ZXZXJN1","B071V54538","B06ZY98NMC","B071H6KKBJ"}
		for _, asin := range asins {
			booksAsins = append(booksAsins, Asin{Type:"main", Asin:asin})
		}
	} else {
		for page := 1; ; page++ {
			url := getUrl(page)
			urlAsins := getUrlAsins(url)
			booksAsins = append(booksAsins, urlAsins...)
			log.Printf("URL:%s Books:%d\n", url, len(urlAsins))
			if len(urlAsins) == 0 {
				break
			}
		}
	}
	booksAsins = validAsins(booksAsins)
	return booksAsins
}

func getXml(asins []Asin) (xmlString string, err error)  {
	var api amazonproduct.AmazonProductAPI
	api.AccessKey = os.Getenv("MANGANOW_AMAZON_ACCESS_KEY")
	api.SecretKey = os.Getenv("MANGANOW_AMAZON_SECRET_KEY")
	api.AssociateTag = os.Getenv("MANGANOW_AMAZON_ASSOCIATE_TAG")
	api.Host = "ecs.amazonaws.jp"
	asinStrings := []string{}
	for _, asin := range asins {
		asinStrings = append(asinStrings, asin.Asin)
	}
	xmlString, err = api.ItemSearch("", map[string]string{
		"Operation": "ItemLookup",
		"IdType": "ASIN",
		"ItemId": strings.Join(asinStrings, ","),
		"ResponseGroup": "Medium",
		"ItemPage": "1",
	})
	time.Sleep(2 * time.Second)
	return xmlString, err
}

func getAsin(asins []Asin, asin string) (Asin, error) {
	// TODO: Change to method of struct
	for _, checkAsin := range asins {
		if checkAsin.Asin == asin {
			return checkAsin, nil
		}
	}
	return Asin{}, errors.New("Asin does not include Asin array.")
}

func get10BooksInfo(asins []Asin) ([]Book, int) {
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
				asin, err := getAsin(asins, book.Asin)
				if err == nil {
					book.TreeType = asin.Type
					book.SubAsinsCol = asin.SubAsinsColString()
					book.Region = "JP"
					// TODO: Get book.Kindle
					attributes := item.Find("ItemAttributes").First()
					if attributes.Length() > 0 {
						book.Title = attributes.Find("Title").Text()
						book.Publisher.Name = attributes.Find("Publisher").Text()
						book.Author.Name = attributes.Find("Author").Text()
						dateStr := attributes.Find("PublicationDate").Text()
					    jst, _ := time.LoadLocation("Asia/Tokyo")
						datePublish, _ := time.ParseInLocation("2006-01-02", dateStr, jst)
						book.DatePublish = datePublish
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
					books = append(books, book)
					//fmt.Printf("----\n")

					// TODO: xmlのテキストを取得する
				}
			})
		} else {
			continue
		}
		break
	}
	return books, countTry + 1
}

func getBooksInfo(asins []Asin) []Book {
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
	for _, book := range books {
		var publisher Publisher
		if db.ORM.Where(&Publisher{Name: book.Publisher.Name}).First(&publisher).RecordNotFound() {
			publisher.Name = book.Publisher.Name
			db.ORM.Create(&publisher)
			//fmt.Printf("[Publisher]Create: %+v\n", publisher)
		} else {
			//fmt.Printf("[Publisher]Exist : %+v\n", publisher)
		}
		book.PublisherID = publisher.ID

		var author Author
		if db.ORM.Where(&Author{Name: book.Author.Name}).First(&author).RecordNotFound() {
			author.Name = book.Author.Name
			db.ORM.Create(&author)
			//fmt.Printf("[Author]Create: %+v\n", author)
		} else {
			//fmt.Printf("[Author]Exist : %+v\n", author)
		}
		book.AuthorID = author.ID

		var existBook Book
		if db.ORM.Where(&Book{Asin: book.Asin}).First(&existBook).RecordNotFound() {
			db.ORM.Set("gorm:save_associations", false).Create(&book)
			//fmt.Printf("[Book]Create: [%s] %s\n", book.Asin, book.Title)
		} else {
			book.ID = existBook.ID
			book.CreatedAt = existBook.CreatedAt
			db.ORM.Set("gorm:save_associations", false).Save(&book)
			//fmt.Printf("[Book]Exist : [%s] %s\n", book.Asin, book.Title)
		}
		//fmt.Printf("%#v\n", book.Asin)
	}
}

func initLog() *os.File {
	filePath := os.Getenv("MANGANOW_UPDATING_BOOK_LOG_FILE")
	f, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Panic("error opening file: %v", err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	return f
}

func main() {
	db.ORM = db.InitDB()
	defer db.ORM.Close()
	//db.ORM.LogMode(true)

	logFile := initLog()
	defer logFile.Close()

	asins := getAsins(false)
	books := getBooksInfo(asins)
	updateDb(books)

	log.Printf("asins: %d\n", len(asins))
	log.Printf("books: %d\n", len(books))
}
