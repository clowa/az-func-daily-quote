package tests

import (
	"slices"
	"testing"

	quoteable "github.com/clowa/az-func-daily-quote/src/lib/quotableClient"
)

func TestGetRandomQuote(t *testing.T) {

	quotableClient := quoteable.NewQuotableClient()

	params := quoteable.GetRandomQuoteQueryParams{
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
		if !slices.Contains(q.Tags, "technology") && !slices.Contains(q.Tags, "famous-quotes") {
			t.Errorf("Expected quote %d to have tag 'technology' or 'famous-quotes', got %v", i, q.Tags)
		}
	}
}

func TestListQuotes(t *testing.T) {

	quotableClient := quoteable.NewQuotableClient()

	params := quoteable.ListQuotesRequestParams{
		Tags:  "technology",
		Limit: 3,
		Page:  1,
	}

	response, err := quotableClient.ListQuotes(params)
	if err != nil {
		t.Fatalf("ListQuotes failed: %s", err)
	}

	if response.Count != 3 {
		t.Errorf("Expected 3 quotes, got %d", response.Count)
	}

	for i, q := range response.Results {
		if !slices.Contains(q.Tags, "technology") {
			t.Errorf("Expected quote %d to have tag 'technology', got %v", i, q.Tags)
		}
	}
}
