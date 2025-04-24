package preset

import (
	"testing"

	midi "github.com/raspidrum-srv/internal/app/mididevice"
	m "github.com/raspidrum-srv/internal/model"
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
	}, nil
}

func Test_augmentFromInstrument(t *testing.T) {
	type args struct {
		pst      *m.KitPreset
		mididevs []midi.MIDIDevice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := augmentFromInstrument(tt.args.pst, tt.args.mididevs); (err != nil) != tt.wantErr {
				t.Errorf("augmentFromInstrument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
