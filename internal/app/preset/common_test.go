package preset

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/goccy/go-yaml"
	m "github.com/raspidrum-srv/internal/model"
	lsampler "github.com/raspidrum-srv/internal/repo/linuxsampler"
	lscp "github.com/raspidrum-srv/libs/liblscp-go"
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

func getDBPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../db/")
	return dir
}

func getProjectPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f), "../../../")
	return dir
}

func connectSampler() (*lsampler.LinuxSampler, error) {
	lsClient := lscp.NewClient("localhost", "8888", "1s")

	err := lsClient.Connect()
	ls := lsampler.LinuxSampler{
		Client: lsClient,
		Engine: "sfz",
	}
	if err != nil {
		return &ls, fmt.Errorf("Failed connect to LinuxSampler: %v", err)
	}
	return &ls, nil
}

func loadPresetFromYAML(t *testing.T, filename string) *m.KitPreset {
	t.Helper()
	absPath, err := filepath.Abs("../../../testdata/preset_load/" + filename)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}
	var preset m.KitPreset
	if err := yaml.Unmarshal(data, &preset); err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	return &preset
}
