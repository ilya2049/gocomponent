package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ilya2049/gocomponent/internal/domain/component"
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
	CustomComponents            []string          `toml:"custom"`
	ComponentColors             map[string]string `toml:"colors"`
}

func (conf *Config) ToComponentGraphConfig() *component.GraphConfig {
	return &component.GraphConfig{
		ExtendComponentIDs:          conf.ExtendComponentIDs,
		IncludeThirdPartyComponents: conf.IncludeThirdPartyComponents,
		ThirdPartyComponentsColor:   component.NewColor(conf.ThirdPartyComponentsColor),
		IncludeParentComponents:     component.NewNamespaces(conf.IncludeParentComponents),
		IncludeChildComponents:      component.NewNamespaces(conf.IncludeChildComponents),
		ExcludeParentComponents:     component.NewNamespaces(conf.ExcludeParentComponents),
		ExcludeChildComponents:      component.NewNamespaces(conf.ExcludeChildComponents),
		CustomComponents:            component.NewNamespaces(conf.CustomComponents),
		ComponentColors:             component.NewNamespaceColorMap(conf.ComponentColors),
	}
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
