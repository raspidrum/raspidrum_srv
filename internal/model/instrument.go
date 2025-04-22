package model

type Instrument struct {
	Id          int64              `yaml:"-"`
	Uid         string             `yaml:"uuid"`
	Key         string             `yaml:"key"`
	Name        string             `yaml:"name"`
	FullName    string             `yaml:"fullName,omitempty"`
	Type        string             `yaml:"type"`
	SubType     string             `yaml:"subtype"`
	Description string             `yaml:"description,omitempty"`
	Copyright   string             `yaml:"copyright,omitempty"`
	Licence     string             `yaml:"licence,omitempty"`
	Credits     string             `yaml:"credits,omitempty"`
	Tags        []string           `yaml:"tags,omitempty"`
	MidiKey     string             `yaml:"midiKey,omitempty"`
	Controls    map[string]Control `yaml:"controls"`
	Layers      map[string]Layer   `yaml:"layers,omitempty"`
}

type Control struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty"`
	Type   string `yaml:"type,omitempty" json:"type,omitempty"`
	CfgKey string `yaml:"key" json:"key"`
}

type Layer struct {
	Name       string             `yaml:"name,omitempty" json:"name,omitempty"`
	CfgMidiKey string             `yaml:"midiKey,omitempty" json:"midiKey,omitempty"`
	Controls   map[string]Control `yaml:"controls,omitempty" json:"controls,omitempty"`
}
