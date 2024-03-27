package handlers

import (
	"net/http"
)

func QuoteRouter(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		getQuoteHandler(w, r)
	case r.Method == "POST":
		createQuoteHandler(w, r)
	}
}
