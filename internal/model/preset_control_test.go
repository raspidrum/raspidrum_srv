package model

import "testing"

// MockSamplerControlSetter implements SamplerControlSetter interface for testing
type MockSamplerControlSetter struct {
	volumeCalled bool
	midiCCCalled bool
	volumeValue  float32
	midiCCValue  float32
	midiCC       int
	channelKey   string
}

func (m *MockSamplerControlSetter) SetChannelVolume(channelKey string, value float32) error {
	m.volumeCalled = true
	m.volumeValue = value
	m.channelKey = channelKey
	return nil
}

func (m *MockSamplerControlSetter) SendChannelMidiCC(channelKey string, cc int, value float32) error {
	m.midiCCCalled = true
	m.midiCCValue = value
	m.midiCC = cc
	m.channelKey = channelKey
	return nil
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
		name           string
		testData       string
		args           args
		controlKey     string
		value          float32
		wantVolumeCall bool
		wantMidiCCCall bool
		wantValue      float32
		wantMidiCC     int
		wantChannelKey string
		wantErr        bool
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
			controlKey:     "c0volume",
			value:          0.75,
			wantVolumeCall: true,
			wantMidiCCCall: false,
			wantValue:      0.75,
			wantChannelKey: "ch1",
			wantErr:        false,
		},
		{
			name:     "set instrument volume with MIDI CC",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey:     "i0volume",
			value:          0.95,
			wantVolumeCall: false,
			wantMidiCCCall: true,
			wantValue:      0.95,
			wantMidiCC:     30,
			wantChannelKey: "ch1",
			wantErr:        false,
		},
		{
			name:     "set instrument pan with MIDI CC",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey:     "i0pan",
			value:          0.54,
			wantVolumeCall: false,
			wantMidiCCCall: true,
			wantValue:      0.54,
			wantMidiCC:     10,
			wantChannelKey: "ch1",
			wantErr:        false,
		},
		{
			name:     "set non-existent control",
			testData: "single_instrument.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey:     "nonexistent",
			value:          0.5,
			wantVolumeCall: false,
			wantMidiCCCall: false,
			wantErr:        true,
		},
		{
			name:     "set layers volume with MIDI CC",
			testData: "single instr_with_layers.yaml",
			args: args{
				mididevs: []MIDIDevice{
					&MockMMIDIDevice{},
				},
			},
			controlKey:     "i0l0volume",
			value:          0.95,
			wantVolumeCall: false,
			wantMidiCCCall: true,
			wantValue:      0.95,
			wantMidiCC:     104,
			wantChannelKey: "ch1",
			wantErr:        false,
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

			// Verify volume call
			if mockSetter.volumeCalled != tt.wantVolumeCall {
				t.Errorf("SetChannelVolume called = %v, want %v", mockSetter.volumeCalled, tt.wantVolumeCall)
			}
			if tt.wantVolumeCall && mockSetter.volumeValue != tt.wantValue {
				t.Errorf("SetChannelVolume value = %v, want %v", mockSetter.volumeValue, tt.wantValue)
			}

			// Verify MIDI CC call
			if mockSetter.midiCCCalled != tt.wantMidiCCCall {
				t.Errorf("SendChannelMidiCC called = %v, want %v", mockSetter.midiCCCalled, tt.wantMidiCCCall)
			}
			if tt.wantMidiCCCall {
				if mockSetter.midiCCValue != tt.wantValue {
					t.Errorf("SendChannelMidiCC value = %v, want %v", mockSetter.midiCCValue, tt.wantValue)
				}
				if mockSetter.midiCC != tt.wantMidiCC {
					t.Errorf("SendChannelMidiCC cc = %v, want %v", mockSetter.midiCC, tt.wantMidiCC)
				}
			}

			// Verify channel key
			if mockSetter.channelKey != tt.wantChannelKey {
				t.Errorf("channel key = %v, want %v", mockSetter.channelKey, tt.wantChannelKey)
			}
		})
	}
}
