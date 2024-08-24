package quote

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/clowa/az-func-daily-quote/src/lib/config"
	quotable "github.com/clowa/az-func-daily-quote/src/lib/quotableSdk"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultContextTimeout = 10 * time.Second
)

// Struct representing a quote
type Quote struct {
	InstanceId   primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Id           string             `bson:"id" json:"id"`
	Content      string             `bson:"content" json:"content"`
	Author       string             `bson:"author" json:"author"`
	AuthorSlug   string             `bson:"authorSlug" json:"authorSlug"`
	Length       int                `bson:"length" json:"length"`
	Tags         []string           `bson:"tags" json:"tags"`
	CreationDate string             `bson:"creationDate" json:"creationDate"`
}

func New() *Quote {
	q := &Quote{}
	q.setId()
	q.setCreationDate()

	return q
}

func (q *Quote) setAll() {
	q.setId()
	q.setCreationDate()
	q.setLength()
	q.setAuthorSlug()
}

// GenerateRandomID generates a random ID of the specified length
func generateRandomId(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Int63n(int64(len(charset)))]
	}
	return string(b)
}

func (q *Quote) SaveToDatabase(client *mongo.Client, config *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()

	q.setCreationDate()
	collection := client.Database(config.MongodbDatabase).Collection(config.MongodbCollection)

	i, err := collection.CountDocuments(ctx, bson.M{"id": q.Id})
	if err != nil {
		return fmt.Errorf("unable to check if quote exists in database: %v", err)
	}

	if i > 0 {
		r, err := collection.ReplaceOne(ctx, bson.M{"id": q.Id}, q)
		if err != nil || r.MatchedCount == 0 {
			return fmt.Errorf("unable to update quote with id %s: %v", q.Id, err)
		}
		log.Infof("Successfully updated quote with %s", q.InstanceId.String())

		return nil
	}

	r, err := collection.InsertOne(ctx, &q)
	if err != nil {
		return err
	}
	log.Infof("Inserted quote with ID %s", r.InsertedID)

	return nil
}

func (q *Quote) setCreationDate() {
	today := time.Now().Format("2006-01-02")
	q.CreationDate = today
}

func (q *Quote) setId() {
	q.Id = generateRandomId(11)
}

func (q *Quote) setLength() {
	q.Length = len(q.Content)
}

func (q *Quote) setAuthorSlug() {
	lower := strings.ToLower(q.Author)
	q.AuthorSlug = strings.ReplaceAll(lower, " ", "-")
}

func (q *Quote) LoadFromQuotable(i *quotable.QuoteResponse) {
	q.Id = i.Id
	q.Content = i.Content
	q.Author = i.Author
	q.AuthorSlug = i.AuthorSlug
	q.Length = i.Length
	q.Tags = i.Tags

	q.setCreationDate()
}

func (q *Quote) LoadFromRequest(r *http.Request) (*Quote, error) {
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		return nil, fmt.Errorf("unable to decode request body into json: %v", err)
	}

	// Set computed properties of the quote
	q.setAll()
	return q, nil
}
