package models

type Config struct {
	Name        string
	Description string
	Version     int
	Author      string
	Tags        []string
}

type Schema struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Tags        string `yaml:"tags"`
}