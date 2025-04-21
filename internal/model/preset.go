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
	Name       string                   `yaml:"name"`
	ChannelKey string                   `yaml:"channelKey"`
	MidiKey    string                   `yaml:"midiKey,omitempty"`
	Controls   map[string]PresetControl `yaml:"controls"`
	Layers     map[string]PresetLayer   `yaml:"layers"`
}

type InstrumentRef struct {
	Id       int64              `yaml:"-"`
	Uid      string             `yaml:"uuid"`
	Key      string             `yaml:"-"`
	Name     string             `yaml:"-"`
	MidiKey  string             `yaml:"-"`
	Controls map[string]Control `yaml:"-"`
	Layers   map[string]Layer   `yaml:"-"`
}

type PresetLayer struct {
	Name     string                   `yaml:"name,omitempty" json:"name,omitempty"`
	MidiKey  string                   `yaml:"midiKey,omitempty" json:"midiKey,omitempty"`
	Controls map[string]PresetControl `yaml:"controls" json:"controls"`
}

type PresetControl struct {
	Name   string  `yaml:"name,omitempty" json:"name,omitempty"`
	Type   string  `yaml:"type" json:"type"`
	MidiCC int     `yaml:"midiCC,omitempty" json:"midiCC,omitempty"`
	Value  float32 `yaml:"value" json:"value"`
}
