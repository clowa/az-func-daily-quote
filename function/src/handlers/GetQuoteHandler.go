package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	quoteOfTheDay, err := getQuoteOfTheDay()
	if quoteOfTheDay.Length == 0 || err != nil {
		http.Error(w, "Failed to fetch new quote", http.StatusInternalServerError)
		return
	}

	log.Infof("Quote of the day: %s by %s", quoteOfTheDay.Content, quoteOfTheDay.Author)

	// Write response
	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(quoteOfTheDay)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes.Bytes())
	w.WriteHeader(http.StatusOK)
}

func getQuoteOfTheDay() (quote.Quote, error) {
	var quoteOfTheDay quote.Quote

	today := time.Now().Format("2006-01-02")
	quoteOfTheDay, err := getQuoteFromDatabase(today)
	if quoteOfTheDay.Length != 0 || err == nil {
		return quoteOfTheDay, nil
	}
	log.Warnf("Error getting quote from database: %s", err)

	quoteOfTheDay, err = getQuoteFromQuotable(true)
	if quoteOfTheDay.Length != 0 || err == nil {
		return quoteOfTheDay, nil
	}
	log.Warnf("Error getting quote from quotable API: %s", err)

	quoteOfTheDay, err = getRandomQuoteFromDatabase()
	if quoteOfTheDay.Length != 0 || err == nil {
		return quoteOfTheDay, nil
	}
	log.Warnf("Error getting random quote from database: %s", err)

	return quote.Quote{}, fmt.Errorf("failed to fetch new quote")
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
	defer results.Close(ctx)

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

// Gets a random quote from the database by picking a random id of all existing quotes
func getRandomQuoteFromDatabase() (quote.Quote, error) {
	config := config.GetConfig()
	client := connect()
	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()
	defer client.Disconnect(ctx)

	// Get all available ids
	// TODO: Limit the number of ids to avoid performance issues
	collection := client.Database(config.MongodbDatabase).Collection(config.MongodbCollection)
	projection := bson.M{"_id": 1}
	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetProjection(projection))
	if err != nil {
		return quote.Quote{}, fmt.Errorf("error getting all ids: %s", err)
	}
	defer cursor.Close(ctx)

	// Store all ids in an array
	var ids []interface{}
	for ok := cursor.Next(ctx); ok; ok = cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Errorf("error decoding id: %s", err)
		}
		ids = append(ids, result["_id"])
	}
	if err := cursor.Err(); err != nil {
		return quote.Quote{}, fmt.Errorf("error iterating over ids: %s", err)
	}

	// Fail if no ids were found
	if len(ids) == 0 {
		return quote.Quote{}, fmt.Errorf("no quotes found")
	}

	// Pick a random id
	randomIndex := rand.Intn(len(ids))
	chosenId := ids[randomIndex]

	var q quote.Quote
	if err := collection.FindOne(ctx, bson.M{"_id": chosenId}).Decode(&q); err != nil {
		return quote.Quote{}, fmt.Errorf("error getting quote by id: %s", err)
	}

	// Update the "creationDate" field of the selected document to the current date
	if err = q.SaveToDatabase(client, config); err != nil {
		log.Warnf("Error updating creation date of quote with id %s: %s", q.Id, err)
	}

	return q, nil
}

func getQuoteFromQuotable(writeToDatabase bool) (quote.Quote, error) {
	var quoteOfTheDay quote.Quote

	quotes, err := quotable.GetRandomQuote(quotable.GetRandomQuoteQueryParams{Limit: 1, Tags: []string{"technology"}})
	if err != nil {
		return quote.Quote{}, fmt.Errorf("error fetching quote from quotable API: %s", err)
	}
	q := quotes[0]
	quoteOfTheDay.LoadFromQuotable(&q)

	if writeToDatabase {
		err = writeQuoteToDatabase(&quoteOfTheDay)
		if err != nil {
			log.Warnf("Error writing quote to database: %s", err)
		}
	}

	return quoteOfTheDay, nil
}

func writeQuoteToDatabase(q *quote.Quote) error {
	config := config.GetConfig()
	client := connect()
	return q.SaveToDatabase(client, config)
}
