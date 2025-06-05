//go:test unit

package model

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

func (i *InstrumentRef) UnmarshalYAML(data []byte) error {
	type alias struct {
		Id         int64              `yaml:"-"`
		Uid        string             `yaml:"uuid"`
		Key        string             `yaml:"-"`
		Name       string             `yaml:"name"`
		CfgMidiKey string             `yaml:"midiKey"`
		Controls   map[string]Control `yaml:"controls"`
		Layers     map[string]Layer   `yaml:"layers"`
	}
	var a alias
	err := yaml.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	*i = InstrumentRef(a)
	return nil
}

func (i *PresetInstrument) UnmarshalYAML(data []byte) error {
	type alias struct {
		Instrument InstrumentRef          `yaml:"instrument"`
		Id         int64                  `yaml:"id"`
		Name       string                 `yaml:"name"`
		ChannelKey string                 `yaml:"channelKey"`
		MidiKey    string                 `yaml:"midiKey,omitempty"`
		MidiNote   int                    `yaml:"-"`
		Controls   ControlMap             `yaml:"controls"`
		Layers     map[string]PresetLayer `yaml:"layers"`
	}
	var a alias
	err := yaml.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	*i = PresetInstrument(a)
	return nil
}

type ExpectedControls map[string]struct {
	Key    string
	Name   string
	Owner  ControlOwner
	MidiCC int
	CfgKey string
	Type   string
	Value  float32
}

// VerifyControlsForTest is a test helper to verify the internal controls state
func VerifyControlsForTest(p *KitPreset, expectedControls ExpectedControls) string {
	if len(p.controls) != len(expectedControls) {
		return fmt.Sprintf("Controls count mismatch: got %d, want %d", len(p.controls), len(expectedControls))
	}

	var differences []string
	for key, expected := range expectedControls {
		ctrlRef, exists := p.controls[key]
		if !exists {
			differences = append(differences, fmt.Sprintf("Control %q not found in preset index", key))
			continue
		}
		if ctrlRef.control.Key != expected.Key {
			differences = append(differences, fmt.Sprintf("Control %q Key mismatch: got %q, want %q", key, ctrlRef.control.Key, expected.Key))
		}
		if ctrlRef.control.MidiCC != expected.MidiCC {
			differences = append(differences, fmt.Sprintf("Control %q MidiCC mismatch: got %d, want %d", key, ctrlRef.control.MidiCC, expected.MidiCC))
		}
		if ctrlRef.control.CfgKey != expected.CfgKey {
			differences = append(differences, fmt.Sprintf("Control %q CfgKey mismatch: got %q, want %q", key, ctrlRef.control.CfgKey, expected.CfgKey))
		}
		if ctrlRef.control.Type != expected.Type {
			differences = append(differences, fmt.Sprintf("Control %q Type mismatch: got %q, want %q", key, ctrlRef.control.Type, expected.Type))
		}
		if ctrlRef.control.Value != expected.Value {
			differences = append(differences, fmt.Sprintf("Control %q Value mismatch: got %f, want %f", key, ctrlRef.control.Value, expected.Value))
		}
	}

	if len(differences) == 0 {
		return ""
	}
	return fmt.Sprintf("Found differences:\n%s", strings.Join(differences, "\n"))
}
