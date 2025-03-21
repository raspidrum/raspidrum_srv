package file

type Kit struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Licence     string   `yaml:"licence"`
	Credits     string   `yaml:"credits"`
	Url         string   `yaml:"url"`
	Tags        []string `yaml:"tags"`
}
