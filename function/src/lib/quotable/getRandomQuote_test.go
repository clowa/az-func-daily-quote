package quotable_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"

	"github.com/clowa/az-func-daily-quote/src/lib/quotable"
)

func TestGetRandomQuote(t *testing.T) {
	const apiEndpoint = "/quotes/random"

	// Create a mock server to simulate the Quotable API.
	// The response is read from a local JSON file for consistency.
	// The server performs basic validation on the incoming request.
	mockingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != apiEndpoint {
			t.Fatalf("Expected request to %s, got %s", apiEndpoint, r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Fatalf("Expected method GET, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.URL.Query().Get("tags") != "technology,famous-quotes" {
			t.Fatalf("Expected tags 'technology,famous-quotes', got %s", r.URL.Query().Get("tags"))
		}
		if r.URL.Query().Get("limit") != "2" {
			t.Fatalf("Expected limit '2', got %s", r.URL.Query().Get("limit"))
		}

		responseFile, err := os.ReadFile("./test_data/quote_random.json")
		if err != nil {
			t.Fatalf("Failed to read response file: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseFile)
	}))
	defer mockingServer.Close()

	//
	quotableClient := quotable.NewQuotableClient(mockingServer.URL)

	params := quotable.GetRandomQuoteQueryParams{
		Tags:  "technology,famous-quotes",
		Limit: 2,
	}

	quotes, err := quotableClient.GetRandomQuote(params)
	if err != nil {
		t.Fatalf("GetRandomQuote failed: %s", err)
	}

	if len(quotes) != 2 {
		t.Errorf("Expected 2 quotes, got %d", len(quotes))
	}

	// check the quotes have the expected tags
	for i, q := range quotes {
		if !slices.Contains(q.Tags, "Technology") && !slices.Contains(q.Tags, "Famous Quotes") {
			t.Errorf("Expected quote %d to have tag 'Technology' or 'Famous Quotes', got %v", i, q.Tags)
		}
	}
}
