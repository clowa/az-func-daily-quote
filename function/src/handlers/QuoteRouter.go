package handlers

import (
	"net/http"
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
	CreationDate string   `json:"creationdate"`
}

func (q *Quote) Load(i *quotable.QuoteResponse) {
	q.Id = i.Id
	q.Content = i.Content
	q.Author = i.Author
	q.AuthorSlug = i.AuthorSlug
	q.Length = i.Length
	q.Tags = i.Tags

	today := time.Now().Format("2006-01-02")
	q.CreationDate = today
}

func QuoteRouter(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		getQuoteHandler(w, r)
	case r.Method == "POST":
		createQuoteHandler(w, r)
	}
}
