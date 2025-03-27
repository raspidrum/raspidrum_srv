package file

type Instrument struct {
	Id            int64
	Uid           string
	InstrumentKey string     `yaml:"instrumentKey"`
	Name          string     `yaml:"name"`
	FullName      string     `yaml:"fullName,omitempty"`
	Type          string     `yaml:"type"`
	SubType       string     `yaml:"subtype"`
	Description   string     `yaml:"description,omitempty"`
	Copyright     string     `yaml:"copyright,omitempty"`
	Licence       string     `yaml:"licence,omitempty"`
	Credits       string     `yaml:"credits,omitempty"`
	Tags          []string   `yaml:"tags,omitempty"`
	MidiKey       string     `yaml:"midiKey,omitempty"`
	Controls      []Controls `yaml:"controls"`
	Layers        []struct {
		Name     string     `yaml:"name"`
		MidiKey  string     `yaml:"midiKey,omitempty"`
		Controls []Controls `yaml:"controls,omitempty"`
	} `yaml:"layers,omitempty"`
}

type Controls struct {
	Name string `yaml:"name"`
	Type string `yaml:"type,omitempty"`
	Key  string `yaml:"key"`
}
