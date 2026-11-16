package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Info struct {
	ID           string
	Name         string
	Endpoint     string
	DefaultModel string
}

type Provider interface {
	Info() Info
	Chat(ctx context.Context, model string, messages []Message) (string, error)
}

const (
	OpenRouterID  = "openrouter"
	HuggingFaceID = "huggingface"
)

var (
	ErrNoToken       = errors.New("token not configured")
	ErrUnknown       = errors.New("unknown provider")
	ErrEmptyResponse = errors.New("empty response from provider")
)

var registry = []Info{
	{
		ID:           OpenRouterID,
		Name:         "OpenRouter",
		Endpoint:     "https://openrouter.ai/api/v1/chat/completions",
		DefaultModel: "openai/gpt-4o-mini",
	},
	{
		ID:           HuggingFaceID,
		Name:         "Hugging Face",
		Endpoint:     "https://router.huggingface.co/v1/chat/completions",
		DefaultModel: "meta-llama/Llama-3.3-70B-Instruct",
	},
}

func All() []Info {
	out := make([]Info, len(registry))
	copy(out, registry)
	return out
}

func Lookup(id string) (Info, bool) {
	for _, info := range registry {
		if info.ID == id {
			return info, true
		}
	}
	return Info{}, false
}

func New(id, token string) (Provider, error) {
	info, ok := Lookup(id)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, id)
	}
	if token == "" {
		return nil, fmt.Errorf("%w: %s", ErrNoToken, info.Name)
	}
	return newCompat(info, token, &http.Client{Timeout: 2 * time.Minute}), nil
}
