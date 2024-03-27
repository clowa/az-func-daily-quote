package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/clowa/az-func-daily-quote/src/lib/config"
	quotable "github.com/clowa/az-func-daily-quote/src/lib/quotableSdk"
	quote "github.com/clowa/az-func-daily-quote/src/lib/quote"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultContextTimeout = 10 * time.Second
)

// Connects to the CosmosDB instance and returns a container client. Configuration in loaded from the environment.
func connect() *mongo.Client {
	config := config.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.MongodbConnectionString)
	c, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Warnf("unable to initialize connection %v", err)
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		log.Warnf("unable to connect %v", err)
	}

	return c
}

func writeQuoteToDatabase(q *quote.Quote) error {
	config := config.GetConfig()
	client := connect()
	ctx := context.Background()
	defer client.Disconnect(ctx)

	collection := client.Database(config.MongodbDatabase).Collection(config.MongodbCollection)
	r, err := collection.InsertOne(ctx, &q)
	if err != nil {
		return err
	}

	log.Infof("Inserted quote with ID %s", r.InsertedID)

	return nil
}

func getQuoteFromDatabase(creationDate string) (quote.Quote, error) {
	config := config.GetConfig()
	client := connect()
	ctx := context.Background()
	defer client.Disconnect(ctx)

	collection := client.Database(config.MongodbDatabase).Collection(config.MongodbCollection)
	filter := bson.D{{Key: "creationdate", Value: creationDate}}
	results, err := collection.Find(ctx, filter)
	if err != nil {
		return quote.Quote{}, err
	}

	var quotes []quote.Quote
	if err = results.All(ctx, &quotes); err != nil {
		return quote.Quote{}, err
	}

	if len(quotes) == 0 {
		return quote.Quote{}, fmt.Errorf("no quotes found for creation date %s", creationDate)
	}

	quote := quotes[0]

	return quote, nil
}

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var quoteOfTheDay quote.Quote

	today := time.Now().Format("2006-01-02")
	quoteOfTheDay, err := getQuoteFromDatabase(today)
	// quoteOfTheDay, err := Quote{}, fmt.Errorf("no quote found")

	if quoteOfTheDay.Length == 0 || err != nil {
		log.Warnf("Error getting quote from database: %s", err)
		log.Info("Fetching quote from quotable API")
		quotes, err := quotable.GetRandomQuote(quotable.GetRandomQuoteQueryParams{Limit: 1, Tags: []string{"technology"}})
		if err != nil {
			log.Warnf("Error fetching quote from quotable API: %s", err)
			http.Error(w, "Failed to fetch new quote", http.StatusInternalServerError)
			return
		}
		quote := quotes[0]

		// Write quote to database
		quoteOfTheDay.LoadFromQuotable(&quote)
		err = writeQuoteToDatabase(&quoteOfTheDay)
		if err != nil {
			log.Warnf("Error writing quote to database: %s", err)
		}
	}

	log.Infof("Quote of the day: %s by %s", quoteOfTheDay.Content, quoteOfTheDay.Author)

	// Write response
	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(quoteOfTheDay)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes.Bytes())
	w.WriteHeader(http.StatusOK)
}
