package model

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/midi"
)

// MapMidiKey maps a MIDI key string to its numeric value using the provided MIDI devices
func MapMidiKey(mkey string, mdevs []midi.MIDIDevice) (int, error) {
	devlist := make([]string, len(mdevs))
	for i, d := range mdevs {
		kmap, err := d.GetKeysMapping()
		if err != nil {
			return 0, fmt.Errorf("failed get MIDI Keys mapping for device %s: %w", d.Name(), err)
		}
		devlist[i] = d.Name()
		midiId, ok := kmap[mkey]
		if ok {
			return midiId, nil
		}
	}
	return 0, fmt.Errorf("MIDI devices %s doen't have mapping for MIDI Key %s", devlist, mkey)
}
