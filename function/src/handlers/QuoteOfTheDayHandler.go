package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/clowa/az-func-daily-quote/src/lib/config"
	quotable "github.com/clowa/az-func-daily-quote/src/lib/quotableSdk"
	log "github.com/sirupsen/logrus"
)

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

// Connects to the CosmosDB instance and returns a container client. Configuration in loaded from the environment.
func connect() *azcosmos.ContainerClient {
	config := config.GetConfig()

	credential, err := azidentity.NewManagedIdentityCredential(nil)
	if err != nil {
		log.Warnf("Error creating managed identity credential: %s", err)
	}

	client, err := azcosmos.NewClient(config.CosmosHost, credential, nil)
	if err != nil {
		log.Warnf("Error creating Cosmos client: %s", err)
	}

	database, err := client.NewDatabase(config.CosmosDatabase)
	if err != nil {
		log.Warnf("Error creating Cosmos database: %s", err)
	}

	container, err := database.NewContainer(config.CosmosContainer)
	if err != nil {
		log.Warnf("Error creating Cosmos container: %s", err)
	}

	return container
}

func writeQuoteToDatabase(q *Quote) error {
	container := connect()

	partitionKey := azcosmos.NewPartitionKeyString(q.AuthorSlug)

	bytes, err := json.Marshal(q)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err = container.UpsertItem(ctx, partitionKey, bytes, nil) // ToDo: change to CreateItem()
	if err != nil {
		return err
	}

	return nil
}

// func getQuoteFromDatabase(creationDate string) (Quote, error) {
// 	config := config.GetConfig()
// 	container := connect()

// 	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
// 	defer cancel()

// 	query := fmt.Sprintf("SELECT * FROM %s q WHERE q.creationDate = \"%s\" OFFSET 0 LIMIT 1", config.CosmosContainer, creationDate)
// 	partitionKey := azcosmos.NewPartitionKeyString("albert-einstein")
// 	opt := azcosmos.QueryOptions{
// 		PageSizeHint: 1,
// 	}
// 	pager := container.NewQueryItemsPager(query, partitionKey, &opt)

// 	var quotes []Quote
// 	for pager.More() {
// 		queryResponse, err := pager.NextPage(ctx)
// 		if err != nil {
// 			log.Warnf("Error querying database: %s", err)
// 		}

// 		for _, item := range queryResponse.Items {
// 			json.Unmarshal(item, &quotes)
// 		}
// 	}
// 	quote := quotes[0]
// 	log.Printf("Query response: %v", quote)

// 	return quote, nil
// }

func QuoteOfTheDayHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := quotable.GetRandomQuote(quotable.GetRandomQuoteQueryParams{Limit: 1, Tags: []string{"technology"}})
	if err != nil {
		handleWarn(w, err)
	}
	quoteOfTheDay := quotes[0]

	log.Infof("Quote of the day: %s by %s", quoteOfTheDay.Content, quoteOfTheDay.Author)

	// Write quote to database
	databaseQuote := Quote{}
	databaseQuote.Load(&quoteOfTheDay)
	err = writeQuoteToDatabase(&databaseQuote)
	if err != nil {
		log.Warnf("Error writing quote to database: %s", err)
	}

	// Write response
	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(quoteOfTheDay)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes.Bytes())
	w.WriteHeader(http.StatusOK)
}
