package audio

type AudioChannelDirection int

const (
	AudioChannelAny = iota
	AudioChannelOutput
	AudioChannelInput
)

type AudioChannel struct {
	Id   string
	Name string
}
