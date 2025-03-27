package file

type Kit struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Licence     string   `yaml:"licence,omitempty"`
	Credits     string   `yaml:"credits,omitempty"`
	Url         string   `yaml:"url,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}
