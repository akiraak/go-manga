package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"os"
	"time"
)

type AsinParam struct {
	Title		string		`json:"title"`
	Publisher	string		`json:"publisher"`
	Author		string		`json:"author"`
	AllText		string		`json:"all_text"`
	DatePublish	time.Time	`json:"date_publish"`
}

func (a *AsinParam)MakeAllText() {
	a.AllText = fmt.Sprintf("%s %s %s", a.Title, a.Publisher, a.Author)
}

type AsinRecord struct {
	AsinParam
	Asin		string
}

func newClient() (context.Context, *elastic.Client, error) {
	endpoint := os.Getenv("MANGANOW_ELASTICSEARCH_ENDPOINT")
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(endpoint),
		elastic.SetSniff(false))
	if err != nil {
		return nil, nil, err
	}
	return ctx, client, nil
}

func BulkAsinIndex(records []AsinRecord) int {
	fmt.Printf("BulkAsinIndex: total=%d\n", len(records))
	updatedIndex := 0
	max := 20000
	for i := 0; ; i++ {
		start := (i * max)
		end := start + max
		if end >= len(records) {
			end = len(records)
		}
		updateRecords := records[start:end]
		for i := 0; i < len(updateRecords); i++ {
			updateRecords[i].MakeAllText()
		}

		ctx, client, err := newClient()
		if err != nil {
			return updatedIndex
		}

		bulkRequest := client.Bulk()
		for _, record := range updateRecords {
			req := elastic.NewBulkIndexRequest().
				Index("asins").
				Type("asin").
				Id(record.Asin).
				Doc(record.AsinParam)
			bulkRequest = bulkRequest.Add(req)
		}

		bulkResponse, err := bulkRequest.Do(ctx)
		if err != nil {
			return updatedIndex
		}
		updatedIndex += len(bulkResponse.Indexed())

		fmt.Printf("%d%% : start:%d end:%d\n", (i * max * 100 / len(records)), start, end)
		if end >= len(records) {
			break
		}
	}

	return updatedIndex
}

func SearchAsins(keyword string, offset int, limit int) ([]AsinRecord, int64) {
	results := []AsinRecord{}
	hitTotal := int64(0)
	ctx, client, err := newClient()
	if err != nil {
		return results, 0
	}
	query := elastic.NewMatchQuery("all_text", keyword).
		Operator("and")
	searchResult, err := client.Search().
		Index("asins").
		Type("asin").
		Query(query).
		Sort("date_publish", false).
		From(offset).Size(limit).
		Do(ctx)
	if err == nil {
		hitTotal = searchResult.Hits.TotalHits
		if searchResult.Hits.TotalHits > 0 {
			for _, hit := range searchResult.Hits.Hits {
				var a AsinParam
				err := json.Unmarshal(*hit.Source, &a)
				if err != nil {
					continue
				}
				results = append(results, AsinRecord{a, hit.Id})
			}
		}
	}
	return results, hitTotal
}

func SearchUserAsins(keywords []string, offset int, limit int) ([]AsinRecord, int64) {
	results := []AsinRecord{}
	hitTotal := int64(0)
	ctx, client, err := newClient()
	if err != nil {
		return results, 0
	}
	matches := []elastic.Query{}
	for _, keyword := range keywords {
		m := elastic.NewMatchPhraseQuery("all_text", keyword)
		matches = append(matches, m)
	}
	query := elastic.NewBoolQuery().Should(matches...)
	searchResult, err := client.Search().
		Index("asins").
		Type("asin").
		Query(query).
		Sort("date_publish", false).
		From(offset).Size(limit).
		Do(ctx)
	if err == nil {
		hitTotal = searchResult.Hits.TotalHits
		if searchResult.Hits.TotalHits > 0 {
			for _, hit := range searchResult.Hits.Hits {
				var a AsinParam
				err := json.Unmarshal(*hit.Source, &a)
				if err != nil {
					continue
				}
				results = append(results, AsinRecord{a, hit.Id})
			}
		}
	}
	return results, hitTotal
}
