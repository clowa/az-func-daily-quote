package handlers

import (
	"net/http"

	quote "github.com/clowa/az-func-daily-quote/src/lib/quote"
	log "github.com/sirupsen/logrus"
)

func createQuoteHandler(w http.ResponseWriter, r *http.Request) {
	// Validate if request is of type application/json
	if r.Header.Get("Content-Type") != "application/json" {
		log.Info("Received invalid Content-Type header")
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Load quote from request
	quote, err := quote.New().LoadFromRequest(r)
	if err != nil {
		log.Warnf("Unable to load quote from request %v", err)
		http.Error(w, "Unable to load quote from request", http.StatusBadRequest)
		return
	}
	log.Infof("Received quote: %s by %s", quote.Content, quote.Author)

	log.Info("Saving quote to database.")
	if err := writeQuoteToDatabase(quote); err != nil {
		log.Warnf("Unable to write quote to database %v", err)
		http.Error(w, "Unable to write quote to database", http.StatusInternalServerError)
		return
	}
}
