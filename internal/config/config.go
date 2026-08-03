package config

type Character struct {
	Name   string `yaml:"name"`
	Realm  string `yaml:"realm"`
	Region string `yaml:"region"`
}

type File struct {
	Characters []Character `yaml:"characters"`
}
