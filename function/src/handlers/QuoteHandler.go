package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/clowa/az-func-daily-quote/src/lib/config"
	quotable "github.com/clowa/az-func-daily-quote/src/lib/quotableSdk"
	log "github.com/sirupsen/logrus"
)

type creationInfo struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// Struct representing the structure returned from the quotable API
type Quote struct {
	Id           string   `json:"id"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	AuthorSlug   string   `json:"authorSlug"`
	Length       int      `json:"length"`
	Tags         []string `json:"tags"`
	CreationDate string   `json:"creationDate"`
}

func (q *Quote) Load(i *quotable.QuoteResponse) {
	q.Id = i.Id
	q.Content = i.Content
	q.Author = i.Author
	q.AuthorSlug = i.AuthorSlug
	q.Length = i.Length
	q.Tags = i.Tags

	now := time.Now()
	q.CreationDate = fmt.Sprintf("%d-%d-%d", now.Year(), int(now.Month()), now.Day())
}

func writeQuoteToDatabase(q *Quote) error {
	config := config.GetConfig()

	client, err := azcosmos.NewClientFromConnectionString(config.CosmosConnectionString, nil)
	if err != nil {
		return err
	}

	database, err := client.NewDatabase(config.CosmosDatabase)
	if err != nil {
		return err
	}

	container, err := database.NewContainer(config.CosmosContainer)
	if err != nil {
		return err
	}

	partitionKey := azcosmos.NewPartitionKeyString(q.AuthorSlug)

	ctx := context.TODO()
	bytes, err := json.Marshal(q)
	if err != nil {
		return err
	}

	_, err = container.UpsertItem(ctx, partitionKey, bytes, nil) // ToDo: change to CreateItem()
	if err != nil {
		return err
	}

	return nil
}

func QuoteHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := quotable.GetRandomQuote(quotable.GetRandomQuoteQueryParams{Limit: 1, Tags: []string{"technology", "famous-quotes"}})
	if err != nil {
		handleWarn(w, err)
	}
	quoteOfTheDay := quotes[0]

	log.Infof("Quote of the day: %s by %s", quoteOfTheDay.Content, quoteOfTheDay.Author)
	databaseQuote := Quote{}
	databaseQuote.Load(&quoteOfTheDay)
	err = writeQuoteToDatabase(&databaseQuote)
	if err != nil {
		handleWarn(w, err)
	}

	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(quoteOfTheDay)

	w.Write(responseBodyBytes.Bytes())
	w.WriteHeader(http.StatusOK)
}
