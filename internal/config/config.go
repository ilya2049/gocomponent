package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/pkg/fs"
)

const defaultConfigFileName = "config.toml"

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
	EnableComponentSize         bool              `toml:"enable_size"`
}

func (c *Config) makeSureProjectDirIsSet() {
	if c.ProjectDirectory == "" {
		c.ProjectDirectory = "./"
	}
}

func newDefault() *Config {
	conf := Config{}

	conf.makeSureProjectDirIsSet()

	return &conf
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
		EnableComponentSize:         conf.EnableComponentSize,
	}
}

type Reader struct {
	fileReader fs.FileReader
}

func NewReader(fileReader fs.FileReader) *Reader {
	return &Reader{
		fileReader: fileReader,
	}
}

func (r *Reader) ReadConfig() (*Config, error) {
	configContents, err := r.fileReader.ReadFile(defaultConfigFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newDefault(), nil
		}

		return nil, fmt.Errorf("read config: %w", err)
	}

	var conf Config
	if _, err := toml.Decode(string(configContents), &conf); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	conf.makeSureProjectDirIsSet()

	return &conf, nil
}
