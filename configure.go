package configure

import (
	"io/ioutil"
	"path/filepath"
)

// ParseConfig takes filename of yml config file, and cli args (e.g. os.Args[1:])
// it calls ParseYAML, ParseEnv, ParseCommandLine
// this is in order to load config in cli -> env -> config order of precedence
func ParseConfig(filename string, args []string) (cfg interface{}, err error) {

	filefullname, _ := filepath.Abs(filename)
	data, err := ioutil.ReadFile(filefullname)
	if err != nil {
		return nil, err
	}

	err = ParseYAML(data, &cfg)
	if err != nil {
		return nil, err
	}

	err = ParseEnv(&cfg)
	if err != nil {
		return nil, err
	}

	err = ParseCommandLine(&cfg, args)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
