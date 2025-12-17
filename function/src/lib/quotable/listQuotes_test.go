package quotable_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"testing"

	"github.com/clowa/az-func-daily-quote/src/lib/quotable"
)

var responseDumps = []string{
	"./test_data/quote_list_page_1.json",
	"./test_data/quote_list_page_2.json",
}

func startMockingServer(t *testing.T, apiEndpoint string, dumpResponsePages ...string) *httptest.Server {

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

		// Fake pagination handling
		pageIndex := 1
		if page := r.URL.Query().Get("page"); page != "" {
			var err error
			pageIndex, err = strconv.Atoi(page)
			if err != nil {
				t.Fatalf("Failed to parse page number: %s", err)
			}

			if pageIndex < 1 || pageIndex > len(dumpResponsePages) {
				t.Fatalf("Invalid page number. Got %s but only have %d pages", page, len(dumpResponsePages))
			}
		}

		pageIndex-- // Convert to zero-based index
		responseFile, err := os.ReadFile(dumpResponsePages[pageIndex])
		if err != nil {
			t.Fatalf("Failed to read response file: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseFile)
	}))

	return mockingServer
}

func TestListQuotes(t *testing.T) {
	const apiEndpoint = "/quotes"

	// Create a mock server to simulate the Quotable API.
	// The response is read from a local JSON file for consistency.
	// The server performs basic validation on the incoming request.
	mockingServer := startMockingServer(t, apiEndpoint, responseDumps...)
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
	const apiEndpoint = "/quotes"

	// Expected total quotes across all pages in the mock responses
	const expectedTotalQuotes = 6

	// Create a mock server to simulate the Quotable API.
	// The response is read from a local JSON file for consistency.
	// The server performs basic validation on the incoming request.
	mockingServer := startMockingServer(t, apiEndpoint, responseDumps...)
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
