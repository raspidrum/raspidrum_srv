package model

// TODO: make maps for Controls and Layers
type Instrument struct {
	Id            int64
	Uid           string    `yaml:"UUID"`
	InstrumentKey string    `yaml:"instrumentKey"`
	Name          string    `yaml:"name"`
	FullName      string    `yaml:"fullName,omitempty"`
	Type          string    `yaml:"type"`
	SubType       string    `yaml:"subtype"`
	Description   string    `yaml:"description,omitempty"`
	Copyright     string    `yaml:"copyright,omitempty"`
	Licence       string    `yaml:"licence,omitempty"`
	Credits       string    `yaml:"credits,omitempty"`
	Tags          []string  `yaml:"tags,omitempty"`
	MidiKey       string    `yaml:"midiKey,omitempty"`
	Controls      []Control `yaml:"controls"`
	Layers        []Layer   `yaml:"layers,omitempty"`
}

type Control struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
	Key  string `yaml:"key" json:"key"`
}

type Layer struct {
	Name     string    `yaml:"name" json:"name"`
	MidiKey  string    `yaml:"midiKey,omitempty" json:"midiKey,omitempty"`
	Controls []Control `yaml:"controls,omitempty" json:"controls,omitempty"`
}
