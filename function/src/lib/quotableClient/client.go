package quotableClient

import "net/http"

const (
	quotableApiUrl = "https://api.quotable.io"
)

type QuotableClient struct {
	baseUrl string
}

func NewQuotableClient() *QuotableClient {
	return &QuotableClient{
		baseUrl: quotableApiUrl,
	}
}

func (c *QuotableClient) do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
