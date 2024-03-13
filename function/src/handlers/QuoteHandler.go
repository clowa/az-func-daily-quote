package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	quotable "github.com/clowa/az-func-daily-quote/src/lib/quotableSdk"
	log "github.com/sirupsen/logrus"
)

// Struct representing the structure returned from the quotable API
type ApiQuote struct {
	ID         string   `json:"_id"`
	Content    string   `json:"content"`
	Author     string   `json:"author"`
	AuthorSlug string   `json:"authorSlug"`
	Length     int      `json:"length"`
	Tags       []string `json:"tags"`
}

func QuoteHandler(w http.ResponseWriter, r *http.Request) {

	quotes, err := quotable.GetRandomQuote(quotable.GetRandomQuoteQueryParams{Limit: 1, Tags: []string{"technology", "famous-quotes"}})
	if err != nil {
		handleWarn(w, err)
	}
	quoteOfTheDay := quotes[0]

	log.Infof("Quote of the day: %s by %s", quoteOfTheDay.Content, quoteOfTheDay.Author)

	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(quoteOfTheDay)

	w.Write(responseBodyBytes.Bytes())
	w.WriteHeader(http.StatusOK)
}
