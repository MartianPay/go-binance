package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MartianPay/go-binance/utils"
)

const (
	BaseURL        = "https://api.binance.com"
	DefaultTimeout = 30 * time.Second
)

type Client struct {
	apiKey     string
	secretKey  string
	baseURL    string
	httpClient *http.Client
	signer     *utils.Signer
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   BaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		signer: utils.NewSigner(apiKey, secretKey),
	}
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

func (c *Client) doRequest(method, endpoint string, params map[string]string, body interface{}, needSign bool) ([]byte, error) {
	url := c.baseURL + endpoint

	var queryString string
	if needSign {
		queryString = c.signer.Sign(params)
	} else if len(params) > 0 {
		queryString = c.signer.BuildQueryString(params)
	}

	if queryString != "" {
		url += "?" + queryString
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	headers := c.signer.GetHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) Get(endpoint string, params map[string]string, needSign bool) ([]byte, error) {
	return c.doRequest(http.MethodGet, endpoint, params, nil, needSign)
}

func (c *Client) Post(endpoint string, params map[string]string, body interface{}, needSign bool) ([]byte, error) {
	return c.doRequest(http.MethodPost, endpoint, params, body, needSign)
}

func (c *Client) Delete(endpoint string, params map[string]string, needSign bool) ([]byte, error) {
	return c.doRequest(http.MethodDelete, endpoint, params, nil, needSign)
}
