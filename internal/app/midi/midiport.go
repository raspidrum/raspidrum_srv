package midi

type MIDIPortState int

const (
	MIDIPortStateUnknown MIDIPortState = iota
	MIDIPortStateConnected
	MIDIPortStateDisconnected
)

type MIDIPortType int

const (
	MIDIPortTypeUnknown MIDIPortType = iota
	MIDIPortTypeHardware
	MIDIPortTypeVirtual
	MIDIPortTypeSoftware
)

type MIDIPortInfo struct {
	Driver     string // "ALSA", "JACK", "COREMIDI"
	DevId      string
	PortId     string
	Name       string
	State      MIDIPortState
	DeviceType MIDIPortType
}

type MIDIPorts map[string]MIDIPortInfo
