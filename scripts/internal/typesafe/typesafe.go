// Package typesafe is a minimal client for the TypeSafe System One API
// (https://docs.typesafe.ai/api), shared by the classification scripts.
package typesafe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"time"
)

const endpoint = "https://api.typesafe.ai/v1/systemone"

// Question is one entry of the request's questions, as the scripts keep them
// in their rules.json files.
type Question struct {
	Type         string         `json:"type"`
	Instructions any            `json:"instructions"`
	Criteria     map[string]any `json:"criteria,omitempty"`
}

// With returns a copy of q whose instructions also carry key: value. Plain
// string instructions become {"question": ...}.
func (q Question) With(key string, value any) Question {
	instructions := map[string]any{}
	switch v := q.Instructions.(type) {
	case map[string]any:
		instructions = maps.Clone(v)
	case string:
		instructions["question"] = v
	}
	instructions[key] = value
	q.Instructions = instructions
	return q
}

// WithCriteria returns a copy of q with extra criteria added to its own.
func (q Question) WithCriteria(extra map[string]any) Question {
	criteria := maps.Clone(q.Criteria)
	if criteria == nil {
		criteria = map[string]any{}
	}
	maps.Copy(criteria, extra)
	q.Criteria = criteria
	return q
}

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
func (c *Client) Ask(state any, questions map[string]Question) (Response, error) {
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
