package config

type Config struct {
	ProjectDirectory      string `toml:"project_directory"`
	ShowThirdPartyImports bool   `toml:"show_third_party_imports"`
}
