package acceptanceTests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/clowa/az-func-daily-quote/src/lib/quote"
	docker "github.com/gruntwork-io/terratest/modules/docker"
)

const hostname = "localhost"

func invokeRestCall(port int, path string, v any) bool {
	url := fmt.Sprintf("http://%s:%d%s", hostname, port, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return false
	}

	return true
}

func invokeGetQuoteApi(port int) (quote.Quote, bool) {
	const path = "/api/quote"
	var quote quote.Quote

	if ok := invokeRestCall(port, path, &quote); !ok {
		return quote, false
	}

	return quote, true
}

func TestQuoteOfTheDay(t *testing.T) {
	const port = 8080
	const composeFile = "../../docker-compose.yaml"

	// spin up end-to-end test environment
	dockerOptions := &docker.Options{
		EnableBuildKit: true,
	}
	docker.RunDockerCompose(t, dockerOptions, "-f", composeFile, "up", "--build", "-d", "--wait")
	defer docker.RunDockerCompose(t, dockerOptions, "-f", composeFile, "down", "--volumes")

	firstQuote, ok := invokeGetQuoteApi(port)
	if !ok {
		t.Errorf("Failed to get first quote")
	}

	secondQuote, ok := invokeGetQuoteApi(port)
	if !ok {
		t.Errorf("Failed to get second quote")
	}

	if reflect.DeepEqual(firstQuote, quote.Quote{}) || reflect.DeepEqual(secondQuote, quote.Quote{}) {
		t.Errorf("Expected a quote, but got an empty quote")
	}

	if !reflect.DeepEqual(firstQuote, secondQuote) {
		t.Errorf("Expected quotes to be the same. Expected: %v, Actual: %v", firstQuote, secondQuote)
	}
}
