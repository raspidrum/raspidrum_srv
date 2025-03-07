package repo

type ParamType interface {
	~int | ~string | ~float64
}

type Param[T ParamType] struct {
	Name  string
	Value T
}

type SamplerRepo interface {
	ConnectAudioOutput(driver string, params []Param[string]) (devId int, err error)
	ConnectMidiInput(driver string, params []Param[string]) (devId int, err error)
	CreateChannel(audioDevId, midiDevId int, instrumentFile string) (channelId int, err error)
}
