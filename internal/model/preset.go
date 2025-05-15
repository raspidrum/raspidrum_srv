package model

type KitPreset struct {
	Uid         string             `yaml:"uuid,omitempty"`
	Kit         KitRef             `yaml:"kit"`
	Name        string             `yaml:"name"`
	Channels    []PresetChannel    `yaml:"channels"`
	Instruments []PresetInstrument `yaml:"instruments"`
}

type KitRef struct {
	Id       int64  `yaml:"-"`
	Uid      string `yaml:"uuid"`
	Name     string `yaml:"-"`
	IsCustom bool   `yaml:"-"`
}

type PresetChannel struct {
	Key      string                   `yaml:"key"`
	Name     string                   `yaml:"name"`
	Controls map[string]PresetControl `yaml:"controls"`
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
}
