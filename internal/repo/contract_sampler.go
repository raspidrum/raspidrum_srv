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

// key (string) - channel key from preset
// value (int) - sampler channel id
type SamplerChannels map[string]int

type SamplerRepo interface {
	ConnectAudioOutput(driver string, params map[int][]Param[string]) (devId int, err error)
	ConnectMidiInput(driver string, params []Param[string]) (devId int, err error)
	CreateChannel(audioDevId, midiDevId int) (channelId int, err error)
	LoadInstrument(instrumentFile string, instrIdx int, channelId int) error
	LoadPreset(audioDevId, midiDevId int, preset *m.KitPreset, fs afero.Fs) (SamplerChannels, error)
	SetChannelVolume(samplerChn int, volume float32) error
	SendMidiCC(samplerChn int, cc int, value float32) error
}
