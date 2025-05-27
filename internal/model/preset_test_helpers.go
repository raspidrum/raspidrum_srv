package model

import (
	"fmt"
	"strings"
)

// VerifyControlsForTest is a test helper to verify the internal controls state
func VerifyControlsForTest(p *KitPreset, expectedControls map[string]struct {
	Key    string
	Owner  ControlOwner
	MidiCC int
	CfgKey string
	Type   string
	Value  float32
}) string {
	if len(p.controls) != len(expectedControls) {
		return fmt.Sprintf("Controls count mismatch: got %d, want %d", len(p.controls), len(expectedControls))
	}

	var differences []string
	for key, expected := range expectedControls {
		control, exists := p.controls[key]
		if !exists {
			differences = append(differences, fmt.Sprintf("Control %q not found", key))
			continue
		}
		if control.MidiCC != expected.MidiCC {
			differences = append(differences, fmt.Sprintf("Control %q MidiCC mismatch: got %d, want %d", key, control.MidiCC, expected.MidiCC))
		}
		if control.CfgKey != expected.CfgKey {
			differences = append(differences, fmt.Sprintf("Control %q CfgKey mismatch: got %q, want %q", key, control.CfgKey, expected.CfgKey))
		}
		if control.Type != expected.Type {
			differences = append(differences, fmt.Sprintf("Control %q Type mismatch: got %q, want %q", key, control.Type, expected.Type))
		}
		if control.Value != expected.Value {
			differences = append(differences, fmt.Sprintf("Control %q Value mismatch: got %f, want %f", key, control.Value, expected.Value))
		}
	}

	if len(differences) == 0 {
		return ""
	}
	return fmt.Sprintf("Found differences:\n%s", strings.Join(differences, "\n"))
}
