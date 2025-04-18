package file

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	m "github.com/raspidrum-srv/internal/model"
)

func ParsePreset(path string) (*m.KitPreset, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}
	var pst m.KitPreset
	if err := yaml.Unmarshal(content, &pst); err != nil {
		return nil, fmt.Errorf("error parsing YAML in %s: %w", path, err)
	}
	return &pst, nil
}
