# Documentation

All URIs are relative to *<https://quotes.clowa.dev>*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetQuote**](DefaultAPI.md#GetQuote) | **Get** /api/quote | Get a quote
[**PostQuote**](DefaultAPI.md#PostQuote) | **Post** /api/quote | Write a quote to the database

## GetQuote

Get a quote

### Example Request

```url
GET https://api.clowa.dev/api/quote HTTP/1.1
Host: api.clowa.dev
Ocp-Apim-Subscription-Key: <Place subscription key here>
```

### Example Response

 ```json
 {
    "content": "Communications tools don't get socially interesting until they get technologically boring.",
    "author": "Clay Shirky",
    "tags": [
        "Technology"
    ]
 }
 ```

### Authorization

[apiKeyHeader](#apikeyheader), [apiKeyQuery](#apikeyquery)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

## PostQuote

Write a quote to the database

### Example Request

```url
POST https://api.clowa.dev/api/quote HTTP/1.1
Host: api.clowa.dev
Content-Type: application/json
Ocp-Apim-Subscription-Key: <Place subscription key here>
```

```json
{
    "id": "01234567890",
    "content": "The only limit to our realization of tomorrow is our doubts of today.",
    "author": "Franklin D. Roosevelt",
    "authorSlug": "franklin-d-roosevelt",
    "length": 69,
    "tags": [
        "inspiration",
        "politics"
    ],
    "creationDate": "2024-09-26"
}
```

### Example Response

- `200`: Success
- `400`: The server failed to parse the request body.
- `401`: No authentication was provided. Please see endpoint documentation for authentication requirements.
- `415`: Request header `Content-Type` is not set to `application/json`.
- `500`: The server failed to write the quote to the database.

### Authorization

[apiKeyHeader](#apikeyheader), [apiKeyQuery](#apikeyquery)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: Not defined

## Documentation For Authorization

Authentication schemes defined for the API:

### apiKeyHeader

- **Type**: API key
- **API key parameter name**: Ocp-Apim-Subscription-Key
- **Location**: HTTP header

### apiKeyQuery

- **Type**: API key
- **API key parameter name**: subscription-key
- **Location**: URL query string

## ToDo

- [ ] Add visualization of the architecture
- [ ] Add visualization of the data flow
