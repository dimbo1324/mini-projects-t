package config

import (
	"mycli/internal/models"
	"os"

	"gopkg.in/yaml.v3"
)

type PathsConfig struct {
	JsonSchemas []models.Schema `yaml:"json_schemas"`
}

func LoadConfig(path string) ([]models.Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg PathsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg.JsonSchemas, nil
}
