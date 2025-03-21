package file

type Instrument struct {
	InstrumentKey string   `yaml:"instrumentKey"`
	Name          string   `yaml:"name"`
	FullName      string   `yaml:"fullName"`
	Type          string   `yaml:"type"`
	SubType       string   `yaml:"subType"`
	Description   string   `yaml:"description"`
	Copyright     string   `yaml:"copyright"`
	Licence       string   `yaml:"licence"`
	Credits       string   `yaml:"credits"`
	Tags          []string `yaml:"tags"`
	MidiKey       string   `yaml:"midiKey"`
	Controls      Controls `yaml:"controls"`
	Layers        struct {
		Name     string   `yaml:"name"`
		MidiKey  string   `yaml:"midiKey"`
		Controls Controls `yaml:"controls"`
	}
}

type Controls struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Key  string `yaml:"key"`
}
