package mididevice

type MIDIDevice interface {
	DevID() string
	Name() string
	GetKeysMapping() (map[string]int, error)
}

type USBMIDIDevice struct {
	// MIDI Device ID
	// Example, for ALSA: "24:0"
	devId string
	// Device Name
	name string
}

func NewUSBMIDIDevice(devId string, name string) USBMIDIDevice {
	return USBMIDIDevice{
		devId: devId,
		name:  name,
	}
}

func (m *USBMIDIDevice) Name() string {
	return m.name
}

func (m *USBMIDIDevice) DevID() string {
	return m.devId
}

// Get MIDI key mapping in format:
// alias : MIDI Key
// eg:
//
//	hihat_close : 42
//	tom1 : 48
func (m *USBMIDIDevice) GetKeysMapping() (map[string]int, error) {
	// TODO: get mapping from repo
	return map[string]int{
		"kick1":            36,
		"snare":            38,
		"snare_rimshot":    39,
		"tom1":             48,
		"tom2":             45,
		"tom3":             43,
		"tom4":             41,
		"hihat_close":      42,
		"hihat_open":       46,
		"hihat_loose":      29,
		"hihat_foot_open":  27,
		"hihat_foot_close": 44,
		"hihat_splash":     28,
		"crash1_edge":      49,
		"ride1_edge":       51,
		"ride1_bell":       53,
	}, nil

}
