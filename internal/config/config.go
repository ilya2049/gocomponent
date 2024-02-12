package config

type Config struct {
	ProjectDirectory                 string   `toml:"project_directory"`
	HideThirdPartyImports            bool     `toml:"hide_third_party_imports"`
	IncludeOnlyNextPackageNamespaces []string `toml:"include_only_next_package_namespaces"`
}
