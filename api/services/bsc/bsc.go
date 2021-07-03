package bsc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://api.bscscan.com/api"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: BaseURL,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func (c *Client) sendRequest(queryValues url.Values, v interface{}) error {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return err
	}

	queryValues.Add("apiKey", c.apiKey)
	base.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", base.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes response
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return err
		}
		if str, ok := errRes.Result.(string); ok {
			return errors.New(str)
		} else {
			return errors.New("wrong format received")
		}
	}

	fullResponse := response{
		Result: v,
	}
	if err := json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
