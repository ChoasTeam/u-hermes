package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ModelConfig struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIBase   string `json:"api_base"`
	APIKey    string `json:"api_key"`
	IsDefault bool   `json:"is_default"`
}

type ChatConfig struct {
	SystemPrompt string `json:"system_prompt"`
}

type Config struct {
	Version int           `json:"version"`
	Models  []ModelConfig `json:"models"`
	Chat    ChatConfig    `json:"chat"`
}

func defaultConfig() *Config {
	return &Config{
		Version: 1,
		Models:  []ModelConfig{},
		Chat: ChatConfig{
			SystemPrompt: "你是一个有用的AI助手",
		},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Chat.SystemPrompt == "" {
		cfg.Chat.SystemPrompt = defaultConfig().Chat.SystemPrompt
	}
	return cfg, nil
}

func Save(path string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("cannot save nil config")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) GetDefaultModel() *ModelConfig {
	for i := range c.Models {
		if c.Models[i].IsDefault {
			return &c.Models[i]
		}
	}
	if len(c.Models) > 0 {
		return &c.Models[0]
	}
	return nil
}
