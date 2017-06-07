package elastic

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
	"os"
)

type AsinParam struct {
	Title		string	`json:"title"`
	Publisher	string	`json:"publisher"`
	Author		string	`json:"author"`
}

type AsinRecord struct {
	Asin		string
	AsinParam	AsinParam
}

func BulkAsinIndex(records []AsinRecord) int {
	endpoint := os.Getenv("MANGANOW_ELASTICSEARCH_ENDPOINT")
	updatedIndex := 0

	max := 20000
	for i := 0; ; i++ {
		start := (i * max)
		end := start + max
		if end >= len(records) {
			end = len(records)
		}
		updateRecords := records[start:end]

		ctx := context.Background()
		client, err := elastic.NewClient(
			elastic.SetURL(endpoint),
			elastic.SetSniff(false))
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

		if end >= len(records) {
			break
		}
	}

	return updatedIndex
}
