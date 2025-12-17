package quotable

import "net/http"

type QuotableClient struct {
	baseUrl string
}

// NewQuotableClient creates a new instance of QuotableClient.
// The host parameter allows to specify the target API endpoint.
func NewQuotableClient(host string) *QuotableClient {
	return &QuotableClient{
		baseUrl: host,
	}
}

func (c *QuotableClient) do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
