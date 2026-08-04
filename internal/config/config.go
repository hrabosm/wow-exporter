package config

type Character struct {
	Name   string `yaml:"name"`
	Realm  string `yaml:"realm"`
	Region string `yaml:"region"`
}

type Config struct {
	CurrentTier string `yaml:"current-tier"`

	Characters []Character `yaml:"characters"`
}
