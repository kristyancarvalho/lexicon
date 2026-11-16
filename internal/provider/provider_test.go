package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistryContainsSupportedProviders(t *testing.T) {
	for _, id := range []string{OpenRouterID, HuggingFaceID} {
		info, ok := Lookup(id)
		if !ok {
			t.Fatalf("provider %s missing", id)
		}
		if info.Endpoint == "" || info.DefaultModel == "" {
			t.Fatalf("provider %s incomplete: %+v", id, info)
		}
	}
	if len(All()) != 2 {
		t.Fatalf("unexpected provider count: %d", len(All()))
	}
}

func TestNewValidatesInput(t *testing.T) {
	if _, err := New("nope", "token"); !errors.Is(err, ErrUnknown) {
		t.Fatalf("expected ErrUnknown, got %v", err)
	}
	if _, err := New(OpenRouterID, ""); !errors.Is(err, ErrNoToken) {
		t.Fatalf("expected ErrNoToken, got %v", err)
	}
	if _, err := New(HuggingFaceID, "hf_x"); err != nil {
		t.Fatal(err)
	}
}

func TestChatSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Errorf("unexpected auth header %q", got)
		}
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
		}
		if req.Model != "m" || len(req.Messages) != 1 || req.Messages[0].Content != "hi" {
			t.Errorf("unexpected request: %+v", req)
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"  hello  "}}]}`))
	}))
	defer server.Close()

	c := newCompat(Info{Endpoint: server.URL, DefaultModel: "d"}, "secret", server.Client())
	got, err := c.Chat(context.Background(), "m", []Message{{Role: RoleUser, Content: "hi"}})
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestChatUsesDefaultModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Model != "d" {
			t.Errorf("expected default model, got %q", req.Model)
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	c := newCompat(Info{Endpoint: server.URL, DefaultModel: "d"}, "t", server.Client())
	if _, err := c.Chat(context.Background(), "", nil); err != nil {
		t.Fatal(err)
	}
}

func TestChatAPIErrorObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid key"}}`))
	}))
	defer server.Close()

	c := newCompat(Info{Endpoint: server.URL}, "t", server.Client())
	_, err := c.Chat(context.Background(), "m", nil)

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.Status != http.StatusUnauthorized || apiErr.Message != "invalid key" {
		t.Fatalf("unexpected error: %+v", apiErr)
	}
}

func TestChatAPIErrorString(t *testing.T) {
