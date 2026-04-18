package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNonexistent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg != nil {
		t.Error("expected nil config for nonexistent file")
	}
}

func TestSaveAndLoad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	want := &Config{
		BaseURL: "https://memos.example.com",
		Token:   "secret-token",
	}

	if err := Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p := filepath.Join(home, ".config", "memos-cli", "config.toml")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.BaseURL != want.BaseURL {
		t.Errorf("BaseURL: got %q, want %q", got.BaseURL, want.BaseURL)
	}
	if got.Token != want.Token {
		t.Errorf("Token: got %q, want %q", got.Token, want.Token)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := &Config{BaseURL: "http://localhost:5230", Token: "tok"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	dir := filepath.Join(home, ".config", "memos-cli")
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat config dir: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected config dir to be a directory")
	}
}

func TestSaveFilePermissions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := Save(&Config{BaseURL: "x", Token: "y"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p := filepath.Join(home, ".config", "memos-cli", "config.toml")
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}
}
