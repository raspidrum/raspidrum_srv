package file

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	m "github.com/raspidrum-srv/internal/model"
)

// Parse kit directory with kit and instrument files
func ParseYAMLDir(dir string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := filepath.WalkDir(dir, func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории и не-YAML файлы
		if d.IsDir() || !(strings.HasSuffix(filepath, ".yaml") || strings.HasSuffix(filepath, ".yml")) {
			return nil
		}

		parsed, err := parseYAMLFile(filepath)
		if err != nil {
			return err
		}

		_, file := path.Split(filepath)
		result[file] = parsed
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
		Instrument *m.Instrument `yaml:"instrument"`
		Kit        *m.Kit        `yaml:"kit"`
	}
	if err := yaml.Unmarshal(content, &probe); err != nil {
		return nil, fmt.Errorf("error parsing YAML in %s: %w", path, err)
	}

	switch {
	case probe.Instrument != nil:
		return probe.Instrument, nil
	case probe.Kit != nil:
		return probe.Kit, nil
	default:
		return nil, fmt.Errorf("unknown YAML format in %s", path)
	}
}
