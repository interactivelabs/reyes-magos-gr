// Package typesafe is a minimal client for the TypeSafe System One API
// (https://docs.typesafe.ai/api), shared by the classification scripts.
package typesafe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const endpoint = "https://api.typesafe.ai/v1/systemone"

type Answer struct {
	Choice        string             `json:"choice"`
	Probabilities map[string]float64 `json:"probabilities"`
	Confidence    float64            `json:"confidence"`
	Noul          float64            `json:"noul"`
}

type Response struct {
	Answers map[string]Answer `json:"answers"`
}

type Client struct {
	apiKey string
	http   *http.Client
}

// NewFromEnv builds a client from TYPESAFE_API_KEY.
func NewFromEnv() (*Client, error) {
	apiKey := os.Getenv("TYPESAFE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("TYPESAFE_API_KEY is not set")
	}
	return &Client{apiKey: apiKey, http: &http.Client{Timeout: 60 * time.Second}}, nil
}

var errRetryable = errors.New("retryable")

// Ask evaluates questions against state with jev-latest, backing off on
// 429/529 as the API docs recommend.
func (c *Client) Ask(state any, questions map[string]any) (Response, error) {
	payload, err := json.Marshal(map[string]any{"model": "jev-latest", "state": state, "questions": questions})
	if err != nil {
		return Response{}, err
	}
	backoff := time.Second
	for attempt := 0; ; attempt++ {
		resp, err := c.post(payload)
		if !errors.Is(err, errRetryable) || attempt == 4 {
			return resp, err
		}
		time.Sleep(backoff)
		backoff *= 2
	}
}

func (c *Client) post(payload []byte) (Response, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}
	switch {
	case res.StatusCode == http.StatusTooManyRequests || res.StatusCode == 529:
		return Response{}, fmt.Errorf("%w: %s", errRetryable, res.Status)
	case res.StatusCode != http.StatusOK:
		return Response{}, fmt.Errorf("%s: %s", res.Status, data)
	}
	var r Response
	return r, json.Unmarshal(data, &r)
}
