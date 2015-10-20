# Configure

`configure` Golang library generates configuration tools based on a struct
definition with tags. It can read a configuration for a struct
from YAML, environment variables and command line.


### Usage 

The latest can be seen if you run 

```
godoc github.com/gravitational/configure
```

But here's a quickstart: Define a sample structure, for examlpe:
```go

	 type Config struct {
	   StringVar   string              `env:"TEST_STRING_VAR" cli:"string" yaml:"string"`
	   BoolVar     bool                `env:"TEST_BOOL_VAR" cli:"bool" yaml:"bool"`
	   IntVar      int                 `env:"TEST_INT_VAR" cli:"int" yaml:"int"`
	   HexVar      hexType             `env:"TEST_HEX_VAR" cli:"hex" yaml:"hex"`
	   MapVar      map[string]string   `env:"TEST_MAP_VAR" cli:"map" yaml:"map,flow"`
	   SliceMapVar []map[string]string `env:"TEST_SLICE_MAP_VAR" cli:"slice" yaml:"slice,flow"`
	}
```

Then you can query the environment and populate that structure like this:

```go
	import (
	   "os"
	   "github.com/gravitational/configure"
	)

	func main() {
	   var cfg Config
	   // parse environment variables
	   err := configure.ParseEnv(&cfg)
	   // parse YAML
	   err = configure.ParseYAML(&cfg)
	   // parse command line arguments
	   err = configure.ParseCommandLine(&cfg, os.Ars[1:])
	}
```
