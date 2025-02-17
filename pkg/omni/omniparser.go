package omni

import (
	"context"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents the API client.
type Client struct {
	BaseURL string
	client  *resty.Client
}

// NewClient creates a new omniparser API client with retry configuration.
func NewClient(baseURL string) *Client {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)

	return &Client{
		BaseURL: baseURL,
		client:  client,
	}
}

// ParseRequest represents the request payload.
type ParseRequest struct {
	Base64Image string `json:"base64_image"`
}

// ParseResponse represents the response from the API.
type ParseResponse struct {
	ImageBase64       string   `json:"som_image_base64"`
	ParsedContentList []string `json:"parsed_content_list"`
	Latency           float64  `json:"latency"`
}

// Parse sends an image for parsing and returns the response.
func (c *Client) Parse(ctx context.Context, base64Image string) (*ParseResponse, error) {
	requestBody := ParseRequest{Base64Image: base64Image}
	response := &ParseResponse{}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(response).
		Post(c.BaseURL + "/parse/")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New("API error: " + resp.Status())
	}

	return response, nil
}

// ProbeResponse represents the response from the probe endpoint.
type ProbeResponse struct {
	Message string `json:"message"`
}

// Probe checks if the API is ready.
func (c *Client) Probe(ctx context.Context) (*ProbeResponse, error) {
	response := &ProbeResponse{}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(response).
		Get(c.BaseURL + "/probe/")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New("API error: " + resp.Status())
	}

	return response, nil
}
