package app

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mycli/internal/db"
	"mycli/internal/models"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func getValueByPath(data map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	var current interface{} = data
	for _, key := range keys {
		currMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("key %s not found or parent is not a map", key)
		}
		val, exists := currMap[key]
		if !exists {
			return nil, fmt.Errorf("key %s not found", key)
		}
		current = val
	}
	return current, nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
type Processor struct {
	Db          *db.DataBase
	JsonSchemas []models.Schema
}

func NewProcessor(database *db.DataBase, schemas []models.Schema) *Processor {
	return &Processor{
		Db:          database,
		JsonSchemas: schemas,
	}
}

func (p *Processor) ProcessDirectory(rootPath string) error {
	var filesCount, insertedCount, updatedCount int
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		var cfg *models.Config
		var parseErr error

		switch ext {
		case ".yaml", ".yml":
			cfg, parseErr = p.processYaml(path)
		case ".json":
			cfg, parseErr = p.processJson(path)
		default:
			return nil
		}

		if parseErr != nil {
			log.Printf("Error processing file %s: %v", path, parseErr)
			return nil
		}

		if cfg != nil {
			filesCount++
			inserted, dbErr := p.Db.UpsertConfig(*cfg)
			if dbErr != nil {
				log.Printf("DB Error for %s: %v", cfg.Name, dbErr)
			} else {
				if inserted {
					insertedCount++
				} else {
					updatedCount++
				}
			}
		}

		return nil
	})
	fmt.Printf("\n--- Statistics ---\n")
	fmt.Printf("Files Processed: %d\n", filesCount)
	fmt.Printf("Records Added:   %d\n", insertedCount)
	fmt.Printf("Records Updated: %d\n", updatedCount)
	return err
}

func (p *Processor) processYaml(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	type YamlFile struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Version     int    `yaml:"version"`
		Metadata    struct {
			Author string   `yaml:"author"`
			Tags   []string `yaml:"tags"`
		} `yaml:"metadata"`
	}

	var yf YamlFile
	if err := yaml.Unmarshal(data, &yf); err != nil {
		return nil, err
	}

	return &models.Config{
		Name:        yf.Name,
		Description: yf.Description,
		Version:     yf.Version,
		Author:      yf.Metadata.Author,
		Tags:        yf.Metadata.Tags,
	}, nil
}

func (p *Processor) processJson(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, err
	}

	for _, schema := range p.JsonSchemas {
		nameVal, err := getValueByPath(rawData, schema.Name)
		if err != nil {
			continue
		}

		nameStr, ok := nameVal.(string)
		if !ok {
			continue
		}

		cfg := &models.Config{Name: nameStr}

		if v, err := getValueByPath(rawData, schema.Description); err == nil {
			if s, ok := v.(string); ok {
				cfg.Description = s
			}
		}

		if v, err := getValueByPath(rawData, schema.Version); err == nil {
			switch val := v.(type) {
			case float64:
				cfg.Version = int(val)
			case int:
				cfg.Version = val
			}
		}

		if v, err := getValueByPath(rawData, schema.Author); err == nil {
			if s, ok := v.(string); ok {
				cfg.Author = s
			}
		}

		if v, err := getValueByPath(rawData, schema.Tags); err == nil {
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					if s, ok := item.(string); ok {
						cfg.Tags = append(cfg.Tags, s)
					}
				}
			}
		}

		return cfg, nil
	}

	return nil, fmt.Errorf("no matching schema found for json file")
}
