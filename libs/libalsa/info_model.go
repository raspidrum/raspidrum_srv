package libalsa

type MidiPortType int

const (
	MidiPortTypeHardware MidiPortType = iota
	MidiPortTypeSoftware
)

// AlsaCard represents ALSA sound card information
type AlsaCard struct {
	ID       int
	Name     string
	LongName string
	Driver   string
	Mixer    string
}

// MidiPortInfo represents ALSA sequencer MIDI port info
type MidiPortInfo struct {
	ClientId   int
	PortId     int
	CardId     int
	ClientName string
	PortName   string
	PortType   MidiPortType
	IsInput    bool // true if port supports input
	IsOutput   bool // true if port supports output
}
