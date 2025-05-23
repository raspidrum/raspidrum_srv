package model

import "fmt"

type ControlOwner interface {
	HandleSetControl(control *PresetControl, value float32) error
}

type ControlIndex map[string]*PresetControl

type KitPreset struct {
	Uid         string             `yaml:"uuid,omitempty"`
	Kit         KitRef             `yaml:"kit"`
	Name        string             `yaml:"name"`
	Channels    []PresetChannel    `yaml:"channels"`
	Instruments []PresetInstrument `yaml:"instruments"`
	controls    ControlIndex
}

type KitRef struct {
	Id       int64  `yaml:"-"`
	Uid      string `yaml:"uuid"`
	Name     string `yaml:"-"`
	IsCustom bool   `yaml:"-"`
}

type PresetChannel struct {
	Key         string                   `yaml:"key"`
	Name        string                   `yaml:"name"`
	Controls    map[string]PresetControl `yaml:"controls"`
	instruments []*PresetInstrument      `yaml:"-"`
}

type PresetInstrument struct {
	Instrument InstrumentRef            `yaml:"instrument"`
	Id         int64                    `yaml:"-"`
	Name       string                   `yaml:"name"`
	ChannelKey string                   `yaml:"channelKey"`
	MidiKey    string                   `yaml:"midiKey,omitempty"`
	MidiNote   int                      `yaml:"-"`
	Controls   map[string]PresetControl `yaml:"controls"`
	Layers     map[string]PresetLayer   `yaml:"layers"`
}

type InstrumentRef struct {
	Id         int64              `yaml:"-"`
	Uid        string             `yaml:"uuid"`
	Key        string             `yaml:"-"`
	Name       string             `yaml:"-"`
	CfgMidiKey string             `yaml:"-"`
	Controls   map[string]Control `yaml:"-"`
	Layers     map[string]Layer   `yaml:"-"`
}

type PresetLayer struct {
	Name       string                   `yaml:"name,omitempty" json:"name,omitempty"`
	MidiKey    string                   `yaml:"midiKey,omitempty" json:"midiKey,omitempty"`
	CfgMidiKey string                   `yaml:"-" json:"-"`
	MidiNote   int                      `yaml:"-"`
	Controls   map[string]PresetControl `yaml:"controls" json:"controls"`
}

// CfgKey - sfz-variable key, same value as Instrument.Controls
type PresetControl struct {
	Name   string  `yaml:"name,omitempty" json:"name,omitempty"`
	Type   string  `yaml:"type" json:"type"`
	MidiCC int     `yaml:"midiCC,omitempty" json:"midiCC,omitempty"`
	CfgKey string  `yaml:"-" json:"-"`
	Value  float32 `yaml:"value" json:"value"`
	owner  ControlOwner
}

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

func (p *KitPreset) GetChannelInstrumentsByIdx(idx int) ([]*PresetInstrument, error) {
	if idx > len(p.Channels)-1 {
		return nil, fmt.Errorf("index %d out of range", idx)
	}
	ins := p.Channels[idx].instruments
	if len(ins) == 0 {
		err := p.indexInstruments()
		if err != nil {
			return nil, err
		}
		ins = p.Channels[idx].instruments
	}
	return ins, nil
}

func (p *KitPreset) GetChannelInstrumentsByKey(key string) ([]*PresetInstrument, error) {
	idx := -1
	for i, v := range p.Channels {
		if v.Key == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, fmt.Errorf("unknown channel key %s", key)
	}
	return p.GetChannelInstrumentsByIdx(idx)
}

func (p *KitPreset) indexInstruments() error {
	chnls := make(map[string]int, len(p.Channels))
	for j, c := range p.Channels {
		chnls[c.Key] = j
	}

	for i, v := range p.Instruments {
		chi, ok := chnls[v.ChannelKey]
		if !ok {
			return fmt.Errorf("instrument '%s' refs to missing channel '%s'", v.Name, v.ChannelKey)
		}
		ch := &p.Channels[chi]
		ch.instruments = append(ch.instruments, &p.Instruments[i])
	}
	return nil
}

