package repo

import (
	midi "github.com/raspidrum-srv/internal/app/mididevice"
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
	CreateChannel(audioDevId, midiDevId int, instrumentFile string) (channelId int, err error)
	LoadPreset(preset *m.KitPreset, mididevs []midi.MIDIDevice, fs afero.Fs) error
}
