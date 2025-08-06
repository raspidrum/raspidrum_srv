package audio

type AudioDeviceProvider interface {
	GetAudioDev() (AudioDevice, error)
}

type AudioDevice interface {
	GetChannels(direction AudioChannelDirection) ([]AudioChannel, error)
	Driver() string
	Name() string
}
