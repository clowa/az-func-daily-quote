package quotableSdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
)

const (
	getRandomQuotePath = "/quotes/random"
)

type GetRandomQuoteQueryParams struct {
	Limit     int                       `url:"limit,omitempty"`
	MaxLength int                       `url:"maxLength,omitempty"`
	MinLength int                       `url:"minLength,omitempty"`
	Tags      commaSeparatedQueryString `url:"tags,omitempty"`
	Author    string                    `url:"author,omitempty"`
	AuthorId  string                    `url:"authorId,omitempty"`
}

type commaSeparatedQueryString []string

func (qp commaSeparatedQueryString) EncodeValues(key string, v *url.Values) error {
	if (qp == nil) || (len(qp) == 0) {
		return nil
	}

	var tags string
	for i, tag := range qp {
		// If we are at the last tag, don't add a comma
		if i == len(qp)-1 {
			tags = tags + tag
			continue
		}
		tags = tags + tag + ","
	}
	v.Set(key, tags)
	return nil
}

type QuoteResponse struct {
	Id         string   `json:"_id"`
	Content    string   `json:"content"`
	Author     string   `json:"author"`
	AuthorSlug string   `json:"authorSlug"`
	Length     int      `json:"length"`
	Tags       []string `json:"tags"`
}

func GetRandomQuote(params GetRandomQuoteQueryParams) ([]QuoteResponse, error) {

	urlValues, err := query.Values(params)
	if err != nil {
		return []QuoteResponse{}, err
	}

	apiEndpoint := quotableApiUrl + getRandomQuotePath + "?" + urlValues.Encode()
	log.Infof("Fetching quote from %s", apiEndpoint)
	req, err := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return []QuoteResponse{}, fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []QuoteResponse{}, fmt.Errorf("failed to make request: %s", err)
	}
	defer resp.Body.Close()

	var quotes []QuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return []QuoteResponse{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	return quotes, nil
}
