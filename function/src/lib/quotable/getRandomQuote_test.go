package quotable_test

import (
	"slices"
	"testing"

	"github.com/clowa/az-func-daily-quote/src/lib/quotable"
)

const randomQuotesApiEndpoint = "/quotes/random"

var randomQuotesResponseDumps = []string{
	"./test_data/quote_random.json",
}

func TestGetRandomQuote(t *testing.T) {
	// Create a mock server to simulate the Quotable API.
	// The response is read from a local JSON file for consistency.
	// The server performs basic validation on the incoming request.
	mockingServer := startMockingServer(t, randomQuotesApiEndpoint, randomQuotesResponseDumps...)
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
