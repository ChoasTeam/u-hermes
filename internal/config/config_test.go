package config

import (
	"path/filepath"
	"testing"
)

func TestLoadConfig_NotExists_ReturnsDefaults(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Models) != 0 {
		t.Errorf("expected 0 models, got %d", len(cfg.Models))
	}
	if cfg.Chat.SystemPrompt != "你是一个有用的AI助手" {
		t.Errorf("unexpected system prompt: %s", cfg.Chat.SystemPrompt)
	}
}

func TestSaveAndLoadConfig_PreservesData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := &Config{
		Version: 1,
		Models: []ModelConfig{
			{ID: "default", Name: "DeepSeek", APIBase: "https://api.deepseek.com", APIKey: "sk-test", IsDefault: true},
		},
		Chat: ChatConfig{SystemPrompt: "hello"},
	}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Verify all fields preserved
	if loaded.Version != 1 {
		t.Errorf("version mismatch: got %d", loaded.Version)
	}
	if len(loaded.Models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(loaded.Models))
	}
	if loaded.Models[0].ID != "default" {
		t.Errorf("ID mismatch: %s", loaded.Models[0].ID)
	}
	if loaded.Models[0].Name != "DeepSeek" {
		t.Errorf("name mismatch: %s", loaded.Models[0].Name)
	}
	if loaded.Models[0].APIBase != "https://api.deepseek.com" {
		t.Errorf("API base mismatch: %s", loaded.Models[0].APIBase)
	}
	if loaded.Models[0].APIKey != "sk-test" {
		t.Errorf("API key mismatch: %s", loaded.Models[0].APIKey)
	}
	if !loaded.Models[0].IsDefault {
		t.Errorf("IsDefault should be true")
	}
	if loaded.Chat.SystemPrompt != "hello" {
		t.Errorf("system prompt mismatch: %s", loaded.Chat.SystemPrompt)
	}
}

func TestConfig_GetDefaultModel_ReturnsDefault(t *testing.T) {
	cfg := &Config{
		Models: []ModelConfig{
			{ID: "a", IsDefault: false},
			{ID: "b", IsDefault: true},
			{ID: "c", IsDefault: false},
		},
	}
	m := cfg.GetDefaultModel()
	if m == nil || m.ID != "b" {
		t.Errorf("expected model b, got %v", m)
	}
}

func TestConfig_GetDefaultModel_NoDefault_ReturnsNil(t *testing.T) {
	cfg := &Config{Models: []ModelConfig{}}
	if m := cfg.GetDefaultModel(); m != nil {
		t.Errorf("expected nil, got %v", m)
	}
}
