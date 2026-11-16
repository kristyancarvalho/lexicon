package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFileReturnsEmptyConfig(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Tokens == nil || cfg.Models == nil {
		t.Fatal("maps must be initialised")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	cfg := New()
	cfg.Active = "openrouter"
	cfg.SetToken("openrouter", "sk-test")

	if err := cfg.Save(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Active != "openrouter" || loaded.Token("openrouter") != "sk-test" {
		t.Fatalf("unexpected config: %+v", loaded)
	}
}

func TestSaveRestrictsPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := New().Save(path); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("expected 0600, got %o", perm)
	}
}

func TestModelFallback(t *testing.T) {
	cfg := New()
	if got := cfg.Model("openrouter", "default"); got != "default" {
		t.Fatalf("got %q", got)
	}
	cfg.Models["openrouter"] = "custom"
	if got := cfg.Model("openrouter", "default"); got != "custom" {
		t.Fatalf("got %q", got)
	}
}
