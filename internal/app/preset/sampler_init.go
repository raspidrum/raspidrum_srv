package preset

import (
	"fmt"

	"github.com/raspidrum-srv/internal/app/audio"
	"github.com/raspidrum-srv/internal/app/midi"
	"github.com/raspidrum-srv/internal/repo"
)

func InitSampler(sampler repo.SamplerRepo, midiDev midi.MIDIDevice, audioDevice audio.AudioDevice) (audioDevId, midiDevId int, err error) {
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
	audioBindings := make(map[int][]repo.Param[string], 0)
	audioDriver := audioDevice.Driver()
	if audioDriver == "JACK" {
		// TODO: get from channels
		audioBindings[0] = []repo.Param[string]{
			{Name: "JACK_BINDINGS", Value: "system:playback_1"},
			{Name: "JACK_BINDINGS", Value: "system:playback_2"},
		}
	}
	audioId, err := sampler.ConnectAudioOutput(audioDriver, audioBindings)
	if err != nil {
		return 0, 0, fmt.Errorf("failed init sampler: %w", err)
	}
	return audioId, midiId, nil
}
