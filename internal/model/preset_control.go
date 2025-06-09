package model

import "math"

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
	var val float32
	if c.Type == CtrlPan {
		val = c.denormalizePan(value)
	} else {
		val = c.denormalizeBase(value)
	}
	return c.owner.HandleControlValue(channelKey, c, val, csetter)
}

func (c *PresetControl) GetNormalizedValue() (val float32, min float32, max float32) {
	if c.Type == CtrlPan {
		return c.normalizePan()
	}
	return c.normalizeBase()
}

func (ctrl *PresetControl) normalizeBase() (val float32, min float32, max float32) {
	if ctrl.MidiCC != 0 {
		// val from 0 to 1 with 3 decimal places
		return roundFloat(float32(ctrl.Value/127), 3), 0, 1
	}
	return roundFloat(float32(ctrl.Value), 3), 0, 1
}

func (ctrl *PresetControl) denormalizeBase(val float32) float32 {
	if ctrl.MidiCC != 0 {
		return roundFloat(float32(val*127), 0)
	}
	return roundFloat(val, 3)
}

func (ctrl *PresetControl) normalizePan() (val float32, min float32, max float32) {
	if ctrl.MidiCC != 0 {
		return roundFloat(float32((ctrl.Value*2/127)-1), 3), -1, 1
	}
	return roundFloat(float32(ctrl.Value), 3), -1, 1
}

func (ctrl *PresetControl) denormalizePan(val float32) float32 {
	if ctrl.MidiCC != 0 {
		return roundFloat(float32((val+1)*127/2), 0)
	}
	return roundFloat(val, 3)
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}
