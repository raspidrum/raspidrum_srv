package model

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/midi"
)

// MapMidiKey maps a MIDI key string to its numeric value using the provided MIDI devices
func MapMidiKey(mkey string, mdev midi.MIDIDevice) (int, error) {
	kmap, err := mdev.GetKeysMapping()
	if err != nil {
		return 0, fmt.Errorf("failed get MIDI Keys mapping for device %s: %w", mdev.Name(), err)
	}
	midiId, ok := kmap[mkey]
	if ok {
		return midiId, nil
	}
	return 0, fmt.Errorf("MIDI devices %s doen't have mapping for MIDI Key %s", mdev.Name(), mkey)
}
