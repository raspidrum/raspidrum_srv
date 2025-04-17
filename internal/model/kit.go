package model

type Kit struct {
	Id          int64
	Uid         string   `yaml:"UUID"`
	Name        string   `yaml:"name"`
	IsCustom    bool     `yaml:"-"`
	Description string   `yaml:"description,omitempty"`
	Copyright   string   `yaml:"copyright,omitempty"`
	Licence     string   `yaml:"licence,omitempty"`
	Credits     string   `yaml:"credits,omitempty"`
	Url         string   `yaml:"url,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}
