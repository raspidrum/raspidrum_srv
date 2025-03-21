package file

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type root struct {
	kit        *Kit        `yaml:"kit"`
	instrument *Instrument `yaml:"instrument"`
}

func ParseYAML(data []byte) (interface{}, error) {
	var root root
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	switch {
	case root.instrument != nil:
		return root.instrument, nil
	case root.kit != nil:
		return root.kit, nil
	default:
		return nil, errors.New("unknown YAML format")
	}
}
