package model

import (
	"fmt"
	"log/slog"
)

type controlRef struct {
	channel *PresetChannel
	control *PresetControl
}

type KitPreset struct {
	Id          int64                 `yaml:"-"`
	Uid         string                `yaml:"uuid,omitempty"`
	Kit         KitRef                `yaml:"kit"`
	Name        string                `yaml:"name"`
	Channels    []PresetChannel       `yaml:"channels"`
	Instruments []PresetInstrument    `yaml:"instruments"`
	controls    map[string]controlRef // key - control.Key
}

type KitRef struct {
	Id       int64  `yaml:"-"`
	Uid      string `yaml:"uuid"`
	Name     string `yaml:"-"`
	IsCustom bool   `yaml:"-"`
}

type PresetChannel struct {
	Key         string              `yaml:"key"`
	Name        string              `yaml:"name"`
	Controls    ControlMap          `yaml:"controls"`
	instruments []*PresetInstrument `yaml:"-"`
}

type PresetInstrument struct {
	Instrument InstrumentRef          `yaml:"instrument"`
	Id         int64                  `yaml:"-"`
	Name       string                 `yaml:"name"`
	ChannelKey string                 `yaml:"channelKey"`
	MidiKey    string                 `yaml:"midiKey,omitempty"`
	MidiNote   int                    `yaml:"-"`
	Controls   ControlMap             `yaml:"controls"`
	Layers     map[string]PresetLayer `yaml:"layers"`
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
	Name       string     `yaml:"name,omitempty" json:"name,omitempty"`
	MidiKey    string     `yaml:"midiKey,omitempty" json:"midiKey,omitempty"`
	CfgMidiKey string     `yaml:"-" json:"-"`
	MidiNote   int        `yaml:"-"`
	Controls   ControlMap `yaml:"controls" json:"controls"`
}

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

// make instrument index for each channel
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

// PrepareToLoad augments preset controls and layers with data from instrument
func (p *KitPreset) PrepareToLoad(mididevs []MIDIDevice) error {
	cnlsIndex := make(map[string]*PresetChannel, len(p.Channels))
	if err := p.indexInstruments(); err != nil {
		return err
	}

	// Initialize control index map if not exists
	if p.controls == nil {
		p.controls = make(map[string]controlRef)
	}

	// Counter for generating unique control IDs
	channelIdx := 0

	// Index channel controls
	for i := range p.Channels {
		ch := &p.Channels[i]
		cnlsIndex[ch.Key] = ch
		instrCount := len(ch.instruments)
		// Index channel controls
		var hasPan bool
		for k, ctrl := range ch.Controls {
			// Link instrument volume or pan control for single instrument with MIDI CC
			if instrCount == 1 && (ctrl.Type == CtrlVolume || ctrl.Type == CtrlPan) {
				if ictrl, ok := ch.instruments[0].Controls.FindControlByType(ctrl.Type); ok {
					if ictrl.MidiCC != 0 {
						ctrl.linkedTo = append(ctrl.linkedTo, ictrl)
						ictrl.linkedWith = ctrl
					}
				}
			}
			// TODO: pan always linked to instrument pan.
			// In case many instruments in channel pan is virtual and linked with pan of all instruments in channel
			// In case one instrument in channel pan regulated in instrument
			if ctrl.Type == CtrlPan {
				hasPan = true
			}
			ctrl.owner = ch
			key := fmt.Sprintf("c%d%s", channelIdx, k)
			ctrl.Key = key
			p.controls[key] = controlRef{channel: ch, control: ctrl}
		}
		if !hasPan && instrCount == 1 {
			// add control for pan linked to instrument pan if not exists in channel controls
			if ictrl, ok := ch.instruments[0].Controls.FindControlByType(CtrlPan); ok {
				if ictrl.MidiCC != 0 {
					ctrl := &PresetControl{
						Name:  ictrl.Name,
						Type:  ictrl.Type,
						owner: ch,
					}
					key := fmt.Sprintf("c%d%s", channelIdx, CtrlPan)
					ctrl.Key = key
					ctrl.linkedTo = append(ctrl.linkedTo, ictrl)
					ictrl.linkedWith = ctrl
					ch.Controls[CtrlPan] = ctrl
					p.controls[key] = controlRef{channel: ch, control: ctrl}
				}
			}
		}
		channelIdx++
	}

	instrumentIdx := 0
	for i := range p.Instruments {
		instr := &p.Instruments[i]
		ch := cnlsIndex[instr.ChannelKey]
		// instrument MIDI Key
		if len(instr.MidiKey) > 0 {
			mkeyid, err := MapMidiKey(instr.MidiKey, mididevs)
			if err != nil {
				return err
			}
			instr.MidiNote = mkeyid
		}

		for k, ctrl := range instr.Controls {
			// find control declaration in instrument
			ctrlMeta, ok := instr.Instrument.Controls[k]
			if !ok {
				return fmt.Errorf("not found control '%s' in instrument '%s'", k, instr.Instrument.Key)
			}
			ctrl.CfgKey = ctrlMeta.CfgKey

			// link with layer controls if instrument has multiple layers
			if ctrl.MidiCC == 0 && (ctrl.Type == CtrlVolume || ctrl.Type == CtrlPan) {
				for _, lr := range instr.Layers {
					// find same type control in layer
					lrCtrl, ok := lr.Controls.FindControlByType(ctrl.Type)
					if ok {
						ctrl.linkedTo = append(ctrl.linkedTo, lrCtrl)
						lrCtrl.linkedWith = ctrl
					}
				}
			}

			ctrl.owner = instr
			// Index instrument controls
			key := fmt.Sprintf("i%d%s", instrumentIdx, k)
			ctrl.Key = key
			p.controls[key] = controlRef{channel: ch, control: ctrl}
		}

		layerIdx := 0
		for lkey, lv := range instr.Layers {
			if len(lv.MidiKey) > 0 {
				mkeyid, err := MapMidiKey(lv.MidiKey, mididevs)
				if err != nil {
					return err
				}
				lv.MidiNote = mkeyid
			}

			lrMeta, ok := instr.Instrument.Layers[lkey]
			if !ok {
				return fmt.Errorf("not found layer '%s' in instrument '%s'", lkey, instr.Instrument.Key)
			}
			lv.CfgMidiKey = lrMeta.CfgMidiKey

			for k, ctrl := range lv.Controls {
				ictrl, ok := lrMeta.Controls[k]
				if !ok {
					return fmt.Errorf("not found control '%s' of layer '%s' in instrument '%s'", k, lkey, instr.Instrument.Key)
				}
				ctrl.CfgKey = ictrl.CfgKey
				ctrl.owner = &lv
				// Index layer controls
				key := fmt.Sprintf("i%dl%d%s", instrumentIdx, layerIdx, k)
				ctrl.Key = key
				p.controls[key] = controlRef{channel: ch, control: ctrl}
			}
			instr.Layers[lkey] = lv
			layerIdx++
		}
		instrumentIdx++
	}
	return nil
}

