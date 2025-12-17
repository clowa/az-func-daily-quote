package quotable_test

import (
	"slices"
	"testing"

	"github.com/clowa/az-func-daily-quote/src/lib/quotable"
)

const listQuotesApiEndpoint = "/quotes"

var listQuotesResponseDumps = []string{
	"./test_data/quote_list_page_1.json",
	"./test_data/quote_list_page_2.json",
}

func TestListQuotes(t *testing.T) {
	// Create a mock server to simulate the Quotable API.
	// The response is read from a local JSON file for consistency.
	// The server performs basic validation on the incoming request.
	mockingServer := startMockingServer(t, listQuotesApiEndpoint, listQuotesResponseDumps...)
	defer mockingServer.Close()

	quotableClient := quotable.NewQuotableClient(mockingServer.URL)

	params := quotable.ListQuotesRequestParams{
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
		if !slices.Contains(q.Tags, "Technology") {
			t.Errorf("Expected quote %d to have tag 'Technology', got %v", i, q.Tags)
		}
	}
}

func TestListAllQuotes(t *testing.T) {
	// Expected total quotes across all pages in the mock responses
	const expectedTotalQuotes = 6

	// Create a mock server to simulate the Quotable API.
	// The response is read from a local JSON file for consistency.
	// The server performs basic validation on the incoming request.
	mockingServer := startMockingServer(t, listQuotesApiEndpoint, listQuotesResponseDumps...)
	defer mockingServer.Close()

	quotableClient := quotable.NewQuotableClient(mockingServer.URL)

	params := quotable.ListQuotesRequestParams{
		Tags:  "technology",
		Limit: 3,
	}

	quotes, err := quotableClient.ListAllQuotes(params)
	if err != nil {
		t.Fatalf("ListAllQuotes failed: %s", err)
	}

	if len(quotes) != expectedTotalQuotes {
		t.Errorf("Expected %d quotes, got %d", expectedTotalQuotes, len(quotes))
	}
}
