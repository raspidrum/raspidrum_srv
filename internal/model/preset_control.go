package model

// PresetControl.Type values MUST match one of the ControlType values
type ControlType int

const (
	CTVolume ControlType = iota
	CTPan
	CTPitch
	CTOther
)

var ControlTypeToString = map[ControlType]string{
	CTVolume: "volume",
	CTPan:    "pan",
	CTPitch:  "pitch",
	CTOther:  "other",
}

var ControlTypeFromString = map[string]ControlType{
	"volume": CTVolume,
	"pan":    CTPan,
	"pitch":  CTPitch,
	"other":  CTOther,
}

var (
	CtrlVolume = ControlTypeToString[CTVolume]
	CtrlPan    = ControlTypeToString[CTPan]
)

type SamplerControlSetter interface {
	SendChannelMidiCC(channelKey string, cc int, value float32) error
	SetChannelVolume(channelKey string, value float32) error
}

type ControlOwner interface {
	HandleControlValue(channelKey string, control *PresetControl, value float32, csetter SamplerControlSetter) error
}

type ControlMap map[string]*PresetControl

// CfgKey - sfz-variable key, same value as Instrument.Controls
// Key - unique id across preset. Used for identification control for communication between srv and ui
// linkedTo - ref to control, example: channel volume control linked to instrument volume control
// linkedWith - ref from control, example: instrument volume control linked from channel volume control
type PresetControl struct {
	Name       string  `yaml:"name,omitempty" json:"name,omitempty"`
	Type       string  `yaml:"type" json:"type"`
	MidiCC     int     `yaml:"midiCC,omitempty" json:"midiCC,omitempty"`
	CfgKey     string  `yaml:"-" json:"-"`
	Value      float32 `yaml:"value" json:"value"`
	Key        string  `yaml:"-" json:"-"`
	owner      ControlOwner
	linkedTo   []*PresetControl
	linkedWith *PresetControl
}

func (c ControlMap) GetControlByKey(key string) (*PresetControl, bool) {
	ctrl, ok := c[key]
	if !ok {
		return nil, false
	}
	return ctrl, true
}

func (c ControlMap) FindControlByType(t string) (*PresetControl, bool) {
	for _, c := range c {
		if t == c.Type {
			return c, true
		}
	}
	return nil, false
}

func (c *PresetControl) SetValue(value float32, channelKey string, csetter SamplerControlSetter) error {
	return c.owner.HandleControlValue(channelKey, c, value, csetter)
}
