package quotable_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// startMockingServer starts a HTTP server to mock the Quotable API.
// It serves responses from the provided dumpResponsePages based on the "page" query parameter.
// The returned httptest.Server should be closed by the caller by calling Close() on it.
// ToDo: Enhance request validation in sense of checking query parameters etc.
// ToDo: Perform actual filtering based on query parameters.
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
