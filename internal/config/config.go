package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ProjectDirectory            string            `toml:"project_directory"`
	ExtendComponentIDs          map[string]int    `toml:"extend_ids"`
	IncludeThirdPartyComponents bool              `toml:"include_third_party"`
	ThirdPartyComponentsColor   string            `toml:"third_party_color"`
	IncludeParentComponents     []string          `toml:"include_parents"`
	IncludeChildComponents      []string          `toml:"include_children"`
	ExcludeParentComponents     []string          `toml:"exclude_parents"`
	ExcludeChildComponents      []string          `toml:"exclude_children"`
	ComponentColors             map[string]string `toml:"colors"`
}

func Read() (*Config, error) {
	configContents, err := os.ReadFile("config.toml")
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var conf Config
	if _, err := toml.Decode(string(configContents), &conf); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &conf, nil
}