func (p *KitPreset) SetControlValue(controlKey string, value float32, csetter SamplerControlSetter) error {
	// find control by key
	if p.controls == nil {
		return fmt.Errorf("controls not initialized")
	}
	ctrl, ok := p.controls[controlKey]
	if !ok {
		return fmt.Errorf("control '%s' not found", controlKey)
	}
	return ctrl.control.SetValue(value, ctrl.channel.Key, csetter)
}

// Volume in channel sets by Sampler API
// Pan in channel virtual (in case many instruments in channel).
// In case one instrument in channel, pan is linked to instrument pan. Pan will be regulated in instrument
// Other controls except volume and pan are not supported in channel
func (c *PresetChannel) HandleControlValue(channelKey string, control *PresetControl, value float32, csetter SamplerControlSetter) error {
	slog.Debug("HandleControlValue", "control", control, "value", value)
	if control.Type == CtrlVolume {
		control.Value = value
		if control.MidiCC == 0 {
			return csetter.SetChannelVolume(c.Key, control.Value)
		} else {
			return csetter.SendChannelMidiCC(c.Key, control.MidiCC, control.Value)
		}
	}
	// TODO: pan handle
	return nil
}

// If channel has linked control (linked to corresponding instrument control)
// then substitute instrument control as channel control, but with channel control key
func (c *PresetChannel) GetControls() func(func(*PresetControl) bool) {
	return func(yield func(*PresetControl) bool) {
		for _, c := range c.Controls {
			var ctrl *PresetControl
			if len(c.linkedTo) > 0 {
				// channel control MAY be linked only with one instrument control
				ctrl = c.linkedTo[0]
			} else {
				ctrl = c
			}
			if !yield(ctrl) {
				return
			}
		}
	}
}

// Volume and pan can be virtual or regulated by MIDI CC accordingly
// If virtual, then regulated by linked layer controls
// Else regulated by MIDI CC
// Other controls always regulated by MIDI CC
func (p *PresetInstrument) HandleControlValue(channelKey string, control *PresetControl, value float32, csetter SamplerControlSetter) error {
	slog.Debug("HandleControlValue", "control", control, "value", value)
	return nil
}

func (c *PresetInstrument) GetControls() func(func(*PresetControl) bool) {
	return func(yield func(*PresetControl) bool) {
		for _, c := range c.Controls {
			// don't yield control if it is linked with channel control
			// it's be yielded by channel
			if c.linkedWith == nil {
				if !yield(c) {
					return
				}
			}
		}
	}
}

func (p *PresetLayer) HandleControlValue(channelKey string, control *PresetControl, value float32, csetter SamplerControlSetter) error {
	slog.Debug("HandleControlValue", "control", control, "value", value)
	if control.Type == CtrlVolume || control.Type == CtrlPan {
		control.Value = value
		if control.MidiCC != 0 {
			return csetter.SendChannelMidiCC(channelKey, control.MidiCC, control.Value)
		}
	}
	return nil
}

// Layer can be in instrument with virtual controls of volume or pan
// In that case its required to calculate correction value by linked instrument control
func (c *PresetLayer) GetControls() func(func(*PresetControl) bool) {
	return func(yield func(*PresetControl) bool) {
		for _, c := range c.Controls {
			if !yield(c) {
				return
			}
		}
	}
}
