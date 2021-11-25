package configure

import (
	"io/ioutil"

	yaml "github.com/hzakher/yaml/v2"
)

const Tag = "config"

// ParseYAML parses yaml-encoded byte string into the struct
// passed to the function.
// EnableTemplating() argument allows to treat configuration file as a template
// for example, it will support {{env "VAR"}} - that will substitute
// environment variable "VAR" and pass it to YAML file parser
func ParseYAML(data []byte, cfg interface{}) error {
	yaml.SetTag(Tag)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}

func SaveYAML(filename string, cfg interface{}) error {
	yaml.SetTag(Tag)
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}
