package mididevice

type MIDIDevice struct {
	// MIDI Device ID
	// Example, for ALSA: "24:0"
	DevId string
	// Device Name
	Name string
}

// Get MIDI key mapping in format:
// alias : MIDI Key
// eg:
//
//	hihat_close : 42
//	tom1 : 48
func (m *MIDIDevice) GetKeysMapping() (map[string]int, error) {
	// TODO: get mapping from repo
	return map[string]int{
		"kick":             36,
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
