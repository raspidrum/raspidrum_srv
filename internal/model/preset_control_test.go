package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type callParam struct {
	VolumeCall bool
	MidiCCCall bool
	Value      float32
	MidiCC     int
	ChannelKey string
}

// MockSamplerControlSetter implements SamplerControlSetter interface for testing
type MockSamplerControlSetter struct {
	CallParams []callParam
}

func (m *MockSamplerControlSetter) SetChannelVolume(channelKey string, value float32) error {
	prm := callParam{
		VolumeCall: true,
		Value:      value,
		ChannelKey: channelKey,
	}
	m.CallParams = append(m.CallParams, prm)
	return nil
}

func (m *MockSamplerControlSetter) SendChannelMidiCC(channelKey string, cc int, value float32) error {
	prm := callParam{
		MidiCCCall: true,
		Value:      value,
		MidiCC:     cc,
		ChannelKey: channelKey,
	}
	m.CallParams = append(m.CallParams, prm)
	return nil
}

func (m *MockSamplerControlSetter) Compare(t *testing.T, wants []callParam) {
	wl := len(wants)
	if len(m.CallParams) != wl {
		if diff := cmp.Diff(wants, m.CallParams); diff != "" {
			t.Errorf("callParams mismatch (-want +got):\n%s", diff)
		}
		return
	}
	// if len > 1 then compare by midiCC
	if wl > 1 {
		for _, prm := range wants {
			// find callParam with same midiCC
			foundIdx := -1
			for ci, cprm := range m.CallParams {
				if cprm.MidiCC == prm.MidiCC {
					foundIdx = ci
					break
				}
			}
			if foundIdx == -1 {
				t.Errorf("callParam with midiCC = %v not found", prm.MidiCC)
			} else {
				if diff := cmp.Diff(prm, m.CallParams[foundIdx]); diff != "" {
					t.Errorf("callParam mismatch (-want +got):\n%s", diff)
				}
			}
		}
	} else {
		if diff := cmp.Diff(wants[0], m.CallParams[0]); diff != "" {
			t.Errorf("callParam mismatch (-want +got):\n%s", diff)
		}
	}
}

// Test doesn't handle case: use instrument control instead of channel control:
//   - Server returns control instrument control in channel controls.
//   - SetControlValues don't expect channel control
func Test_SetControlValue(t *testing.T) {
	type args struct {
		preset   *KitPreset
		mididevs []MIDIDevice
	}
	type testCase struct {
		name       string
		testData   string
		args       args
		controlKey string
		value      float32
		wants      []callParam
		wantErr    bool
	}
	tests := []testCase{
		{
			name:     "set channel volume",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "c0volume",
			value:      0.75,
			wants: []callParam{
				{
					VolumeCall: true,
					MidiCCCall: false,
					Value:      0.75,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set instrument volume with MIDI CC",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "i0volume",
			value:      0.95,
			wants: []callParam{
				{
					VolumeCall: false,
					MidiCCCall: true,
					Value:      121.0,
					MidiCC:     30,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set instrument pan with MIDI CC",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "i0pan",
			value:      0.54,
			wants: []callParam{
				{
					VolumeCall: false,
					MidiCCCall: true,
					Value:      98.0,
					MidiCC:     10,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set non-existent control",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "nonexistent",
			value:      0.5,
			wants: []callParam{
				{
					VolumeCall: false,
					MidiCCCall: false,
				},
			},
			wantErr: true,
		},
		{
			name:     "set pitch (tunes) control in instrument",
			testData: "single_instrument_with_tunes.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "i0pitch",
			value:      0.75,
			wants: []callParam{
				{
					VolumeCall: false,
					MidiCCCall: true,
					Value:      95,
					MidiCC:     11,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set pitch (tunes) control in layer",
			testData: "layer_with_tune.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "i0bellpitch",
			value:      0.75,
			wants: []callParam{
				{
					VolumeCall: false,
					MidiCCCall: true,
					Value:      95,
					MidiCC:     16,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set instrument virtual vol",
			testData: "single instr_with_layers.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "i0volume",
			value:      0.5,
			wants: []callParam{
				{
					MidiCCCall: true,
					Value:      40,
					MidiCC:     104,
					ChannelKey: "ch1",
				},
				{
					MidiCCCall: true,
					Value:      45,
					MidiCC:     103,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set layer vol with instrument correction",
			testData: "single instr_with_layers.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "i0bellvolume",
			value:      1.0,
			wants: []callParam{
				{
					MidiCCCall: true,
					Value:      121,
					MidiCC:     104,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
		{
			name:     "set channel virtual pan",
			testData: "two_instruments_channel_virtual_pan.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey: "c0pan",
			value:      0.50,
			wants: []callParam{
				{
					MidiCCCall: true,
					Value:      27,
					MidiCC:     10,
					ChannelKey: "ch1",
				},
				{
					MidiCCCall: true,
					Value:      43,
					MidiCC:     11,
					ChannelKey: "ch1",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load preset from YAML
			preset := loadPresetFromYAML(t, tt.testData)
			tt.args.preset = preset

			// Prepare the preset
			if err := preset.PrepareToLoad(tt.args.mididevs); err != nil {
				t.Fatalf("PrepareToLoad() error = %v", err)
			}

			// Create mock control setter
			mockSetter := &MockSamplerControlSetter{}

			// Call SetControlValue
			err := preset.SetControlValue(tt.controlKey, tt.value, mockSetter)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("SetControlValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			mockSetter.Compare(t, tt.wants)

		})
	}
}
