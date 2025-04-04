package model

type KitPreset struct {
	Kit struct {
		Uid string `yaml:"uuid"`
	} `yaml:"kit"`
	Name        string             `yaml:"name"`
	Channels    []PresetChannel    `yaml:"channels"`
	Instruments []PresetInstrument `yaml:"instruments"`
}

type PresetChannel struct {
	Key      string          `yaml:"key"`
	Name     string          `yaml:"name"`
	Controls []PresetControl `yaml:"controls"`
}

type PresetInstrument struct {
	Instrument struct {
		Uid string `yaml:"uuid"`
	} `yaml:"instrument"`
	Name       string          `yaml:"name"`
	ChannelKey string          `yaml:"channelKey"`
	MidiKey    string          `yaml:"midiKey,omitempty"`
	Controls   []PresetControl `yaml:"controls"`
	Layers     []PresetLayer   `yaml:"layers"`
}

type PresetLayer struct {
	Name     string          `yaml:"name"`
	MidiKey  string          `yaml:"midiKey,omitempty"`
	Controls []PresetControl `yaml:"controls"`
}

type PresetControl struct {
	Name   string  `yaml:"name"`
	Type   string  `yaml:"type"`
	MidiCC int     `yaml:"midiCC,omitempty"`
	Value  float32 `yaml:"value"`
}
