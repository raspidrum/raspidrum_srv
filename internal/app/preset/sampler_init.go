package preset

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/midi"
	"github.com/raspidrum-srv/internal/repo"
)

// TODO: remove hardcoded value
const audioDriver = "COREAUDIO"

// TODO: может сделать тип Sampler, в который сохранять созданные идентификаторы устройств и каналов

func InitSampler(sampler repo.SamplerRepo, midiDev midi.MIDIDevice) (audioDevId, midiDevId int, err error) {
	// Init MIDI
	var midiBindings repo.Param[string]
	mport, err := midiDev.GetOutPort()
	if err != nil {
		return 0, 0, fmt.Errorf("failed get MIDI out port: %w", err)
	}
	if mport == nil {
		return 0, 0, fmt.Errorf("no MIDI out port available for device %s", midiDev.Name())
	}
	switch mport.Driver {
	case "COREMIDI":
		midiBindings.Name = "CORE_MIDI_BINDINGS"
	case "ALSA":
		midiBindings.Name = "ALSA_SEQ_BINDINGS"
	}
	midiBindings.Value = mport.PortId

	midiId, err := sampler.ConnectMidiInput(mport.Driver, []repo.Param[string]{midiBindings})
	if err != nil {
		return 0, 0, fmt.Errorf("failed init sampler: %w", err)
	}

	// Init Audio
	audioId, err := sampler.ConnectAudioOutput(audioDriver, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed init sampler: %w", err)
	}
	return audioId, midiId, nil
}
