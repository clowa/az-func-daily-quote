package quote

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	quotable "github.com/clowa/az-func-daily-quote/src/lib/quotableSdk"
)

// Struct representing a quote
type Quote struct {
	Id           string   `json:"id"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	AuthorSlug   string   `json:"authorSlug"`
	Length       int      `json:"length"`
	Tags         []string `json:"tags"`
	CreationDate string   `json:"creationDate"`
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

func (q *Quote) setId() {
	q.Id = generateRandomId(11)
}

func (q *Quote) setCreationDate() {
	today := time.Now().Format("2006-01-02")
	q.CreationDate = today
}

func (q *Quote) setLength() {
	q.Length = len(q.Content)
}

func (q *Quote) setAuthorSlug() {
	lower := strings.ToLower(q.Author)
	q.AuthorSlug = strings.ReplaceAll(lower, " ", "-")
}

func (q *Quote) LoadFromQuotable(i *quotable.QuoteResponse) *Quote {
	q.Id = i.Id
	q.Content = i.Content
	q.Author = i.Author
	q.AuthorSlug = i.AuthorSlug
	q.Length = i.Length
	q.Tags = i.Tags

	q.setCreationDate()
	return q
}

func (q *Quote) LoadFromRequest(r *http.Request) (*Quote, error) {
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		return nil, fmt.Errorf("unable to decode request body into json: %v", err)
	}

	// Set computed properties of the quote
	q.setAll()
	return q, nil
}
