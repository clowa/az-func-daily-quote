package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Quote struct {
	ID         string   `json:"_id"`
	Content    string   `json:"content"`
	Author     string   `json:"author"`
	AuthorSlug string   `json:"authorSlug"`
	Length     int      `json:"length"`
	Tags       []string `json:"tags"`
}

func QuoteHandler(w http.ResponseWriter, r *http.Request) {
	quoteOfTheDay := getQuoteOfTheDay()

	fmt.Printf("Quote of the day: %s by %s", quoteOfTheDay.Content, quoteOfTheDay.Author)

	responseBodyBytes := new(bytes.Buffer)
	json.NewEncoder(responseBodyBytes).Encode(quoteOfTheDay)

	w.Write(responseBodyBytes.Bytes())
	w.WriteHeader(http.StatusOK)
}

func getQuoteOfTheDay() Quote {
	resp, err := http.Get("https://api.quotable.io/quotes/random?tags=technology,famous-quotes&limit=1")
	if err != nil {
		fmt.Printf("Failed to fetch quote. Request failed with error: %s", err)
	}
	defer resp.Body.Close()

	var quote []Quote
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		msg := "Error decoding Quote"
		log.Printf("%s: %v", msg, err)
	}

	return quote[0]
}
