package repo

import (
	m "github.com/raspidrum-srv/internal/model"
	"github.com/spf13/afero"
)

type ParamType interface {
	~int | ~string | ~float64
}

type Param[T ParamType] struct {
	Name  string
	Value T
}

type SamplerRepo interface {
	ConnectAudioOutput(driver string, params map[int][]Param[string]) (devId int, err error)
	ConnectMidiInput(driver string, params []Param[string]) (devId int, err error)
	CreateChannel(audioDevId, midiDevId int) (channelId int, err error)
	LoadInstrument(instrumentFile string, instrIdx int, channelId int) error
	LoadPreset(audioDevId, midiDevId int, preset *m.KitPreset, fs afero.Fs) error
	SetChannelVolume(samplerChn int, volume float64) error
}
