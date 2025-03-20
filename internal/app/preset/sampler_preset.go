package preset

import (
	"fmt"

	"github.com/raspidrum-srv/internal/repo"
)

// TODO: убрать хардкод
const audioDriver = "COREAUDIO"
const midiDriver = "COREMIDI"

// TODO: Добавить инициализацию repo Linuxsampler и заинжектить в него lscp.Client

// TODO: может сделать тип Sampler, в который сохранять созданные идентификаторы устройств и каналов

func InitSampler(sampler repo.SamplerRepo) (audioDevId, midiDevId int, err error) {
	audioId, err := sampler.ConnectAudioOutput(audioDriver, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed init sampler: %w", err)
	}

	// TODO: убрать хардкод
	midiBindings := repo.Param[string]{
		Name:  "CORE_MIDI_BINDINGS",
		Value: "vmpk vmpk out",
	}
	midiId, err := sampler.ConnectMidiInput(midiDriver, []repo.Param[string]{midiBindings})
	if err != nil {
		return 0, 0, fmt.Errorf("failed init sampler: %w", err)
	}
	return audioId, midiId, nil
}

func LoadPresetToSampler(sampler repo.SamplerRepo, audDevId, midiDevId int, instrumentFile string) (chnl int, err error) {
	chnl, err = sampler.CreateChannel(audDevId, midiDevId, instrumentFile)
	if err != nil {
		return chnl, fmt.Errorf("failed load preset: %w", err)
	}
	return
}
