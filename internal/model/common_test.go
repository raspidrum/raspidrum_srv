package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
)

type MockMMIDIDevice struct{}

func (m *MockMMIDIDevice) Name() string {
	return "Dummy"
}

func (m *MockMMIDIDevice) DevID() string {
	return "0:0"
}

func (m *MockMMIDIDevice) GetKeysMapping() (map[string]int, error) {
	return map[string]int{
		"kick1":      36,
		"snare":      38,
		"ride1_edge": 51,
		"ride1_bell": 53,
		"tom1":       48,
	}, nil
}

func loadPresetFromYAML(t *testing.T, filename string) *KitPreset {
	t.Helper()
	absPath, err := filepath.Abs("../../testdata/preset_load/" + filename)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}
	var preset KitPreset
	if err := yaml.Unmarshal(data, &preset); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	return &preset
}
