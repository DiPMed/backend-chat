package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPServer HTTPServerConfig `yaml:"http_server"`
	LLM        LLMConfig        `yaml:"llm"`
}

type HTTPServerConfig struct {
	Address string `yaml:"address"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"-"`
	Model    string `yaml:"model"`
}

func (c *Config) Validate() error {
	if c.HTTPServer.Address == "" {
		return fmt.Errorf("http_server.address is required")
	}
	if c.LLM.APIKey == "" {
		return fmt.Errorf("llm.api_key is required")
	}
	return nil
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := &Config{
		HTTPServer: HTTPServerConfig{Address: ":8080"},
		LLM: LLMConfig{
			Provider: "gemini",
			Model:    "gemini-2.5-flash",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	cfg.LLM.APIKey = os.Getenv("LLM_API_KEY")

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}
