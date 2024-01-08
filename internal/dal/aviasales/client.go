package aviasales

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClientWithProxy *resty.Client
	httpClient          *resty.Client
}

func NewClient(proxy string) (*Client, error) {
	userAgent := uarand.GetRandom()

	httpClient := resty.New()
	httpClient.SetHeader("User-Agent", userAgent)
	httpClient.SetRedirectPolicy(resty.NoRedirectPolicy())
	httpClient.SetCookieJar(nil)

	httpClientWithProxy := resty.New()
	httpClientWithProxy.SetHeader("User-Agent", userAgent)
	httpClientWithProxy.SetRedirectPolicy(resty.NoRedirectPolicy())
	httpClientWithProxy.SetProxy(proxy)
	httpClientWithProxy.SetCookieJar(nil)

	return &Client{
		httpClient:          httpClient,
		httpClientWithProxy: httpClientWithProxy,
	}, nil
}

func (ac *Client) getOriginCookie() (string, error) {
	resp, err := ac.httpClient.R().Get("https://www.aviasales.ru/")
	if err != nil {
		return "", fmt.Errorf("can't do request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("server sent code %d, but expected 200: %s", resp.StatusCode(), resp.Body())
	}

	cookieStrings := make([]string, 0)
	for _, cookie := range resp.Cookies() {
		cookieStrings = append(cookieStrings, cookie.String())
	}

	return strings.Join(cookieStrings, "; "), nil
}

func (ac *Client) StartSearch(reqBody SearchStartRequestBody) (*SearchStartResponseBody, error) {
	originCookie, err := ac.getOriginCookie()
	if err != nil {
		return nil, fmt.Errorf("can't get auid cookie: %w", err)
	}

	resp, err := ac.httpClientWithProxy.R().
		SetBody(reqBody).
		SetHeader("x-origin-cookie", originCookie).
		SetHeader("x-client-type", "mobile").
		Post("https://tickets-api.aviasales.ru/search/v2/start")
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server sent code %d, but expected 200: %s", resp.StatusCode(), resp.Body())
	}

	var respBody SearchStartResponseBody
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal response body: %w", err)
	}

	return &respBody, nil
}

func (ac *Client) GetSearchResults(resultsURL string, reqBody SearchResultsRequestBody) (SearchResultsResponseBody, bool, error) {
	resp, err := ac.httpClient.R().
		SetBody(reqBody).
		Post(fmt.Sprintf("https://%s/search/v3/results", resultsURL))
	if err != nil {
		return nil, false, fmt.Errorf("can't send request: %w", err)
	}

	isLast := false
	if resp.Header().Get("X-Stop-Marker") != "" {
		isLast = true
	}

	if resp.StatusCode() == 304 {
		return SearchResultsResponseBody{}, isLast, nil
	}

	if resp.StatusCode() != 200 {
		return nil, isLast, fmt.Errorf("server sent code %d, but expected 200 or 304: %s", resp.StatusCode(), resp.Body())
	}

	var respBody SearchResultsResponseBody
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, isLast, fmt.Errorf("can't unmarshal response body: %w", err)
	}

	return respBody, isLast, nil
}

func (ac *Client) SearchIATA(filter string) (SearchIATAResponseBody, error) {
	queryParams := url.Values{}
	queryParams.Add("locale", "en_US")
	queryParams.Add("max", "5")
	queryParams.Add("term", filter)
	queryParams.Add("types", "city")
	queryParams.Add("types", "airport")
	queryParams.Add("types", "country")

	resp, err := ac.httpClient.R().
		SetQueryParamsFromValues(queryParams).
		Get("https://suggest.aviasales.com/v2/places.json")
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("server sent code %d, but expected 200: %s", resp.StatusCode(), resp.Body())
	}

	var respBody SearchIATAResponseBody
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal response body: %w", err)
	}

	return respBody, nil
}
