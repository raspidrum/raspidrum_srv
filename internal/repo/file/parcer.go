package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse kit directory with kit and instrument files
func parseYAMLDir(dir string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории и не-YAML файлы
		if d.IsDir() || !(strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			return nil
		}

		parsed, err := parseYAMLFile(path)
		if err != nil {
			return err
		}

		result[path] = parsed
		return nil
	})

	return result, err
}

// Parse yaml file with kit or instrument
func parseYAMLFile(path string) (interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	// Определяем тип по первому ключу
	var probe struct {
		Instrument *Instrument `yaml:"instrument"`
		Kit        *Kit        `yaml:"kit"`
	}
	if err := yaml.Unmarshal(content, &probe); err != nil {
		return nil, fmt.Errorf("error parsing YAML in %s: %w", path, err)
	}

	switch {
	case probe.Instrument != nil:
		var result Instrument
		err = yaml.Unmarshal(content, &result)
		return result, err
	case probe.Kit != nil:
		var result Kit
		err = yaml.Unmarshal(content, &result)
		return result, err
	default:
		return nil, fmt.Errorf("unknown YAML format in %s", path)
	}
}
