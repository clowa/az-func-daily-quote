package quotableClient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
)

type ListQuotesRequestParams struct {
	MaxLength int `url:"maxLength,omitempty"`
	MinLength int `url:"minLength,omitempty"`
	// Tags can be a comma-separated list of tags and supports logical OR (|) and AND (,) operators
	Tags     string `url:"tags,omitempty"`
	Author   string `url:"author,omitempty"`
	AuthorId string `url:"authorId,omitempty"`
	SortBy   string `url:"sortBy,omitempty"`
	Order    string `url:"order,omitempty"`
	Limit    int    `url:"limit,omitempty"`
	Page     int    `url:"page,omitempty"`
}

type ListQuotesResponse struct {
	Count         int             `json:"count"`
	TotalCount    int             `json:"totalCount"`
	Page          int             `json:"page"`
	TotalPages    int             `json:"totalPages"`
	LastItemIndex int             `json:"lastItemIndex"`
	Results       []QuoteResponse `json:"results"`
}

func (c *QuotableClient) ListQuotes(params ListQuotesRequestParams) (ListQuotesResponse, error) {
	const listQuotesPath = "/quotes"

	urlValues, err := query.Values(params)
	if err != nil {
		return ListQuotesResponse{}, err
	}

	apiEndpoint := c.baseUrl + listQuotesPath + "?" + urlValues.Encode()
	log.Infof("Listing quotes from %s", apiEndpoint)
	req, err := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return ListQuotesResponse{}, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return ListQuotesResponse{}, fmt.Errorf("failed to make request: %s", err)
	}
	defer resp.Body.Close()

	var listQuotesResponse ListQuotesResponse
	err = json.NewDecoder(resp.Body).Decode(&listQuotesResponse)
	if err != nil {
		return ListQuotesResponse{}, fmt.Errorf("failed to decode response body: %s", err)
	}

	return listQuotesResponse, nil
}

func (c *QuotableClient) ListAllQuotes(params ListQuotesRequestParams) ([]QuoteResponse, error) {
	var allQuotes []QuoteResponse
	currentPage := 1
	for {
		params.Page = currentPage
		listQuotesResponse, err := c.ListQuotes(params)
		if err != nil {
			return nil, err
		}

		allQuotes = append(allQuotes, listQuotesResponse.Results...)

		if currentPage >= listQuotesResponse.TotalPages {
			break
		}
		currentPage++
	}

	return allQuotes, nil
}