func (p *PresetChannel) HandleSetControl(control *PresetControl, value float32) error {
	return fmt.Errorf("unimplemented")
}

func (p *PresetInstrument) HandleSetControl(control *PresetControl, value float32) error {
	return fmt.Errorf("unimplemented")
}

func (p *PresetLayer) HandleSetControl(control *PresetControl, value float32) error {
	return fmt.Errorf("unimplemented")
}

// PrepareToLoad augments preset controls and layers with data from instrument
func (p *KitPreset) PrepareToLoad(mididevs []MIDIDevice) error {
	// Initialize control index map if not exists
	if p.controls == nil {
		p.controls = make(ControlIndex)
	}

	// Counter for generating unique control IDs
	controlId := 0

	// Index channel controls
	for i := range p.Channels {
		ch := &p.Channels[i]
		for k, ctrl := range ch.Controls {
			tmpCtrl := ctrl // Create a copy to take address of
			tmpCtrl.owner = ch
			key := fmt.Sprintf("c%d", controlId)
			p.controls[key] = &tmpCtrl
			ch.Controls[k] = tmpCtrl // Update the original map
			controlId++
		}
	}

	for i, v := range p.Instruments {
		// instrument MIDI Key
		if len(v.MidiKey) > 0 {
			mkeyid, err := MapMidiKey(v.MidiKey, mididevs)
			if err != nil {
				return err
			}
			p.Instruments[i].MidiNote = mkeyid
		}

		// Index instrument controls
		instr := &p.Instruments[i]
		for k, ctrl := range v.Controls {
			ictrl, ok := v.Instrument.Controls[k]
			if !ok {
				return fmt.Errorf("not found control '%s' in instrument '%s'", k, v.Instrument.Key)
			}
			tmpCtrl := ctrl // Create a copy to take address of
			tmpCtrl.CfgKey = ictrl.CfgKey
			tmpCtrl.owner = instr
			key := fmt.Sprintf("i%d", controlId)
			p.controls[key] = &tmpCtrl
			instr.Controls[k] = tmpCtrl // Update the original map
			controlId++
		}

		// instrument layers
		for lkey, lv := range v.Layers {
			if len(lv.MidiKey) > 0 {
				mkeyid, err := MapMidiKey(lv.MidiKey, mididevs)
				if err != nil {
					return err
				}
				lv.MidiNote = mkeyid
			}

			ilrs, ok := v.Instrument.Layers[lkey]
			if !ok {
				return fmt.Errorf("not found layer '%s' in instrument '%s'", lkey, v.Instrument.Key)
			}
			lv.CfgMidiKey = ilrs.CfgMidiKey

			// Index layer controls
			for k, ctrl := range lv.Controls {
				ictrl, ok := ilrs.Controls[k]
				if !ok {
					return fmt.Errorf("not found control '%s' of layer '%s' in instrument '%s'", k, lkey, v.Instrument.Key)
				}
				tmpCtrl := ctrl // Create a copy to take address of
				tmpCtrl.CfgKey = ictrl.CfgKey
				tmpCtrl.owner = &lv
				ctrlKey := fmt.Sprintf("l%d", controlId)
				p.controls[ctrlKey] = &tmpCtrl
				lv.Controls[k] = tmpCtrl // Update the original map
				controlId++
			}
			instr.Layers[lkey] = lv
		}
	}
	return nil
}

// MapMidiKey maps a MIDI key string to its numeric value using the provided MIDI devices
func MapMidiKey(mkey string, mdevs []MIDIDevice) (int, error) {
	devlist := make([]string, len(mdevs))
	for i, d := range mdevs {
		kmap, err := d.GetKeysMapping()
		if err != nil {
			return 0, fmt.Errorf("failed get MIDI Keys mapping for device %s: %w", d.Name(), err)
		}
		devlist[i] = d.Name()
		midiId, ok := kmap[mkey]
		if ok {
			return midiId, nil
		}
	}
	return 0, fmt.Errorf("MIDI devices %s doen't have mapping for MIDI Key %s", devlist, mkey)
}

// MIDIDevice interface defines methods required for MIDI device operations
type MIDIDevice interface {
	Name() string
	GetKeysMapping() (map[string]int, error)
}
