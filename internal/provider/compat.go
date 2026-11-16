package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxBodyBytes = 4 << 20

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (%d): %s", e.Status, e.Message)
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error json.RawMessage `json:"error"`
}

type compat struct {
	info   Info
	token  string
	client *http.Client
}

func newCompat(info Info, token string, client *http.Client) *compat {
	return &compat{info: info, token: token, client: client}
}

func (c *compat) Info() Info {
	return c.info
}

func (c *compat) Chat(ctx context.Context, model string, messages []Message) (string, error) {
	if model == "" {
		model = c.info.DefaultModel
	}

	body, err := json.Marshal(chatRequest{Model: model, Messages: messages})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.info.Endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Title", "lexicon")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return "", err
	}

	var out chatResponse
	parseErr := json.Unmarshal(raw, &out)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := errorMessage(out.Error)
		if msg == "" {
			msg = truncate(strings.TrimSpace(string(raw)), 200)
		}
		return "", &APIError{Status: resp.StatusCode, Message: msg}
	}
	if parseErr != nil {
		return "", parseErr
	}
	if len(out.Choices) == 0 {
		return "", ErrEmptyResponse
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

func errorMessage(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var obj struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil && obj.Message != "" {
		return obj.Message
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
