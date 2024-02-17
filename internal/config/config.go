package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ProjectDirectory                 string            `toml:"project_directory"`
	IncludeThirdPartyComponents      bool              `toml:"include_third_party_components"`
	IncludeOnlyNextPackageNamespaces []string          `toml:"include_only_next_package_namespaces"`
	NamespaceColors                  map[string]string `toml:"namespace_colors"`
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
